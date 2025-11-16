# Go Fundamentals

Master the basics of Go programming language.

## Prerequisites
- Basic programming knowledge in any language
- Go installed on your system

## Learning Objectives
By the end of this section, you will:
- Understand Go syntax and basic data types
- Write simple Go programs
- Handle input/output operations
- Perform basic mathematical operations
- Set up a proper Go development environment

## Modules

### 1. [Hello World](./hello/)
Your first Go program demonstrating:
- Basic program structure
- Package declaration
- Importing standard libraries
- The `main` function

**Time**: 30 minutes  
**Files**: `hello.go`, `README.md`

### 2. [Mathematics](./math/)
Basic mathematical operations in Go:
- Numeric data types
- Arithmetic operations
- Function definitions
- Testing with Go's testing package

**Time**: 1 hour  
**Files**: `math.go`, `math_test.go`

## Running the Examples

```bash
# Navigate to any module directory
cd hello/

# Run the program
go run hello.go

# For modules with tests
go test -v
```

## Key Concepts Covered
- Variables and constants
- Data types (int, float, string, bool)
- Functions and return values
- Basic testing
- Go module structure

## Next Steps
Once comfortable with these fundamentals, proceed to [Intermediate](../02-intermediate/) concepts including structs, interfaces, and composition patterns.

## Practice Examples
Apply these concepts with practical examples:
- [Beginner Examples](../../examples/01-beginner/) - Calculator and Hello Concurrent
- [Getting Started Guide](../../docs/tutorials/getting-started.md) - Environment setup and first programs

## Common Gotchas
- Go is statically typed - variable types must be known at compile time
- Unused imports and variables cause compilation errors
- The `main` package and `main()` function are required for executable programs
- Go uses package-level visibility (capitalized = public, lowercase = private)