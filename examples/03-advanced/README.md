# Advanced Examples

Production-ready implementations demonstrating enterprise-grade Go development practices.

## Prerequisites
- ✅ Completed [Intermediate Examples](../intermediate/) or equivalent experience  
- ✅ Strong understanding of interfaces, HTTP services, and testing
- ✅ Familiar with design patterns and clean architecture principles
- ✅ Completed [Learning Advanced](../../learning/advanced/) (recommended)

## Examples in This Directory

### 🏗️ [Production API](./production-api/)
**Estimated Time**: 3-4 weeks  
**Concepts**: Hexagonal architecture, dependency injection, production deployment

A complete, production-ready API implementation featuring:

#### Architecture
- **Hexagonal (Clean) Architecture**: Clear separation of concerns
- **Dependency Injection**: Using uber-go/dig container
- **Domain-Driven Design**: Rich domain models and use cases
- **SOLID Principles**: Applied throughout the codebase

#### Technical Implementation
- **Modern Tooling**: Cobra CLI, Viper configuration, structured logging
- **Comprehensive Testing**: Unit, integration, and E2E tests with mocks
- **Quality Assurance**: GolangCI-Lint, test coverage, cyclomatic complexity limits
- **Database Integration**: Repository pattern with database abstraction
- **HTTP Layer**: RESTful API with proper error handling and validation

#### Production Features
- **Configuration Management**: Environment-based configuration
- **Health Checks**: Kubernetes-ready health probes
- **Containerization**: Docker support with multi-stage builds
- **CI/CD Integration**: GitHub Actions, automated testing and deployment
- **Observability**: Structured logging, metrics, and monitoring hooks
- **Security**: Input validation, error sanitization, secure defaults

**What You'll Learn**:
- Production-grade application architecture
- Hexagonal/Clean architecture implementation
- Advanced testing strategies and patterns
- CI/CD pipeline design and implementation
- Containerization and deployment practices
- Observability and monitoring patterns
- Security best practices for Go APIs

```bash
cd production-api/

# View comprehensive documentation
cat README.md
cat TESTING.md

# Run with Docker
docker-compose up

# Run locally
make run

# Run comprehensive tests
make test
make test-integration
make test-coverage

# Check code quality
make lint
make security-scan
```

## Learning Objectives

### Architecture Mastery
By completing this example, you'll understand:
- 🎯 **Clean Architecture**: Organizing code for maintainability and testability
- 🎯 **Dependency Injection**: Managing dependencies and enabling testing
- 🎯 **Domain-Driven Design**: Modeling business logic effectively
- 🎯 **Hexagonal Architecture**: Ports and adapters pattern
- 🎯 **SOLID Principles**: Applied in Go context

### Production Skills
- 🎯 **Testing Strategies**: Unit, integration, E2E, and performance testing
- 🎯 **Quality Engineering**: Code quality, security, and maintainability
- 🎯 **Deployment Practices**: Containerization, orchestration, and CI/CD
- 🎯 **Observability**: Logging, monitoring, and debugging in production
- 🎯 **Performance**: Optimization, profiling, and scalability considerations

### Enterprise Patterns
- 🎯 **Configuration Management**: Environment-based configuration
- 🎯 **Error Handling**: Structured error handling and logging
- 🎯 **Security**: Authentication, authorization, and input validation
- 🎯 **API Design**: RESTful APIs with proper HTTP semantics
- 🎯 **Database Patterns**: Repository pattern and data access layers

## Key Architecture Patterns

### Hexagonal Architecture
```
              External World
                     |
            ┌─────────────────────┐
            │   Adapters (HTTP)   │
            └─────────────────────┘
                     |
            ┌─────────────────────┐
            │   Application       │ ← Use Cases
            │   (Business Logic)  │
            └─────────────────────┘
                     |
            ┌─────────────────────┐
            │     Domain          │ ← Core Business Rules
            └─────────────────────┘
                     |
            ┌─────────────────────┐
            │  Infrastructure     │ ← Database, External APIs
            └─────────────────────┘
```

### Dependency Injection Pattern
```go
// Domain defines interfaces
type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// Application implements use cases
type UserService struct {
    repo UserRepository
    log  Logger
}

// Infrastructure implements interfaces
type PostgreSQLUserRepository struct {
    db *sql.DB
}

// DI container wires everything together
container.Register(func() UserRepository {
    return &PostgreSQLUserRepository{db: db}
})
```

### Testing Architecture
```
Production Code          Test Code
├── Domain              ├── Unit Tests
├── Application         ├── Integration Tests  
├── Infrastructure      ├── E2E Tests
└── Adapters           └── Performance Tests
```

## Real-World Applications

### When You'd Use This Architecture

