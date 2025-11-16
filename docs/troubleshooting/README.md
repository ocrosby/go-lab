# Troubleshooting & Debugging Guide 🔧

Comprehensive guide to solving common issues and debugging Go applications effectively.

## 📋 Quick Reference

| Problem Category | Common Issues | Quick Solutions |
|------------------|---------------|-----------------|
| [**Environment**](#environment-issues) | Go not found, GOPATH issues | Installation and path configuration |
| [**Compilation**](#compilation-errors) | Syntax errors, import issues | Code fixes and module management |
| [**Runtime**](#runtime-errors) | Panics, nil pointers, deadlocks | Debugging techniques and prevention |
| [**Testing**](#testing-issues) | Test failures, coverage problems | Test debugging and organization |
| [**Performance**](#performance-issues) | Memory leaks, slow code | Profiling and optimization |
| [**Dependencies**](#dependency-issues) | Module problems, version conflicts | Module management and resolution |

## 🚀 Environment Issues

### Go Installation Problems

#### Problem: "go: command not found"
```bash
# Check if Go is installed
go version
# Error: go: command not found
```

**Solutions:**
```bash
# 1. Install Go from official site
# Visit https://golang.org/dl/

# 2. Check if Go is in PATH
echo $PATH | grep -o "/usr/local/go/bin"

# 3. Add Go to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# 4. Verify installation
go version
```

#### Problem: Wrong Go Version
```bash
go version
# go version go1.18.x darwin/amd64
# But you need 1.21+
```

**Solutions:**
```bash
# 1. Remove old version
sudo rm -rf /usr/local/go

# 2. Install new version from golang.org

# 3. Update PATH if needed
export PATH=/usr/local/go/bin:$PATH

# 4. Verify
go version
```

### GOPATH and Module Issues

#### Problem: "cannot find package" in legacy code
```bash
go run main.go
# Error: cannot find package "github.com/example/package"
```

**Solutions:**
```bash
# Modern approach - use Go modules
go mod init your-project
go mod tidy

# Legacy approach - set GOPATH (not recommended for new projects)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## ⚠️ Compilation Errors

### Import and Module Errors

#### Problem: Import path issues
```go
// ❌ Common import mistakes
import "github.com/nonexistent/package"  // Package doesn't exist
import "./local/package"                 // Relative imports not allowed
import "fmt "                           // Extra space
```

**Solutions:**
```go
// ✅ Correct imports
import "fmt"                                    // Standard library
import "github.com/gorilla/mux"                // External package
import "github.com/youruser/yourproject/pkg"   // Local package
```

#### Problem: Module not found
```bash
go mod tidy
# Error: module github.com/example/nonexistent@latest: not found
```

**Solutions:**
```bash
# 1. Check module name spelling
go list -m all

# 2. Use correct version
go get github.com/example/package@v1.2.3

# 3. Clean module cache if corrupted
go clean -modcache
go mod download
```

### Syntax Errors

#### Problem: Common syntax mistakes
```go
// ❌ Common errors
func main() {
    if x := 5 {        // Missing condition
        fmt.Println("Error")
    }
    
    var x int = "hello"  // Type mismatch
    
    fmt.Println("Hello"  // Missing closing parenthesis
}
```

**Solutions:**
```go
// ✅ Correct syntax
func main() {
    if x := 5; x > 0 {   // Proper condition
        fmt.Println("Correct")
    }
    
    var x int = 5        // Correct type
    // or
    var x = "hello"      // Type inference
    
    fmt.Println("Hello") // Proper parentheses
}
```

## 💥 Runtime Errors

### Panic and Recovery

#### Problem: Nil pointer dereference
```go
// ❌ Common panic cause
var ptr *int
fmt.Println(*ptr)  // Panic: nil pointer dereference
```

**Solutions:**
```go
// ✅ Safe approaches
var ptr *int
if ptr != nil {
    fmt.Println(*ptr)
} else {
    fmt.Println("Pointer is nil")
}

// Or initialize properly
value := 42
ptr = &value
fmt.Println(*ptr)  // Safe
```

#### Problem: Index out of bounds
```go
// ❌ Panic cause
slice := []int{1, 2, 3}
fmt.Println(slice[5])  // Panic: index out of range
```

**Solutions:**
```go
// ✅ Safe access
slice := []int{1, 2, 3}
index := 5

if index < len(slice) {
    fmt.Println(slice[index])
} else {
    fmt.Println("Index out of bounds")
}

// Or use range
for i, v := range slice {
    fmt.Printf("Index %d: %v\n", i, v)
}
```

### Concurrency Issues

#### Problem: Deadlock
```go
// ❌ Deadlock example
func main() {
    ch := make(chan int)
    ch <- 1  // Blocks forever - no receiver
    fmt.Println(<-ch)
}
```

**Solutions:**
```go
// ✅ Solutions
// 1. Buffered channel
func main() {
    ch := make(chan int, 1)  // Buffer size 1
    ch <- 1
    fmt.Println(<-ch)
}

// 2. Goroutine
func main() {
    ch := make(chan int)
    go func() {
        ch <- 1
    }()
    fmt.Println(<-ch)
}

// 3. Select with default
func main() {
    ch := make(chan int)
    select {
    case ch <- 1:
        fmt.Println("Sent")
    default:
        fmt.Println("Would block")
    }
}
```

#### Problem: Race conditions
```go
// ❌ Race condition
var counter int

func increment() {
    counter++  // Not atomic
}

func main() {
    for i := 0; i < 1000; i++ {
        go increment()
    }
    time.Sleep(time.Second)
    fmt.Println(counter)  // Unpredictable result
}
```

**Solutions:**
```go
// ✅ Thread-safe solutions
import "sync"

// 1. Mutex
var (
    counter int
    mutex   sync.Mutex
)

func increment() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}

// 2. Atomic operations
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

// 3. Channel-based approach
func main() {
    counterCh := make(chan int)
    
    // Worker goroutine
    go func() {
        count := 0
        for range counterCh {
            count++
        }
        fmt.Println("Final count:", count)
    }()
    
    // Send increments
    for i := 0; i < 1000; i++ {
        go func() {
            counterCh <- 1
        }()
    }
}
```

## 🧪 Testing Issues

### Test Discovery and Execution

#### Problem: Tests not running
```bash
go test
# no test files
```

**Solutions:**
```bash
# 1. Check file naming - must end with _test.go
ls *_test.go

# 2. Check function naming - must start with Test
grep "func Test" *_test.go

# 3. Run specific tests
go test -run TestSpecificFunction

# 4. Run tests in all subdirectories
go test ./...
```

#### Problem: Test failures with unclear messages
```go
// ❌ Unclear test
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Error("Test failed")  // Unclear what failed
    }
}
```

**Solutions:**
```go
// ✅ Clear test messages
func TestAdd(t *testing.T) {
    a, b := 2, 3
    expected := 5
    result := Add(a, b)
    
    if result != expected {
        t.Errorf("Add(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// ✅ Table-driven tests for multiple cases
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 5, 5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Mock and Testing Setup Issues

#### Problem: Mock generation failures
```bash
mockgen -source=interface.go -destination=mock.go
# Error: failed to load package
```

**Solutions:**
```bash
# 1. Ensure interface file compiles
go build interface.go

# 2. Use correct mockgen syntax
mockgen -source=interface.go -destination=mock.go -package=main

# 3. Generate from import path
mockgen -destination=mock.go package_path InterfaceName

# 4. Install/update mockgen
go install go.uber.org/mock/mockgen@latest
```

## 🐌 Performance Issues

### Memory Problems

#### Problem: Memory leaks
```bash
# Signs of memory leaks
go tool pprof http://localhost:6060/debug/pprof/heap
# Shows increasing memory usage
```

**Debugging steps:**
```go
// 1. Add pprof to your application
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Your application code
}
```

```bash
# 2. Profile memory usage
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# 3. Look for growing allocations
# 4. Check for:
#    - Unclosed resources (files, connections)
#    - Growing slices/maps
#    - Goroutine leaks
```

**Common fixes:**
```go
// ✅ Close resources properly
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Always close
    
    // Process file
    return nil
}

// ✅ Limit slice growth
func processItems(items []string) {
    // Pre-allocate if size is known
    results := make([]string, 0, len(items))
    
    for _, item := range items {
        if shouldProcess(item) {
            results = append(results, processItem(item))
        }
    }
}

// ✅ Clean up goroutines
func worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case job := <-jobs:
            processJob(job)
        case <-ctx.Done():
            return  // Exit when context is cancelled
        }
    }
}
```

### CPU Performance Issues

#### Problem: Slow performance
```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

**Common optimizations:**
```go
// ❌ Inefficient string concatenation
func inefficient(items []string) string {
    result := ""
    for _, item := range items {
        result += item  // Creates new string each time
    }
    return result
}

// ✅ Efficient string building
func efficient(items []string) string {
    var builder strings.Builder
    builder.Grow(estimatedSize)  // Pre-allocate if possible
    
    for _, item := range items {
        builder.WriteString(item)
    }
    return builder.String()
}

// ❌ Inefficient map lookups
func inefficient(data map[string]int, keys []string) {
    for _, key := range keys {
        if _, exists := data[key]; exists {
            value := data[key]  // Double lookup
            process(value)
        }
    }
}

// ✅ Single lookup
func efficient(data map[string]int, keys []string) {
    for _, key := range keys {
        if value, exists := data[key]; exists {
            process(value)
        }
    }
}
```

## 📦 Dependency Issues

### Module Resolution Problems

#### Problem: Version conflicts
```bash
go mod tidy
# Error: module A requires B@v1.0.0, but module C requires B@v2.0.0
```

**Solutions:**
```bash
# 1. Check dependency tree
go mod graph | grep problematic-package

# 2. Force specific version
go mod edit -require=package@version

# 3. Use replace directive for testing
go mod edit -replace=old-package=new-package@version

# 4. Update dependencies
go get -u ./...
```

#### Problem: Vendor directory issues
```bash
go build
# Error: package is not in vendor/ but should be
```

**Solutions:**
```bash
# 1. Rebuild vendor directory
go mod vendor

# 2. Ensure vendor is up to date
go mod tidy
go mod vendor

# 3. Or remove vendor and use modules
rm -rf vendor/
# Edit go.mod to remove vendor references
```

## 🛠️ Debugging Techniques

### Using Delve Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a program
dlv debug main.go

# Debug tests
dlv test

# Attach to running process
dlv attach <pid>
```

**Common Delve commands:**
```
(dlv) break main.main     # Set breakpoint
(dlv) continue            # Continue execution
(dlv) next               # Step over
(dlv) step               # Step into
(dlv) print variable     # Print variable value
(dlv) stack              # Show stack trace
(dlv) goroutines         # List goroutines
(dlv) goroutine 1 bt     # Backtrace for goroutine 1
```

### Print Debugging

```go
// Simple debugging
fmt.Printf("Debug: variable = %+v\n", variable)

// Structured logging
import "log"

log.Printf("Processing item %d: %+v", i, item)

// JSON pretty printing
import "encoding/json"

func debugJSON(v interface{}) {
    b, _ := json.MarshalIndent(v, "", "  ")
    fmt.Println(string(b))
}
```

### Runtime Debugging

```go
import "runtime"

// Get stack trace
func printStackTrace() {
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    fmt.Printf("Stack trace:\n%s", buf[:n])
}

// Get goroutine info
func goroutineInfo() {
    fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}

// Memory stats
func memoryStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("Allocated memory: %d KB\n", m.Alloc/1024)
    fmt.Printf("Total allocations: %d\n", m.TotalAlloc)
    fmt.Printf("Number of GC cycles: %d\n", m.NumGC)
}
```

## 🎯 IDE and Editor Issues

### VS Code Go Extension Problems

#### Problem: "Go tools are missing"
**Solution:**
```bash
# Install Go tools
Ctrl+Shift+P -> "Go: Install/Update Tools"
# Select all tools and install
```

#### Problem: Import suggestions not working
**Solutions:**
1. Ensure `gopls` is installed and updated
2. Check VS Code Go extension settings
3. Restart VS Code and Go language server
4. Clear module cache: `go clean -modcache`

### Module-aware editing issues

#### Problem: Red underlines on valid imports
**Solutions:**
1. Ensure you're in a module directory (`go.mod` exists)
2. Run `go mod tidy`
3. Restart language server
4. Check if imports are accessible from current module

## 📚 Getting Help

### Documentation and Resources

```bash
# Built-in help
go help                    # General help
go help build             # Help for specific command
go doc fmt.Println        # Documentation for function
go doc -all fmt           # All documentation for package

# Online resources
# https://golang.org/doc/
# https://pkg.go.dev/
# https://forum.golangbridge.org/
```

### Community Support

1. **Go Forum**: https://forum.golangbridge.org/
2. **Reddit**: https://reddit.com/r/golang
3. **Gophers Slack**: https://gophers.slack.com/
4. **Stack Overflow**: Tag questions with `go` or `golang`
5. **GitHub Issues**: For specific project issues

### Reporting Bugs

When reporting issues, include:
1. **Go version**: `go version`
2. **Operating system**: `go env GOOS GOARCH`
3. **Minimal reproducible example**
4. **Expected vs actual behavior**
5. **Full error messages**
6. **Steps to reproduce**

---

## 🔄 Prevention Tips

### Best Practices to Avoid Common Issues

1. **Use `go mod tidy` regularly** to keep dependencies clean
2. **Write tests early** to catch issues during development  
3. **Use `go vet`** to catch common mistakes
4. **Enable race detection** during testing: `go test -race`
5. **Use linters** like `golangci-lint` for code quality
6. **Handle errors explicitly** - don't ignore them
7. **Close resources** with `defer` statements
8. **Use context for cancellation** in long-running operations

### Debugging Checklist

When encountering issues:
- [ ] Can you reproduce the issue consistently?
- [ ] Have you checked the Go version and environment?
- [ ] Are all dependencies up to date?
- [ ] Do the tests pass?
- [ ] Have you run `go vet` and linters?
- [ ] Have you checked for race conditions?
- [ ] Are you handling errors properly?
- [ ] Have you checked the documentation?

Remember: Every Go developer faces these issues. The key is systematic debugging and learning from each problem! 🚀