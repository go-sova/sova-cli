package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-sova/sova-cli/pkg/utils"
)

//go:embed cli/* api/*
var TemplateFS embed.FS

// TemplateLoader handles loading templates from the embedded filesystem
type TemplateLoader struct {
	fs     fs.FS
	logger *utils.Logger
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader() *TemplateLoader {
	return &TemplateLoader{
		fs:     TemplateFS,
		logger: utils.NewLoggerWithPrefix(utils.Info, "TemplateLoader"),
	}
}

func (l *TemplateLoader) SetLogger(logger *utils.Logger) {
	l.logger = logger
}

// LoadTemplate loads a template by name from the embedded filesystem
func (l *TemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	// If the template name already includes a category prefix (e.g. "api/env.tpl"),
	// try loading it directly
	content, err := fs.ReadFile(l.fs, name)
	if err == nil {
		tmpl, err := template.New(filepath.Base(name)).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		return tmpl, nil
	}

	// If direct loading fails, try each category as a fallback
	categories := []string{"cli", "api"}
	for _, category := range categories {
		if tmpl, err := l.LoadTemplateFromCategory(category, name); err == nil {
			return tmpl, nil
		}
	}

	return nil, fmt.Errorf("template not found: %s", name)
}

// LoadTemplateFromCategory loads a template from a specific category
func (l *TemplateLoader) LoadTemplateFromCategory(category, name string) (*template.Template, error) {
	templatePath := filepath.Join(category, name)
	content, err := fs.ReadFile(l.fs, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(name)).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	return tmpl, nil
}

// FileGenerator handles generating files from templates
type FileGenerator struct {
	loader *TemplateLoader
	logger *utils.Logger
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(loader *TemplateLoader) *FileGenerator {
	return &FileGenerator{
		loader: loader,
		logger: utils.NewLoggerWithPrefix(utils.Info, "FileGenerator"),
	}
}

func (g *FileGenerator) SetLogger(logger *utils.Logger) {
	g.logger = logger
}

// GenerateFile generates a file from a template
func (g *FileGenerator) GenerateFile(templateName, outputPath string, data interface{}) error {
	g.logger.Debug("Generating file %s from template %s", outputPath, templateName)

	// Create the directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Load the template
	tmpl, err := g.loader.LoadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// GetTemplateFS returns the embedded filesystem containing all templates
func GetTemplateFS() fs.FS {
	return TemplateFS
}

// GetTemplatePath returns the path to a specific template within the embedded filesystem
func GetTemplatePath(category, name string) string {
	return filepath.Join(category, name)
}
