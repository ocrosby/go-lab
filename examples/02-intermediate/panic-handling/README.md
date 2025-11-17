# Panic Handling in Goroutines

This example demonstrates the critical importance of proper panic handling in spawned goroutines in Go applications.

## The Problem

In Go, panics in goroutines do NOT propagate to the calling goroutine. When a goroutine panics without recovery:
1. The goroutine terminates immediately
2. The panic does NOT crash the main program (unless it's the main goroutine)
3. Resources may not be cleaned up properly
4. Other goroutines continue running unaware of the failure
5. Silent failures occur with no error reporting

This is particularly dangerous in production systems where:
- Worker pools process background jobs
- HTTP handlers spawn goroutines for async operations
- Long-running services use goroutines for concurrent tasks

## Project Structure

```
panic-handling/
├── before/              # ❌ Incorrect implementation
│   ├── worker.go        # Worker pool without panic recovery
│   └── worker_test.go   # Tests demonstrating the failure
├── after/               # ✅ Correct implementation
│   ├── worker.go        # Worker pool with proper panic recovery
│   └── worker_test.go   # Tests showing graceful handling
└── README.md           # This file
```

## Before: The Unsafe Pattern

### Key Issues

The `before/` implementation demonstrates common mistakes:

1. **No panic recovery in worker goroutines**
   ```go
   go func(workerID int) {
       for job := range jobs {
           processJob(workerID, job) // If this panics, goroutine dies
       }
   }(i)
   ```

2. **Silent failures** - Panics are not captured or logged
3. **Resource leaks** - Worker goroutines terminate without cleanup
4. **No error propagation** - Callers never know about the panic
5. **Reduced capacity** - Dead workers reduce pool capacity

### Running the Tests

```bash
cd before
go test -v
```

You'll see that when a job causes a panic:
- The test hangs or times out
- Not all results are received
- The worker goroutine dies silently

## After: The Safe Pattern

### Key Improvements

The `after/` implementation shows best practices:

1. **Deferred recovery in each goroutine**
   ```go
   func (wp *WorkerPool) processJobWithRecovery(workerID int, job Job) (result JobResult) {
       defer func() {
           if r := recover(); r != nil {
               wp.panicHandler(r, workerID, job)
               result = JobResult{
                   Job:       job,
                   IsPanic:   true,
                   PanicInfo: fmt.Sprintf("%v", r),
                   Err:       fmt.Errorf("panic recovered: %v", r),
               }
           }
       }()
       return wp.processJob(workerID, job)
   }
   ```

2. **Custom panic handlers** - Configurable logging and alerting
3. **Error propagation** - Panics converted to error results
4. **Worker survival** - Workers continue processing after panic
5. **WaitGroup tracking** - Proper goroutine lifecycle management
6. **Helper function** - `SafeGo()` for ad-hoc goroutines

### Running the Tests

```bash
cd after
go test -v
```

All tests pass, demonstrating:
- Normal job processing
- Graceful error handling
- Panic recovery with custom handlers
- Multiple panic recovery
- All results received even after panics

## Design Patterns Used

### 1. Dependency Injection (DI)
Custom panic handlers are injected into the worker pool:
```go
wp := NewWorkerPool(3, customPanicHandler)
```

This allows:
- Different panic handling strategies per use case
- Easy testing with mock handlers
- Separation of concerns

### 2. Strategy Pattern
`PanicHandler` is a strategy that can be swapped:
```go
type PanicHandler func(recovered interface{}, workerID int, job Job)
```

Different strategies might:
- Log to different destinations
- Send alerts to monitoring systems
- Retry failed jobs
- Implement circuit breakers

### 3. Template Method Pattern
The recovery logic is standardized while allowing customization:
```go
func (wp *WorkerPool) processJobWithRecovery(workerID int, job Job) (result JobResult) {
    defer func() {
        if r := recover(); r != nil {
            wp.panicHandler(r, workerID, job) // Customizable
            result = buildPanicResult(r, job)  // Standardized
        }
    }()
    return wp.processJob(workerID, job)
}
```

## Best Practices

### ✅ DO

1. **Always use `defer recover()` in spawned goroutines**
   ```go
   go func() {
       defer func() {
           if r := recover(); r != nil {
               log.Printf("Recovered panic: %v", r)
           }
       }()
       // Your goroutine logic
   }()
   ```

2. **Create a helper function for launching goroutines**
   ```go
   SafeGo(func() {
       // Your logic
   }, panicHandler)
   ```

3. **Log panic information with context**
   - Stack traces
   - Worker/goroutine ID
   - Job/request data
   - Timestamps

4. **Convert panics to errors when possible**
   - Return error results instead of panicking
   - Propagate errors to callers
   - Enable graceful degradation

5. **Use WaitGroups for lifecycle management**
   ```go
   wp.wg.Add(1)
   go func() {
       defer wp.wg.Done()
       // Worker logic
   }()
   ```

6. **Inject panic handlers for testability**
   - Makes panic handling observable in tests
   - Allows different handling strategies
   - Follows dependency injection principle

### ❌ DON'T

1. **Never spawn goroutines without panic recovery**
   ```go
   // BAD - Panic will crash this goroutine silently
   go func() {
       doWork() // What if this panics?
   }()
   ```

2. **Don't ignore panics in worker pools**
   - Reduces pool capacity
   - Creates silent failures
   - Difficult to debug

3. **Don't assume panics propagate across goroutines**
   - Each goroutine needs its own recovery
   - Parent goroutines won't catch child panics

4. **Don't panic in library code**
   - Return errors instead
   - Let callers decide how to handle failures
   - Panics are for truly unrecoverable situations

## Testing Strategies

### Test for Expected Panics
```go
func TestWorkerPool_PanicRecovery(t *testing.T) {
    panicHandlerCalled := false
    handler := func(r interface{}, workerID int, job Job) {
        panicHandlerCalled = true
    }
    
    wp := NewWorkerPool(3, handler)
    // Submit panic-inducing job
    // Verify panic was handled
    
    if !panicHandlerCalled {
        t.Error("Expected panic to be caught")
    }
}
```

### Test Resource Cleanup
```go
func TestWorkerPool_Wait(t *testing.T) {
    wp := NewWorkerPool(2, nil)
    // Submit jobs
    wp.Close()
    
    done := make(chan bool)
    go func() {
        wp.Wait() // Should complete when all workers done
        done <- true
    }()
    
    select {
    case <-done:
        // Success
    case <-time.After(timeout):
        t.Fatal("Workers didn't complete")
    }
}
```

## Real-World Applications

### HTTP Server Middleware
```go
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic in HTTP handler: %v", err)
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Background Job Processor
```go
func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.numWorkers; i++ {
        wp.wg.Add(1)
        SafeGo(func() {
            defer wp.wg.Done()
            wp.worker(ctx, i)
        }, wp.panicHandler)
    }
}
```

### Long-Running Service
```go
func runService(ctx context.Context) {
    SafeGo(func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                processEvents() // Might panic
            }
        }
    }, func(r interface{}) {
        alertMonitoring(r)
        log.Printf("Service panic: %v", r)
    })
}
```

## Comparison Summary

| Aspect | Before (Unsafe) | After (Safe) |
|--------|----------------|--------------|
| Panic Recovery | ❌ None | ✅ Every goroutine |
| Error Propagation | ❌ Silent failures | ✅ Panic → Error |
| Worker Survival | ❌ Dies on panic | ✅ Continues working |
| Logging | ❌ No panic info | ✅ Custom handlers |
| Testability | ❌ Tests hang | ✅ Fully testable |
| Resource Cleanup | ❌ Leaks | ✅ WaitGroup tracking |
| Production Ready | ❌ No | ✅ Yes |

## Key Takeaways

1. **Panics in goroutines don't propagate** - Each goroutine needs its own recovery
2. **Use `defer recover()` in every spawned goroutine** - This is not optional
3. **Create helper functions** - `SafeGo()` makes safe launching easy
4. **Inject panic handlers** - Enables logging, alerting, and testing
5. **Convert panics to errors** - Allows graceful degradation
6. **Test panic scenarios** - Verify your recovery logic works
7. **Use WaitGroups** - Track goroutine lifecycles properly

## Further Reading

- [Go Blog: Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- [Effective Go: Recover](https://go.dev/doc/effective_go#recover)
- [Uber Go Style Guide: Goroutines](https://github.com/uber-go/guide/blob/master/style.md#goroutine-lifetimes)
