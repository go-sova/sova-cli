package cli

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
	Use:   "cli [project-name]",
	Short: "Initialize a new Go CLI project",
	Long: `Initialize a new Go CLI project with a clean architecture structure.
This command will create a new directory with the project name and set up all necessary files and directories.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		projectDir := filepath.Join(".", projectName)

		if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
			fmt.Printf("Error: directory %s already exists\n", projectDir)
			return
		}

		if err := os.MkdirAll(projectDir, 0755); err != nil {
			fmt.Printf("Error: failed to create project directory: %v\n", err)
			return
		}

		answers := &questions.ProjectAnswers{
			ProjectName: projectName,
			UseZap:      useZap,
			UsePostgres: usePostgres,
			UseRedis:    useRedis,
			UseRabbitMQ: useRabbitMQ,
		}

		generator := NewCLIProjectGenerator(projectName, projectDir, answers)

		files, dirs, err := generator.Generate()
		if err != nil {
			fmt.Printf("Error: failed to generate project files: %v\n", err)
			return
		}

		for _, dir := range dirs {
			dirPath := filepath.Join(projectDir, dir)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Printf("Error: failed to create directory %s: %v\n", dir, err)
				return
			}
			fmt.Printf("Created directory: %s\n", dirPath)
		}

		if err := generator.WriteFiles(files); err != nil {
			fmt.Printf("Error: failed to write files: %v\n", err)
			return
		}

		// Initialize and tidy up the Go module
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = projectDir
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			fmt.Printf("Error: failed to run go mod tidy: %v\nOutput: %s\n", err, output)
			return
		}

		fmt.Printf("\nProject %s created successfully!\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf(" cd %s\n", projectName)
		fmt.Println("go mod tidy")
		fmt.Println("go run main.go")
		fmt.Println("\nTry your CLI commands:")
		fmt.Printf("   ./%s command1\n", projectName)
		fmt.Printf("   ./%s command2\n", projectName)
	},
}

func init() {
	InitCmd.Flags().BoolVar(&useZap, "use-zap", false, "Use zap logger")
	InitCmd.Flags().BoolVar(&usePostgres, "use-postgres", false, "Use PostgreSQL")
	InitCmd.Flags().BoolVar(&useRedis, "use-redis", false, "Use Redis")
	InitCmd.Flags().BoolVar(&useRabbitMQ, "use-rabbitmq", false, "Use RabbitMQ")
}
