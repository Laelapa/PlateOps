#  PlateOps

<img align="left" width="320" src="PlateOpsMascot.png">

PlateOps is a registry of known foods & recipes, and a food stock management and daily nutrition tracking service. It provides an API with user authentication and it is built to serve as a backend for consumer-facing applications.


---

## Table of Contents

- [Features](#features)
- [Technical Highlights](#technical-highlights)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
- [Getting Started](#getting-started)
- [Security Considerations](#security-considerations)
- [Development Notes](#development-notes)
- [License](#license)

---

## Features

### User Features

- **User Registration & Authentication**
  - Secure signup and login endpoints with password hashing (bcrypt).
  - JWT-based stateless authentication for accessing protected endpoints.
  - Refresh token mechanism for secure session renewal.
  - Logout and token revocation support.
 
- **OpenAPI Documentation**
  - Self-hosted OpenAPI JSON and Swagger UI for API exploration.

- **Food Registry**
  - Create, update, and search food entries.
  - Retrieve food by ID or GTIN (barcode).
  - Fuzzy search for foods by name with tolerance for typos.

- **Inventory Management \***
  - Track owned food items, quantities, and expiration dates.
  - Track package opening and usage to monitor freshness.

- **Recipe Knowledge Base \***
  - Create and manage recipes with flexible ingredient lists.
  - Attach foods from the registry as ingredients.

- **Role Management \***
  - Basic role support (`user`, `admin`, `limited`) for future access control expansions.

\* *Migrations and queries implemented, endpoints soon*


---

## Technical Highlights

### Security

- **JWT Authentication**: Stateless, signed tokens with configurable expiry.
- **Refresh Token Rotation**: Long-lived, database-backed tokens with fingerprinting (user agent & IP).
- **Rate Limiting (Planned)**: Middleware-ready structure for future rate limiting.
- **Role-based Access (Planned)**: Role model implemented for future granular permissions.

### Data Integrity & Validation

- **Centralized Validation**: All user input is validated.
- **SQLC Typed Queries**: All database access is generated and type-checked using [sqlc](https://sqlc.dev/), reducing runtime query errors.
- **Migration-based Schema**: Database schema is managed via [goose](https://github.com/pressly/goose) migrations.

### Observability

- **Structured Logging**: All requests and errors are logged with [zap](https://github.com/uber-go/zap).
- **Composable Middleware**: Security, logging, and cache headers are applied via middleware.

### Extensibility

- **Kafka Integration**: Events like food updates and user signups can be published to Kafka for downstream consumers.
- **Modular Design**: Clean package boundaries for authentication, database, middleware, routing, and business logic.

### Performance

- **Connection Pooling**: Uses pgxpool for efficient PostgreSQL access.
- **Optimized Queries**: Generated type-safe compile-ready Go code instead of ORM.

---

## Project Structure

```
cmd/api/                 # Main entrypoint (API server)
internal/app/            # Core application bootstrapping
internal/repository/     # SQLC-generated DB access (queries, models)
internal/migrations/     # Database migration files
internal/routes/         # HTTP routing and handlers
internal/services/       # Business logic modules (auth, etc)
internal/middleware/     # HTTP middleware (security, logging)
internal/env/            # Environment/config parsing
util/                    # Supporting utilities (validation, net, parsing, etc)
auth/tokenauthority/     # JWT and refresh token issuing & validation service
docs/                    # OpenAPI schema & documentation
```

---

## API Overview

PlateOps exposes a RESTful API. The core endpoints include:

- `POST /signup` — Register a new user.
- `POST /login` — Authenticate and receive JWT/refresh tokens.
- `POST /refresh` — Exchange a refresh token for a new JWT.
- `POST /logout` — Invalidate a refresh token.
- `POST /food` — Create a new food entry (JWT required).
- `GET /food/id/{id}` — Retrieve a food entry by internal ID.
- `GET /food/gtin/{gtin}` — Retrieve a food entry by GTIN.
- `PATCH /food/id/{id}` — Update an existing food entry (JWT required).
- `GET /foods/name/{name}` — Fuzzy search foods by name.
- `GET /health` — Health check endpoint.
- `GET /openapi.json` — OpenAPI schema.
- `GET /docs` — Swagger UI.

See the [OpenAPI spec](docs/openapi.json) for request and response details.

---

## Getting Started

### Prerequisites

- Go
- [goose](https://github.com/pressly/goose)
- [sqlc](https://sqlc.dev/) (optional, if extending/modifying/regenerating queries is needed)
- PostgreSQL
- Kafka/Redpanda (optional, for event streaming)
- Docker (optional, for containerization/deployment)

### Environment Variables

See the `dotenv.example` file for required configuration (documentation in the github wiki pages planned).
Key variables include:

- `ENVIRONMENT` — Sets logging and other functionalities to developer friendly or production configuration
- `DB_URL` — PostgreSQL connection string.
- `POSTGRES_*` — USER, PASSWORD, PORT, DB, basic PostgreSQL env variables.
- `JWT_SECRET` — Secret for signing JWTs (min 16 chars, ideally 32+).
- `SERVICE_NAME` — Used as JWT issuer and Kafka client ID.
- `SERVER_PORT` — Port to bind the API server.
- `STATIC_DIR` — Directory for static file (i.e favicon) serving (default: `static/`).
- `KAFKA_BROKERS` — (optional) Comma-separated Kafka broker URLs.

### Migrations

Apply database migrations with [goose](https://github.com/pressly/goose):

*If you have set up the goose environment variables, otherwise define information in args.*
```sh
goose up
```


### Running

You can compile and run the Go code normally:

```sh
go run ./cmd/api
```

or you can build a ready for deployment docker image through the provided `Dockerfile`.

---

## Security Considerations

- All passwords are hashed using bcrypt before storage.
- Refresh tokens are invalidated upon logout, expiration, or server-side revocation.
- Refresh tokens are bound to user agent and IP.
- JWT lifetime should be set to something short (i.e. 3 minutes) for production environments. Control of JWT lifetime via environment variable coming soon.
- All endpoints validate input rigorously to prevent injection and malformed data.
- Database schema uses foreign keys and check constraints for referential integrity.

---

## Development Notes

- **Tests**: (To be expanded) Unit and integration tests are planned for all major components.
- **Error Handling**: All errors are logged; user-facing errors avoid broadcasting sensitive details.
- **Extensibility**: The codebase is designed for adding new food/inventory/recipe endpoints in code, synergistic microservices through pub/sub and api, or external integrations such as front-ends to the api.
- **Contributions**: Please open issues, discussions, or PRs to discuss features, bugs, or improvements.

---
