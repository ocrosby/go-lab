# Calculator Examples

Mathematical operations implemented with iterative development showing testing evolution.

## Versions

### 📊 [Version 1 (v1)](./v1/)
Basic calculator implementation with fundamental testing:
- **Files**: `calculator.go`, `calculator_test.go`
- **Testing**: Standard Go testing package
- **Operations**: Add, Subtract, Multiply, Divide
- **Focus**: Core functionality and basic test coverage

```bash
cd v1/
go test -v
```

### 🎯 [Version 2 (v2)](./v2/)
Enhanced calculator with BDD testing framework:
- **Files**: `calculator.go`, `calculator_test.go`, `calculator_suite_test.go`
- **Testing**: Ginkgo & Gomega BDD framework
- **Operations**: Same core operations with enhanced testing
- **Focus**: Behavior-driven development, expressive test syntax

```bash
cd v2/
go test -v
# or with ginkgo
ginkgo -v
```

## Evolution Demonstrated

### Testing Framework Progression
1. **v1**: Standard `testing` package with table-driven tests
2. **v2**: Ginkgo/Gomega BDD with descriptive test scenarios

### Key Differences

**v1 Testing Style:**
```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}
```

**v2 BDD Style:**
```go
Describe("Calculator", func() {
    Context("when adding two numbers", func() {
        It("should return the sum", func() {
            result := Add(2, 3)
            Expect(result).To(Equal(5))
        })
    })
})
```

## Learning Objectives

### From v1 (Fundamentals)
- Basic function implementation
- Standard Go testing patterns
- Test-driven development basics
- Error handling in mathematical operations

### From v2 (Intermediate)
- BDD testing methodology
- Expressive test organization
- Advanced testing frameworks
- Test suite management

## Running Examples

```bash
# Version 1 - Standard testing
cd v1/
go run calculator.go
go test -v

# Version 2 - BDD testing  
cd v2/
go run calculator.go
go test -v
ginkgo -v  # if ginkgo is installed
```

## Key Concepts

### Mathematical Operations
- **Addition**: Sum of two numbers
- **Subtraction**: Difference between numbers
- **Multiplication**: Product of two numbers
- **Division**: Quotient with zero-division handling

### Error Handling
Both versions demonstrate:
- Division by zero protection
- Proper error return patterns
- Test coverage for error conditions

### Testing Patterns
- **Positive Cases**: Valid operations
- **Negative Cases**: Error conditions
- **Edge Cases**: Boundary value testing

## Installing Ginkgo (for v2)

### Step 1. Setup Environment
Add to your shell configuration:
```shell
export GOROOT=~/go
export GOBIN=${GOROOT}/bin
export PATH=${GOBIN}:${PATH}
```

### Step 2. Install Libraries
```bash
go get github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/...
```

### Step 3. Install Ginkgo CLI
```bash
go install github.com/onsi/ginkgo/v2/ginkgo
```

### Step 4. Bootstrap Test Suite
```bash
cd v2/
ginkgo bootstrap  # Creates calculator_suite_test.go
ginkgo generate calculator  # Creates calculator_test.go
```

## Next Steps

After mastering these calculator examples:
1. Explore [HTTP Services](../http-services/) for web-based calculators
2. Study [Advanced Patterns](../../learning/advanced/) for production architectures
3. Review [Testing Strategies](../../testing/) for comprehensive test approaches

