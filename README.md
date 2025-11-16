# Go Laboratory 🚀

A comprehensive learning and experimentation repository containing 108+ Go source files that demonstrate modern Go development practices, design patterns, and architectural approaches. The project serves as both a reference implementation and hands-on learning resource for Go developers at various skill levels.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 What You'll Learn

- **Go Fundamentals**: Complete coverage of Go fundamentals including data types, control structures, functions, constants, and variables
- **Advanced Architectural Patterns**: Full implementations of Gang of Four patterns (Creational: Builder, Prototype, Singleton; Structural: Adapter) and SOLID principles
- **Web Development**: HTTP services, REST APIs, middleware, and JSONPlaceholder integration for REST API experimentation
- **Testing Strategies**: Comprehensive testing with Ginkgo, Gomega, GoConvey, mock generation using uber-go/mock, and integration testing
- **Production Skills**: Hexagonal architecture, clean code, dependency injection using uber-go/dig, Kubernetes-ready health probes, and containerization support

## 🗂️ Repository Structure

```
go-lab/
├── 📚 learning/                    # Structured learning path
│   ├── 01-fundamentals/           # Basic Go concepts (1-2 weeks)
│   ├── 02-intermediate/           # OOP and composition (2-3 weeks)
│   └── 03-advanced/               # Concurrency and patterns (4-6 weeks)
├── 💡 examples/                   # Practical, skill-based examples
│   ├── 01-beginner/               # Entry-level projects
│   ├── 02-intermediate/           # Real-world applications  
│   └── 03-advanced/               # Production-ready systems
├── 🧪 testing/                    # Comprehensive testing strategies
├── 📖 docs/                       # Tutorials and architecture guides
│   ├── tutorials/                 # Step-by-step guides
│   ├── architecture/              # Design patterns and decisions
│   ├── api-reference/             # API documentation standards
│   └── troubleshooting/           # Debugging and problem-solving
├── 🛠️ templates/                 # Project scaffolding
├── 🚀 deployment/                 # Infrastructure and deployment patterns
├── 📊 scripts/                    # Automation and assessment tools
├── 🤖 .github/workflows/          # CI/CD and quality automation
└── 📋 Root documentation          # GitHub standard files (README, CONTRIBUTING, etc.)
```

## 🚀 Quick Start

