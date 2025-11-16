# Complete Guide to Testing in Go

Master testing strategies from basic unit tests to advanced mocking patterns.

## Table of Contents
1. [Testing Fundamentals](#testing-fundamentals)
2. [Table-Driven Tests](#table-driven-tests)
3. [BDD Testing with Ginkgo](#bdd-testing-with-ginkgo)
4. [Mocking and Interfaces](#mocking-and-interfaces)
5. [Testing HTTP Services](#testing-http-services)
6. [Advanced Testing Patterns](#advanced-testing-patterns)

## Testing Fundamentals

### Basic Test Structure
```go
package calculator

import "testing"

func TestAdd(t *testing.T) {
    // Arrange
    a, b := 2, 3
    expected := 5
    
    // Act
    result := Add(a, b)
    
    // Assert
    if result != expected {
        t.Errorf("Add(%d, %d) = %d; want %d", a, b, result, expected)
    }
}
```

### Running Tests
```bash
go test                    # Current package
go test -v                 # Verbose output
go test ./...              # All packages
go test -run TestAdd       # Specific test
go test -cover             # With coverage
```

### Test File Conventions
- Test files end with `_test.go`
- Test functions start with `Test`
- Test functions take `*testing.T` parameter
- Place tests in the same package as code being tested

## Table-Driven Tests

Efficient way to test multiple scenarios:

```go
func TestCalculatorOperations(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        op       string
        expected int
        hasError bool
    }{
        {"add positive", 2, 3, "add", 5, false},
        {"subtract", 5, 3, "subtract", 2, false},
        {"multiply", 4, 3, "multiply", 12, false},
        {"divide", 10, 2, "divide", 5, false},
        {"divide by zero", 10, 0, "divide", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var result int
            var err error
            
            switch tt.op {
            case "add":
                result = Add(tt.a, tt.b)
            case "divide":
                result, err = Divide(tt.a, tt.b)
            }
            
            if tt.hasError && err == nil {
                t.Errorf("Expected error but got none")
            }
            
            if !tt.hasError && result != tt.expected {
                t.Errorf("Expected %d, got %d", tt.expected, result)
            }
        })
    }
}
```

## BDD Testing with Ginkgo

### Installation
```bash
go get github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/...
go install github.com/onsi/ginkgo/v2/ginkgo
```

### Bootstrap Test Suite
```bash
cd your-package/
ginkgo bootstrap        # Creates package_suite_test.go
ginkgo generate calc    # Creates calc_test.go
```

### BDD Test Structure
```go
package calculator_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    . "your-package/calculator"
)

var _ = Describe("Calculator", func() {
    Context("when adding numbers", func() {
        It("should return the sum of positive numbers", func() {
            result := Add(2, 3)
            Expect(result).To(Equal(5))
        })
        
        It("should handle negative numbers", func() {
            result := Add(-2, 3)
            Expect(result).To(Equal(1))
        })
    })
    
    Context("when dividing numbers", func() {
        It("should return the quotient", func() {
            result, err := Divide(10, 2)
            Expect(err).NotTo(HaveOccurred())
            Expect(result).To(Equal(5))
        })
        
        It("should return error for division by zero", func() {
            _, err := Divide(10, 0)
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("division by zero"))
        })
    })
})
```

### Running Ginkgo Tests
```bash
ginkgo -v                   # Verbose output
ginkgo -r                   # Recursive (all packages)
ginkgo --focus="adding"     # Run specific tests
ginkgo --skip="slow"        # Skip certain tests
```

## Mocking and Interfaces

### Define Interfaces for Testability
```go
// Define interface
type UserService interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// Real implementation
type userService struct {
    db Database
}

func (s *userService) GetUser(id int) (*User, error) {
    return s.db.FindUser(id)
}
```

### Generate Mocks with uber-go/mock
```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mocks
mockgen -source=user_service.go -destination=mocks/mock_user_service.go
```

### Using Mocks in Tests
```go
func TestUserHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockUserService(ctrl)
    handler := NewUserHandler(mockService)
    
    // Setup expectations
    expectedUser := &User{ID: 1, Name: "John"}
    mockService.EXPECT().
        GetUser(1).
        Return(expectedUser, nil).
        Times(1)
    
    // Execute test
    user, err := handler.HandleGetUser(1)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
}
```

## Testing HTTP Services

### Testing HTTP Handlers
```go
func TestUserHandler(t *testing.T) {
    // Create request
    req := httptest.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    
    // Setup handler with mock service
    mockService := &MockUserService{}
    handler := NewUserHandler(mockService)
    
    // Execute
    handler.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var user User
    err := json.Unmarshal(w.Body.Bytes(), &user)
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### Testing HTTP Clients
```go
func TestHTTPClient(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(User{ID: 1, Name: "John"})
    }))
    defer server.Close()
    
    // Test client
    client := NewAPIClient(server.URL)
    user, err := client.GetUser(1)
    
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

## Advanced Testing Patterns

### Subtests for Organization
```go
func TestUserService(t *testing.T) {
    t.Run("GetUser", func(t *testing.T) {
        t.Run("existing user", func(t *testing.T) {
            // Test existing user logic
        })
        
        t.Run("non-existing user", func(t *testing.T) {
            // Test error case
        })
    })
    
    t.Run("CreateUser", func(t *testing.T) {
        // Test user creation
    })
}
```

### Benchmarking
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

func BenchmarkComplexOperation(b *testing.B) {
    data := setupBenchmarkData()
    b.ResetTimer() // Don't count setup time
    
    for i := 0; i < b.N; i++ {
        ComplexOperation(data)
    }
}
```

### Example Tests (Documentation)
```go
func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleCalculator_Divide() {
    calc := &Calculator{}
    result, err := calc.Divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
    // Output: 5
}
```

### Testing with External Dependencies
```go
func TestWithDatabase(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test with real database
    service := NewUserService(db)
    user, err := service.CreateUser(&User{Name: "John"})
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

## Best Practices

### Test Organization
1. **AAA Pattern**: Arrange, Act, Assert
2. **One Assertion Per Test**: Focus on single behavior
3. **Descriptive Names**: Test names should explain what's being tested
4. **Independent Tests**: Tests shouldn't depend on each other

### Mocking Guidelines
1. **Mock External Dependencies**: Not internal logic
2. **Use Interfaces**: Enable easy mocking
3. **Verify Interactions**: When behavior matters more than state
4. **Don't Over-Mock**: Mock only what you need

### Coverage and Quality
```bash
# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Coverage by function
go test -cover -v

# Set minimum coverage
go test -cover | grep -E '^coverage:' | awk '{print $2}' | grep -E '^[0-9]+\.[0-9]+%$'
```

## Integration with Examples

### Practical Examples
- **Basic Testing**: [`/learning/fundamentals/math/`](../../learning/fundamentals/math/)
- **BDD Testing**: [`/examples/calculator/v2/`](../../examples/calculator/v2/)
- **Mocking**: [`/testing/mocking/`](../../testing/mocking/)
- **HTTP Testing**: [`/examples/http-services/`](../../examples/http-services/)

### Next Steps
1. Practice with different testing frameworks
2. Learn property-based testing with [gopter](https://github.com/leanovate/gopter)
3. Explore mutation testing
4. Study the production API tests in [`/projects/api/`](../../projects/api/)

## Common Pitfalls

1. **Testing Implementation, Not Behavior**: Focus on what, not how
2. **Fragile Tests**: Tests that break with minor refactoring
3. **Slow Tests**: Keep unit tests fast, integration tests separate
4. **Poor Test Data**: Use realistic but minimal test data
5. **Ignoring Edge Cases**: Test boundary conditions and error paths

Remember: Good tests serve as documentation and safety nets for refactoring!