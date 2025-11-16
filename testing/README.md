# Go Testing Examples

Comprehensive testing strategies, patterns, and frameworks in Go.

## Testing Approaches Covered

### 🎭 [Mocking](./mocking/)
Mock generation and testing with dependencies:
- **Cars & Trucks**: Mock implementations for automotive domain
- **Generated Mocks**: Using uber-go/mock for automatic mock generation
- **Interface Testing**: Testing through contracts rather than implementations
- Demonstrates: Dependency injection, interface-based testing, mock strategies

### 📊 [Test Patterns](./test/)
Specialized testing patterns and frameworks:
- **Channels**: Testing concurrent Go code with channels
- **Ginkgo/Gomega**: BDD-style testing examples
- Demonstrates: Concurrent testing, behavior-driven development, test organization

## Testing Frameworks Demonstrated

### Standard Testing
- Built-in `testing` package
- Table-driven tests
- Benchmarks and examples

### Ginkgo & Gomega
- Behavior-driven development (BDD)
- Expressive test syntax
- Advanced matchers and assertions

### Mock Generation
- uber-go/mock integration
- Interface-based mocking
- Dependency injection testing

## Key Testing Concepts

### Test Organization
```go
func TestSomething(t *testing.T) {
    // Arrange
    // Act  
    // Assert
}
```

### Table-Driven Tests
```go
tests := []struct {
    name     string
    input    int
    expected int
}{
    {"positive", 5, 5},
    {"negative", -3, 3},
}
```

### Mock Usage
```go
mockObj := mocks.NewMockInterface(ctrl)
mockObj.EXPECT().Method().Return(value)
```

## Running Tests

```bash
# Standard tests
go test ./...

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Specific package
cd mocking/cars/honda/
go test -v

# Ginkgo tests
cd test/channels/
ginkgo -v
```

## Best Practices Demonstrated

### Test Structure
- **AAA Pattern**: Arrange, Act, Assert
- **Given-When-Then**: BDD style with Ginkgo
- **Table-Driven**: Multiple test cases efficiently

### Mock Usage
- Mock external dependencies, not internal logic
- Use interfaces to enable mocking
- Test behavior, not implementation details
- Verify interactions when necessary

### Test Coverage
- Aim for high coverage but focus on critical paths
- Test edge cases and error conditions
- Integration tests for end-to-end flows
- Unit tests for individual components

## Integration with Learning Path

- **Fundamentals**: Basic testing with `testing` package
- **Intermediate**: Interface-based testing and simple mocks
- **Advanced**: Complex mocking strategies, BDD, concurrent testing

## Production Testing

For production-ready testing strategies, see:
- [`/projects/api/`](../projects/api/) - Complete testing suite
- Quality engineering practices
- CI/CD integration examples