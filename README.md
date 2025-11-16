# Go Laboratory 🚀

A comprehensive learning repository for mastering Go development, from fundamentals to production-ready applications.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 What You'll Learn

- **Go Fundamentals**: Syntax, data types, functions, and control structures
- **Advanced Concepts**: Concurrency, interfaces, and design patterns  
- **Web Development**: HTTP services, REST APIs, and middleware
- **Testing Strategies**: Unit testing, BDD, mocking, and integration testing
- **Production Skills**: Clean architecture, CI/CD, deployment, and monitoring

## 🗂️ Repository Structure

```
go-lab/
├── 📚 learning/              # Structured learning path
│   ├── fundamentals/         # Basic Go concepts (1-2 weeks)
│   ├── intermediate/         # OOP and composition (2-3 weeks)
│   └── advanced/            # Concurrency and patterns (4-6 weeks)
├── 💡 examples/             # Practical, skill-based examples
│   ├── beginner/            # Entry-level projects
│   ├── intermediate/        # Real-world applications  
│   └── advanced/           # Production-ready systems
├── 🧪 testing/              # Comprehensive testing strategies
├── 📖 docs/                 # Tutorials and architecture guides
├── 🛠️ templates/           # Project scaffolding
└── 📜 scripts/             # Utility scripts
```

## 🚀 Quick Start

### New to Go?
1. **Install Go**: [Download from golang.org](https://golang.org/dl/)
2. **Get Started**: Read our [Getting Started Guide](docs/tutorials/getting-started.md)
3. **Follow the Path**: Begin with [Learning Fundamentals](learning/fundamentals/)

### Have Go Experience?
1. **Assess Your Level**: Check prerequisites in each directory
2. **Jump to Examples**: Try [skill-based examples](examples/)
3. **Explore Patterns**: Study [advanced concepts](learning/advanced/)

### Want Production Patterns?
1. **Study the API**: Explore [Production API](examples/advanced/production-api/)
2. **Use Templates**: Scaffold projects with [templates](templates/)
3. **Apply Patterns**: Implement in your own projects

## 🎓 Learning Path

### 📚 [Learning Track](learning/)
**Structured, comprehensive Go education**

```
Fundamentals (1-2 weeks)     Intermediate (2-3 weeks)     Advanced (4-6 weeks)
├── Basic Syntax            ├── Interfaces & Composition   ├── Concurrency Patterns
├── Functions & Types       ├── Error Handling            ├── Design Patterns  
├── Testing Basics          ├── Package Organization      ├── Clean Architecture
└── Development Setup       └── Object-Oriented Patterns  └── Production Practices
```

### 💡 [Examples Track](examples/)
**Practical, project-based learning**

```
Beginner                    Intermediate                  Advanced
├── Calculator (v1 & v2)    ├── HTTP Services            ├── Production API
├── Hello Concurrent        ├── REST API Client          ├── Hexagonal Architecture
└── Basic Testing           ├── JSON Processing          ├── Comprehensive Testing
                           └── Service Architecture      └── CI/CD & Deployment
```

## 🛠️ Key Features

### Comprehensive Testing Examples
- **Unit Testing**: Standard library and table-driven tests
- **BDD Testing**: Ginkgo/Gomega framework examples
- **Mocking**: uber-go/mock integration and patterns
- **Integration Testing**: HTTP services and database testing

### Production-Ready Patterns
- **Clean Architecture**: Hexagonal architecture implementation
- **Design Patterns**: Gang of Four patterns adapted for Go
- **Dependency Injection**: Interface-based design and DI containers
- **Observability**: Structured logging, metrics, and monitoring

### Modern Development Practices
- **CI/CD**: GitHub Actions workflows and deployment pipelines
- **Containerization**: Docker and Docker Compose configurations
- **Code Quality**: GolangCI-Lint, security scanning, and best practices
- **Documentation**: Comprehensive guides and API documentation

## 📊 Project Statistics

- **108+ Go source files** demonstrating modern practices
- **Complete test coverage** with multiple testing frameworks
- **Production-ready examples** with deployment configurations
- **Comprehensive documentation** with tutorials and guides

## 🎯 Use Cases

### For Learning
- **Bootcamp Curriculum**: Structured learning from basics to advanced
- **Self-Study**: Progressive skill development with clear milestones
- **Team Training**: Onboarding new Go developers effectively

### For Development
- **Reference Implementation**: Production-ready patterns and practices
- **Project Scaffolding**: Templates for quick project initialization
- **Best Practices**: Code quality, testing, and architecture examples

### For Teaching
- **Course Material**: Ready-to-use examples and exercises
- **Workshop Content**: Hands-on learning materials
- **Assessment**: Progressive skill evaluation opportunities

## 🔧 Development Tools

### Recommended Setup
```bash
# Essential tools
go install golang.org/x/tools/cmd/gofmt@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing frameworks
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Mock generation
go install go.uber.org/mock/mockgen@latest
```

### IDE Configuration
- **VS Code**: Install Go extension with enhanced features
- **GoLand**: JetBrains Go IDE with full toolchain integration
- **Vim/Neovim**: LSP configuration for Go development

## 🎨 Architecture Highlights

### Hexagonal Architecture
Demonstrated in the production API with clear separation:
- **Domain**: Core business logic and entities
- **Application**: Use cases and orchestration
- **Infrastructure**: External integrations and adapters

### Testing Strategy
Multi-layered testing approach:
- **Unit Tests**: Fast, isolated component testing
- **Integration Tests**: Component interaction verification
- **E2E Tests**: Complete user journey validation

### Quality Engineering
- **Cyclomatic Complexity**: Limited to ≤7 for maintainability
- **Code Coverage**: Comprehensive test coverage metrics
- **Security Scanning**: Automated vulnerability detection
- **Performance**: Benchmarking and optimization examples

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. **Issues**: Report bugs or suggest features
2. **Examples**: Add new learning examples or improve existing ones
3. **Documentation**: Enhance guides, tutorials, and API docs
4. **Templates**: Create project scaffolding for common use cases

## 📈 Next Steps

After exploring this repository:

1. **Build Your Own Projects**: Apply learned patterns to real applications
2. **Contribute to Open Source**: Use skills in Go community projects
3. **Share Knowledge**: Teach others using these materials
4. **Stay Current**: Follow Go evolution and update practices

## 📚 External Resources

### Official Go Resources
- [Go Documentation](https://golang.org/doc/) - Official language documentation
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices guide
- [Go Blog](https://blog.golang.org/) - Latest developments and insights

### Community Resources  
- [Go Forum](https://forum.golangbridge.org/) - Community discussions
- [Gophers Slack](https://gophers.slack.com/) - Real-time community chat
- [Awesome Go](https://github.com/avelino/awesome-go) - Curated Go resources

---

**Ready to start your Go journey?** Begin with the [Getting Started Guide](docs/tutorials/getting-started.md) or jump directly to [Learning Fundamentals](learning/fundamentals/).

*This repository represents a comprehensive investment in Go development excellence, providing both educational value and practical implementation guidance for building robust, maintainable Go applications.*
