package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-sova/sova-cli/pkg/questions"
	"github.com/spf13/cobra"
)

var (
	useZap      bool
	usePostgres bool
	useRedis    bool
	useRabbitMQ bool
)

var InitCmd = &cobra.Command{
	Use:   "api [project-name]",
	Short: "Initialize a new Go API project",
	Long: `Initialize a new Go API project with a clean architecture structure.
This command will create a new directory with the project name and set up all necessary files and directories.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectDir := filepath.Join(".", projectName)

		if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
			return fmt.Errorf("directory %s already exists", projectDir)
		}

		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}

		answers := &questions.ProjectAnswers{
			ProjectName: projectName,
			UseZap:      useZap,
			UsePostgres: usePostgres,
			UseRedis:    useRedis,
			UseRabbitMQ: useRabbitMQ,
		}

		generator := NewAPIProjectGenerator(projectName, projectDir, answers)

		files, dirs, err := generator.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate project files: %v", err)
		}

		for _, dir := range dirs {
			dirPath := filepath.Join(projectDir, dir)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dir, err)
			}
			fmt.Printf("Created directory: %s\n", dirPath)
		}

		if err := generator.WriteFiles(files); err != nil {
			return fmt.Errorf("failed to write files: %v", err)
		}

		// Initialize and tidy up the Go module
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectDir
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run go mod tidy: %v\nOutput: %s", err, output)
		}

		// Get required dependencies
		getCmd := exec.Command("go", "get", "github.com/gin-gonic/gin", "github.com/joho/godotenv")
		getCmd.Dir = projectDir
		if output, err := getCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to get dependencies: %v\nOutput: %s", err, output)
		}

		// Run go mod tidy again to clean up dependencies
		tidyCmd = exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectDir
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run go mod tidy: %v\nOutput: %s", err, output)
		}

		fmt.Printf("\nProject %s created successfully!\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("cd %s\n", projectName)
		fmt.Println("go mod tidy")
		fmt.Println("docker compose up -d")
		fmt.Println("go run cmd/main.go")
		fmt.Println("\nYour API will be available at http://localhost:8080")
		fmt.Println("Test the ping endpoint: curl http://localhost:8080/api/ping")

		return nil
	},
}

func init() {
	InitCmd.Flags().BoolVar(&useZap, "use-zap", false, "Use zap logger")
	InitCmd.Flags().BoolVar(&usePostgres, "use-postgres", false, "Use PostgreSQL")
	InitCmd.Flags().BoolVar(&useRedis, "use-redis", false, "Use Redis")
	InitCmd.Flags().BoolVar(&useRabbitMQ, "use-rabbitmq", false, "Use RabbitMQ")
}
