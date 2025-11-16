# Project Templates

Scaffolding and templates for quickly starting new Go projects based on patterns from this repository.

## Available Templates

### 🚀 Basic CLI Tool (Coming Soon)
**Use Case**: Command-line utilities, scripts, and automation tools
**Based On**: Patterns from [Learning Fundamentals](../learning/01-fundamentals/) and [Intermediate Examples](../examples/02-intermediate/)

**Features**:
- Cobra CLI framework setup
- Viper configuration management  
- Structured logging with logrus
- Basic testing structure
- Makefile for common tasks

**When to Use**:
- Building command-line tools
- Automation scripts
- DevOps utilities
- System administration tools

### 🌐 [HTTP Service](./http-service/)
**Use Case**: REST APIs, web services, microservices
**Based On**: [Intermediate HTTP Services Examples](../examples/02-intermediate/http-services/)

**Features**:
- HTTP server with routing
- JSON request/response handling
- Middleware support (logging, CORS, etc.)
- Basic error handling
- Health check endpoints
- Configuration management
- Docker support

**When to Use**:
- REST API development
- Microservices
- Web application backends
- Integration endpoints

### 🏗️ Production API Template (Coming Soon)
**Use Case**: Enterprise-grade APIs with full production features
**Based On**: [Advanced Production API Example](../examples/03-advanced/production-api/)

**Features**:
- Hexagonal architecture
- Dependency injection with dig
- Comprehensive testing suite
- Database integration
- CI/CD pipeline
- Observability (logging, metrics)
- Security best practices
- Docker & Kubernetes deployment

**When to Use**:
- Production applications
- Enterprise systems
- High-reliability services
- Team-developed projects

### 📚 Go Library (Coming Soon)
**Use Case**: Reusable packages and libraries
**Based On**: Package design patterns from [Learning Advanced](../learning/03-advanced/) and throughout the repository

**Features**:
- Package structure best practices
- Comprehensive testing
- Example usage
- Documentation generation
- CI/CD for library publishing
- Semantic versioning

**When to Use**:
- Creating reusable packages
- Open source libraries
- Internal company libraries
- Shared utilities

## How to Use Templates

### Method 1: Direct Copy
```bash
# Copy template to new project location
cp -r templates/http-service/ ../my-new-service/
cd ../my-new-service/

# Initialize as new Go module
go mod init github.com/yourusername/my-new-service

# Install dependencies
go mod tidy

# Customize the template
# - Update README.md
# - Modify package names
# - Add your business logic
```

### Method 2: Template Generator (Future Enhancement)
```bash
# Future: Template generation tool
go-lab create --template=http-service --name=my-service --module=github.com/me/my-service
```

## Template Structure

### Common Structure
All templates follow this basic organization:
```
template-name/
├── README.md              # Template-specific documentation
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── Makefile              # Common tasks automation
├── Dockerfile            # Container definition
├── .gitignore            # Git ignore rules
├── .golangci.yml         # Linting configuration
├── cmd/                  # CLI commands (if applicable)
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP handlers (if applicable)
│   └── services/         # Business logic
├── pkg/                  # Public library code
├── tests/                # Test files and utilities
└── docs/                 # Additional documentation
```

### Configuration Files Included

#### Makefile
```makefile
.PHONY: build test lint run clean

build:
	go build -o bin/app ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run

run:
	go run ./cmd/app

clean:
	rm -rf bin/
```

#### .golangci.yml
Pre-configured linting rules based on repository standards:
- Cyclomatic complexity limits
- Code quality checks
- Security scanning
- Best practices enforcement

#### Dockerfile
Multi-stage builds for efficient containerization:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Runtime stage  
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## Customization Guide

### Step-by-Step Customization

1. **Initialize Project**
   ```bash
   # Copy template
   cp -r templates/your-choice/ ../your-project/
   cd ../your-project/
   
   # Update Go module
   go mod init your-module-name
   ```

