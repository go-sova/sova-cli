# {{.ProjectName}}

{{.ProjectDescription}}

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for running dependencies)
{{if .UsePostgres}}- PostgreSQL{{end}}
{{if .UseRedis}}- Redis{{end}}
{{if .UseRabbitMQ}}- RabbitMQ{{end}}

### Installing

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd {{.ProjectName}}
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Start the dependencies:
   ```bash
   docker compose up -d
   ```

5. Run the application:
   ```bash
   go run cmd/main.go
   ```

The API will be available at http://localhost:8080

## API Endpoints

- `GET /api/ping` - Health check endpoint
- `GET /api/v1/...` - Your API endpoints

## Project Structure

```
.
├── cmd/            # Application entry points
├── internal/       # Private application code
│   ├── handlers/   # HTTP handlers
│   ├── middleware/ # HTTP middleware
│   ├── routes/     # Route definitions
│   ├── server/     # Server setup
│   └── service/    # Business logic
├── pkg/            # Public library code
└── api/            # API documentation
```

## Built With

* [Go](https://golang.org/) - The programming language used
* [Gin](https://gin-gonic.com/) - HTTP web framework
{{if .UseZap}}* [Zap](https://github.com/uber-go/zap) - Structured logging{{end}}
{{if .UsePostgres}}* [PostgreSQL](https://www.postgresql.org/) - Database{{end}}
{{if .UseRedis}}* [Redis](https://redis.io/) - Cache and message broker{{end}}
{{if .UseRabbitMQ}}* [RabbitMQ](https://www.rabbitmq.com/) - Message broker{{end}}

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details 