package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLICommands(t *testing.T) {
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

	tempDir, err := os.MkdirTemp("", "sova-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name          string
		args          []string
		expectedOut   string
		expectedError bool
	}{
		{
			name:          "Version command",
			args:          []string{"version"},
			expectedOut:   "Sova CLI vdev",
			expectedError: false,
		},
		{
			name:          "Help command",
			args:          []string{"help"},
			expectedOut:   "Available Commands:",
			expectedError: false,
		},
		{
			name:          "Init command with project name",
			args:          []string{"init", "test-project", "--type", "cli", "--use-zap=false", "--use-postgres=false", "--use-redis=false", "--use-rabbitmq=false"},
			expectedOut:   "Project initialized successfully",
			expectedError: false,
		},
		{
			name:          "Invalid command",
			args:          []string{"invalid-command"},
			expectedOut:   "",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tc.args...)
			cmd.Dir = tempDir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !strings.Contains(output, tc.expectedOut) {
				t.Errorf("Expected output containing %q, got %q", tc.expectedOut, output)
			}
		})
	}
}

func TestCLIFlags(t *testing.T) {
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

	tempDir, err := os.MkdirTemp("", "sova-cli-flags-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name          string
		args          []string
		flags         []string
		expectedOut   string
		expectedError bool
	}{
		{
			name:          "Init with type flag",
			args:          []string{"init", "test-project"},
			flags:         []string{"--type", "cli", "--use-zap=false", "--use-postgres=false", "--use-redis=false", "--use-rabbitmq=false"},
			expectedOut:   "Project initialized successfully",
			expectedError: false,
		},
		{
			name:          "Init with invalid type",
			args:          []string{"init", "test-project"},
			flags:         []string{"--type", "nonexistent"},
			expectedOut:   "Error: unsupported project type: nonexistent",
			expectedError: true,
		},
		{
			name:          "Version with json flag",
			args:          []string{"version"},
			flags:         []string{"--json"},
			expectedOut:   "{",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmdArgs := append(tc.args, tc.flags...)
			cmd := exec.Command(binaryPath, cmdArgs...)
			cmd.Dir = tempDir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tc.expectedError {
				if err == nil && !strings.Contains(output, tc.expectedOut) {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !strings.Contains(output, tc.expectedOut) {
				t.Errorf("Expected output containing %q, got %q", tc.expectedOut, output)
			}
		})
	}
}