2. **Update Package Names**
   ```bash
   # Find and replace template package names
   find . -name "*.go" -exec sed -i 's/template-package/your-package/g' {} \;
   ```

3. **Customize Configuration**
   - Update `README.md` with your project details
   - Modify configuration files in `internal/config/`
   - Update Docker image names and ports
   - Customize Makefile targets

4. **Add Your Logic**
   - Implement business logic in `internal/services/`
   - Add HTTP handlers in `internal/handlers/`
   - Create domain models as needed
   - Add tests alongside your code

5. **Test and Validate**
   ```bash
   # Run tests
   make test
   
   # Check linting
   make lint
   
   # Build and run
   make build
   make run
   ```

## Template Development

### Creating New Templates

When adding new templates to this repository:

1. **Base on Existing Examples**: Use proven patterns from `/examples/` and `/learning/`

2. **Follow Template Standards**:
   ```
   templates/new-template/
   ├── README.md           # Template documentation
   ├── TEMPLATE.md         # Customization instructions
   ├── .template-config    # Template metadata
   └── [project files]     # Actual template code
   ```

3. **Include Documentation**:
   - Clear use case description
   - Customization instructions
   - Example usage
   - Architecture decisions

4. **Test Templates**:
   - Verify templates work out-of-the-box
   - Test customization instructions
   - Ensure all dependencies are included

### Template Metadata (.template-config)
```yaml
name: "HTTP Service"
description: "REST API service with middleware and configuration"
difficulty: "intermediate"
based_on: "examples/intermediate/http-services"
features:
  - "HTTP routing"
  - "JSON handling"
  - "Configuration management"
  - "Docker support"
prerequisites:
  - "Basic Go knowledge"
  - "Understanding of HTTP"
estimated_setup_time: "30 minutes"
```

## Integration with Learning Path

Templates are designed to work with the repository's learning progression:

### Beginner → CLI Tool Template
After completing [Learning Fundamentals](../learning/01-fundamentals/), use CLI tool template to:
- Practice Go syntax in a real project
- Learn configuration management
- Understand project structure
- Apply concepts from [Beginner Examples](../examples/01-beginner/)

### Intermediate → HTTP Service Template  
After [HTTP Services Examples](../examples/02-intermediate/http-services/), use this template to:
- Build real web services
- Practice middleware patterns
- Implement proper error handling
- Apply [Intermediate Learning](../learning/02-intermediate/) concepts

### Advanced → Production API Template
After mastering [Advanced Concepts](../learning/03-advanced/), use this template to:
- Build enterprise-grade applications
- Implement clean architecture from [Production API Example](../examples/03-advanced/production-api/)
- Practice production deployment with [Deployment Guides](../deployment/)

## Future Enhancements

### Planned Features
- **Interactive Template Generator**: CLI tool for template customization
- **Template Variants**: Different configurations for same template type
- **IDE Integration**: VS Code and GoLand template support
- **Template Registry**: Searchable template catalog

### Community Templates
We welcome community contributions of templates for:
- Specific frameworks (Gin, Echo, Fiber)
- Database integrations (PostgreSQL, MongoDB, Redis)
- Cloud platforms (AWS, GCP, Azure)
- Specialized use cases (gRPC services, GraphQL APIs)

## Getting Help

### Template Issues
1. **Verify Prerequisites**: Ensure Go version and dependencies
2. **Check Documentation**: Read template-specific README
3. **Test Minimal Case**: Start with basic template usage
4. **Check Examples**: Reference the original examples the template is based on

### Contributing Templates
1. **Discuss First**: Open an issue to discuss the template idea
2. **Follow Standards**: Use existing templates as guides
3. **Include Tests**: Ensure templates work correctly
4. **Document Thoroughly**: Provide clear instructions and examples

Remember: Templates are starting points, not finished products. Customize them to fit your specific needs and requirements!