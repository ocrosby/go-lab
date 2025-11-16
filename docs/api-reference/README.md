# API Reference

Auto-generated and curated API documentation for Go Laboratory packages.

## Overview

This directory contains detailed API documentation for the various packages and modules in the Go Laboratory repository.

## Documentation Structure

### Package Documentation
- **Automatically Generated**: Using `go doc` and `godoc`
- **Code Examples**: Practical usage examples for each package
- **Cross-References**: Links to related packages and concepts

### How to Generate Documentation

#### Local Documentation Server
```bash
# Install godoc (if not already installed)
go install golang.org/x/tools/cmd/godoc@latest

# Start local documentation server
godoc -http=:8080

# View in browser: http://localhost:8080
```

#### Command Line Documentation
```bash
# Package-level documentation
go doc package-name

# Function-specific documentation  
go doc package-name.FunctionName

# Examples
go doc fmt.Println
go doc testing.T
```

## Package Index

### Learning Packages

#### Fundamentals
- **math**: Basic mathematical operations with comprehensive testing
  ```bash
  go doc ./learning/fundamentals/math
  ```

#### Intermediate  
- **composition**: Object composition and embedding patterns
  ```bash
  go doc ./learning/intermediate/composition
  ```

#### Advanced
- **concurrency**: Goroutines, channels, and synchronization primitives
- **dependency-injection**: Clean architecture and DI patterns  
- **patterns**: Design pattern implementations

### Example Packages

#### Calculator
- **v1**: Basic calculator with standard testing
- **v2**: Enhanced calculator with BDD testing

#### HTTP Services
- **jsonplaceholder**: REST API client with comprehensive models
- **server**: HTTP server implementation patterns

### Production Packages

#### API Project
- **application**: Business logic and use cases
- **domain**: Core domain entities and interfaces
- **infrastructure**: External service adapters
- **validation**: Input validation and sanitization

## Documentation Standards

### Package Documentation
Each package should include:

```go
// Package calculator implements basic mathematical operations
// with comprehensive error handling and testing support.
//
// This package demonstrates fundamental Go concepts including
// function definitions, error handling, and testing patterns.
//
// Example usage:
//
//     result := calculator.Add(2, 3)
//     fmt.Println(result) // Output: 5
//
//     quotient, err := calculator.Divide(10, 2) 
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println(quotient) // Output: 5
package calculator
```

### Function Documentation
```go
// Add returns the sum of two integers.
//
// This function performs basic addition with no error conditions.
// Both positive and negative integers are supported.
//
// Example:
//     sum := Add(5, 3)  // Returns 8
//     sum := Add(-2, 7) // Returns 5
func Add(a, b int) int {
    return a + b
}

// Divide returns the quotient of two integers and an error if division by zero.
//
// The function returns an error if the divisor is zero to prevent
// division by zero panics. The result is truncated to an integer.
//
// Example:
//     result, err := Divide(10, 2)
//     if err != nil {
//         // Handle error
//     }
//     // result == 5
func Divide(dividend, divisor int) (int, error) {
    if divisor == 0 {
        return 0, errors.New("division by zero")
    }
    return dividend / divisor, nil
}
```

### Type Documentation
```go
// User represents a system user with authentication capabilities.
//
// User contains all necessary information for user management
// including authentication, authorization, and profile data.
type User struct {
    // ID is the unique identifier for the user
    ID int `json:"id"`
    
    // Username is the unique username for authentication  
    Username string `json:"username"`
    
    // Email is the user's email address
    Email string `json:"email"`
    
    // CreatedAt is the timestamp when the user was created
    CreatedAt time.Time `json:"created_at"`
}

// String returns a string representation of the User.
// Implements the fmt.Stringer interface.
func (u User) String() string {
    return fmt.Sprintf("User{ID: %d, Username: %s}", u.ID, u.Username)
}
```

## Browsing Documentation

### Online Documentation
For published packages, documentation is available at:
- [pkg.go.dev](https://pkg.go.dev/) for public packages
- Local godoc server for private/local development

### IDE Integration
Most Go IDEs provide built-in documentation:
- **VS Code**: Hover over functions/types
- **GoLand**: Ctrl+Q (Quick Documentation)
- **Vim/Neovim**: Using LSP plugins

### Command Line Access
```bash
# Quick function lookup
go doc fmt.Printf

# Package overview
go doc net/http

# All package functions
go doc -all package-name

# Source code
go doc -src package-name.Function
```

## Generating Custom Documentation

### For This Repository
```bash
# Generate HTML documentation for all packages
godoc -html > docs/api-reference/generated.html

# Generate for specific package
go doc -all ./learning/fundamentals/math > docs/api-reference/math.md
```

### Documentation Comments Best Practices

1. **Start with the name**: "Add returns..." not "This function adds..."
2. **Be concise but complete**: Explain what, when, and why
3. **Include examples**: Show typical usage patterns
4. **Document edge cases**: Explain error conditions
5. **Use proper formatting**: Follow godoc conventions

### Example Documentation
See these well-documented packages:
- [`/learning/fundamentals/math/`](../../learning/fundamentals/math/) - Basic documentation
- [`/examples/http-services/jsonplaceholder/`](../../examples/http-services/jsonplaceholder/) - Service documentation
- [`/projects/api/`](../../projects/api/) - Production-level documentation

## Automated Documentation Generation

### Using go generate
```go
//go:generate go doc -all . > api-docs.txt
```

### CI/CD Integration
```yaml
# Example GitHub Actions step
- name: Generate Documentation
  run: |
    go install golang.org/x/tools/cmd/godoc@latest
    godoc -html > docs/api-reference/generated.html
```

## Contributing to API Documentation

When adding new packages:
1. Include comprehensive package-level documentation
2. Document all exported functions and types
3. Provide practical examples
4. Follow established documentation patterns
5. Update this index with new packages

Remember: Good documentation is as important as good code!