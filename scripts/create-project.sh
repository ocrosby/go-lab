#!/bin/bash

# Go Laboratory Project Generator
# Creates new Go projects using repository templates and best practices

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
TEMPLATE=""
PROJECT_NAME=""
MODULE_PATH=""
OUTPUT_DIR=""
VERBOSE=false

# Available templates
declare -a TEMPLATES=("http-service" "cli-tool" "library")

# Print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Print usage information
usage() {
    cat << EOF
Go Laboratory Project Generator

Usage: $0 [OPTIONS]

OPTIONS:
    -t, --template TEMPLATE     Template to use (${TEMPLATES[*]})
    -n, --name NAME            Project name
    -m, --module MODULE        Go module path (e.g., github.com/user/project)
    -o, --output DIR          Output directory (default: current directory)
    -v, --verbose             Verbose output
    -h, --help                Show this help message

EXAMPLES:
    # Create HTTP service
    $0 -t http-service -n my-api -m github.com/myuser/my-api

    # Create CLI tool
    $0 -t cli-tool -n my-tool -m github.com/myuser/my-tool -o ~/projects

    # Create library
    $0 -t library -n my-lib -m github.com/myuser/my-lib

TEMPLATES:
    http-service    REST API service with middleware and testing
    cli-tool        Command-line application with Cobra CLI
    library         Reusable Go library with comprehensive testing

EOF
}

# Validate template exists
validate_template() {
    local template=$1
    for tmpl in "${TEMPLATES[@]}"; do
        if [[ "$tmpl" == "$template" ]]; then
            return 0
        fi
    done
    return 1
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Using Go version: $go_version"
}

# Check if template exists
check_template() {
    local template_dir="templates/$TEMPLATE"
    if [[ ! -d "$template_dir" ]]; then
        print_error "Template '$TEMPLATE' not found in $template_dir"
        exit 1
    fi
}

# Validate project name
validate_project_name() {
    if [[ ! "$PROJECT_NAME" =~ ^[a-zA-Z][a-zA-Z0-9_-]*$ ]]; then
        print_error "Invalid project name. Use letters, numbers, hyphens, and underscores only."
        exit 1
    fi
}

# Validate module path
validate_module_path() {
    if [[ ! "$MODULE_PATH" =~ ^[a-zA-Z][a-zA-Z0-9._/-]*$ ]]; then
        print_error "Invalid module path. Use standard Go module naming conventions."
        exit 1
    fi
}

# Create project directory
create_project_dir() {
    local target_dir="$OUTPUT_DIR/$PROJECT_NAME"
    
    if [[ -d "$target_dir" ]]; then
        print_error "Directory '$target_dir' already exists"
        exit 1
    fi
    
    mkdir -p "$target_dir"
    print_success "Created project directory: $target_dir"
    echo "$target_dir"
}

