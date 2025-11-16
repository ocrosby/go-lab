# Go Workspaces

Go workspaces provide a way to work with multiple Go modules simultaneously within a single directory tree. Introduced in Go 1.18, workspaces allow you to develop across multiple modules without needing to publish them or use replace directives.

## What are Go Workspaces?

A Go workspace is a collection of modules that are developed and built together. The workspace is defined by a `go.work` file at the root of the workspace directory, which specifies the modules that are part of the workspace.

## Key Benefits

- **Multi-module development**: Work on multiple related modules simultaneously
- **Local development**: Make changes across modules without publishing intermediate versions
- **Simplified dependency management**: Avoid complex replace directives in go.mod files
- **Consistent builds**: Ensure all modules in the workspace use compatible versions

## Creating a Workspace

### Initialize a workspace
```bash
go work init
```

### Add modules to the workspace
```bash
go work use ./module1 ./module2
```

### Initialize with modules
```bash
go work init ./module1 ./module2
```

## The go.work File

The `go.work` file defines the workspace structure:

```go
go 1.21

use (
    ./api
    ./shared
    ./client
)

replace example.com/old => example.com/new v1.2.3
```

### Sections:
- `go`: Specifies the Go version
- `use`: Lists directories containing modules to include in the workspace
- `replace`: Module replacements (similar to go.mod replace directives)

## Common Commands

### View workspace status
```bash
go work status
```

### Sync workspace dependencies
```bash
go work sync
```

### Edit workspace
```bash
go work edit -use=./newmodule
go work edit -dropuse=./oldmodule
```

## Example Workspace Structure

```
myproject/
├── go.work
├── api/
│   ├── go.mod
│   ├── main.go
│   └── handlers/
├── shared/
│   ├── go.mod
│   └── models/
└── client/
    ├── go.mod
    └── main.go
```

## Best Practices

1. **Use workspaces for related modules**: Keep modules that are developed together in the same workspace
2. **Version compatibility**: Ensure all modules use compatible Go versions
3. **Clean boundaries**: Maintain clear module boundaries and interfaces
4. **CI/CD considerations**: Remember that workspace files are typically not committed to version control
5. **Team coordination**: Document workspace setup for team members

## Workspace vs Other Approaches

| Approach | Use Case | Pros | Cons |
|----------|----------|------|------|
| Workspaces | Multi-module local development | Clean, no go.mod changes | Go 1.18+ only |
| Replace directives | Module replacement | Works with older Go versions | Clutters go.mod |
| Vendor directories | Dependency management | Full control over deps | Large repository size |

## Common Pitfalls

- **Version conflicts**: Different modules requiring incompatible versions
- **Circular dependencies**: Modules depending on each other
- **Build context**: Commands run from workspace root vs module directories
- **Publishing**: Remember to test modules independently before publishing

## Integration with IDEs

Most modern Go IDEs support workspaces:
- **VS Code**: Install the Go extension and open the workspace root
- **GoLand**: Open the directory containing go.work
- **Vim/Neovim**: Use gopls language server with workspace support

## Conclusion

Go workspaces simplify multi-module development by providing a clean, standardized way to work with related modules. They eliminate the need for complex replace directives and make local development more straightforward while maintaining the benefits of modular architecture.