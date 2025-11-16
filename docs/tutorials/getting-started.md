# Getting Started with Go Laboratory

A step-by-step guide to begin your Go learning journey.

## Prerequisites

- Basic programming knowledge in any language
- Computer with internet access
- Text editor or IDE

## Step 1: Install Go

### Option A: Official Installation
1. Visit [golang.org/dl/](https://golang.org/dl/)
2. Download the installer for your operating system
3. Follow the installation instructions

### Option B: Using Package Managers

**macOS (Homebrew):**
```bash
brew install go
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**Windows (Chocolatey):**
```bash
choco install golang
```

### Verify Installation
```bash
go version
# Should output: go version go1.21.x darwin/amd64 (or similar)
```

## Step 2: Environment Setup

### Modern Go (Go 1.11+)
With Go modules, you no longer need to set `GOPATH`. However, you may want to customize:

```bash
# Optional: Set where Go installs binaries
export GOBIN=$HOME/go/bin
export PATH=$PATH:$GOBIN
```

### Legacy Go Workspace (Pre-1.11)
If working with older Go code:

```bash
export GOROOT=/usr/local/go  # Where Go is installed
export GOPATH=$HOME/go       # Your workspace
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Step 3: Clone the Go Laboratory

```bash
# Clone the repository
git clone https://github.com/your-username/go-lab.git
cd go-lab

# Initialize Go module (if needed)
go mod tidy
```

## Step 4: Your First Go Program

### Run Hello World
```bash
cd learning/fundamentals/hello/
go run hello.go
```

Expected output:
```
Hello, World!
```

### Understanding the Code
```go
package main          // Executable package

import "fmt"          // Import standard library

func main() {         // Entry point
    fmt.Println("Hello, World!")
}
```

## Step 5: Follow the Learning Path

### Recommended Progression
1. **[Fundamentals](../../learning/01-fundamentals/)** (1-2 weeks)
   - Basic syntax and data types
   - Functions and control structures  
   - Simple applications
   - **Practice with**: [Beginner Examples](../../examples/01-beginner/)

2. **[Intermediate](../../learning/02-intermediate/)** (2-3 weeks)
   - Structs and methods
   - Interfaces and composition
   - Object-oriented patterns
   - **Practice with**: [Intermediate Examples](../../examples/02-intermediate/)

3. **[Advanced](../../learning/03-advanced/)** (4-6 weeks)
   - Concurrency and channels
   - Design patterns
   - Production architectures
   - **Practice with**: [Advanced Examples](../../examples/03-advanced/)

## Step 6: Practice with Examples

### Beginner Examples
```bash
cd examples/01-beginner/calculator/v1/
go run calculator.go
go test -v
```

### Test Your Understanding
Try modifying the calculator to:
- Add more operations (power, modulo)
- Handle floating-point numbers
- Add input validation

## Step 7: Set Up Your Development Environment

### Recommended Editors
- **VS Code**: Install the Go extension
- **GoLand**: JetBrains Go IDE
- **Vim/Neovim**: With Go plugins

### Essential Tools
```bash
# Code formatting
go install golang.org/x/tools/cmd/gofmt@latest

# Imports management
go install golang.org/x/tools/cmd/goimports@latest

# Linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing framework (optional)
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

## Common Commands

### Running Code
```bash
go run main.go           # Run a program
go build main.go         # Compile to binary
go install               # Build and install to $GOBIN
```

### Testing
```bash
go test                  # Run tests in current package
go test -v               # Verbose output
go test ./...            # Test all packages
go test -cover           # With coverage report
```

### Module Management
```bash
go mod init module-name  # Initialize new module
go mod tidy              # Add missing/remove unused dependencies
go mod download          # Download dependencies
```

## Next Steps

1. **Complete Fundamentals**: Work through all examples in [Learning Fundamentals](../../learning/01-fundamentals/)
2. **Follow the Full Learning Path**: See the [Complete Learning Roadmap](../../LEARNING_ROADMAP.md)
3. **Join the Community**: 
   - [Go Forum](https://forum.golangbridge.org/)
   - [r/golang](https://reddit.com/r/golang)
   - [Gophers Slack](https://gophers.slack.com/)
4. **Build Something**: Create your own project using [Project Templates](../../templates/)
5. **Study Production Examples**: Analyze the [Production API](../../examples/03-advanced/production-api/)
6. **Contribute**: Add examples or fix issues in this repository following the [Contributing Guide](../../CONTRIBUTING.md)

## Troubleshooting

### "go: command not found"
- Verify Go is installed: `which go`
- Check your PATH includes Go binary directory
- Restart your terminal after installation

### Import Issues
- Ensure you're in the correct directory
- Run `go mod tidy` to resolve dependencies
- Check for typos in import paths

### Module Problems
- Initialize module if needed: `go mod init`
- Update dependencies: `go get -u ./...`
- Clear module cache: `go clean -modcache`

## Getting Help

1. **Official Documentation**: [golang.org/doc](https://golang.org/doc/)
2. **This Repository**: Check [Architecture Docs](../architecture/) for design patterns and [Troubleshooting Guide](../troubleshooting/)
3. **Community**: Ask questions in Go forums and Slack channels
4. **Stack Overflow**: Tag questions with `go` or `golang`

Ready to start coding? Head to the [Learning Path](../../learning/README.md)!