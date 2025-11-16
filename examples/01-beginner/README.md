# Beginner Examples

Perfect starting examples for those new to Go programming.

## Prerequisites
- Basic programming knowledge in any language
- Go installed and working (see [Getting Started Guide](../../docs/tutorials/getting-started.md))
- Completed [Learning Fundamentals](../../learning/fundamentals/) (recommended)

## Examples in This Directory

### 🧮 [Calculator](./calculator/)
**Estimated Time**: 2-3 hours  
**Concepts**: Functions, testing, error handling, package organization

Mathematical operations showing evolution from basic testing to BDD:
- **v1/**: Standard Go testing with table-driven tests
- **v2/**: Ginkgo/Gomega BDD testing framework

**What You'll Learn**:
- Function definition and calling
- Error handling with multiple return values
- Writing and running tests
- Table-driven test patterns
- Testing framework evolution

```bash
cd calculator/v1/
go run calculator.go
go test -v

cd ../v2/
go test -v
ginkgo -v  # if ginkgo installed
```

### ⚡ [Hello Concurrent](./hello-concurrent/)
**Estimated Time**: 1-2 hours  
**Concepts**: Goroutines, basic concurrency, program flow

Introduction to Go's concurrency model:
- Basic goroutine usage
- Understanding concurrent execution
- Simple synchronization patterns

**What You'll Learn**:
- Creating and running goroutines
- Concurrent vs sequential execution
- Basic synchronization concepts
- Go's approach to concurrency

```bash
cd hello-concurrent/
go run main.go
```

## Learning Path

### Recommended Order
1. **Calculator v1**: Master basic Go syntax and testing
2. **Hello Concurrent**: Understand Go's concurrency basics
3. **Calculator v2**: Learn advanced testing frameworks

### Skills Developed
By completing these examples, you'll understand:
- ✅ Go syntax and basic data types
- ✅ Function definition and error handling
- ✅ Testing with the standard library
- ✅ Basic concurrent programming
- ✅ Package organization
- ✅ BDD testing frameworks

## Common Challenges & Solutions

### "go: command not found"
- Verify Go installation: `go version`
- Check PATH includes Go binary directory
- See [Environment Setup](../../docs/tutorials/getting-started.md#step-2-environment-setup)

### Import Path Issues
- Ensure you're in the correct directory
- Run `go mod tidy` to resolve dependencies
- Initialize module if needed: `go mod init example`

### Test Failures
- Read error messages carefully
- Check function signatures match expected parameters
- Verify return types and error handling

### Ginkgo Not Found (Calculator v2)
```bash
# Install Ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Verify installation
ginkgo version
```

## Key Concepts Reinforced

### Go Fundamentals
- **Package Declaration**: Every Go file starts with `package`
- **Import Statements**: Bringing in standard library and external packages
- **Function Syntax**: Parameters, return types, and naming conventions
- **Error Handling**: Go's explicit error handling pattern

### Testing Patterns
- **Test File Naming**: `*_test.go` files
- **Test Function Naming**: `TestXxx(t *testing.T)`
- **Table-Driven Tests**: Efficient multiple test case handling
- **BDD Style**: Descriptive test organization with Ginkgo

### Concurrency Basics
- **Goroutines**: Lightweight threads with `go` keyword
- **Program Flow**: Understanding when programs exit
- **Resource Sharing**: Basic concepts of concurrent access

## Next Steps

After mastering these examples:

### Immediate Next Steps
1. **Modify Examples**: Change functionality to test understanding
2. **Add Features**: Extend calculator with more operations
3. **Error Handling**: Add robust error handling to concurrent example

### Ready for Intermediate?
Check your readiness:
- ✅ Can write functions with proper error handling
- ✅ Understand how to structure and run tests
- ✅ Comfortable with basic goroutine usage
- ✅ Know package organization principles

If yes, proceed to [Intermediate Examples](../intermediate/)

### Still Need Practice?
- Review [Learning Fundamentals](../../learning/fundamentals/)
- Try the [Testing Guide](../../docs/tutorials/testing-guide.md)
- Experiment more with the current examples

## Getting Help

### Within This Repository
- [Documentation](../../docs/) - Comprehensive guides
- [Learning Path](../../learning/) - Structured theoretical content
- [Testing Examples](../../testing/) - Advanced testing patterns

### External Resources
- [Go Tour](https://tour.golang.org/) - Interactive Go tutorial
- [Go by Example](https://gobyexample.com/) - Code examples
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices

### Community
- [Go Forum](https://forum.golangbridge.org/)
- [r/golang](https://reddit.com/r/golang) 
- [Gophers Slack](https://gophers.slack.com/)

Remember: Take your time with these examples. Understanding the fundamentals well will make advanced concepts much easier to grasp!