### New to Go?
1. **Install Go**: [Download from golang.org](https://golang.org/dl/)
2. **Get Started**: Read our [Getting Started Guide](docs/tutorials/getting-started.md)
3. **Follow the Path**: Begin with [Learning Fundamentals](learning/01-fundamentals/)
4. **Assess Progress**: Use `./scripts/assess-skills.sh` to track learning

### Have Go Experience?
1. **Assess Your Level**: Run `./scripts/assess-skills.sh` for personalized recommendations
2. **Jump to Examples**: Try [skill-based examples](examples/) based on your level
3. **Explore Patterns**: Study [advanced concepts](learning/03-advanced/)
4. **Build Projects**: Use `./scripts/create-project.sh` to scaffold new applications

### Want Production Patterns?
1. **Study the API**: Explore [Production API](examples/03-advanced/production-api/)
2. **Use Templates**: Scaffold projects with [templates](templates/)
3. **Deploy Applications**: Follow [deployment guides](deployment/)
4. **Monitor Performance**: Use [benchmarking examples](examples/03-advanced/performance-benchmarks/)

## 🎓 Learning Path

### 📚 [Learning Track](learning/) - Theory First
**Structured, comprehensive Go education**

```
01-Fundamentals (1-2 weeks)  02-Intermediate (2-3 weeks)  03-Advanced (4-6 weeks)
├── Hello World & Basics    ├── Interfaces & Composition  ├── Concurrency Patterns
├── Functions & Types       ├── Error Handling           ├── Design Patterns  
├── Testing Introduction    ├── Package Organization     ├── Clean Architecture
└── Development Setup       └── Object-Oriented Patterns └── Production Practices
```

### 💡 [Examples Track](examples/) - Practice First
**Practical, project-based learning**

```
01-Beginner                 02-Intermediate               03-Advanced
├── Calculator (v1 & v2)    ├── HTTP Services            ├── Production API
├── Hello Concurrent        ├── REST API Client          ├── Hexagonal Architecture
└── Basic Testing           ├── JSON Processing          ├── Performance Benchmarks
                           └── Service Architecture      └── CI/CD & Deployment
```

### 🎯 Navigation Tools
- **[INDEX.md](docs/INDEX.md)**: Complete reference to all concepts and locations
- **[LEARNING_ROADMAP.md](docs/LEARNING_ROADMAP.md)**: Detailed career progression guide
- **Assessment Tool**: `./scripts/assess-skills.sh` for personalized recommendations

## 🛠️ How to Use This Repository

### 🎓 For Learning (Self-Study)
```bash
# 1. Assess your current level
./scripts/assess-skills.sh

# 2. Follow personalized recommendations
# Start with learning/01-fundamentals/ if new to Go
cd learning/01-fundamentals/hello
go run hello.go

# 3. Practice with examples
cd ../../examples/01-beginner/calculator/v1
go test -v

# 4. Track progress regularly
./scripts/assess-skills.sh  # Re-run periodically
```

### 👥 For Team Training
```bash
# 1. Set up team development environment
make setup-dev  # Install recommended tools

# 2. Use as curriculum
# - learning/ for theoretical foundations
# - examples/ for hands-on practice
# - testing/ for testing strategies

# 3. Create team projects from templates
./scripts/create-project.sh -t http-service -n team-api -m github.com/company/team-api

# 4. Track team progress
# Use scripts/assess-skills.sh for individual assessments
```

### 🏢 For Corporate Training
- **Curriculum**: Use [LEARNING_ROADMAP.md](docs/LEARNING_ROADMAP.md) for structured courses
- **Assessment**: Built-in skill assessment and progress tracking
- **Projects**: Templates for real-world team projects from [templates/](templates/)
- **Quality Standards**: Automated code quality and CI/CD examples in [examples/03-advanced/](examples/03-advanced/)

### 🛠️ For Development Projects
```bash
# 1. Generate new project from templates
./scripts/create-project.sh -t http-service -n my-api -m github.com/mycompany/my-api

# 2. Study production patterns
cd examples/03-advanced/production-api/
# Review architecture, testing, and deployment

# 3. Apply patterns to your project
# Use deployment/ for infrastructure patterns
# Reference docs/architecture/ for design decisions
```

## ⚙️ Key Features & Capabilities

### 🤖 Automated Learning Tools
- **Skill Assessment**: Interactive progress tracking with `./scripts/assess-skills.sh`
- **Project Generation**: Automated scaffolding with `./scripts/create-project.sh`
- **Quality Gates**: Comprehensive CI/CD with educational focus
- **Progress Validation**: Automated compilation and test verification

### 📚 Comprehensive Educational Content
- **Structured Progression**: Numbered paths (01-fundamentals → 03-advanced)
- **Multiple Learning Styles**: Theory-first vs. practice-first approaches
- **Real-World Examples**: Production-ready patterns and architectures
- **Testing Mastery**: Unit, BDD, mocking, integration, and E2E testing

### 🚀 Production-Ready Patterns
- **Hexagonal Architecture**: Clean separation using ports and adapters pattern with uber-go/dig dependency injection
- **Modern Tooling**: Integration with Cobra CLI, Viper configuration, and comprehensive quality assurance
- **Design Patterns**: Gang of Four patterns adapted for Go idioms with working examples
- **Specialized Modules**: Concurrency primitives, pipelines, dependency injection with concrete/abstract/safety patterns
- **Performance Optimization**: Benchmarking, profiling, and optimization with cyclomatic complexity limits (≤7)
- **Kubernetes Ready**: Built-in health probes and containerization support with multi-stage Docker builds

### 🔧 Development Automation & Quality Engineering
- **Code Quality**: Educational-focused linting with GolangCI-Lint and automated formatting
- **Testing Frameworks**: Examples using Ginkgo, Gomega, GoConvey, and automated mock creation
- **CI/CD Pipeline**: Comprehensive quality checks with coverage reporting and security scanning
- **Documentation**: Auto-generated API docs and comprehensive architectural guides
- **Community Features**: Contribution guidelines and recognition system

## 📊 Project Statistics & Technical Stack

### Project Scale
- **108+ Go source files** demonstrating modern practices
- **Complete test coverage** with multiple testing frameworks
- **Production-ready examples** with deployment configurations
- **Comprehensive documentation** with tutorials and guides

### Technical Stack
- **Language**: Go 1.19+
- **Testing**: Ginkgo/Gomega, GoConvey, uber-go/mock
- **Web Framework**: Beego v2, standard HTTP ServeMux
- **Build System**: Go-Task for automated build processes
- **Containerization**: Docker support with multi-stage builds
- **Code Quality**: GolangCI-Lint, automated formatting

## 🎯 Use Cases & Strategic Benefits

### For Learning
- **Bootcamp Curriculum**: Structured learning from basics to advanced with step-by-step progression
- **Self-Study**: Progressive skill development with clear milestones and automated assessment
- **Team Training**: Onboarding new Go developers effectively with comprehensive quality benchmarks

### For Development
- **Reference Implementation**: Production-ready patterns following industry best practices
- **Development Accelerator**: Reusable components and patterns for rapid Go application development
- **Project Scaffolding**: Templates for quick project initialization with modern tooling
- **Best Practices**: Code quality, testing strategies, and architecture examples

### For Teaching & Enterprise
- **Course Material**: Ready-to-use examples and exercises with educational resource progression
- **Workshop Content**: Hands-on learning materials with practical implementation guidance
- **Assessment**: Progressive skill evaluation opportunities and CI/CD integration
- **Enterprise Applications**: Demonstrates modern Go development workflows and tooling integration

### Strategic Benefits
- Reduces Go development learning curve through practical examples
- Provides tested, production-ready code patterns for enterprise applications  
- Serves as a foundation for scaling Go development teams and practices
- Demonstrates testing strategies, code quality enforcement, and modern development workflows

## 🔧 Development Tools

### One-Command Setup
```bash
# Clone and setup (recommended)
git clone https://github.com/ocrosby/go-lab.git
cd go-lab
make setup  # Installs all recommended tools

# Or manual setup
go install golang.org/x/tools/cmd/gofmt@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install go.uber.org/mock/mockgen@latest
```

### Available Commands
```bash
# Learning and assessment
./scripts/assess-skills.sh          # Evaluate your Go skills
./scripts/create-project.sh --help  # Generate new projects

# Development workflow
make test                           # Run all tests
make lint                          # Code quality checks
make docs                          # Generate documentation
make validate                      # Validate learning paths

# Quality assurance
make test-coverage                 # Test with coverage report
make benchmark                     # Run performance benchmarks
make security-scan                 # Security vulnerability scan
```

### IDE Configuration
- **VS Code**: Install Go extension + use .vscode/settings.json (provided)
- **GoLand**: JetBrains Go IDE with full toolchain integration
- **Vim/Neovim**: Use provided LSP configuration examples

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

## 🔄 How to Extend & Customize

### 📚 Adding New Learning Content
```bash
# 1. Choose appropriate level
learning/01-fundamentals/    # Basic concepts
learning/02-intermediate/    # Applied concepts  
learning/03-advanced/       # Production patterns

# 2. Follow naming convention
your-topic/
├── README.md              # Prerequisites, objectives, concepts
├── example.go             # Working code example
├── example_test.go        # Comprehensive tests
└── docs/                  # Additional documentation

# 3. Update cross-references
# - Add to docs/INDEX.md
# - Reference in docs/LEARNING_ROADMAP.md
# - Link from related examples
```

### 💡 Adding New Examples
```bash
# 1. Determine skill level
examples/01-beginner/       # Entry-level projects
examples/02-intermediate/   # Real-world applications
examples/03-advanced/      # Production systems

# 2. Create complete example
your-example/
├── README.md              # Usage instructions and concepts
├── main.go                # Entry point
├── *_test.go             # Test coverage
├── Makefile              # Build automation
└── docs/                 # Architecture documentation

# 3. Ensure quality
make lint                  # Code quality
make test                  # Test coverage
./scripts/assess-skills.sh # Validate integration
```

### 🛠️ Creating New Templates
```bash
# 1. Design template structure
templates/your-template/
├── README.md              # Customization guide
├── main.go                # Placeholder code
├── go.mod                 # Module definition
├── Makefile              # Build targets
├── Dockerfile            # Containerization
└── .template-config       # Template metadata

# 2. Use placeholders
# TEMPLATE_PROJECT_NAME    → replaced with actual project name
# TEMPLATE_MODULE_PATH     → replaced with module path
# template-project         → replaced with kebab-case name

# 3. Test template generation
./scripts/create-project.sh -t your-template -n test-project -m github.com/test/test-project
```

### 📖 Enhancing Documentation
```bash
# Documentation structure
docs/
├── tutorials/            # Step-by-step guides
├── architecture/         # Design decisions  
├── api-reference/       # API documentation
└── troubleshooting/     # Problem solving

# Follow documentation standards
# - Clear prerequisites and objectives
# - Working code examples
# - Cross-references to related content
# - Troubleshooting common issues
```

## 🤝 Contributing Guidelines

### Quick Contribution
1. **Fork & Clone**: Standard GitHub workflow
2. **Install Tools**: Run `make setup` for development environment
3. **Quality Checks**: Use `make lint && make test` before committing
4. **Documentation**: Update relevant README files and cross-references

### Contribution Types
- 🐛 **Bug Fixes**: Code corrections and improvements
- 📚 **Learning Content**: New examples and educational materials
- 📖 **Documentation**: Guides, tutorials, and API documentation
- 🛠️ **Templates**: Project scaffolding for common use cases
- 🚀 **Infrastructure**: CI/CD, deployment, and automation improvements

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines, coding standards, and review process.

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

## 🎯 What Makes This Repository Special

### 🏆 Educational Excellence
- **Progressive Learning**: Clear skill-level progression with validation
- **Multiple Learning Styles**: Theory-first and practice-first approaches
- **Real-World Relevance**: Production patterns used in actual Go projects
- **Assessment Tools**: Built-in progress tracking and skill evaluation

### 🤖 Automation & Quality
- **Self-Validating**: Automated testing of all learning paths and examples
- **Quality Enforcement**: Educational-focused linting and security scanning  
- **Continuous Integration**: Comprehensive CI/CD with educational quality gates
- **Community-Ready**: Professional contribution workflows and recognition

### 🌟 Production-Ready
- **Enterprise Patterns**: Hexagonal architecture, dependency injection, clean code
- **Deployment Infrastructure**: Docker, Kubernetes, cloud deployment examples
- **Performance Engineering**: Benchmarking, profiling, and optimization guides
- **Monitoring & Observability**: Structured logging, metrics, and debugging

---

## 🚀 Get Started Now

### For Absolute Beginners
```bash
git clone https://github.com/ocrosby/go-lab.git
cd go-lab
./scripts/assess-skills.sh  # Get personalized recommendations
```

### For Experienced Developers
```bash
# Jump straight to advanced patterns
cd examples/03-advanced/production-api/
make run

# Or generate a new project
./scripts/create-project.sh -t http-service -n my-api -m github.com/me/my-api
```

### For Teams & Organizations
- Use [LEARNING_ROADMAP.md](docs/LEARNING_ROADMAP.md) for structured training curriculum
- Implement [quality standards](.golangci.yml) in your projects
- Deploy using [production patterns](deployment/) and infrastructure examples

---

**Ready to master Go development?** 

🎯 **Start Here**: [Getting Started Guide](docs/tutorials/getting-started.md) → [Skill Assessment](scripts/assess-skills.sh) → [Learning Path](learning/)

📖 **Reference**: [Complete Index](docs/INDEX.md) | [Learning Roadmap](docs/LEARNING_ROADMAP.md) | [Troubleshooting](docs/troubleshooting/)

🛠️ **Build**: [Project Templates](templates/) | [Production Examples](examples/03-advanced/) | [Deployment Guides](deployment/)

*The Go Laboratory project represents a significant investment in Go development excellence, providing both educational value and practical implementation guidance for building robust, maintainable Go applications. This repository transforms Go learning from scattered tutorials into a comprehensive, automated educational ecosystem.*

[![Star this repo](https://img.shields.io/github/stars/ocrosby/go-lab?style=social)](https://github.com/ocrosby/go-lab) [![Fork this repo](https://img.shields.io/github/forks/ocrosby/go-lab?style=social)](https://github.com/ocrosby/go-lab/fork)
