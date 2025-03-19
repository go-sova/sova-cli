package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-sova/sova-cli/cmd"
	"github.com/go-sova/sova-cli/templates"
)

func TestTemplateRetrievalExtended(t *testing.T) {
	loader := templates.NewTemplateLoader()
	tests := []struct {
		name         string
		category     string
		templateName string
		wantErr      bool
	}{
		{
			name:         "Valid template",
			category:     "cli",
			templateName: "main.tpl",
			wantErr:      false,
		},
		{
			name:         "Invalid category",
			category:     "nonexistent",
			templateName: "main.tpl",
			wantErr:      true,
		},
		{
			name:         "Invalid template name",
			category:     "cli",
			templateName: "nonexistent.tpl",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := templates.GetTemplatePath(tt.category, tt.templateName)
			tmpl, err := loader.LoadTemplate(templatePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tmpl == nil {
				t.Error("LoadTemplate() returned nil template for valid path")
			}
		})
	}
}

func TestProjectInitializationExtended(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sova-project-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name       string
		projectDir string
		template   string
		wantFiles  []string
		wantErr    bool
	}{
		{
			name:       "Basic project initialization",
			projectDir: "test-project",
			template:   "cli",
			wantFiles:  []string{"main.go", "go.mod", "README.md"},
			wantErr:    false,
		},
		{
			name:       "Project with existing directory",
			projectDir: "existing-project",
			template:   "cli",
			wantFiles:  []string{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath := filepath.Join(tempDir, tt.projectDir)

			if tt.name == "Project with existing directory" {
				if err := os.MkdirAll(projectPath, 0755); err != nil {
					t.Fatalf("Failed to create existing directory: %v", err)
				}
			}

			// Test project initialization
			err := cmd.InitializeProject(projectPath, tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if expected files exist
				for _, file := range tt.wantFiles {
					path := filepath.Join(projectPath, file)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("Expected file %s does not exist", file)
					}
				}
			}
		})
	}
}

func TestTemplateValidationExtended(t *testing.T) {
	tests := []struct {
		name     string
		template string
		wantErr  bool
	}{
		{
			name:     "Valid template name",
			template: "cli",
			wantErr:  false,
		},
		{
			name:     "Empty template name",
			template: "",
			wantErr:  true,
		},
		{
			name:     "Invalid template name",
			template: "../invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.ValidateTemplate(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
