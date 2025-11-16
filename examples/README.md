# Go Examples

Practical, real-world examples organized by skill level and complexity.

## Skill-Based Organization

### 🚀 [Beginner Examples](./01-beginner/)
Perfect for those starting with Go:
- **Calculator**: Basic arithmetic with testing evolution (v1 → v2)
- **Hello Concurrent**: Introduction to goroutines and concurrency
- **Prerequisites**: Basic programming knowledge, [Getting Started Guide](../docs/tutorials/getting-started.md)
- **Theory Background**: [Learning Fundamentals](../learning/01-fundamentals/)
- **Time**: 1-2 weeks

### 🎯 [Intermediate Examples](./02-intermediate/)  
Building on Go fundamentals:
- **HTTP Services**: REST API clients and servers
- **Prerequisites**: Complete [Learning Fundamentals](../learning/01-fundamentals/)
- **Theory Background**: [Learning Intermediate](../learning/02-intermediate/)
- **Time**: 2-3 weeks

### 🏆 [Advanced Examples](./03-advanced/)
Production-ready implementations:
- **Production API**: Complete hexagonal architecture with testing, CI/CD, and deployment
- **Performance Benchmarks**: Optimization and profiling examples
- **Prerequisites**: Complete [Learning Intermediate](../learning/02-intermediate/)
- **Theory Background**: [Learning Advanced](../learning/03-advanced/)
- **Time**: 4-6 weeks

## How to Use

### Follow Your Skill Level
1. **Assess Prerequisites**: Each level builds on previous knowledge
2. **Start Appropriate**: Don't skip levels - each introduces important concepts
3. **Practice Actively**: Modify examples to test your understanding

### Running Examples
```bash
# Beginner: Calculator v1
cd 01-beginner/calculator/v1/
go run calculator.go
go test -v

# Intermediate: HTTP Services
cd 02-intermediate/http-services/jsonplaceholder/
go run main.go

# Advanced: Production API
cd 03-advanced/production-api/
make run
make test
```

## Learning Path Integration

Examples are designed to complement the [structured learning path](../learning/):

```
📚 Learning Theory          🛠️ Practical Examples
├── Fundamentals       →     ├── Beginner
├── Intermediate       →     ├── Intermediate  
└── Advanced           →     └── Advanced
```

### Recommended Flow
1. **Study Theory**: Read concepts in `/learning/`
2. **Apply Knowledge**: Try corresponding examples
3. **Build Understanding**: Modify and extend examples
4. **Test Mastery**: Move to next level

## Example Progression

### Mathematical Operations (Beginner)
- **v1**: Basic functions + standard testing ([calculator/v1/](01-beginner/calculator/v1/))
- **v2**: Same logic + BDD testing framework ([calculator/v2/](01-beginner/calculator/v2/))
- **Learning**: Testing evolution and framework adoption
- **Related Theory**: [Mathematics Learning](../learning/01-fundamentals/math/)

### Web Services (Intermediate)
- **Client**: REST API consumption patterns ([http-services/jsonplaceholder/](02-intermediate/http-services/jsonplaceholder/))
- **Server**: HTTP handler implementation ([http-services/server/](02-intermediate/http-services/server/))
- **Learning**: Service architecture and HTTP patterns
- **Related Theory**: [Composition Patterns](../learning/02-intermediate/composition/)

### Production Systems (Advanced)
- **Architecture**: Hexagonal/Clean architecture ([production-api/](03-advanced/production-api/))
- **Testing**: Comprehensive testing strategies (see [Testing Guide](../docs/tutorials/testing-guide.md))
- **Deployment**: Production-ready patterns ([deployment/](../deployment/))
- **Learning**: Enterprise-grade Go development
- **Related Theory**: [Advanced Patterns](../learning/03-advanced/patterns/)

## Contributing Examples

When adding new examples:

### Structure Requirements
```
new-example/
├── README.md           # Comprehensive guide
├── main.go             # Entry point
├── *_test.go           # Test files
└── docs/               # Additional documentation
```

### Documentation Standards
1. **Prerequisites**: What knowledge is required
2. **Learning Objectives**: What concepts are demonstrated
3. **Running Instructions**: Clear setup and execution steps
4. **Key Concepts**: What patterns/principles are shown
5. **Next Steps**: Where to go after mastering this example

### Quality Standards
- **Comprehensive Testing**: Unit, integration, and example tests
- **Go Best Practices**: Follow established conventions
- **Real-World Relevance**: Solve actual problems
- **Progressive Complexity**: Build on previous examples

## Quick Start Guide

### New to Go?
1. Start with [Getting Started Tutorial](../docs/tutorials/getting-started.md)
2. Work through [Learning Fundamentals](../learning/01-fundamentals/)
3. Try [Beginner Examples](./01-beginner/)
4. Follow the [Complete Learning Roadmap](../docs/LEARNING_ROADMAP.md)

### Have Go Experience?
1. Assess your level using the prerequisites in each directory
2. Jump to appropriate examples directory ([01-beginner/](01-beginner/), [02-intermediate/](02-intermediate/), or [03-advanced/](03-advanced/))
3. Focus on patterns and architectures you haven't seen
4. Use [Project Templates](../templates/) to apply patterns to your own projects

### Want Production Patterns?
1. Go directly to [Advanced Examples](./03-advanced/)
2. Study the [Production API Implementation](./03-advanced/production-api/)
3. Review [Architecture Documentation](../docs/architecture/)
4. Adapt patterns to your own projects using [deployment guides](../deployment/)

Remember: The goal is not just to run the code, but to understand the principles and adapt them to your own projects!