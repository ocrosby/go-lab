# HTTP Service Template

A starter template for building REST APIs and web services in Go.

## Quick Start

```bash
# Copy template to your project
cp -r templates/http-service/ ../my-api/
cd ../my-api/

# Initialize Go module
go mod init github.com/yourusername/my-api

# Install dependencies
go mod tidy

# Run the service
make run

# Test the endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/users
```

## What's Included

### 🌐 HTTP Server
- **Router**: Gorilla Mux for routing
- **Middleware**: Logging, CORS, recovery
- **JSON Handling**: Request/response processing
- **Error Handling**: Structured error responses

### 📋 Configuration
- **Environment-based**: Different configs for dev/prod
- **Viper Integration**: YAML, JSON, and environment variables
- **Default Values**: Sensible defaults for quick start

### 🔍 Observability
- **Structured Logging**: JSON-formatted logs with logrus
- **Health Checks**: Kubernetes-ready health endpoints
- **Metrics**: Basic request metrics

### 🐳 Deployment
- **Docker**: Multi-stage build for efficiency
- **Docker Compose**: Local development environment
- **Makefile**: Common development tasks

## Project Structure

```
http-service/
├── main.go                 # Application entry point
├── go.mod                  # Dependencies
├── Makefile               # Build and run tasks
├── Dockerfile             # Container definition
├── docker-compose.yml     # Local development
├── .gitignore            # Git ignore rules
├── .golangci.yml         # Linting configuration
├── config/
│   ├── config.go         # Configuration management
│   └── config.yaml       # Default configuration
├── handlers/
│   ├── health.go         # Health check handlers
│   ├── users.go          # User API handlers
│   └── middleware.go     # HTTP middleware
├── models/
│   ├── user.go           # Data models
│   └── response.go       # API response models
├── services/
│   └── user_service.go   # Business logic
└── tests/
    ├── integration_test.go  # Integration tests
    └── handlers_test.go     # Handler unit tests
```

## API Endpoints

### Health Check
- `GET /health` - Basic health check
- `GET /health/ready` - Readiness probe
- `GET /health/live` - Liveness probe

### Users API
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

## Configuration

### Environment Variables
```bash
# Server configuration
HTTP_PORT=8080
HTTP_HOST=localhost

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Database (if you add database support)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapi
```

### Configuration File (config.yaml)
```yaml
server:
  port: 8080
  host: "localhost"
  read_timeout: "30s"
  write_timeout: "30s"

logging:
  level: "info"
  format: "json"

cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE"]
```

## Development Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run linting
make lint

# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up

# Clean build artifacts
make clean
```

## Testing

### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./handlers/
```

### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./tests/

# Run with test database
make test-integration
```

### Manual Testing
```bash
# Health check
curl http://localhost:8080/health

# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get user
curl http://localhost:8080/api/users/1
```

## Customization Guide

### 1. Update Module Name
```bash
# Replace template module with your module
find . -name "*.go" -exec sed -i 's|github.com/template/http-service|github.com/yourusername/your-service|g' {} \;
```

### 2. Add Your Domain Models
```go
// models/your_model.go
type YourModel struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 3. Create Business Logic
```go
// services/your_service.go  
type YourService struct {
    // Add dependencies (database, external APIs, etc.)
}

func (s *YourService) CreateYourModel(model *YourModel) error {
    // Implement your business logic
    return nil
}
```

### 4. Add HTTP Handlers
```go
// handlers/your_handlers.go
func (h *YourHandler) CreateYourModel(w http.ResponseWriter, r *http.Request) {
    // Parse request, call service, return response
}
```

### 5. Register Routes
```go
// main.go or router setup
r.HandleFunc("/api/your-models", handlers.CreateYourModel).Methods("POST")
r.HandleFunc("/api/your-models", handlers.ListYourModels).Methods("GET")
```

## Adding Database Support

### PostgreSQL Example
```go
// Add to go.mod
// github.com/lib/pq v1.10.7

// config/database.go
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
}

// services/user_service.go
type UserService struct {
    db *sql.DB
}

func (s *UserService) CreateUser(user *User) error {
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    return s.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
}
```

## Adding Authentication

### JWT Example
```go
// Add to go.mod
// github.com/golang-jwt/jwt/v4 v4.4.3

// middleware/auth.go
func JWTAuthentication(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        // Validate JWT token
        next.ServeHTTP(w, r)
    }
}

// Apply to protected routes
r.HandleFunc("/api/protected", JWTAuthentication(handlers.ProtectedHandler))
```

## Production Deployment

### Docker
```bash
# Build production image
docker build -t my-api:latest .

# Run container
docker run -p 8080:8080 my-api:latest
```

### Kubernetes
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-api
  template:
    metadata:
      labels:
        app: my-api
    spec:
      containers:
      - name: my-api
        image: my-api:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
```

## Monitoring & Observability

### Structured Logging
```go
import "github.com/sirupsen/logrus"

log.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_user",
    "duration": duration,
}).Info("User created successfully")
```

### Metrics (Prometheus)
```go
// Add prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

## Security Best Practices

### Input Validation
```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

func validateInput(req *CreateUserRequest) error {
    validate := validator.New()
    return validate.Struct(req)
}
```

### Error Handling
```go
// Don't expose internal errors
if err != nil {
    log.WithError(err).Error("Database error")
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
}
```

## Next Steps

After setting up this template:
1. **Add Database Integration**: Choose your database and add persistence
2. **Implement Authentication**: Add JWT or OAuth2 support
3. **Add Monitoring**: Integrate Prometheus metrics and distributed tracing
4. **Write More Tests**: Increase test coverage and add integration tests
5. **Set Up CI/CD**: Add GitHub Actions or similar for automated deployment

## Based On

This template is based on patterns from:
- [Intermediate HTTP Services](../../examples/intermediate/http-services/)
- [Production API Example](../../examples/advanced/production-api/)
- [Architecture Documentation](../../docs/architecture/)

For more advanced features, see the [Production API Template](../production-api/).