# HTTP Services Examples

Web service implementations demonstrating REST APIs, HTTP clients, and service architectures.

## Services

### 🌐 [JSONPlaceholder Client](./jsonplaceholder/)
REST API client demonstrating modern Go HTTP patterns:
- **Domain Models**: User, Post, Comment, Album, Photo, Todo
- **Service Layer**: Clean separation of HTTP concerns
- **HTTP Client**: Custom client with proper error handling
- **Testing**: Mock HTTP client for unit testing
- Demonstrates: REST API consumption, JSON marshaling, service architecture

```bash
cd jsonplaceholder/
go run main.go
go test ./...
```

### 🖥️ [HTTP Server](./server/)
Basic HTTP server implementation:
- **Server Setup**: Standard HTTP server patterns
- **Route Handling**: Basic request routing
- **Testing**: Server testing strategies
- Demonstrates: HTTP server basics, handler patterns, testing servers

```bash
cd server/
go run server.go
go test -v
```

## Architecture Patterns

### Service Layer Pattern
```
main.go
├── models/          # Domain entities
├── services/        # Business logic
└── pkg/http/        # HTTP infrastructure
```

### HTTP Client Design
- Interface-based HTTP client for testability
- Proper error handling and timeouts
- JSON request/response handling
- Mock implementations for testing

### Server Patterns
- Handler function organization
- Middleware concepts
- Request/response patterns
- Graceful error handling

## Key Concepts Demonstrated

### REST API Client
- **HTTP Verbs**: GET, POST, PUT, DELETE
- **JSON Handling**: Marshal/Unmarshal patterns
- **Error Handling**: HTTP status codes and error responses
- **Testing**: Mocking external dependencies

### HTTP Server
- **Request Routing**: URL pattern matching
- **Handler Functions**: Request processing
- **Response Writing**: JSON and text responses
- **Middleware**: Cross-cutting concerns

### Testing Strategies
- **Unit Testing**: Individual service methods
- **Integration Testing**: End-to-end HTTP flows
- **Mock Testing**: Simulating external services
- **Server Testing**: HTTP server testing patterns

## Running Examples

### JSONPlaceholder Client
```bash
cd jsonplaceholder/
# Run the example
go run main.go

# Run tests with coverage
go test -cover ./...

# Test specific service
go test -v ./services/
```

### HTTP Server
```bash
cd server/
# Start the server
go run server.go

# In another terminal, test the server
curl http://localhost:8080/

# Run tests
go test -v
```

## Integration with Learning Path

### Prerequisites
- **Fundamentals**: Basic Go syntax, functions, structs
- **Intermediate**: Interfaces, error handling, testing

### Skills Developed
- HTTP client/server programming
- JSON handling and APIs
- Service-oriented architecture
- Testing HTTP services
- Error handling in distributed systems

## Production Considerations

These examples demonstrate foundational concepts. For production systems, consider:

- **Authentication/Authorization**: JWT, OAuth, API keys
- **Rate Limiting**: Protecting against abuse
- **Logging/Monitoring**: Observability and debugging
- **Configuration**: Environment-based configuration
- **Graceful Shutdown**: Proper server lifecycle management
- **TLS/HTTPS**: Secure communications

See [`/projects/api/`](../../projects/api/) for a production-ready implementation with these concerns addressed.

## Next Steps

1. **Enhance the Examples**: Add authentication, logging, middleware
2. **Build Your Own API**: Create a REST API for a domain you're interested in
3. **Study Production Code**: Explore the complete API in `/projects/api/`
4. **Learn Advanced Patterns**: Microservices, event-driven architecture