# Testing Guide

This document provides comprehensive information about the testing strategy and implementation for the User Management API.

## Test Structure

The test suite is organized following the project's hexagonal architecture:

### Unit Tests
- **Domain Layer**: `internal/domain/*_test.go`
  - User model validation tests
  - Email validation tests
  - Business logic validation

- **Application Layer**: `internal/application/*_test.go` 
  - User service comprehensive tests
  - Error handling tests
  - Business rules validation

- **Infrastructure Layer**: `internal/infrastructure/**/.*_test.go`
  - HTTP handler tests
  - Repository implementation tests
  - External adapter tests

- **Configuration**: `internal/config/*_test.go`
  - Configuration validation tests
  - Environment variable handling tests

- **Health Package**: `pkg/health/*_test.go`
  - Health checker functionality tests
  - HTTP health endpoint tests
  - Concurrent access tests

### Integration Tests
- **Full API Lifecycle**: `integration_test.go`
  - End-to-end user management workflow
  - Error handling scenarios
  - Pagination testing

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race condition detection
go test -race ./...

# Run only unit tests (skip integration)
go test -short ./...

# Run specific test package
go test ./internal/application

# Run specific test
go test -run TestUserService_CreateUser ./internal/application
```

### Using Makefile

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Generate coverage report
make test-coverage

# Show coverage in terminal
make test-coverage-cli
```

### Docker Testing

```bash
# Run tests in Docker container
make docker-test

# Generate coverage report in Docker
make docker-coverage
```

## Test Coverage

The test suite aims for comprehensive coverage across all layers:

- **Domain Layer**: 100% - All business logic and validations
- **Application Layer**: 95%+ - All service methods and error paths
- **Infrastructure Layer**: 90%+ - All handlers, repositories, and adapters
- **Integration Tests**: Key user workflows and error scenarios

### Coverage Reports

Generate HTML coverage report:
```bash
make test-coverage
open coverage.html
```

View coverage summary:
```bash
make test-coverage-cli
```

## Testing Patterns

### 1. Unit Test Structure

Each unit test follows the Arrange-Act-Assert pattern:

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo, logger)
    
    // Set expectations
    mockRepo.EXPECT().GetByEmail(ctx, email).Return(nil, domain.ErrUserNotFound)
    mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
    
    // Act
    user, err := service.CreateUser(ctx, email, name)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, email, user.Email)
}
```

### 2. Mock Usage

The project uses GoMock for mocking dependencies:

- **Repository Mocks**: For testing service layer
- **Service Mocks**: For testing HTTP handlers
- **Interface Mocks**: For testing external dependencies

### 3. Test Data Management

Use the `testutil` package for consistent test data:

```go
// Create test users
user := testutil.CreateTestUser("1", "test@example.com", "Test User")
users := testutil.CreateTestUsers(5)

// Assert user equality
testutil.AssertUserEqual(t, expected, actual)
testutil.AssertUsersEqual(t, expectedSlice, actualSlice)
```

### 4. HTTP Testing

HTTP handlers are tested using `httptest`:

```go
func TestUserHandler_CreateUser(t *testing.T) {
    // Setup
    mockService := mocks.NewMockUserService(ctrl)
    handler := NewUserHandler(mockService, logger)
    
    // Prepare request
    reqBody := CreateUserRequest{Email: "test@example.com", Name: "Test User"}
    jsonBody, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
    w := httptest.NewRecorder()
    
    // Execute
    handler.createUser(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Test Categories

### 1. Happy Path Tests
- Valid inputs and expected outputs
- Successful operations
- Proper data flow

### 2. Error Path Tests
- Invalid inputs
- External service failures
- Edge cases and boundary conditions

### 3. Integration Tests
- Full request/response cycles
- Database integration
- Multi-component interactions

### 4. Concurrent Access Tests
- Race condition detection
- Thread safety validation
- Performance under load

## Mock Generation

The project uses `go generate` to create mocks:

```bash
# Generate all mocks
go generate ./...

# Or use the Makefile
make generate-mocks
```

Mock files are generated in the `mocks/` directory and follow the naming pattern `*_mocks.go`.

## Best Practices

### 1. Test Naming
- Use descriptive test names: `TestUserService_CreateUser_InvalidEmail`
- Include the component and scenario being tested
- Use table-driven tests for multiple scenarios

### 2. Test Organization
- One test file per source file (`user_service.go` → `user_service_test.go`)
- Group related tests using subtests
- Keep tests focused and independent

### 3. Assertions
- Use specific assertions rather than generic ones
- Include meaningful error messages
- Test both positive and negative cases

### 4. Test Data
- Use the `testutil` package for consistent test data
- Avoid hardcoded values in test assertions
- Create minimal test data that focuses on the scenario

### 5. Mocking
- Mock only external dependencies
- Set clear expectations on mock calls
- Verify mock expectations are met

## Continuous Integration

The test suite is designed to run in CI/CD pipelines:

```bash
# CI command that runs all checks
make ci

# This runs:
# - Linting (golangci-lint)
# - Tests with coverage
# - Race condition detection
# - Build verification
```

## Debugging Tests

### Running Individual Tests
```bash
# Run specific test with verbose output
go test -v -run TestUserService_CreateUser ./internal/application

# Run tests with debugging
go test -v -race -count=1 ./internal/application
```

### Common Issues

1. **Mock Expectations**: Ensure all mock expectations are set correctly
2. **Race Conditions**: Use `-race` flag to detect concurrent access issues
3. **Test Isolation**: Ensure tests don't depend on each other
4. **Resource Cleanup**: Use `defer` for proper cleanup

## Performance Testing

While not included in this suite, consider adding:
- Benchmark tests for critical paths
- Load testing for HTTP endpoints
- Memory usage profiling

```bash
# Run benchmarks (if implemented)
go test -bench=. ./...

# Profile memory usage
go test -memprofile mem.prof ./...
```

## Future Enhancements

Consider adding:
- Contract testing for external APIs
- Property-based testing for complex business logic
- End-to-end testing with real database
- Performance regression testing