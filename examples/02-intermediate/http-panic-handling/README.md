# Panic Handling in HTTP Servers

This example demonstrates proper panic handling in HTTP servers, including the special case of `http.ErrAbortHandler`.

## The Problem

HTTP servers need robust panic handling because:
1. Panics in handlers crash the entire server if not recovered
2. Panics in background goroutines spawned by handlers can silently crash those goroutines
3. The `http.ErrAbortHandler` sentinel value requires special treatment
4. Poor error responses confuse clients when panics occur

## Project Structure

```
http-panic-handling/
├── before/              # ❌ Incorrect implementation
│   ├── server.go        # HTTP server without panic recovery
│   └── server_test.go   # Tests showing server crashes
├── after/               # ✅ Correct implementation
│   ├── server.go        # HTTP server with proper panic recovery
│   └── server_test.go   # Tests showing graceful handling
└── README.md           # This file
```

## Before: The Unsafe Pattern

### Key Issues

The `before/` implementation demonstrates common mistakes:

1. **No panic recovery in HTTP handlers**
   ```go
   func (s *Server) handlePanic(w http.ResponseWriter, r *http.Request) {
       panic("intentional panic in handler") // Crashes the server
   }
   ```

2. **No recovery in background goroutines**
   ```go
   go func() {
       time.Sleep(100 * time.Millisecond)
       panic("panic in background goroutine") // Silent goroutine crash
   }()
   ```

3. **No error responses** - Clients receive no meaningful response
4. **Server crashes** - Entire server can be brought down by a single panic
5. **No logging** - Panics are not captured for debugging

### Running the Tests

```bash
cd before
go test -v
```

You'll see that:
- Tests demonstrate panics crash the server
- Background goroutine panics cause test failures
- No graceful error handling

## After: The Safe Pattern

### Key Improvements

The `after/` implementation shows best practices:

1. **Recovery middleware wrapping all handlers**
   ```go
   func (s *Server) withPanicRecovery(handler http.HandlerFunc) http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
           defer func() {
               rec := recover()
               if rec == nil {
                   return
               }

               // Special handling for http.ErrAbortHandler
               if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
                   panic(rec) // Re-panic to allow http package to handle it
               }

               s.panicHandler(rec, r)

               w.Header().Set("Content-Type", "application/json")
               w.WriteHeader(http.StatusInternalServerError)
               json.NewEncoder(w).Encode(ErrorResponse{
                   Error:  "Internal Server Error",
                   Status: "error",
               })
           }()

           handler(w, r)
       }
   }
   ```

2. **Safe goroutine launcher with recovery**
   ```go
   func safeGo(fn func(), panicHandler func(interface{})) {
       go func() {
           defer func() {
               rec := recover()
               if rec == nil {
                   return
               }

               // Special handling for http.ErrAbortHandler
               if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
                   panic(rec)
               }

               if panicHandler != nil {
                   panicHandler(rec)
                   return
               }

               log.Printf("PANIC RECOVERED in goroutine: %v\n%s", rec, debug.Stack())
           }()
           fn()
       }()
   }
   ```

3. **Proper handling of http.ErrAbortHandler**
   - This special sentinel value tells the http package to abort the connection
   - Must be re-panicked to allow the http package to handle it properly
   - Used for scenarios like HTTP/2 connection hijacking

4. **Custom panic handlers** - Configurable logging and alerting
5. **Error responses** - Clients receive proper 500 status with JSON error
6. **Stack traces** - Full debugging information captured

### Running the Tests

```bash
cd after
go test -v
```

All tests pass, demonstrating:
- Handlers recover from panics gracefully
- Background goroutines handle panics without crashing
- `http.ErrAbortHandler` is properly re-panicked
- Concurrent requests with panics are all handled correctly
- Proper error responses sent to clients

## Understanding http.ErrAbortHandler

### What is http.ErrAbortHandler?

`http.ErrAbortHandler` is a sentinel error defined by the `net/http` package:

```go
var ErrAbortHandler = errors.New("net/http: abort Handler")
```

### When is it used?

The http package uses this to signal that a handler has intentionally aborted processing and wants to prevent any further processing, including:
- Writing response headers
- Writing response body
- Logging the request
- Running response middleware

### Why must it be re-panicked?

If you catch `http.ErrAbortHandler` and don't re-panic it, the http server will:
1. Continue processing the request
2. Potentially write duplicate headers/body
3. Not properly close connections in HTTP/2 scenarios

**Correct pattern:**
```go
defer func() {
    if r := recover(); r != nil {
        // Special case: re-panic http.ErrAbortHandler
        if err, ok := r.(error); ok && errors.Is(err, http.ErrAbortHandler) {
            panic(r)
        }
        
        // Handle all other panics
        handlePanic(r)
    }
}()
```

## Design Patterns Used

### 1. Middleware Pattern
Recovery logic is implemented as middleware:
```go
func (s *Server) setupRoutes() {
    s.mux.HandleFunc("/health", s.withPanicRecovery(s.handleHealth))
    s.mux.HandleFunc("/panic", s.withPanicRecovery(s.handlePanic))
}
```

### 2. Dependency Injection (DI)
Custom panic handlers are injected:
```go
server := NewServer(customPanicHandler)
```

This allows:
- Different panic handling strategies per environment
- Easy testing with mock handlers
- Separation of concerns

### 3. Template Method Pattern
Recovery logic is standardized while allowing customization:
```go
defer func() {
    if r := recover(); r != nil {
        s.panicHandler(r, req)           // Customizable
        sendErrorResponse(w)              // Standardized
    }
}()
```

## Best Practices

### ✅ DO

