# Go Intermediate Concepts

Build on Go fundamentals with object-oriented programming patterns and interface design.

## Prerequisites
- Complete [Fundamentals](../fundamentals/) section
- Comfortable with Go syntax and basic functions
- Understanding of Go's type system

## Learning Objectives
By the end of this section, you will:
- Design and implement custom types using structs
- Understand Go's approach to object-oriented programming
- Use composition over inheritance effectively
- Work with interfaces for flexible design
- Apply embedding patterns

## Modules

### 1. [Composition](./composition/)
Go's approach to object-oriented design through composition:
- Struct definitions and methods
- Embedding and interface implementation
- Polymorphism through interfaces
- Real-world examples with `Human`, `Ninja`, and `SeniorNinja`

**Time**: 2-3 hours  
**Files**: `human.go`, `ninja.go`, `senior_ninja.go`

## Running the Examples

```bash
# Navigate to the composition directory
cd composition/

# Run individual files or create a main.go to test interactions
go run *.go
```

## Key Concepts Covered
- Struct types and method receivers
- Interface design and implementation
- Composition over inheritance
- Method sets and interface satisfaction
- Embedding anonymous fields
- Polymorphic behavior

## Design Principles Demonstrated
- **Composition over Inheritance**: Go doesn't have traditional inheritance; instead use embedding
- **Interface Segregation**: Small, focused interfaces are preferred
- **Implicit Interface Implementation**: Types implement interfaces automatically
- **Zero Values**: Structs have useful zero values by design

## Real-World Applications
These patterns appear frequently in:
- Web service handlers
- Database models
- Plugin architectures
- Middleware chains

## Next Steps
Ready for advanced concepts? Move on to [Advanced](../advanced/) topics including:
- Concurrency patterns
- Design pattern implementations  
- Dependency injection
- Production-ready architectures

## Common Patterns
- **Builder Pattern**: Use method chaining with pointer receivers
- **Strategy Pattern**: Implement through interfaces
- **Decorator Pattern**: Use embedding to extend functionality
- **Factory Pattern**: Return interfaces, not concrete types