# Documentation

Comprehensive documentation for the Go Laboratory project.

## 📚 Documentation Structure

### [Tutorials](./tutorials/)
Step-by-step learning guides for specific topics:
- Getting started with Go development
- Setting up your development environment  
- Building your first Go application
- Advanced patterns and best practices

### [Architecture](./architecture/)
Design decisions and architectural patterns:
- Project structure rationale
- Design pattern implementations
- Clean architecture principles
- Production system design

### [API Reference](./api-reference/)
Detailed API documentation:
- Package documentation
- Function signatures and usage
- Code examples and best practices
- Generated documentation

## 🎯 Quick Start Guides

### For Beginners
1. Read [Getting Started Tutorial](./tutorials/getting-started.md)
2. Follow the [Learning Path](../learning/README.md)
3. Try [Beginner Examples](../examples/beginner/)

### For Intermediate Developers
1. Review [Architecture Decisions](./architecture/)
2. Explore [Intermediate Examples](../examples/intermediate/)
3. Study [Testing Strategies](../testing/)

### For Advanced Developers
1. Analyze [Production API](../projects/api/)
2. Study [Advanced Patterns](../learning/advanced/)
3. Review [Architecture Documentation](./architecture/)

## 📖 External Resources

### Official Go Documentation
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Community Resources
- [Go Blog](https://blog.golang.org/)
- [Go Wiki](https://github.com/golang/go/wiki)
- [Awesome Go](https://github.com/avelino/awesome-go)

## 🤝 Contributing to Documentation

When adding new documentation:
1. Follow the existing structure and formatting
2. Include practical examples and code snippets
3. Cross-reference related topics
4. Keep explanations clear and concise
5. Update the index when adding new sections

## 📋 Documentation Standards

### Structure
- Use clear headings and subheadings
- Include table of contents for long documents
- Add prerequisites and learning objectives
- Provide next steps and related resources

### Code Examples
- Include working, tested code
- Add comments explaining key concepts
- Show both positive and error cases
- Include expected output where helpful

### Cross-References
- Link to related documentation
- Reference relevant examples and tests
- Point to external resources when appropriate
- Maintain bidirectional navigation

---

## 🔧 Legacy Content: Environment Setup

### Setting GOROOT to Define the Go Binary Location

The operating system needs to know how to find the Go installation. In most instances, if you've installed Go in the 
default path, such as /usr/local/go on a Unix system, you don't have to take any action. However, in the event that
you've chosen to install Go in a nonstandard path or are installing Go on Windows, you'll need to tell the operating
system where to find the Go binary.

You can do this from your command line by setting the reserved `GOROOT` environment variable to the location of your
binary.

### Setting GOPATH to Define the Go Workspace

Unlike setting your `GOROOT` environment variable, which is optional, you must set your `GOPATH` environment
to instruct the Go toolset where your source code, third party libraries, and compiled programs will exist.

The `GOPATH` environment variable is a list of paths that point to your Go workspace. The Go workspace is a 
directory that contains your Go source code, third party libraries, and compiled Go programs. The `GOPATH` 
contains three subdirectories within: bin, pkg, and src.

The bin directory will contain your compiled and installed Go executable binaries. Binaries that are built and 
installed will be automatically placed in this location.

The pkg directory stores various package objects including third-party Go dependencies that your code might rely on.

The src directory will contain all the source code you'll write.

### SOLID Principles in Go

The SOLID principles are essentially a set of rules for helping you write clean
and maintainable object-oriented code:

* **Single responsibility**: Each type should have one reason to change
* **Open/closed**: Open for extension, closed for modification
* **Liskov substitution**: Interfaces should be substitutable
* **Interface segregation**: Many specific interfaces are better than one general-purpose interface
* **Dependency inversion**: Depend on abstractions, not concretions

Go is an object-based programming language that supports these principles through interfaces and composition.

For detailed examples of SOLID principles in Go, see [Architecture Documentation](./architecture/).

