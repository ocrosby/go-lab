# Go API Example

A comprehensive Go API example demonstrating modern Go development practices and architectural patterns.

## Features

- ✅ **Hexagonal Architecture (Ports and Adapters)** - Clean separation of concerns
- ✅ **Built-in HTTP ServeMux** - Using Go's standard library HTTP multiplexer
- ✅ **Kubernetes Health Probes** - Liveness, readiness, and startup probes
- ✅ **Dependency Injection** - Using uber-go/dig for clean dependency management
- ✅ **Cobra CLI & Viper Configuration** - Command-line interface and configuration management  
- ✅ **Mock Generation** - Using uber-go/mock for comprehensive testing
- ✅ **Interface-Driven Design** - Interfaces everywhere for better testability
- ✅ **Gang of Four Design Patterns** - Factory and Observer patterns implemented
- ✅ **Low Cyclomatic Complexity** - All functions maintain complexity ≤ 7
- ✅ **Go-Task Build System** - Modern task runner for build automation
- ✅ **GolangCI-Lint** - Comprehensive linting configuration

## Architecture

```
├── cmd/api/                    # Application entrypoint
├── internal/
│   ├── application/           # Application services (use cases)
│   ├── domain/               # Domain entities and interfaces
│   ├── infrastructure/
│   │   └── adapters/
│   │       ├── http/         # HTTP handlers and server
│   │       └── repository/   # Data persistence implementations
│   ├── config/              # Configuration management
│   ├── di/                  # Dependency injection container
│   └── patterns/            # Design pattern implementations
├── pkg/
│   └── health/              # Health check functionality
└── mocks/                   # Generated mock interfaces
```

## Quick Start

### Prerequisites

- Go 1.21+
- Task (go-task.github.io/task)
- GolangCI-Lint (optional)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd api

# Install dependencies
go mod tidy

# Install task runner (if not already installed)
go install github.com/go-task/task/v3/cmd/task@latest
```

### Running the Application

```bash
# Using task
task run

# Or directly
go run cmd/api/main.go server
```

The API will be available at:
- Main API: `http://localhost:8080`
- Health endpoints: `http://localhost:8081`

### Available Tasks

```bash
task                    # Show all available tasks
task build             # Build the application
task run               # Run the application
task test              # Run tests
task test:coverage     # Run tests with coverage
task lint              # Run linting
task fmt               # Format code
task clean             # Clean build artifacts
task generate:mocks    # Generate mock files
task docker:build      # Build Docker image
task docker:run        # Run in Docker
task ci                # Full CI pipeline
```

## API Endpoints

### Users API
- `GET /users` - List all users (with pagination)
- `POST /users` - Create a new user  
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Health Probes (Port 8081)
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe  
- `GET /startupz` - Startup probe

## Design Patterns Used

### Factory Pattern
Located in `internal/patterns/factory.go`, provides a factory for creating different repository implementations.

### Observer Pattern  
Located in `internal/patterns/observer.go`, implements an event system for user-related actions.

### Dependency Injection Pattern
Implemented using uber-go/dig in `internal/di/container.go` for clean dependency management.

## Configuration

Configuration is handled by Viper and can be provided via:
- Environment variables
- Configuration files (config.yaml)
- Command-line flags

Default configuration:
- Server port: 8080
- Health port: 8081
- Host: 0.0.0.0

## Testing

Run the complete test suite:

```bash
task test              # Basic tests
task test:coverage     # Tests with coverage report
```

Generate new mocks:

```bash
task generate:mocks
```

## Linting and Code Quality

The project enforces a maximum cyclomatic complexity of 7 and includes comprehensive linting:

```bash
task lint             # Run GolangCI-Lint
task fmt              # Format and organize imports  
```

## Docker Support

Build and run with Docker:

```bash
task docker:build     # Build Docker image
task docker:run       # Run container
```

Or manually:

```bash
docker build -t api:latest .
docker run -p 8080:8080 -p 8081:8081 api:latest
```

## Development

### Adding New Features

1. Define domain interfaces in `internal/domain/`
2. Implement business logic in `internal/application/`
3. Create adapters in `internal/infrastructure/adapters/`
4. Register dependencies in `internal/di/container.go`
5. Generate mocks: `task generate:mocks`
6. Write tests
7. Run quality checks: `task check`

### Code Style

- All functions maintain cyclomatic complexity ≤ 7
- Interfaces are used throughout for better testability
- Dependencies are injected via the DI container
- Error handling follows Go best practices
- Tests include both unit tests and integration tests

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the established patterns
4. Run `task ci` to ensure all checks pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.