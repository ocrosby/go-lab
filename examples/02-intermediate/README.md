# Intermediate Examples

Building on Go fundamentals with real-world service patterns and architectures.

## Prerequisites
- ✅ Completed [Beginner Examples](../beginner/) or equivalent experience
- ✅ Comfortable with Go syntax, functions, and basic testing
- ✅ Understanding of interfaces and basic concurrency
- ✅ Completed [Learning Intermediate](../../learning/intermediate/) (recommended)

## Examples in This Directory

### 🌐 [HTTP Services](./http-services/)
**Estimated Time**: 1-2 weeks  
**Concepts**: REST APIs, HTTP clients/servers, JSON handling, service architecture

Real-world web service implementations demonstrating:

#### JSONPlaceholder Client
Professional REST API client with:
- **Domain Models**: User, Post, Comment, Album, Photo, Todo
- **Service Layer**: Clean separation of HTTP concerns  
- **HTTP Client**: Custom client with proper error handling
- **Mock Testing**: HTTP client testing with mocks

**What You'll Learn**:
- REST API consumption patterns
- JSON marshaling and unmarshaling
- HTTP client design and testing
- Service layer architecture
- Error handling in distributed systems

#### HTTP Server  
Basic server implementation covering:
- **Route Handling**: Request routing patterns
- **Handler Functions**: Processing HTTP requests
- **Response Writing**: JSON and text responses
- **Server Testing**: Testing HTTP handlers

**What You'll Learn**:
- HTTP server fundamentals
- Handler pattern implementation
- Request/response processing
- Server testing strategies

```bash
cd http-services/

# Run JSONPlaceholder client
cd jsonplaceholder/
go run main.go
go test -v ./...

# Run HTTP server
cd ../server/  
go run server.go

# In another terminal, test the server
curl http://localhost:8080/
```

## Learning Objectives

### Technical Skills
By completing these examples, you'll master:
- 🎯 **HTTP Programming**: Client and server implementation
- 🎯 **JSON Processing**: Marshaling, unmarshaling, and validation
- 🎯 **Service Architecture**: Clean separation of concerns
- 🎯 **Error Handling**: Robust error handling in network code
- 🎯 **Testing Strategies**: Mocking external dependencies
- 🎯 **Interface Design**: Creating testable, flexible interfaces

### Architecture Patterns
- **Service Layer Pattern**: Business logic separation
- **Repository Pattern**: Data access abstraction  
- **Dependency Injection**: Interface-based design
- **Clean Architecture**: Organizing code for maintainability

## Key Concepts Reinforced

### HTTP Programming
```go
// Client-side HTTP request
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    return nil, fmt.Errorf("failed to fetch users: %w", err)
}
defer resp.Body.Close()

// Server-side handler
func userHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### JSON Handling
```go
// Struct tags for JSON mapping
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// Marshaling and error handling
data, err := json.Marshal(user)
if err != nil {
    return fmt.Errorf("failed to marshal user: %w", err)
}
```

### Interface-Based Design
```go
// Define interfaces for testability
type HTTPClient interface {
    Get(url string) (*http.Response, error)
    Post(url string, body io.Reader) (*http.Response, error)
}

// Use interfaces in services
type UserService struct {
    client HTTPClient
    baseURL string
}
```

### Testing HTTP Code
```go
// Testing HTTP handlers
func TestUserHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    handler.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}

// Testing HTTP clients with mocks
func TestUserService_GetUser(t *testing.T) {
    mockClient := &MockHTTPClient{}
    service := NewUserService(mockClient, "https://api.example.com")
    
    // Setup mock expectations and test
}
```

## Real-World Applications

### When You'd Use These Patterns

#### REST API Clients
- Consuming third-party APIs (GitHub, Stripe, AWS)
- Microservice communication
- Data synchronization between services
- Integration with external systems

#### HTTP Servers
- Web APIs and microservices
- Webhook handlers
- Health check endpoints
- Admin interfaces

#### Service Architecture
- Separating business logic from HTTP concerns
- Making code testable and maintainable
- Supporting multiple transport protocols
- Enabling dependency injection

## Common Challenges & Solutions

### HTTP Client Issues
```go
// Problem: Default client has no timeout
client := &http.Client{
    Timeout: 30 * time.Second,
}

// Problem: Not handling response body properly
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
```

### JSON Handling Issues
```go
// Problem: Not handling JSON errors
var user User
if err := json.Unmarshal(data, &user); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}

// Problem: Empty values in JSON
type User struct {
    Name string `json:"name,omitempty"`
}
```

### Testing Issues
```go
// Problem: Testing real HTTP calls
// Solution: Use httptest.Server for integration tests
server := httptest.NewServer(handler)
defer server.Close()

// Problem: Hard to mock dependencies  
// Solution: Use interfaces and dependency injection
```

## Performance Considerations

### HTTP Clients
- **Connection Pooling**: Reuse HTTP clients
- **Timeouts**: Set appropriate timeouts
- **Context**: Use context for cancellation
- **Rate Limiting**: Respect API limits

### HTTP Servers
- **Graceful Shutdown**: Handle shutdown signals
- **Middleware**: Implement logging, metrics
- **Resource Cleanup**: Close database connections
- **Error Handling**: Don't expose internal errors

## Next Steps

After mastering these examples:

### Extend the Examples
1. **Add Authentication**: JWT tokens, API keys
2. **Add Middleware**: Logging, metrics, rate limiting  
3. **Add Databases**: Persistence layer integration
4. **Add Configuration**: Environment-based config
5. **Add Monitoring**: Health checks, metrics

### Check Your Readiness for Advanced
- ✅ Can design and implement REST APIs
- ✅ Understand service layer architecture
- ✅ Comfortable with HTTP client/server patterns
- ✅ Can test HTTP code effectively
- ✅ Understand JSON processing and validation
- ✅ Know interface-based design principles

If yes, proceed to [Advanced Examples](../advanced/)

### Still Need Practice?
- Extend the HTTP services with new endpoints
- Add database integration to the server
- Implement authentication and authorization
- Study the [Testing Guide](../../docs/tutorials/testing-guide.md)

## Integration with Learning Path

These examples reinforce concepts from:
- [Learning Intermediate](../../learning/intermediate/) - Interfaces and composition
- [Testing Examples](../../testing/) - Advanced testing patterns  
- [Architecture Docs](../../docs/architecture/) - Design patterns

## Production Readiness

While these examples demonstrate important patterns, production systems need:
- **Security**: Authentication, authorization, input validation
- **Observability**: Logging, metrics, tracing
- **Reliability**: Circuit breakers, retries, timeouts
- **Scalability**: Load balancing, caching, database optimization

See [Advanced Examples](../advanced/production-api/) for production-ready implementations.

## Getting Help

### Debugging HTTP Issues
```bash
# Use curl to test endpoints
curl -v http://localhost:8080/endpoint

# Check server logs for errors
go run server.go 2>&1 | grep ERROR

# Use httputil for request debugging
dump, _ := httputil.DumpRequest(req, true)
fmt.Printf("%s", dump)
```

### Common HTTP Status Codes
- **200 OK**: Success
- **400 Bad Request**: Client error  
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Server error

Understanding these patterns will prepare you for building robust, scalable web services in Go!