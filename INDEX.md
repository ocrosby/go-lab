# Go Laboratory - Comprehensive Index 📚

A complete reference to all concepts, patterns, and examples in this repository.

## 🗂️ Quick Navigation

| Category | Location | Description |
|----------|----------|-------------|
| [**Learning Path**](#learning-path) | `/learning/` | Structured educational content |
| [**Practical Examples**](#practical-examples) | `/examples/` | Project-based learning |
| [**Testing Strategies**](#testing-strategies) | `/testing/` | Comprehensive testing approaches |
| [**Documentation**](#documentation) | `/docs/` | Guides, tutorials, and references |
| [**Project Templates**](#project-templates) | `/templates/` | Scaffolding for new projects |

## 📚 Learning Path

### 01-Fundamentals (1-2 weeks)
**Location**: [`learning/01-fundamentals/`](learning/01-fundamentals/)
**Practical Examples**: [`examples/01-beginner/`](examples/01-beginner/)

#### Core Concepts
- **Hello World** [`hello/`](learning/01-fundamentals/hello/)
  - Program structure and `package main`
  - Import statements and standard library
  - The `main()` function as entry point

- **Mathematics** [`math/`](learning/01-fundamentals/math/)
  - Function definitions and parameters
  - Return values and error handling
  - Basic unit testing patterns
  - Table-driven tests

#### Skills Developed
- ✅ Go syntax and basic data types
- ✅ Function definition and calling
- ✅ Error handling patterns
- ✅ Unit testing with standard library
- ✅ Package organization basics

### 02-Intermediate (2-3 weeks)  
**Location**: [`learning/02-intermediate/`](learning/02-intermediate/)
**Practical Examples**: [`examples/02-intermediate/`](examples/02-intermediate/)

#### Core Concepts
- **Composition** [`composition/`](learning/02-intermediate/composition/)
  - Struct types and method receivers
  - Interface design and implementation  
  - Embedding and composition patterns
  - Polymorphism through interfaces

#### Skills Developed
- ✅ Object-oriented design in Go
- ✅ Interface-based programming
- ✅ Composition over inheritance
- ✅ Method sets and interface satisfaction

### 03-Advanced (4-6 weeks)
**Location**: [`learning/03-advanced/`](learning/03-advanced/)
**Practical Examples**: [`examples/03-advanced/`](examples/03-advanced/)

#### Core Concepts
- **Concurrency** [`concurrency/`](learning/03-advanced/concurrency/)
  - Goroutines and channels
  - Select statements and timeouts
  - Pipeline patterns
  - Synchronization primitives

- **Dependency Injection** [`dependency-injection/`](learning/03-advanced/dependency-injection/)
  - Constructor injection patterns
  - Interface-based design
  - Service locator pattern
  - Clean architecture principles

- **Design Patterns** [`patterns/`](learning/03-advanced/patterns/)
  - **Creational Patterns**
    - Builder: [`patterns/creational/builder/`](learning/03-advanced/patterns/creational/builder/)
    - Prototype: [`patterns/creational/prototype/`](learning/03-advanced/patterns/creational/prototype/)
    - Singleton: [`patterns/creational/singleton/`](learning/03-advanced/patterns/creational/singleton/)
  - **Structural Patterns**
    - Adapter: [`patterns/structural/adapter/`](learning/03-advanced/patterns/structural/adapter/)

#### Skills Developed
- ✅ Concurrent programming mastery
- ✅ Design pattern implementation
- ✅ Clean architecture design
- ✅ Production-ready code patterns

## 💡 Practical Examples

### 01-Beginner Examples
**Location**: [`examples/01-beginner/`](examples/01-beginner/)

#### Projects
- **Calculator** [`calculator/`](examples/01-beginner/calculator/)
  - **v1/**: Basic arithmetic with standard testing
  - **v2/**: Enhanced with Ginkgo/Gomega BDD testing
  - **Evolution**: Demonstrates testing framework progression

- **Hello Concurrent** [`hello-concurrent/`](examples/01-beginner/hello-concurrent/)
  - Basic goroutine usage
  - Concurrent execution patterns
  - Program synchronization

#### Learning Objectives
- Apply Go fundamentals in practical projects
- Experience testing evolution (standard → BDD)
- Understand basic concurrency concepts

### 02-Intermediate Examples
**Location**: [`examples/02-intermediate/`](examples/02-intermediate/)

#### Projects  
- **HTTP Services** [`http-services/`](examples/02-intermediate/http-services/)
  - **JSONPlaceholder Client**: REST API consumption patterns
  - **HTTP Server**: Basic server implementation
  - Service architecture and testing strategies

#### Learning Objectives
- Master HTTP programming in Go
- Implement service-oriented architectures
- Design testable, maintainable code

### 03-Advanced Examples
**Location**: [`examples/03-advanced/`](examples/03-advanced/)

#### Projects
- **Production API** [`production-api/`](examples/03-advanced/production-api/)
  - Hexagonal (Clean) architecture
  - Comprehensive testing strategies
  - CI/CD and deployment practices
  - Production monitoring and observability

#### Learning Objectives
- Implement enterprise-grade architectures
- Master production deployment practices
- Apply comprehensive quality engineering

## 🧪 Testing Strategies

### Testing Approaches
**Location**: [`testing/`](testing/)

#### Mocking Strategies [`mocking/`](testing/mocking/)
- **Cars Example**: [`cars/`](testing/mocking/cars/)
  - Interface-based mocking
  - uber-go/mock integration
  - Dependency injection testing

- **Generated Mocks**: [`mocks/`](testing/mocking/mocks/)
  - Automated mock generation
  - Mock verification patterns
  - Integration testing strategies

#### Test Patterns [`test/`](testing/test/)
- **Channels Testing**: [`channels/`](testing/test/channels/)
  - Concurrent code testing
  - Channel communication patterns
  - Timeout and cancellation testing

#### Testing Frameworks
- **Standard Library**: `testing` package patterns
- **BDD Testing**: Ginkgo & Gomega integration
- **Mock Generation**: uber-go/mock usage
- **Integration Testing**: HTTP and database testing

## 📖 Documentation

### Tutorials [`docs/tutorials/`](docs/tutorials/)
- **Getting Started**: [`getting-started.md`](docs/tutorials/getting-started.md)
  - Environment setup and installation
  - First Go program walkthrough
  - Development tool configuration

- **Testing Guide**: [`testing-guide.md`](docs/tutorials/testing-guide.md)
  - Complete testing strategies
  - Framework comparisons
  - Best practices and patterns

### Architecture [`docs/architecture/`](docs/architecture/)
- **Project Structure**: [`project-structure.md`](docs/architecture/project-structure.md)
  - Repository organization rationale
  - Scalability considerations
  - Decision records and trade-offs

- **Design Patterns**: [`design-patterns.md`](docs/architecture/design-patterns.md)
  - Go-specific pattern implementations
  - Classical patterns adapted for Go
  - Modern Go patterns and idioms

### API Reference [`docs/api-reference/`](docs/api-reference/)
- Documentation generation standards
- Package documentation examples
- API design best practices

## 🛠️ Project Templates

### Available Templates
**Location**: [`templates/`](templates/)

#### HTTP Service Template [`http-service/`](templates/http-service/)
- REST API starter template
- Middleware and configuration setup
- Docker and deployment ready
- Testing framework integration

#### Future Templates (Planned)
- **CLI Tool Template**: Command-line application scaffolding
- **Production API Template**: Enterprise-grade API template
- **Library Template**: Reusable package template

## 🔍 Concept Cross-Reference

### Core Go Concepts

#### Language Fundamentals
| Concept | Learning Location | Example Location | Documentation |
|---------|------------------|------------------|---------------|
| Functions | `01-fundamentals/math/` | `01-beginner/calculator/` | `tutorials/getting-started.md` |
| Interfaces | `02-intermediate/composition/` | `02-intermediate/http-services/` | `architecture/design-patterns.md` |
| Goroutines | `03-advanced/concurrency/` | `01-beginner/hello-concurrent/` | `tutorials/testing-guide.md` |
| Error Handling | `01-fundamentals/math/` | All examples | `architecture/project-structure.md` |

#### Testing Strategies
| Strategy | Location | Example Usage | Complexity |
|----------|----------|---------------|------------|
| Unit Testing | `01-fundamentals/math/` | `01-beginner/calculator/v1/` | Beginner |
| Table-Driven Tests | `01-fundamentals/math/` | `01-beginner/calculator/v1/` | Beginner |
| BDD Testing | `testing/` | `01-beginner/calculator/v2/` | Intermediate |
| Mocking | `testing/mocking/` | `02-intermediate/http-services/` | Intermediate |
| Integration Testing | `testing/` | `03-advanced/production-api/` | Advanced |

#### Architecture Patterns
| Pattern | Theory Location | Implementation | Template |
|---------|----------------|----------------|----------|
| Clean Architecture | `docs/architecture/` | `03-advanced/production-api/` | `templates/http-service/` |
| Dependency Injection | `03-advanced/dependency-injection/` | `03-advanced/production-api/` | `templates/http-service/` |
| Repository Pattern | `docs/architecture/` | `03-advanced/production-api/` | - |
| Service Layer | `docs/architecture/` | `02-intermediate/http-services/` | `templates/http-service/` |

### Design Patterns Reference

#### Creational Patterns
| Pattern | Location | Use Case | Difficulty |
|---------|----------|----------|------------|
| Builder | `03-advanced/patterns/creational/builder/` | Complex object construction | Intermediate |
| Singleton | `03-advanced/patterns/creational/singleton/` | Single instance management | Beginner |
| Factory | `docs/architecture/design-patterns.md` | Object creation abstraction | Intermediate |
| Prototype | `03-advanced/patterns/creational/prototype/` | Object cloning | Intermediate |

#### Structural Patterns  
| Pattern | Location | Use Case | Difficulty |
|---------|----------|----------|------------|
| Adapter | `03-advanced/patterns/structural/adapter/` | Interface compatibility | Intermediate |
| Decorator | `docs/architecture/design-patterns.md` | Dynamic behavior addition | Advanced |
| Facade | `02-intermediate/http-services/` | Interface simplification | Beginner |

#### Behavioral Patterns
| Pattern | Location | Use Case | Difficulty |
|---------|----------|----------|------------|
| Strategy | `docs/architecture/design-patterns.md` | Algorithm selection | Intermediate |
| Observer | `docs/architecture/design-patterns.md` | Event notification | Intermediate |
| Command | `03-advanced/production-api/` | Request encapsulation | Advanced |

## 🎯 Learning Paths by Goal

### Goal: Web API Development
**Recommended Path:**
1. `01-fundamentals/` → Basic Go syntax
2. `01-beginner/calculator/` → Testing foundations  
3. `02-intermediate/http-services/` → HTTP programming
4. `templates/http-service/` → Project scaffolding
5. `03-advanced/production-api/` → Production patterns

### Goal: System Programming
**Recommended Path:**
1. `01-fundamentals/` → Language basics
2. `03-advanced/concurrency/` → Concurrent programming
3. `01-beginner/hello-concurrent/` → Practical concurrency
4. `03-advanced/patterns/` → Advanced patterns
5. Custom CLI tools using templates

### Goal: Testing Mastery
**Recommended Path:**
1. `01-fundamentals/math/` → Basic testing
2. `01-beginner/calculator/v1/` → Standard library testing
3. `01-beginner/calculator/v2/` → BDD testing
4. `testing/mocking/` → Mock strategies
5. `docs/tutorials/testing-guide.md` → Complete guide

### Goal: Architecture Design
**Recommended Path:**
1. `02-intermediate/composition/` → Interface design
2. `03-advanced/dependency-injection/` → DI patterns
3. `03-advanced/patterns/` → Design patterns
4. `docs/architecture/` → Architecture theory
5. `03-advanced/production-api/` → Implementation

## 🔗 External References

### Official Resources
- **Go Documentation**: https://golang.org/doc/
- **Go Tour**: https://tour.golang.org/
- **Effective Go**: https://golang.org/doc/effective_go.html
- **Go Blog**: https://blog.golang.org/

### Community Resources
- **Go by Example**: https://gobyexample.com/
- **Awesome Go**: https://github.com/avelino/awesome-go
- **Go Forum**: https://forum.golangbridge.org/
- **Gophers Slack**: https://gophers.slack.com/

---

**Navigation Tip**: Use Ctrl+F (or Cmd+F) to search for specific concepts, patterns, or file locations in this index.

**Last Updated**: This index reflects the current repository structure and will be updated as new content is added.