#### Enterprise APIs
- Microservices in large organizations
- APIs requiring high reliability and maintainability
- Systems with complex business logic
- Applications requiring extensive testing

#### Production Considerations
- **Scalability**: Handling increasing load and data
- **Reliability**: Fault tolerance and graceful degradation  
- **Security**: Protecting against common vulnerabilities
- **Maintainability**: Supporting team development and evolution

## Advanced Concepts Demonstrated

### Clean Code Principles
- **Single Responsibility**: Each component has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Dependency Inversion**: Depend on abstractions, not concretions
- **Interface Segregation**: Small, focused interfaces

### Testing Strategies
```go
// Unit Test - Tests business logic in isolation
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, logger)
    
    // Test business logic without external dependencies
}

// Integration Test - Tests component interactions
func TestUserAPI_CreateUser(t *testing.T) {
    db := setupTestDatabase()
    server := setupTestServer(db)
    
    // Test actual HTTP endpoints with real database
}

// E2E Test - Tests complete user journeys  
func TestUserWorkflow(t *testing.T) {
    // Test complete user registration → login → usage flow
}
```

### Production Deployment
```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - DB_HOST=postgres
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user  
      POSTGRES_PASSWORD: password
```

## Performance & Scalability

### Optimization Techniques
- **Database**: Connection pooling, query optimization, indexing
- **HTTP**: Keep-alive connections, compression, caching headers
- **Memory**: Object pooling, garbage collection tuning
- **Concurrency**: Worker pools, rate limiting, circuit breakers

### Monitoring & Observability
```go
// Structured logging
log.WithFields(log.Fields{
    "user_id": userID,
    "action": "create_user",
    "duration": duration,
}).Info("User created successfully")

// Metrics collection
metrics.Counter("api.requests.total").
    WithTags("method", "POST", "endpoint", "/users").
    Increment()

// Distributed tracing
span := tracing.StartSpan("user.service.create")
defer span.Finish()
```

## Security Considerations

### Common Vulnerabilities Addressed
- **SQL Injection**: Parameterized queries, ORM usage
- **XSS**: Input sanitization and output encoding  
- **CSRF**: Token validation and SameSite cookies
- **Authentication**: Proper session management and JWT handling
- **Authorization**: Role-based access control

### Security Best Practices
```go
// Input validation
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
}

// Error handling (don't leak information)
if err != nil {
    log.WithError(err).Error("Database error")
    return errors.New("internal server error")
}
```

## Next Steps

After mastering this example:

### Apply to Your Projects
1. **Adapt Architecture**: Use hexagonal architecture in your applications
2. **Improve Testing**: Implement comprehensive testing strategies
3. **Production Deployment**: Deploy with proper CI/CD and monitoring
4. **Scale Up**: Handle increasing load and complexity

### Advanced Topics to Explore
- **Microservices Architecture**: Service mesh, distributed systems
- **Event-Driven Architecture**: Message queues, event sourcing
- **Advanced Security**: OAuth2, OpenID Connect, zero-trust architecture
- **Performance Optimization**: Profiling, optimization, and scaling
- **Observability**: Distributed tracing, advanced monitoring

### Contribute to Open Source
- Apply these patterns to open-source Go projects
- Share your learnings with the community
- Build libraries that implement these patterns

## Getting Help

### Debugging Production Issues
```bash
# Application logs
docker logs production-api

# Database connectivity
docker exec -it production-db psql -U user -d myapp

# Performance profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Load testing
hey -n 1000 -c 10 http://localhost:8080/api/users
```

### Common Production Issues
- **Memory Leaks**: Monitor memory usage, use pprof
- **Database Locks**: Analyze slow queries, optimize indexes
- **High Latency**: Profile critical paths, add caching
- **Error Rates**: Improve error handling and monitoring

## Integration with Learning Path

This example synthesizes concepts from:
- [Learning Advanced](../../learning/advanced/) - Design patterns and concurrency
- [Testing Strategies](../../testing/) - Comprehensive testing approaches
- [Architecture Documentation](../../docs/architecture/) - Design principles

## Recognition of Mastery

You've mastered advanced Go development when you can:
- ✅ Design and implement clean, maintainable architectures
- ✅ Write comprehensive test suites with high coverage
- ✅ Deploy and monitor production-ready applications
- ✅ Handle security, performance, and scalability concerns
- ✅ Debug and optimize production systems
- ✅ Lead technical decisions on Go projects

**Congratulations!** You now have the skills to build and maintain production-grade Go applications. Consider mentoring others, contributing to open source, or taking on technical leadership roles in Go projects.

The journey continues with emerging patterns, new tools, and evolving best practices in the Go ecosystem. Stay curious and keep building!