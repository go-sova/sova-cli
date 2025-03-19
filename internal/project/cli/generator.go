package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sova/sova-cli/pkg/questions"
	"github.com/go-sova/sova-cli/pkg/utils"
	"github.com/go-sova/sova-cli/templates"
)

type CLIProjectGenerator struct {
	ProjectName    string
	ProjectDir     string
	Answers        *questions.ProjectAnswers
	templateLoader *templates.TemplateLoader
	fileGenerator  *templates.FileGenerator
	logger         *utils.Logger
}

func NewCLIProjectGenerator(projectName, projectDir string, answers *questions.ProjectAnswers) *CLIProjectGenerator {
	loader := templates.NewTemplateLoader()
	return &CLIProjectGenerator{
		ProjectName:    projectName,
		ProjectDir:     projectDir,
		Answers:        answers,
		templateLoader: loader,
		fileGenerator:  templates.NewFileGenerator(loader),
		logger:         utils.NewLoggerWithPrefix(utils.Info, "CLIProjectGenerator"),
	}
}

func (g *CLIProjectGenerator) SetLogger(logger *utils.Logger) {
	g.logger = logger
	g.templateLoader.SetLogger(logger)
	g.fileGenerator.SetLogger(logger)
}

func (g *CLIProjectGenerator) Generate() (map[string]string, []string, error) {
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
		"docs",
		"scripts",
		"test",
		"cmd/root",
		"internal/commands",
		"internal/config",
	}

	fileTemplates := map[string]string{
		"main.go":                   "cli/main.tpl",
		"go.mod":                    "cli/go-mod.tpl",
		"README.md":                 "cli/readme.tpl",
		"cmd/root.go":               "cli/root.tpl",
		"cmd/version.go":            "cli/version.tpl",
		"internal/commands/cmd.go":  "cli/commands.tpl",
		"internal/config/config.go": "cli/config.tpl",
		"internal/utils/utils.go":   "cli/utils.tpl",
		".gitignore":                "cli/gitignore.tpl",
	}

	files := make(map[string]string)
	for filePath, templateName := range fileTemplates {
		files[filePath] = templateName
	}

	return files, dirs, nil
}

func (g *CLIProjectGenerator) WriteFiles(files map[string]string) error {
	for filePath, templateName := range files {
		fullPath := filepath.Join(g.ProjectDir, filePath)

		data := map[string]interface{}{
			"ProjectName":        g.ProjectName,
			"ProjectDescription": "A CLI application with clean architecture",
			"ModuleName":         g.ProjectName,
			"GoVersion":          "1.21",
		}

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}

		if err := g.fileGenerator.GenerateFile(templateName, fullPath, data); err != nil {
			return fmt.Errorf("failed to generate file %s from template %s: %v", filePath, templateName, err)
		}

		fmt.Printf("Created file: %s\n", fullPath)
	}

	return nil
}
