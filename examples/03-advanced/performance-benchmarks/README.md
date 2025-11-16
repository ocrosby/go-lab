# Performance Benchmarks & Optimization 🚀

Learn to measure, analyze, and optimize Go application performance.

## Prerequisites
- ✅ Completed [Advanced Examples](../production-api/) or equivalent experience
- ✅ Understanding of Go's runtime and memory management
- ✅ Familiarity with profiling tools and methodologies

## Overview

Performance optimization in Go requires understanding:
- **Benchmarking**: Measuring code performance accurately
- **Profiling**: Identifying bottlenecks and resource usage
- **Memory Management**: Understanding allocations and GC behavior
- **Concurrency**: Optimizing parallel and concurrent code
- **I/O Operations**: Optimizing network and disk operations

## Benchmarking Examples

### Basic Benchmarking
```go
// benchmark_test.go
package performance

import (
    "testing"
)

func BenchmarkStringConcatenation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        result := ""
        for j := 0; j < 1000; j++ {
            result += "a"
        }
    }
}

func BenchmarkStringBuilderConcatenation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 1000; j++ {
            builder.WriteString("a")
        }
        _ = builder.String()
    }
}
```

### Memory Allocation Benchmarking
```go
func BenchmarkSliceAppend(b *testing.B) {
    b.ReportAllocs() // Report memory allocations
    
    for i := 0; i < b.N; i++ {
        var slice []int
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}

func BenchmarkSlicePrealloc(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        slice := make([]int, 0, 1000) // Pre-allocate capacity
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}
```

### Parallel Benchmarking
```go
func BenchmarkParallelWork(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Work that can be done in parallel
            _ = expensiveOperation()
        }
    })
}
```

## Profiling Examples

### CPU Profiling
```go
// profile_example.go
package main

import (
    "log"
    "os"
    "runtime/pprof"
)

func main() {
    // CPU profiling
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    if err := pprof.StartCPUProfile(f); err != nil {
        log.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    
    // Your application code here
    doWork()
}

func doWork() {
    // Simulate CPU-intensive work
    for i := 0; i < 1000000; i++ {
        _ = fibonacci(30)
    }
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

### Memory Profiling
```go
// memory_profile.go
package main

import (
    "log"
    "os"
    "runtime"
    "runtime/pprof"
)

func main() {
    // Memory profiling
    f, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    // Your application code here
    doMemoryIntensiveWork()
    
    runtime.GC() // Force garbage collection
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Fatal(err)
    }
}

func doMemoryIntensiveWork() {
    // Simulate memory-intensive work
    data := make([][]byte, 1000)
    for i := range data {
        data[i] = make([]byte, 1024*1024) // 1MB per slice
    }
    
    // Use the data
    for _, slice := range data {
        slice[0] = 1
    }
}
```

### HTTP Server Profiling
```go
// server_with_profiling.go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof" // Import for side effects
)

func main() {
    // Start profiling server on separate port
    go func() {
        log.Println("Profiling server starting on :6060")
        log.Println(http.ListenAndServe(":6060", nil))
    }()
    
    // Your actual HTTP server
    http.HandleFunc("/", handler)
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Simulate some work
    result := 0
    for i := 0; i < 100000; i++ {
        result += i
    }
    
    w.Write([]byte("Hello, World!"))
}
```

## Optimization Strategies

### String Operations
```go
// Inefficient: String concatenation
func inefficientStringConcat(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ","
    }
    return result
}

// Efficient: Using strings.Builder
func efficientStringConcat(items []string) string {
    var builder strings.Builder
    for _, item := range items {
        builder.WriteString(item)
        builder.WriteString(",")
    }
    return builder.String()
}

// Most efficient: Using strings.Join
func mostEfficientStringConcat(items []string) string {
    return strings.Join(items, ",")
}
```

### Memory Management
```go
// Inefficient: Growing slice dynamically
func inefficientSliceUsage() []int {
    var result []int
    for i := 0; i < 10000; i++ {
        result = append(result, i)
    }
    return result
}

// Efficient: Pre-allocating slice
func efficientSliceUsage() []int {
    result := make([]int, 0, 10000) // Pre-allocate capacity
    for i := 0; i < 10000; i++ {
        result = append(result, i)
    }
    return result
}

// Most efficient: Direct indexing
func mostEfficientSliceUsage() []int {
    result := make([]int, 10000) // Pre-allocate length
    for i := 0; i < 10000; i++ {
        result[i] = i
    }
    return result
}
```

### Concurrent Processing
```go
// Sequential processing
func sequentialProcessing(items []int) []int {
    results := make([]int, len(items))
    for i, item := range items {
        results[i] = expensiveOperation(item)
    }
    return results
}

// Concurrent processing with worker pool
func concurrentProcessing(items []int, numWorkers int) []int {
    jobs := make(chan int, len(items))
    results := make(chan result, len(items))
    
    // Start workers
    for w := 0; w < numWorkers; w++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for i, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Collect results
    output := make([]int, len(items))
    for i := 0; i < len(items); i++ {
        res := <-results
        output[res.index] = res.value
    }
    
    return output
}

