# Chirpy HTTP Server

A simple lightweight HTTP server project built with Go.

## Overview

Chirpy is a simple yet functional HTTP server implementation that demonstrates core concepts of web development in Go. This project covers routing, middleware, request handling, and RESTful API design patterns.

## Features

- [x] HTTP/1.1 compliant server
- [x] RESTful API endpoints
- [x] Static file serving
- [x] Request validation
- [x] Error handling
- [x] Structured logging
- [x] CORS support (if applicable)

## Project Structure

```
chirpy-http-server/
├── README.md           # Project documentation
├── go.mod             # Go module definition
├── go.sum             # Dependency checksums
├── main.go            # Application entry point
├── handler/           # HTTP request handlers
```

## Installation

### Prerequisites

- Go 1.22 or higher
- Git
- sqlc

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/rangaroo/chirpy-http-server
cd chirpy-http-server

# Download dependencies
go mod download

# Build the project
go build -o chirpy
```

## Usage

### Running the Server

```bash
# Run directly with Go
go run .

```

The server will start on `http://localhost:8080` by default.

### Environment Variables

```bash
PORT=8080              # Server port (default: 8080)
DEBUG=true             # Enable debug mode
```

## API Endpoints

### Health Check

```
GET /api/healthz
```

Returns server health status.

**Response:**
```json
{
    "status": "ok"
}
```

### Chirps

```
GET    /api/chirps                 # Get all chirps
GET    /api/chirps/{chirpID}       # Get chirp by ID
POST   /api/chirps                 # Create new chirp
DELETE /api/chirps/{chirpID}       # Delete chirp
```

### Static Files

```
GET /                        # Serves index.html
GET /assets/*                # Serves static assets
```

## Configuration

Configuration can be managed through environment variables or a config file (if implemented).

## Testing

```bash
# Run all tests
go test ./...

```

## Acknowledgments

- Built as part of [Boot.dev](https://boot.dev) backend development course