package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-sova/sova-cli/internal/project/api"
	"github.com/go-sova/sova-cli/internal/project/cli"
	"github.com/go-sova/sova-cli/pkg/questions"
	"github.com/go-sova/sova-cli/templates"
	"github.com/spf13/cobra"
)

var (
	projectType string
	useZap      bool
	usePostgres bool
	useRedis    bool
	useRabbitMQ bool
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project",
	Long: `Initialize a new project with the specified name.
This command will guide you through the project setup process.
If you don't provide a project name, you'll be prompted to enter one.
You can choose between different project types:
  - api: A Go API project with clean architecture
  - cli: A Go CLI project with clean architecture`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		var err error

		if len(args) > 0 {
			projectName = args[0]
		} else {
			projectName, err = questions.AskProjectName()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		// Validate project type early
		switch projectType {
		case "api", "cli":
			// Valid project type, continue
		case "":
			// If project type is not provided via flag, ask for it
			projectType, err = questions.AskProjectType()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		default:
			fmt.Printf("Error: unsupported project type: %s\n", projectType)
			return
		}

		// If flags are not provided, ask for project configuration
		if !cmd.Flags().Changed("use-zap") &&
			!cmd.Flags().Changed("use-postgres") &&
			!cmd.Flags().Changed("use-redis") &&
			!cmd.Flags().Changed("use-rabbitmq") {
			answers, err := questions.AskProjectQuestions(projectType)
			if err != nil {
				fmt.Printf("Error: failed to get project configuration: %v\n", err)
				return
			}
			useZap = answers.UseZap
			usePostgres = answers.UsePostgres
			useRedis = answers.UseRedis
			useRabbitMQ = answers.UseRabbitMQ
		}

		// Initialize project based on type
		switch projectType {
		case "api":
			apiCmd := api.InitCmd
			args := []string{projectName}
			if cmd.Flags().Changed("use-zap") {
				args = append(args, "--use-zap="+fmt.Sprintf("%v", useZap))
			}
			if cmd.Flags().Changed("use-postgres") {
				args = append(args, "--use-postgres="+fmt.Sprintf("%v", usePostgres))
			}
			if cmd.Flags().Changed("use-redis") {
				args = append(args, "--use-redis="+fmt.Sprintf("%v", useRedis))
			}
			if cmd.Flags().Changed("use-rabbitmq") {
				args = append(args, "--use-rabbitmq="+fmt.Sprintf("%v", useRabbitMQ))
			}
			apiCmd.SetArgs(args)
			err = apiCmd.Execute()
		case "cli":
			cliCmd := cli.InitCmd
			args := []string{projectName}
			if cmd.Flags().Changed("use-zap") {
				args = append(args, "--use-zap="+fmt.Sprintf("%v", useZap))
			}
			if cmd.Flags().Changed("use-postgres") {
				args = append(args, "--use-postgres="+fmt.Sprintf("%v", usePostgres))
			}
			if cmd.Flags().Changed("use-redis") {
				args = append(args, "--use-redis="+fmt.Sprintf("%v", useRedis))
			}
			if cmd.Flags().Changed("use-rabbitmq") {
				args = append(args, "--use-rabbitmq="+fmt.Sprintf("%v", useRabbitMQ))
			}
			cliCmd.SetArgs(args)
			err = cliCmd.Execute()
		default:
			err = fmt.Errorf("unsupported project type: %s", projectType)
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Project initialized successfully\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&projectType, "type", "t", "", "Project type (api or cli)")
	initCmd.Flags().BoolVar(&useZap, "use-zap", false, "Use zap logger")
	initCmd.Flags().BoolVar(&usePostgres, "use-postgres", false, "Use PostgreSQL")
	initCmd.Flags().BoolVar(&useRedis, "use-redis", false, "Use Redis")
	initCmd.Flags().BoolVar(&useRabbitMQ, "use-rabbitmq", false, "Use RabbitMQ")
}

// InitializeProject creates a new project at the specified path using the given template
func InitializeProject(projectPath string, templateName string) error {
	if err := ValidateTemplate(templateName); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Check if directory already exists
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s already exists", projectPath)
	}

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize template loader and file generator
	loader := templates.NewTemplateLoader()
	generator := templates.NewFileGenerator(loader)

	// Define template files to generate
	templateFiles := map[string]string{
		"main.tpl":      "main.go",
		"go-mod.tpl":    "go.mod",
		"readme.tpl":    "README.md",
		"root.tpl":      "cmd/root.go",
		"version.tpl":   "cmd/version.go",
		"gitignore.tpl": ".gitignore",
	}

	// Generate each file from template
	for tpl, outFile := range templateFiles {
		templatePath := templates.GetTemplatePath(templateName, tpl)
		outputPath := filepath.Join(projectPath, outFile)

		data := map[string]interface{}{
			"ProjectName": filepath.Base(projectPath),
			"ModulePath":  "github.com/example/" + filepath.Base(projectPath),
		}

		if err := generator.GenerateFile(templatePath, outputPath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", outFile, err)
		}
	}

	return nil
}

// ValidateTemplate checks if the template name is valid
func ValidateTemplate(templateName string) error {
	if templateName == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(templateName, "..") {
		return fmt.Errorf("invalid template name: contains path traversal")
	}

	// Check if template exists by trying to load a common template file
	loader := templates.NewTemplateLoader()
	templatePath := templates.GetTemplatePath(templateName, "main.tpl")
	if _, err := loader.LoadTemplate(templatePath); err != nil {
		return fmt.Errorf("template %s does not exist", templateName)
	}

	return nil
}
