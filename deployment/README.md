# Deployment Examples & Infrastructure 🚀

Production deployment patterns and infrastructure examples for Go applications.

## Overview

This directory contains real-world deployment configurations and infrastructure-as-code examples for deploying Go applications in various environments.

## Deployment Strategies

### Local Development
- **Docker Compose**: Multi-service local development
- **Hot Reload**: Development workflow optimization
- **Environment Configuration**: Local environment setup

### Container Platforms
- **Docker**: Containerization best practices
- **Kubernetes**: Production-grade orchestration
- **Helm Charts**: Kubernetes package management

### Cloud Platforms
- **AWS**: ECS, EKS, Lambda, and EC2 deployments
- **Google Cloud**: GKE, Cloud Run, and App Engine
- **Azure**: AKS, Container Instances, and App Service

### CI/CD Pipelines
- **GitHub Actions**: Automated testing and deployment
- **GitLab CI/CD**: Complete DevOps pipeline
- **Jenkins**: Traditional CI/CD approach

## Directory Structure

```
deployment/
├── README.md                    # This file
├── docker/
│   ├── Dockerfile.dev           # Development container
│   ├── Dockerfile.prod          # Production container
│   ├── docker-compose.yml       # Local development stack
│   └── docker-compose.prod.yml  # Production-like local stack
├── kubernetes/
│   ├── namespace.yaml           # Namespace definition
│   ├── deployment.yaml          # Application deployment
│   ├── service.yaml             # Service definition
│   ├── ingress.yaml             # Ingress configuration
│   ├── configmap.yaml           # Configuration management
│   ├── secrets.yaml             # Secrets management
│   └── hpa.yaml                 # Horizontal Pod Autoscaler
├── helm/
│   ├── Chart.yaml               # Helm chart definition
│   ├── values.yaml              # Default values
│   ├── values-dev.yaml          # Development values
│   ├── values-prod.yaml         # Production values
│   └── templates/               # Helm templates
├── cloud/
│   ├── aws/                     # AWS deployment examples
│   ├── gcp/                     # Google Cloud examples
│   └── azure/                   # Azure deployment examples
├── ci-cd/
│   ├── github-actions/          # GitHub Actions workflows
│   ├── gitlab-ci/               # GitLab CI configuration
│   └── jenkins/                 # Jenkins pipeline
└── monitoring/
    ├── prometheus/              # Prometheus configuration
    ├── grafana/                 # Grafana dashboards
    └── alertmanager/            # Alert configuration
```

## Docker Deployment

### Development Dockerfile
```dockerfile
# docker/Dockerfile.dev
FROM golang:1.21-alpine AS development

# Install hot reload tool
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Use air for hot reload
CMD ["air"]
```

### Production Dockerfile
```dockerfile
# docker/Dockerfile.prod
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
```

### Docker Compose for Development
```yaml
# docker/docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: ..
      dockerfile: docker/Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - ../:/app
      - /app/vendor
    environment:
      - ENV=development
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

## Kubernetes Deployment

### Application Deployment
```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
  namespace: go-lab
  labels:
    app: go-app
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
        version: v1
    spec:
      containers:
      - name: go-app
        image: go-lab/app:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db.host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: db.password
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

### Service Definition
```yaml
# kubernetes/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: go-app-service
  namespace: go-lab
  labels:
    app: go-app
spec:
  selector:
    app: go-app
  ports:
  - name: http
    port: 80
    targetPort: 8080
  type: ClusterIP
```

### Ingress Configuration
```yaml
# kubernetes/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-app-ingress
  namespace: go-lab
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.yourapp.com
    secretName: go-app-tls
  rules:
  - host: api.yourapp.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-app-service
            port:
              number: 80
```

## Helm Charts

### Chart Definition
```yaml
# helm/Chart.yaml
apiVersion: v2
name: go-app
description: A Helm chart for Go application
type: application
version: 0.1.0
appVersion: "1.0.0"
```

### Values Configuration
```yaml
# helm/values.yaml
replicaCount: 3

image:
  repository: go-lab/app
  pullPolicy: IfNotPresent
  tag: "latest"

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.yourapp.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: go-app-tls
      hosts:
        - api.yourapp.com

resources:
  requests:
    memory: 64Mi
    cpu: 250m
  limits:
    memory: 128Mi
    cpu: 500m

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

config:
  env: production
  logLevel: info

secrets:
  dbPassword: "your-secret-password"
```

