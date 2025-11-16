# Go Examples

Practical, real-world examples organized by skill level and complexity.

## Skill-Based Organization

### 🚀 [Beginner Examples](./beginner/)
Perfect for those starting with Go:
- **Calculator**: Basic arithmetic with testing evolution (v1 → v2)
- **Hello Concurrent**: Introduction to goroutines and concurrency
- **Prerequisites**: Basic programming knowledge
- **Time**: 1-2 weeks

### 🎯 [Intermediate Examples](./intermediate/)  
Building on Go fundamentals:
- **HTTP Services**: REST API clients and servers
- **Prerequisites**: Complete [Fundamentals](../learning/fundamentals/)
- **Time**: 2-3 weeks

### 🏆 [Advanced Examples](./advanced/)
Production-ready implementations:
- **Production API**: Complete hexagonal architecture with testing, CI/CD, and deployment
- **Prerequisites**: Complete [Intermediate](../learning/intermediate/)
- **Time**: 4-6 weeks

## How to Use

### Follow Your Skill Level
1. **Assess Prerequisites**: Each level builds on previous knowledge
2. **Start Appropriate**: Don't skip levels - each introduces important concepts
3. **Practice Actively**: Modify examples to test your understanding

### Running Examples
```bash
# Beginner: Calculator v1
cd beginner/calculator/v1/
go run calculator.go
go test -v

# Intermediate: HTTP Services
cd intermediate/http-services/jsonplaceholder/
go run main.go

# Advanced: Production API
cd advanced/production-api/
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
- **v1**: Basic functions + standard testing
- **v2**: Same logic + BDD testing framework
- **Learning**: Testing evolution and framework adoption

### Web Services (Intermediate)
- **Client**: REST API consumption patterns
- **Server**: HTTP handler implementation  
- **Learning**: Service architecture and HTTP patterns

### Production Systems (Advanced)
- **Architecture**: Hexagonal/Clean architecture
- **Testing**: Comprehensive testing strategies
- **Deployment**: Production-ready patterns
- **Learning**: Enterprise-grade Go development

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
2. Work through [Learning Fundamentals](../learning/fundamentals/)
3. Try [Beginner Examples](./beginner/)

### Have Go Experience?
1. Assess your level using the prerequisites
2. Jump to appropriate examples directory
3. Focus on patterns and architectures you haven't seen

### Want Production Patterns?
1. Go directly to [Advanced Examples](./advanced/)
2. Study the production API implementation
3. Adapt patterns to your own projects

Remember: The goal is not just to run the code, but to understand the principles and adapt them to your own projects!