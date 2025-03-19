package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sova-utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("Directory operations", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "test-dir")
		nestedDirPath := filepath.Join(dirPath, "nested")

		err := os.MkdirAll(nestedDirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create nested directory: %v", err)
		}

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("Directory was not created")
		}

		if _, err := os.Stat(nestedDirPath); os.IsNotExist(err) {
			t.Error("Nested directory was not created")
		}
	})

	t.Run("File operations", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test.txt")
		content := []byte("Hello, World!")

		err := os.WriteFile(filePath, content, 0644)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(readContent) != string(content) {
			t.Errorf("File content mismatch. Want %q, got %q", content, readContent)
		}
	})

	t.Run("File permissions", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "permissions.txt")
		content := []byte("Test permissions")

		err := os.WriteFile(filePath, content, 0600)
		if err != nil {
			t.Fatalf("Failed to create file with permissions: %v", err)
		}

		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}

		if info.Mode().Perm() != 0600 {
			t.Errorf("File permissions mismatch. Want %o, got %o", 0600, info.Mode().Perm())
		}
	})

	t.Run("File deletion", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "to-delete.txt")
		content := []byte("Delete me")

		err := os.WriteFile(filePath, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}

		err = os.Remove(filePath)
		if err != nil {
			t.Fatalf("Failed to delete file: %v", err)
		}

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Error("File was not deleted")
		}
	})
}

func TestPathOperations(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		wantDir       string
		wantBase      string
		wantExtension string
	}{
		{
			name:          "Simple file",
			path:          "file.txt",
			wantDir:       ".",
			wantBase:      "file",
			wantExtension: ".txt",
		},
		{
			name:          "Nested file",
			path:          "dir/subdir/file.go",
			wantDir:       "dir/subdir",
			wantBase:      "file",
			wantExtension: ".go",
		},
		{
			name:          "Hidden file",
			path:          ".config",
			wantDir:       ".",
			wantBase:      ".config",
			wantExtension: "",
		},
		{
			name:          "Multiple extensions",
			path:          "archive.tar.gz",
			wantDir:       ".",
			wantBase:      "archive.tar",
			wantExtension: ".gz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Dir(tc.path)
			if dir != tc.wantDir {
				t.Errorf("Directory mismatch for %s. Want %s, got %s", tc.path, tc.wantDir, dir)
			}

			base := filepath.Base(tc.path)
			if tc.name == "Hidden file" {
				if base != tc.wantBase {
					t.Errorf("Base name mismatch for %s. Want %s, got %s", tc.path, tc.wantBase, base)
				}
				if tc.wantExtension != "" {
					t.Errorf("Extension mismatch for %s. Want %s, got %s", tc.path, tc.wantExtension, "")
				}
			} else {
				ext := filepath.Ext(tc.path)
				baseWithoutExt := base[:len(base)-len(ext)]

				if baseWithoutExt != tc.wantBase {
					t.Errorf("Base name mismatch for %s. Want %s, got %s", tc.path, tc.wantBase, baseWithoutExt)
				}

				if ext != tc.wantExtension {
					t.Errorf("Extension mismatch for %s. Want %s, got %s", tc.path, tc.wantExtension, ext)
				}
			}
		})
	}
}
