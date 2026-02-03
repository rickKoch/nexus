# Nexus

Nexus is a segment management service built with Go. It provides a REST API for
creating, reading, updating, and deleting segments with optional TTL
(time-to-live) settings.

## Getting Started

### Running with Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
docker compose up -d
```

This will start:
- **PostgreSQL** database on port `5432`
- **Segment Service** API on port `8080`

To stop the services:

```bash
docker compose down
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `POSTGRES_HOST` | PostgreSQL host | `postgres` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | PostgreSQL username | `nexus` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `nexus` |
| `POSTGRES_DB` | PostgreSQL database name | `nexus` |
| `POSTGRES_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | `*` |

## API Reference

The service exposes a REST API for managing segments.

### Base URL

```
http://localhost:8080
```

### Endpoints

#### List Segments

```http
GET /segment
```

**Response:**

```json
[
  {
    "id": 1,
    "name": "premium-users",
    "ttl_seconds": 3600,
    "created_at": "2026-02-03T10:00:00Z",
    "updated_at": "2026-02-03T10:00:00Z"
  }
]
```

#### Get Segment

```http
GET /segment/:id
```

**Response:**

```json
{
  "id": 1,
  "name": "premium-users",
  "ttl_seconds": 3600,
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

#### Create Segment

```http
POST /segment
Content-Type: application/json

{
  "name": "premium-users",
  "ttl_seconds": 3600
}
```

**Response:** `201 Created`

```json
{
  "id": 1,
  "name": "premium-users",
  "ttl_seconds": 3600,
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

#### Update Segment

```http
PUT /segment/:id
Content-Type: application/json

{
  "name": "vip-users",
  "ttl_seconds": 7200
}
```

**Response:** `200 OK`

```json
{
  "id": 1,
  "name": "vip-users",
  "ttl_seconds": 7200,
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T12:00:00Z"
}
```

#### Delete Segment

```http
DELETE /segment/:id
```

**Response:** `204 No Content`

## Development

### Project Structure

```
├── internal/
│   └── segments/           # Segment service module
│       ├── adapters/       # Database adapters
│       ├── app/            # Application layer (use cases)
│       ├── domain/         # Domain entities
│       ├── port/           # HTTP handlers
│       └── service/        # Service configuration
├── pkg/                    # Shared packages
├── scripts/                # Build and lint scripts
├── sql/                    # Database schema
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

### Available Make Commands

```bash
# Run linter
make lint

# Format code
make fmt

# Run tests
make test
```

### Running Tests

```bash
make test
```

### Running Locally

1. Start the PostgreSQL database:

   ```bash
   docker compose up -d postgres
   ```

2. Run the service:

   ```bash
   go run ./internal/segments
   ```

## Database Schema

The `segments` table schema:

```sql
CREATE TABLE segments (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  ttl_seconds INT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);
```