type result struct {
    index int
    value int
}

func worker(jobs <-chan int, results chan<- result) {
    for job := range jobs {
        results <- result{
            index: job,
            value: expensiveOperation(job),
        }
    }
}
```

## Running Performance Tests

### Benchmarking Commands
```bash
# Run all benchmarks
go test -bench=.

# Run benchmarks with memory stats
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkStringConcatenation

# Run benchmarks multiple times for accuracy
go test -bench=. -count=5

# Generate CPU profile during benchmarking
go test -bench=. -cpuprofile=cpu.prof

# Generate memory profile during benchmarking
go test -bench=. -memprofile=mem.prof
```

### Profiling Analysis
```bash
# Analyze CPU profile
go tool pprof cpu.prof

# Analyze memory profile
go tool pprof mem.prof

# Generate visual profile (requires graphviz)
go tool pprof -png cpu.prof > cpu.png

# Web interface for profile analysis
go tool pprof -http=:8080 cpu.prof
```

### Live Profiling
```bash
# Profile running application (assuming pprof endpoint)
go tool pprof http://localhost:6060/debug/pprof/profile

# Memory profile of running application
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine analysis
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## Performance Testing Scenarios

### Database Operations
```go
func BenchmarkDatabaseInsert(b *testing.B) {
    db := setupTestDB(b)
    defer db.Close()
    
    b.ResetTimer() // Don't count setup time
    
    for i := 0; i < b.N; i++ {
        insertRecord(db, generateTestData())
    }
}

func BenchmarkDatabaseBatchInsert(b *testing.B) {
    db := setupTestDB(b)
    defer db.Close()
    
    batchSize := 100
    data := make([]TestRecord, batchSize)
    for i := range data {
        data[i] = generateTestData()
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        batchInsert(db, data)
    }
}
```

### HTTP Client Performance
```go
func BenchmarkHTTPClient(b *testing.B) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }))
    defer server.Close()
    
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        resp, err := client.Get(server.URL)
        if err != nil {
            b.Fatal(err)
        }
        resp.Body.Close()
    }
}

func BenchmarkHTTPClientWithConnectionReuse(b *testing.B) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }))
    defer server.Close()
    
    // Configure client for connection reuse
    client := &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        resp, err := client.Get(server.URL)
        if err != nil {
            b.Fatal(err)
        }
        resp.Body.Close()
    }
}
```

## Best Practices

### Benchmarking Best Practices
1. **Use b.ResetTimer()**: Exclude setup time from measurements
2. **Report Allocations**: Use `b.ReportAllocs()` to track memory usage
3. **Multiple Runs**: Use `-count=5` for more accurate results
4. **Realistic Data**: Use representative data sizes and patterns
5. **Stable Environment**: Run on consistent hardware/OS conditions

### Optimization Guidelines
1. **Measure First**: Always profile before optimizing
2. **Focus on Hotspots**: Optimize the most expensive operations first
3. **Memory vs Speed**: Consider trade-offs between memory usage and speed
4. **Avoid Premature Optimization**: Don't optimize until you have metrics
5. **Test Performance**: Include performance tests in CI/CD

### Common Performance Pitfalls
1. **String Concatenation**: Use `strings.Builder` for multiple concatenations
2. **Slice Growth**: Pre-allocate slices when size is known
3. **Map Lookups**: Consider sync.Map for concurrent access
4. **Interface Conversions**: Minimize boxing/unboxing operations
5. **Reflection**: Avoid reflection in hot paths

## Integration with Production API

### Adding Benchmarks to Production API
```go
// internal/benchmarks/user_service_bench_test.go
func BenchmarkUserService_CreateUser(b *testing.B) {
    service := setupUserService(b)
    user := &domain.User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := service.CreateUser(context.Background(), user)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Performance Monitoring in Production
```go
// Add to your HTTP middleware
func performanceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        
        // Log slow requests
        if duration > 100*time.Millisecond {
            log.Printf("Slow request: %s %s took %v", r.Method, r.URL.Path, duration)
        }
        
        // Record metrics (Prometheus, etc.)
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
    })
}
```

## Next Steps

After mastering performance optimization:

1. **Production Integration**: Add benchmarks to your production API
2. **Monitoring Setup**: Implement continuous performance monitoring
3. **Load Testing**: Use tools like `hey`, `wrk`, or `k6` for load testing
4. **Advanced Profiling**: Learn flame graphs and advanced profiling techniques
5. **Distributed Tracing**: Implement tracing for microservices performance

## Tools and Resources

### Essential Tools
- **go tool pprof**: Built-in profiling tool
- **go-torch**: Flame graph generation (deprecated, use pprof)
- **graphviz**: Visual profile generation
- **hey**: HTTP load testing tool
- **wrk**: Modern HTTP benchmarking tool

### Advanced Tools
- **Jaeger**: Distributed tracing
- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **DataDog**: Application performance monitoring

Understanding performance characteristics is crucial for building scalable Go applications. Use these tools and techniques to ensure your applications perform well under load!