## CI/CD Pipelines

### GitHub Actions Workflow
```yaml
# ci-cd/github-actions/deploy.yml
name: Build and Deploy

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Run linting
      uses: golangci/golangci-lint-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Log in to registry
      uses: docker/login-action@v2
      with:
        registry: ${{ secrets.REGISTRY }}
        username: ${{ secrets.REGISTRY_USERNAME }}
        password: ${{ secrets.REGISTRY_PASSWORD }}
    
    - name: Build and push
      uses: docker/build-push-action@v3
      with:
        context: .
        file: docker/Dockerfile.prod
        push: true
        tags: ${{ secrets.REGISTRY }}/go-app:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up kubectl
      uses: azure/setup-kubectl@v3
    
    - name: Deploy to Kubernetes
      env:
        KUBE_CONFIG: ${{ secrets.KUBE_CONFIG }}
      run: |
        echo "$KUBE_CONFIG" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig
        
        # Update image tag
        kubectl set image deployment/go-app go-app=${{ secrets.REGISTRY }}/go-app:${{ github.sha }} -n go-lab
        
        # Wait for rollout
        kubectl rollout status deployment/go-app -n go-lab
```

## Cloud Platform Examples

### AWS ECS Task Definition
```json
{
  "family": "go-app",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "go-app",
      "image": "your-registry/go-app:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:db-password"
        }
      ],
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3
      },
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/go-app",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Run Configuration
```yaml
# cloud/gcp/cloudrun.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: go-app
  annotations:
    run.googleapis.com/ingress: all
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "100"
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containerConcurrency: 80
      containers:
      - image: gcr.io/project-id/go-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: production
        resources:
          limits:
            cpu: "1"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
```

## Monitoring and Observability

### Prometheus Configuration
```yaml
# monitoring/prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'go-app'
    static_configs:
      - targets: ['go-app-service:80']
    metrics_path: /metrics
    scrape_interval: 30s

rule_files:
  - "go-app-rules.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

### Grafana Dashboard
```json
{
  "dashboard": {
    "title": "Go Application Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "5xx errors"
          }
        ]
      }
    ]
  }
}
```

## Best Practices

### Security
- **Image Scanning**: Scan container images for vulnerabilities
- **Secrets Management**: Use proper secrets management (K8s secrets, AWS Secrets Manager, etc.)
- **Network Policies**: Implement network segmentation
- **RBAC**: Use role-based access control
- **Security Contexts**: Run containers with non-root users

### Performance
- **Resource Limits**: Set appropriate CPU and memory limits
- **Health Checks**: Implement proper liveness and readiness probes
- **Autoscaling**: Configure horizontal pod autoscaling
- **Load Balancing**: Use appropriate load balancing strategies
- **Caching**: Implement caching layers where appropriate

### Reliability
- **Multi-Zone Deployment**: Deploy across multiple availability zones
- **Rolling Updates**: Use rolling deployment strategies
- **Circuit Breakers**: Implement circuit breaker patterns
- **Graceful Shutdown**: Handle SIGTERM signals properly
- **Backup and Recovery**: Plan for disaster recovery

## Getting Started

1. **Choose Your Platform**: Select deployment target (local, K8s, cloud)
2. **Review Examples**: Study relevant configuration files
3. **Customize Configuration**: Adapt to your specific needs
4. **Test Locally**: Use Docker Compose for local testing
5. **Deploy and Monitor**: Deploy to target environment and set up monitoring

Each deployment strategy has its own trade-offs in terms of complexity, cost, and scalability. Choose the approach that best fits your requirements and team expertise.

## Integration with Examples

These deployment configurations work with:
- [Production API](../examples/03-advanced/production-api/) - Complete application ready for deployment
- [HTTP Service Template](../templates/http-service/) - Template with deployment configurations
- [Performance Benchmarks](../examples/03-advanced/performance-benchmarks/) - Load testing deployed applications

Master these deployment patterns to successfully run Go applications in production environments!