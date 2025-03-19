package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectInitialization(t *testing.T) {
	// Create a temporary directory for the binary
	tempBinDir := filepath.Join(os.TempDir(), "sova-cli-test-bin")
	if err := os.MkdirAll(tempBinDir, 0755); err != nil {
		t.Fatalf("Failed to create binary directory: %v", err)
	}
	defer os.RemoveAll(tempBinDir)

	// Build the binary
	binaryPath := filepath.Join(tempBinDir, "sova")
	buildCmd := exec.Command("go", "build", "-o", binaryPath)
	buildCmd.Dir = ".." // Run from the project root
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "sova-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name        string
		projectName string
		projectType string
		wantDirs    []string
		wantFiles   []string
		wantErr     bool
	}{
		{
			name:        "Basic project creation",
			projectName: "test-project",
			projectType: "cli",
			wantDirs: []string{
				"cmd",
			},
			wantFiles: []string{
				"main.go",
				"go.mod",
				"README.md",
				"cmd/root.go",
				"cmd/version.go",
				".gitignore",
			},
			wantErr: false,
		},
		{
			name:        "API project creation",
			projectName: "api-project",
			projectType: "api",
			wantDirs: []string{
				"cmd",
				"internal",
				"pkg",
				"api",
			},
			wantFiles: []string{
				"main.go",
				"go.mod",
				"README.md",
			},
			wantErr: false,
		},
		{
			name:        "Invalid project type",
			projectName: "invalid-project",
			projectType: "nonexistent",
			wantDirs:    []string{},
			wantFiles:   []string{},
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectDir := filepath.Join(tempDir, tc.projectName)
			cmd := exec.Command(binaryPath, "init", tc.projectName, "--type", tc.projectType,
				"--use-zap=false", "--use-postgres=false", "--use-redis=false", "--use-rabbitmq=false")
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()

			if tc.wantErr {
				if err == nil && !strings.Contains(string(output), "Error: unsupported project type") {
					t.Errorf("Expected error but got none. Output: %s", string(output))
				}
				return
			} else if err != nil {
				t.Fatalf("Failed to run command: %v\nOutput: %s", err, string(output))
			}

			for _, dir := range tc.wantDirs {
				dirPath := filepath.Join(projectDir, dir)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					t.Errorf("Expected directory %s does not exist", dir)
				}
			}

			for _, file := range tc.wantFiles {
				filePath := filepath.Join(projectDir, file)
				info, err := os.Stat(filePath)
				if os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", file)
					continue
				}
				if info.Size() == 0 {
					t.Errorf("File %s exists but is empty", file)
				}
			}

			if !tc.wantErr {
				buildCmd := exec.Command("go", "build", "./...")
				buildCmd.Dir = projectDir
				if output, err := buildCmd.CombinedOutput(); err != nil {
					t.Errorf("Project failed to build: %v\nOutput: %s", err, string(output))
				}
			}
		})
	}
}
