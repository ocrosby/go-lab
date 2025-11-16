# Contributing to Go Laboratory 🤝

Thank you for your interest in contributing to the Go Laboratory! This project aims to provide world-class Go education, and we welcome contributions from developers of all skill levels.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Types of Contributions](#types-of-contributions)
- [Development Setup](#development-setup)
- [Contribution Process](#contribution-process)
- [Style Guidelines](#style-guidelines)
- [Quality Standards](#quality-standards)
- [Documentation Standards](#documentation-standards)
- [Testing Requirements](#testing-requirements)
- [Review Process](#review-process)

## 📜 Code of Conduct

This project adheres to a code of conduct that ensures a welcoming environment for all contributors. By participating, you agree to:

- **Be respectful and inclusive** to all community members
- **Focus on constructive feedback** in code reviews and discussions
- **Help newcomers** learn and contribute effectively
- **Maintain high educational standards** in all contributions
- **Give credit** where credit is due

## 🚀 Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Go 1.21 or later** installed
- **Git** configured with your name and email
- **Basic Go knowledge** (at minimum completed our fundamentals)
- **Familiarity with GitHub workflow** (fork, clone, pull request)

### First Time Contributors

1. **Explore the repository** structure using our [comprehensive index](docs/INDEX.md)
2. **Read the learning roadmap** to understand the project's educational goals
3. **Try the examples** to understand the teaching approach
4. **Look for "good first issue" labels** in GitHub issues
5. **Join discussions** in issues before starting major changes

## 🎯 Types of Contributions

We welcome various types of contributions:

### 📚 Learning Content
- **New learning examples** demonstrating Go concepts
- **Additional exercises** for existing modules
- **Alternative explanations** of complex topics
- **Real-world use cases** and applications

### 💡 Code Examples
- **Beginner-friendly examples** with clear educational value
- **Advanced patterns** demonstrating production practices
- **Performance optimizations** with benchmarks and analysis
- **Testing strategies** and comprehensive test suites

### 📖 Documentation
- **Tutorial improvements** with clearer explanations
- **Architecture documentation** for complex patterns
- **Troubleshooting guides** for common issues
- **API documentation** enhancements

### 🛠️ Infrastructure
- **CI/CD improvements** for better automation
- **Template enhancements** for project scaffolding
- **Tool integrations** for development workflow
- **Deployment examples** for various platforms

### 🐛 Bug Fixes
- **Code corrections** in examples or templates
- **Documentation fixes** for accuracy and clarity
- **Link repairs** and reference updates
- **Performance improvements** in existing code

## 💻 Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/go-lab.git
cd go-lab

# Add upstream remote
git remote add upstream https://github.com/ocrosby/go-lab.git
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install golang.org/x/tools/cmd/gofmt@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### 3. Verify Setup

```bash
# Run tests to ensure everything works
make test

# Check code quality
make lint

# Verify examples compile
make check-examples
```

## 🔄 Contribution Process

### 1. Plan Your Contribution

- **Check existing issues** for similar ideas or problems
- **Open an issue** to discuss significant changes before starting
- **Follow the roadmap** and maintain consistency with existing content
- **Consider the educational value** of your contribution

### 2. Create a Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-contribution-name
```

### 3. Make Your Changes

- **Follow our style guidelines** (see below)
- **Write comprehensive tests** for new code
- **Update documentation** as needed
- **Ensure examples compile and run**

### 4. Test Your Changes

```bash
# Run full test suite
make test-all

# Check code quality
make lint

# Validate learning path integrity
make validate-learning-path

# Test documentation links
make check-docs
```

### 5. Commit Your Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "feat: add advanced concurrency patterns example

- Implements worker pool pattern with context cancellation
- Adds comprehensive tests and benchmarks
- Includes detailed documentation and usage examples
- Integrates with existing advanced learning path"
```

### 6. Submit Pull Request

```bash
# Push to your fork
git push origin feature/your-contribution-name

# Open pull request on GitHub
```

## 🎨 Style Guidelines

### Go Code Style

```go
// ✅ Good: Clear, educational code with comments
// CalculateAverage computes the arithmetic mean of a slice of integers.
// This function demonstrates error handling and slice operations.
//
// Example usage:
//     numbers := []int{1, 2, 3, 4, 5}
//     avg, err := CalculateAverage(numbers)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Printf("Average: %.2f\n", avg)
func CalculateAverage(numbers []int) (float64, error) {
    if len(numbers) == 0 {
        return 0, errors.New("cannot calculate average of empty slice")
    }

    var sum int
    for _, num := range numbers {
        sum += num
    }

    return float64(sum) / float64(len(numbers)), nil
}

// ❌ Bad: Unclear, no comments, poor educational value
func calc(n []int) float64 {
    s := 0
    for _, v := range n {
        s += v
    }
    return float64(s) / float64(len(n))
}
```

### Documentation Style

```markdown
# ✅ Good: Clear structure with learning objectives

# Advanced Concurrency Patterns

Learn advanced Go concurrency patterns for production applications.

## Prerequisites
- ✅ Completed [Intermediate Examples](../02-intermediate/)
- ✅ Understanding of goroutines and channels
- ✅ Familiarity with context package

## Learning Objectives
By completing this module, you will:
- Master worker pool patterns
- Implement graceful shutdown mechanisms  
- Apply context for cancellation and timeouts
- Optimize concurrent performance

## Examples

### Worker Pool Pattern
Demonstrates efficient task distribution across multiple goroutines...
```

### Commit Message Format

```
type(scope): brief description

Detailed explanation of what this change does and why.

- Lists specific changes made
- Explains the educational value
- References any issues addressed

Fixes #123
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## ⚡ Quality Standards

### Educational Value

All contributions must have clear educational value:

- **Progressive Learning**: Build on previous concepts appropriately
- **Clear Examples**: Code should be self-explanatory with comments
- **Real-World Relevance**: Examples should demonstrate practical applications
- **Multiple Approaches**: Show different ways to solve problems when beneficial

### Code Quality

- **Cyclomatic Complexity**: Keep functions under complexity 7
- **Function Length**: Limit functions to 80 lines for readability
- **Test Coverage**: Maintain >80% test coverage for new code
- **Error Handling**: Demonstrate proper Go error handling patterns

### Performance

- **Benchmarks**: Include benchmarks for performance-critical code
- **Memory Efficiency**: Consider memory allocations and cleanup
- **Concurrency Safety**: Ensure thread-safe code where applicable
- **Resource Management**: Properly close resources (files, connections, etc.)

## 📖 Documentation Standards

### README Requirements

Every directory must include a README.md with:

```markdown
# Module Name

Brief description of what this module teaches.

## Prerequisites
- List required knowledge
- Reference previous modules

## Learning Objectives  
- What skills will be developed
- What concepts will be understood

## Examples
- Brief description of each example
- How to run the examples
- Expected outcomes

## Key Concepts
- Important concepts demonstrated
- Best practices shown
- Common pitfalls avoided

## Next Steps
- Where to go after mastering this module
- Related concepts to explore
```

### Code Documentation

- **Package Comments**: Every package needs a clear description
- **Function Comments**: Explain purpose, parameters, return values, and examples
- **Complex Logic**: Comment non-obvious algorithms or business logic
- **Educational Notes**: Explain "why" not just "what" for learning purposes

## 🧪 Testing Requirements

### Test Coverage Standards

- **Unit Tests**: All public functions must have tests
- **Table-Driven Tests**: Use for multiple input scenarios  
- **Error Cases**: Test all error conditions
- **Edge Cases**: Include boundary value testing

### Example Test Structure

```go
func TestCalculateAverage(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected float64
        hasError bool
    }{
        {
            name:     "valid numbers",
            input:    []int{1, 2, 3, 4, 5},
            expected: 3.0,
            hasError: false,
        },
        {
            name:     "empty slice",
            input:    []int{},
            expected: 0,
            hasError: true,
        },
        {
            name:     "single number",
            input:    []int{42},
            expected: 42.0,
            hasError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateAverage(tt.input)
            
            if tt.hasError {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Tests

For complex examples, include integration tests that verify:
- Complete workflows function correctly
- External dependencies are handled properly
- Error recovery mechanisms work
- Performance meets expectations

## 🔍 Review Process

### Automated Checks

All pull requests undergo automated validation:

- **Code Quality**: GolangCI-Lint with strict educational standards
- **Test Coverage**: Minimum 80% coverage for new code
- **Documentation**: Link validation and completeness checks
- **Learning Path**: Validation that examples integrate properly
- **Security**: Vulnerability scanning and dependency audits

### Manual Review

Maintainers will review for:

- **Educational Value**: Does this improve learning outcomes?
- **Code Quality**: Is the code exemplary for students?
- **Documentation**: Are explanations clear and helpful?
- **Integration**: Does this fit well with existing content?
- **Consistency**: Does this match the repository's style and approach?

### Review Timeline

- **Small fixes**: 1-3 days
- **New examples**: 3-7 days  
- **Major features**: 1-2 weeks
- **Architecture changes**: 2-4 weeks

We aim to provide feedback quickly, but thorough review takes time to ensure educational quality.

## 🎖️ Recognition

Contributors will be recognized through:

- **Contributors section** in README
- **Commit attribution** preserved in git history
- **Changelog entries** for significant contributions
- **Special mentions** for outstanding contributions

## ❓ Getting Help

### Questions and Discussions

- **GitHub Issues**: For bugs, feature requests, and substantial discussions
- **GitHub Discussions**: For questions, ideas, and community chat
- **Pull Request Comments**: For specific code review discussions

### Finding Your First Contribution

Great places to start:
- **Documentation improvements**: Fix typos, clarify explanations
- **Example enhancements**: Add error handling, tests, or comments
- **New beginner examples**: Simple, focused learning examples
- **Issue triage**: Help categorize and reproduce reported issues

Look for these labels:
- `good first issue`: Perfect for new contributors
- `help wanted`: We'd appreciate community assistance
- `documentation`: Documentation improvements needed
- `beginner friendly`: Suitable for those new to Go

## 📈 Growing as a Contributor

### Contribution Ladder

1. **First-time Contributor**: Fix documentation, add tests
2. **Regular Contributor**: New examples, feature enhancements  
3. **Core Contributor**: Major features, architectural decisions
4. **Maintainer**: Code review, project direction, community management

### Skills You'll Develop

Contributing to Go Laboratory helps you:
- **Master Go**: Work with idiomatic, production-ready Go code
- **Learn Teaching**: Develop skills in explaining complex concepts
- **Practice Code Review**: Learn from experienced developers
- **Build Community**: Help grow the Go learning community
- **Advance Career**: Gain recognition in the Go ecosystem

---

**Thank you for contributing to Go Laboratory!** Your efforts help developers worldwide learn Go more effectively. Together, we're building the premier Go learning resource.

Questions? Open an issue or start a discussion. We're here to help! 🚀