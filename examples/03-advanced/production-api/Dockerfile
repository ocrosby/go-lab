FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Install swag for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate Swagger docs
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/api .

EXPOSE 8080 8081

CMD ["./api", "server"]