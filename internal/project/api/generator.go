package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sova/sova-cli/pkg/questions"
	"github.com/go-sova/sova-cli/pkg/utils"
	"github.com/go-sova/sova-cli/templates"
)

type APIProjectGenerator struct {
	ProjectName    string
	ProjectDir     string
	Answers        *questions.ProjectAnswers
	templateLoader *templates.TemplateLoader
	fileGenerator  *templates.FileGenerator
	logger         *utils.Logger
}

func NewAPIProjectGenerator(projectName, projectDir string, answers *questions.ProjectAnswers) *APIProjectGenerator {
	loader := templates.NewTemplateLoader()
	return &APIProjectGenerator{
		ProjectName:    projectName,
		ProjectDir:     projectDir,
		Answers:        answers,
		templateLoader: loader,
		fileGenerator:  templates.NewFileGenerator(loader),
		logger:         utils.NewLoggerWithPrefix(utils.Info, "APIProjectGenerator"),
	}
}

func (g *APIProjectGenerator) SetLogger(logger *utils.Logger) {
	g.logger = logger
	g.templateLoader.SetLogger(logger)
	g.fileGenerator.SetLogger(logger)
}

func (g *APIProjectGenerator) Generate() (map[string]string, []string, error) {
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
		"api",
		"internal/server",
		"internal/service",
		"internal/handlers",
		"internal/middleware",
		"internal/routes",
	}

	fileTemplates := map[string]string{
		"main.go":                       "api/main.tpl",
		"README.md":                     "api/readme.tpl",
		"cmd/main.go":                   "api/main.tpl",
		"internal/server/server.go":     "api/server.tpl",
		"internal/routes/routes.go":     "api/routes.tpl",
		"internal/service/service.go":   "api/service-init.tpl",
		"internal/handlers/handlers.go": "api/handlers.tpl",
		"internal/middleware/auth.go":   "api/middleware.tpl",
		".env":                          "api/env.tpl",
		"docker-compose.yml":            "api/docker-compose.tpl",
		"Dockerfile":                    "api/dockerfile.tpl",
		"go.mod":                        "api/go-mod.tpl",
		".gitignore":                    "api/gitignore.tpl",
	}

	if g.Answers.UseZap {
		fileTemplates["internal/middleware/logging.go"] = "api/logging.tpl"
	}

	if g.Answers.UsePostgres {
		fileTemplates["internal/service/postgres.go"] = "api/postgres.tpl"
	}

	if g.Answers.UseRedis {
		fileTemplates["internal/service/redis.go"] = "api/redis.tpl"
	}

	if g.Answers.UseRabbitMQ {
		fileTemplates["internal/service/rabbitmq.go"] = "api/rabbitmq.tpl"
	}

	files := make(map[string]string)
	for filePath, templateName := range fileTemplates {
		files[filePath] = templateName
	}

	return files, dirs, nil
}

func (g *APIProjectGenerator) WriteFiles(files map[string]string) error {
	for filePath, templateName := range files {
		fullPath := filepath.Join(g.ProjectDir, filePath)

		data := map[string]interface{}{
			"ProjectName":        g.ProjectName,
			"ProjectDescription": "A Go API with clean architecture",
			"ModuleName":         g.ProjectName,
			"GoVersion":          "1.21",
			"UsePostgres":        g.Answers.UsePostgres,
			"UseRedis":           g.Answers.UseRedis,
			"UseRabbitMQ":        g.Answers.UseRabbitMQ,
			"UseZap":             g.Answers.UseZap,
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