1. **Always wrap HTTP handlers with panic recovery**
   ```go
   func RecoveryMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           defer func() {
               if err := recover(); err != nil {
                   // Handle panic
               }
           }()
           next.ServeHTTP(w, r)
       })
   }
   ```

2. **Re-panic http.ErrAbortHandler**
   ```go
   if err, ok := r.(error); ok && errors.Is(err, http.ErrAbortHandler) {
       panic(r)
   }
   ```

3. **Log panics with full stack traces**
   ```go
   log.Printf("Panic: %v\n%s", recovered, debug.Stack())
   ```

4. **Send proper error responses to clients**
   ```go
   w.WriteHeader(http.StatusInternalServerError)
   json.NewEncoder(w).Encode(ErrorResponse{
       Error: "Internal Server Error",
   })
   ```

5. **Use safe goroutine launchers for background work**
   ```go
   safeGo(func() {
       // Background work
   }, panicHandler)
   ```

6. **Test panic scenarios**
   - Verify recovery works
   - Check error responses
   - Ensure http.ErrAbortHandler propagates

### ❌ DON'T

1. **Never let panics crash the HTTP server**
   ```go
   // BAD - No recovery
   func handleRequest(w http.ResponseWriter, r *http.Request) {
       doWork() // What if this panics?
   }
   ```

2. **Don't swallow http.ErrAbortHandler**
   ```go
   // BAD - Breaks HTTP/2 connection handling
   defer func() {
       if r := recover(); r != nil {
           log.Printf("Panic: %v", r) // Should re-panic ErrAbortHandler!
       }
   }()
   ```

3. **Don't spawn unprotected goroutines in handlers**
   ```go
   // BAD - Panic in goroutine will crash it silently
   go func() {
       processInBackground() // What if this panics?
   }()
   ```

4. **Don't expose panic details to clients**
   ```go
   // BAD - Exposes internal implementation
   http.Error(w, fmt.Sprintf("Panic: %v", recovered), 500)
   ```

## Testing Strategies

### Test Normal Panic Recovery
```go
func TestServer_PanicRecovery(t *testing.T) {
    var panicHandlerCalled bool
    panicHandler := func(recovered interface{}, r *http.Request) {
        panicHandlerCalled = true
    }
    
    server := NewServer(panicHandler)
    req := httptest.NewRequest(http.MethodGet, "/panic", nil)
    w := httptest.NewRecorder()
    
    server.ServeHTTP(w, req)
    
    if !panicHandlerCalled {
        t.Error("Expected panic to be caught")
    }
    
    if w.Code != http.StatusInternalServerError {
        t.Errorf("Expected 500 status")
    }
}
```

### Test http.ErrAbortHandler Propagation
```go
func TestServer_AbortHandler(t *testing.T) {
    server := NewServer(nil)
    req := httptest.NewRequest(http.MethodGet, "/abort", nil)
    w := httptest.NewRecorder()
    
    defer func() {
        r := recover()
        if r != http.ErrAbortHandler {
            t.Error("Expected http.ErrAbortHandler to be re-panicked")
        }
    }()
    
    server.ServeHTTP(w, req)
}
```

### Test Concurrent Panic Handling
```go
func TestServer_MultipleConcurrentPanics(t *testing.T) {
    server := NewServer(nil)
    
    concurrency := 10
    for i := 0; i < concurrency; i++ {
        go func() {
            req := httptest.NewRequest(http.MethodGet, "/panic", nil)
            w := httptest.NewRecorder()
            server.ServeHTTP(w, req)
        }()
    }
    
    // Verify all requests handled gracefully
}
```

## Real-World Applications

### Production HTTP Server
```go
func main() {
    panicHandler := func(recovered interface{}, r *http.Request) {
        log.Printf("Panic in %s %s: %v\n%s", 
            r.Method, r.URL.Path, recovered, debug.Stack())
        
        // Alert monitoring system
        alerting.SendAlert("HTTP Panic", fmt.Sprintf("%v", recovered))
    }
    
    server := NewServer(panicHandler)
    log.Fatal(server.Start(":8080"))
}
```

### Middleware Chain
```go
func NewServer() *http.Server {
    mux := http.NewServeMux()
    
    // Apply recovery middleware to all routes
    handler := RecoveryMiddleware(
        LoggingMiddleware(
            AuthMiddleware(mux)))
    
    return &http.Server{
        Handler: handler,
        Addr:    ":8080",
    }
}
```

## Comparison Summary

| Aspect | Before (Unsafe) | After (Safe) |
|--------|----------------|--------------|
| Handler Panics | ❌ Crashes server | ✅ Recovered gracefully |
| Background Goroutines | ❌ Silent crashes | ✅ Safe with recovery |
| http.ErrAbortHandler | ❌ Not handled | ✅ Properly re-panicked |
| Error Responses | ❌ None | ✅ Proper 500 JSON |
| Logging | ❌ None | ✅ With stack traces |
| Client Impact | ❌ Connection drops | ✅ Proper error response |
| Testability | ❌ Hard to test | ✅ Fully testable |
| Production Ready | ❌ No | ✅ Yes |

## Key Takeaways

1. **Always use recovery middleware** - Wrap all HTTP handlers with panic recovery
2. **Re-panic http.ErrAbortHandler** - This is critical for proper HTTP handling
3. **Protect background goroutines** - Use safe launchers with recovery
4. **Send proper error responses** - Don't leave clients hanging
5. **Log with stack traces** - Essential for debugging production issues
6. **Test panic scenarios** - Verify recovery works correctly
7. **Inject panic handlers** - Enables different strategies per environment

## Further Reading

- [Go Blog: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- [net/http Package Documentation](https://pkg.go.dev/net/http)
- [http.ErrAbortHandler Source](https://github.com/golang/go/blob/master/src/net/http/server.go)
