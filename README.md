# Chirpy HTTP Server

A RESTful API server for a Twitter-like social media platform built with Go.  Chirpy allows users to create accounts, post chirps (tweets), and manage their profiles with JWT-based authentication. 

## Notes

This is a guided project I compeleted on [BootDev](https://boot.dev)'s Golang backend course.

## Features

-  User authentication with JWT tokens
-  Create, read, and delete chirps
-  Token refresh mechanism
-  Premium user upgrades via webhook integration
-  Admin metrics and monitoring
-  PostgreSQL database with sqlc for type-safe queries

## Installation

### 1. Install Go 1.25 or later

[The webi installer](https://webinstall.dev/golang/) is the simplest way for most people.  Just run this in your terminal:

```bash
curl -sS https://webi.sh/golang | sh
```

_Read the output of the command and follow any instructions._

### 2. Install PostgreSQL

Chirpy uses PostgreSQL as its database.  You'll need to have PostgreSQL installed and running on your system.

**For Mac OS**: 
```bash
brew install postgresql
brew services start postgresql
```

**For Linux/WSL**:
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo service postgresql start
```

Create a database for Chirpy:
```bash
createdb chirpy
```

### 3. Clone and Setup the Project

Clone the repository:
```bash
git clone https://github.com/rangaroo/chirpy-http-server.git
cd chirpy-http-server
```

Install dependencies:
```bash
go mod download
```

### 4. Configure Environment Variables

Create a `.env` file in the project root with the following variables:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
TOKEN_SECRET=your-secret-key-here
POLKA_KEY=your-polka-api-key
```

**Important:**
- `DB_URL`: Your PostgreSQL connection string
- `PLATFORM`: Environment (e.g., "dev", "prod")
- `TOKEN_SECRET`: Secret key for JWT token signing (use a strong random string)
- `POLKA_KEY`: API key for Polka webhook authentication

### 5. Run Database Migrations

Apply the SQL migrations to set up your database schema:
```bash
# Navigate to the sql directory and apply migrations
# The exact command depends on your migration tool
```

### 6. Start the Server

Run the server:
```bash
go build -o chirpy-http-server
./chirpy-http-server
```

The server will start on `http://localhost:8080`.

## API Endpoints

### Health & Metrics

- `GET /api/healthz` - Check server health
- `GET /admin/metrics` - View server metrics (admin only)
- `POST /admin/reset` - Reset server metrics (admin only)

### Users

- `POST /api/users` - Create a new user
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```

- `PUT /api/users` - Update user information (requires authentication)
  ```json
  {
    "email": "newemail@example.com",
    "password": "newpassword"
  }
  ```

### Authentication

- `POST /api/login` - Login and receive JWT tokens
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```

- `POST /api/refresh` - Refresh access token using refresh token

- `POST /api/revoke` - Revoke refresh token (logout)

### Chirps

- `POST /api/chirps` - Create a new chirp (requires authentication)
  ```json
  {
    "body": "This is my chirp!"
  }
  ```

- `GET /api/chirps` - Get all chirps

- `GET /api/chirps/{chirpID}` - Get a specific chirp by ID

- `DELETE /api/chirps/{chirpID}` - Delete a chirp (requires authentication, must be owner)

### Webhooks

- `POST /api/polka/webhooks` - Polka webhook for user upgrades (requires API key)

## Authentication

Chirpy uses JWT (JSON Web Tokens) for authentication: 

1. **Login** to receive an access token and refresh token
2. **Access token** is used for authenticated requests (add to `Authorization:  Bearer <token>` header)
3. **Refresh token** is used to get a new access token when it expires
4. **Revoke** refresh tokens to log out

## Architecture

Chirpy is built with: 
- **Go** - Core language using standard library `net/http`
- **PostgreSQL** - Database for storing users and chirps
- **sqlc** - Type-safe SQL query generation
- **JWT** - Token-based authentication (`golang-jwt/jwt/v5`)
- **Argon2** - Password hashing (`alexedwards/argon2id`)
- **godotenv** - Environment variable management

The project structure includes:
- Handler files for each API endpoint
- Database queries generated with sqlc
- Internal packages for database access
- SQL migrations for database schema
- Middleware for metrics and request logging

## Development

### Running Tests

```bash
go test ./...
```
