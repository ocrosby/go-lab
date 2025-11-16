# Swagger Documentation

This API automatically generates Swagger/OpenAPI documentation from code annotations.

## Accessing the Documentation

Once the API is running, you can access the Swagger UI at:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json

## Regenerating Documentation

To regenerate the Swagger documentation after making changes to API annotations:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go
```

## API Endpoints

The following endpoints are documented:

### Users
- `GET /users` - List users with pagination
- `POST /users` - Create a new user
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Health Checks
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe
- `GET /startupz` - Startup probe

## Docker

The Dockerfile automatically generates Swagger documentation during the build process.

## Development

When adding new endpoints or modifying existing ones:
1. Add proper Swagger annotations to your handler functions
2. Run `swag init -g cmd/api/main.go` to regenerate documentation
3. Test the endpoint in Swagger UI