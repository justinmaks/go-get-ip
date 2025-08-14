# go-get-ip

A lightweight, production-ready Go API service that returns the client's IP address. Designed to work behind reverse proxies like Nginx and Cloudflare with IP detection fallback strategies.

## API Endpoints

### `GET /`
Returns the client's IP address (fastest detection - IPv4 or IPv6).

**Response:**
```json
{
  "ip": "192.168.1.100"
}
```

### `GET /ipv4`
Returns only IPv4 address if available.

**Response:**
```json
{
  "ip": "192.168.1.100"
}
```

**Error Response (404):**
```json
{
  "error": "No IPv4 address found"
}
```

### `GET /ipv6`
Returns only IPv6 address if available.

**Response:**
```json
{
  "ip": "2001:db8::1"
}
```

**Error Response (404):**
```json
{
  "error": "No IPv6 address found"
}
```

### `GET /health`
Health check endpoint for monitoring and load balancers.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1692000000
}
```

## IP Detection Strategy

The service checks headers in this priority order:
1. `CF-Connecting-IP` (Cloudflare)
2. `X-Forwarded-For` (takes first IP from comma-separated list)
3. `X-Real-IP` (Nginx)
4. `X-Client-IP`
5. `X-Forwarded`
6. `X-Cluster-Client-IP`
7. `Forwarded-For`
8. `Forwarded`
9. Falls back to `RemoteAddr`

## Quick Start

### Using Docker Compose (Recommended)

1. Clone this repository
2. Build and run:
```bash
docker-compose up --build -d
```

The service will be available at `http://localhost:3000`

### Using Go Directly

1. Install dependencies:
```bash
go mod download
```

2. Run the service:
```bash
go run main.go
```

### Using Docker

```bash
# Build
docker build -t go-get-ip .

# Run
docker run -p 3000:3000 go-get-ip
```

## Configuration

### Environment Variables

- `GIN_MODE`: Set to `release` for production (default in Docker)
- `PORT`: Port number (defaults to 3000)

### Docker Compose Configuration

The included `docker-compose.yml` is pre-configured for Nginx Proxy Manager:

- **Network**: `npm-network` (external)
- **Port**: 3000 (internal and external)
- **Health Check**: Configured with 30s intervals
- **Restart Policy**: `unless-stopped`

## Production Deployment

### With Nginx Proxy Manager

1. Ensure the `npm-network` exists:
```bash
docker network create npm-network
```

2. Deploy the service:
```bash
docker-compose up -d
```

3. In Nginx Proxy Manager:
   - Create a new proxy host
   - Set destination to `go-get-ip-app-npmnetwork:3000`
   - Configure your domain and SSL

### Health Monitoring

The service includes comprehensive health monitoring:

- **Health Endpoint**: `/health` returns service status
- **Docker Health Check**: Automatic container health monitoring
- **Graceful Shutdown**: 5-second timeout for clean shutdowns
- **Request Logging**: All requests logged in production

### Production Features

- **Graceful Shutdown**: Handles SIGINT and SIGTERM signals
- **Request Recovery**: Automatic panic recovery middleware
- **Minimal Attack Surface**: Alpine-based container (< 10MB)
- **No Root User**: Runs as non-root for security
- **Health Checks**: Built-in monitoring endpoints

## Development

### Testing the Service

```bash
# Test root endpoint
curl http://localhost:3000/

# Test IPv4 endpoint
curl http://localhost:3000/ipv4

# Test IPv6 endpoint  
curl http://localhost:3000/ipv6

# Test health check
curl http://localhost:3000/health
```

### Building for Production

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Build Docker image
docker build -t go-get-ip:latest .
```

## Architecture

- **Language**: Go 1.21
- **Framework**: Gin (high-performance HTTP framework)
- **Container**: Multi-stage Docker build with Alpine Linux
- **Size**: ~10MB final image
- **Performance**: Sub-millisecond response times

## Security Considerations

- No external dependencies beyond Gin framework
- Input validation on all IP addresses
- No logging of sensitive information
- Minimal container attack surface
- Non-root container execution

## License

This project is open source and available under the MIT License.