# Copy template files
copy_template() {
    local template_dir="templates/$TEMPLATE"
    local target_dir=$1
    
    print_info "Copying template files from $template_dir..."
    
    # Copy all files except README.md (will be customized)
    find "$template_dir" -type f -not -name "README.md" -not -name ".template-config" | while read -r file; do
        local rel_path=${file#$template_dir/}
        local target_file="$target_dir/$rel_path"
        
        # Create directory if needed
        mkdir -p "$(dirname "$target_file")"
        
        # Copy file
        cp "$file" "$target_file"
        
        if [[ "$VERBOSE" == true ]]; then
            print_info "  Copied: $rel_path"
        fi
    done
}

# Customize files with project-specific values
customize_files() {
    local target_dir=$1
    
    print_info "Customizing files with project-specific values..."
    
    # Find all Go files and customize them
    find "$target_dir" -name "*.go" -o -name "go.mod" -o -name "*.yml" -o -name "*.yaml" -o -name "Dockerfile" -o -name "Makefile" | while read -r file; do
        # Replace template placeholders
        sed -i.bak \
            -e "s/TEMPLATE_PROJECT_NAME/$PROJECT_NAME/g" \
            -e "s|TEMPLATE_MODULE_PATH|$MODULE_PATH|g" \
            -e "s/template-project/$PROJECT_NAME/g" \
            -e "s|github.com/template/|$MODULE_PATH|g" \
            "$file"
        
        # Remove backup file
        rm -f "$file.bak"
        
        if [[ "$VERBOSE" == true ]]; then
            print_info "  Customized: ${file#$target_dir/}"
        fi
    done
}

# Create custom README
create_readme() {
    local target_dir=$1
    local template_readme="templates/$TEMPLATE/README.md"
    local target_readme="$target_dir/README.md"
    
    if [[ -f "$template_readme" ]]; then
        # Create customized README
        cat > "$target_readme" << EOF
# $PROJECT_NAME

Generated from Go Laboratory $TEMPLATE template.

## Description

Add your project description here.

## Prerequisites

- Go 1.21 or later
- [Add any other prerequisites]

## Installation

\`\`\`bash
git clone [your-repo-url]
cd $PROJECT_NAME
go mod tidy
\`\`\`

## Usage

\`\`\`bash
# Add usage examples here
go run main.go
\`\`\`

## Development

\`\`\`bash
# Run tests
make test

# Run linting
make lint

# Build
make build
\`\`\`

## Based On

This project was generated using the Go Laboratory [$TEMPLATE template](https://github.com/ocrosby/go-lab/tree/main/templates/$TEMPLATE).

For more advanced features, see the [Go Laboratory documentation](https://github.com/ocrosby/go-lab).
EOF
        
        print_success "Created customized README.md"
    fi
}

# Initialize Go module
init_go_module() {
    local target_dir=$1
    
    print_info "Initializing Go module..."
    
    cd "$target_dir"
    
    # Initialize module
    go mod init "$MODULE_PATH"
    
    # Download dependencies
    go mod tidy
    
    print_success "Initialized Go module: $MODULE_PATH"
}

# Run quality checks
run_quality_checks() {
    local target_dir=$1
    
    print_info "Running quality checks..."
    
    cd "$target_dir"
    
    # Check if code compiles
    if ! go build -o /tmp/test-build ./...; then
        print_error "Generated code does not compile"
        return 1
    fi
    
    # Run tests if any exist
    if find . -name "*_test.go" | grep -q .; then
        if ! go test ./...; then
            print_error "Tests are failing"
            return 1
        fi
        print_success "All tests passing"
    fi
    
    # Run linting if golangci-lint is available
    if command -v golangci-lint &> /dev/null; then
        if golangci-lint run --timeout 5m; then
            print_success "Code quality checks passed"
        else
            print_warning "Some linting issues found (fixable)"
        fi
    fi
    
    print_success "Quality checks completed"
}

# Main execution
main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--template)
                TEMPLATE="$2"
                shift 2
                ;;
            -n|--name)
                PROJECT_NAME="$2"
                shift 2
                ;;
            -m|--module)
                MODULE_PATH="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    # Set defaults
    OUTPUT_DIR=${OUTPUT_DIR:-$(pwd)}
    
    # Validate required arguments
    if [[ -z "$TEMPLATE" ]]; then
        print_error "Template is required. Use -t or --template."
        usage
        exit 1
    fi
    
    if [[ -z "$PROJECT_NAME" ]]; then
        print_error "Project name is required. Use -n or --name."
        usage
        exit 1
    fi
    
    if [[ -z "$MODULE_PATH" ]]; then
        print_error "Module path is required. Use -m or --module."
        usage
        exit 1
    fi
    
    # Validate inputs
    if ! validate_template "$TEMPLATE"; then
        print_error "Invalid template. Available templates: ${TEMPLATES[*]}"
        exit 1
    fi
    
    validate_project_name
    validate_module_path
    
    # Check prerequisites
    check_go
    check_template
    
    print_info "Creating project '$PROJECT_NAME' using template '$TEMPLATE'..."
    print_info "Module path: $MODULE_PATH"
    print_info "Output directory: $OUTPUT_DIR"
    
    # Create project
    local target_dir
    target_dir=$(create_project_dir)
    
    copy_template "$target_dir"
    customize_files "$target_dir"
    create_readme "$target_dir"
    
    # Initialize Go module
    init_go_module "$target_dir"
    
    # Run quality checks
    if ! run_quality_checks "$target_dir"; then
        print_warning "Quality checks had issues, but project was created successfully"
    fi
    
    # Success message
    print_success "Project '$PROJECT_NAME' created successfully!"
    echo
    print_info "Next steps:"
    echo "  1. cd $PROJECT_NAME"
    echo "  2. Open in your editor and start coding!"
    echo "  3. Run 'make test' to verify everything works"
    echo "  4. See README.md for more information"
    echo
    print_info "Happy coding! 🚀"
}

# Execute main function
main "$@"