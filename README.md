# Chirpy

A RESTful microblogging backend in Go, inspired by Twitter. Supports user accounts, short posts (chirps), JWT-based authentication with token refresh/revocation, and webhook event handling. Built with PostgreSQL via `sqlc` for type-safe database access.

## Features

- Create, read, and delete chirps (140-character posts)
- User registration and login with bcrypt password hashing
- JWT access tokens + refresh token rotation
- Token revocation (logout)
- Webhook endpoint for external event integration
- Type-safe SQL queries generated with `sqlc`

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Query generation:** sqlc
- **Auth:** JWT (access + refresh tokens)

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (optional, for regenerating queries)

### Setup

1. Clone the repository

```bash
git clone https://github.com/Ocidemus/chirpy.git
cd chirpy
```

2. Install dependencies

```bash
go mod download
```

3. Configure environment

```bash
cp .env.example .env
```

Update `.env` with your values:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy
JWT_SECRET=your_secret_here
POLKA_KEY=your_webhook_key_here
PORT=8080
```

4. Run the server

```bash
go run .
```

## API Overview

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create a new user |
| PUT | `/api/users` | Update user details |

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/login` | Login and receive tokens |
| POST | `/api/refresh` | Refresh access token |
| POST | `/api/revoke` | Revoke refresh token |

### Chirps

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/chirps` | Create a chirp |
| GET | `/api/chirps` | List all chirps |
| GET | `/api/chirps/{id}` | Get a chirp by ID |
| DELETE | `/api/chirps/{id}` | Delete a chirp |

### Webhooks

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/polka/webhooks` | Receive webhook events |

## Project Structure

```
chirpy/
├── main.go                  — server setup and routing
├── handler_chirps.go        — chirp CRUD handlers
├── handler_users_create.go  — user registration
├── handler_login.go         — login and token issuance
├── handler_refresh.go       — token refresh
├── handler_revoke.go        — token revocation
├── handler_update.go        — user update
├── handler_webhooks.go      — webhook handler
├── json.go                  — JSON response helpers
├── reset.go                 — dev reset endpoint
├── internal/                — db queries and auth logic
├── sql/                     — SQL schema and migrations
└── sqlc.yaml                — sqlc configuration
```

## Security Note

JWT secrets and database credentials are stored in `.env` and never committed. Rotate your `JWT_SECRET` if exposed.
