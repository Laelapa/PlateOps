#  PlateOps

[![CodeFactor](https://www.codefactor.io/repository/github/laelapa/plateops/badge)](https://www.codefactor.io/repository/github/laelapa/plateops) [![Go Report Card](https://goreportcard.com/badge/github.com/Laelapa/PlateOps)](https://goreportcard.com/report/github.com/Laelapa/PlateOps)

<img align="right" width="320" src="PlateOpsMascot.png">

PlateOps is a comprehensive diet management platform. Currently it provides: 
- a registry of foods with fuzzy search capability,
- RESTful API endpoints,
- a self-rolled authentication system that supports hashed passwords, stateless and short-lived JWTs, database-backed & fingerprinted refresh tokens with logout/revocation support,
- structured logging, 
- event publishing capability to Kafka topics, 
- type-safe database queries through sqlc instead of an ORM,
- validation of all user input. 

Built as an API service, PlateOps empowers developers to create food-focused applications on top of it.

In its final version it will enable users to manage their pantry/inventory, track their nutritional intake, and receive recipe suggestions. It will recommend recipes based on the users' food stocks, their macronutrient goals for the day, and ingredient expiration dates & open packaging, minimizing food spoilage and waste.

---
Whenever I come across anything interesting during the development process I write about it in my [blog](https://laelapa.dev/blog).
- https://laelapa.dev/blog/001-the-three-value-problem--a-journey-from-json-to-go-pointers-to-postgres-native-types : Here you can read about how I tried to overcome some patterns that emerged in my code when I tried to add partial updates via PATCH. 

## Features

- **API**
  - OpenAPI JSON and Swagger UI for API exploration.
  - Rigorous input validation on all endpoints guarding against injections and malformed data.

- **Self-Rolled User Registration & Authentication**
  - Secure signup and login endpoints with password hashing (bcrypt).
  - JWT-based stateless authentication for accessing protected endpoints.
  - Refresh token mechanism for secure session renewal.
  - Logout and token revocation support.
 
- **Food Registry**
  - Create, update, and search food entries.
  - Retrieve food by ID or GTIN (barcode).
  - Fuzzy search for foods by name.

- **Database**
  - All database access code hand-typed SQL wrapped in type-safe Go functions using [sqlc](https://sqlc.dev/) affording us the following benefits:
    - compile-time verification & type safety,
    - LSP support,
    - no runtime overhead,
    - total transparency on what query is being executed (because they are all hand-typed & then wrapped)
  - Migration-based Schema managed with [goose](https://github.com/pressly/goose).
  - Uses pgxpool connection pooling.

- **Observability**
  - Structured Logging of requests, errors, and application state with [uber/zap](https://github.com/uber-go/zap).


- **Structure & Extensibility**
  - Composable Middleware: Security, logging, and cache headers are applied via middleware.
  - Kafka Integration: Events like food updates can be published to Kafka for downstream consumers.
  - Modular Design: Clean package boundaries for authentication, database, middleware, and routing.

- **Inventory Management \***
  - Track owned food items, quantities, and expiration dates.
  - Track packaging opening to monitor freshness.

- **Recipe Knowledge Base \***
  - Create and manage recipes with flexible ingredient lists.
  - Attach foods from the registry as ingredients.

- **Role Management \***
  - Basic role support (`user`, `admin`, `limited`) for future access control expansions.
  

\* *Migrations and queries implemented, endpoints and business logic are in the oven*

## Project Structure

| Directory | Description |
| ---- | ---- |
|[cmd/api/](cmd/api/) | API server, main entrypoint |
|[internal/app/](internal/app/) | Core application bootstrapping |
|[internal/queries/](internal/queries/) | Hand-typed SQL |
|[internal/repository/](internal/repository/) | sqlc-generated DB access layer |
|[internal/migrations/](internal/migrations/) | Database migration files |
|[internal/routes/](internal/routes/) | HTTP routing and handlers |
|[internal/services/](internal/services/) | Business logic modules (auth etc) |
|[internal/middleware/](internal/middleware/) | HTTP middleware (security, logging) |
|[internal/env/](internal/env/) | Environment & config parsing |
|[util/](util/) | misc utilites (data parsing & validation, net, type converters, etc) |
|[auth/tokenauthority/](auth/tokenauthority/) | JWT and refresh token issuing & validation service |
|[docs/](docs/) | OpenAPI schema & documentation | 

---

## API Overview

PlateOps exposes a RESTful API. The core endpoints include:

- `POST /signup`: Register a new user.
- `POST /login`: Authenticate and receive JWT/refresh tokens.
- `POST /refresh`: Exchange a refresh token for a new JWT.
- `POST /logout`: Invalidate a refresh token.
- `POST /food`: Create a new food entry (JWT required).
- `GET /food/id/{id}`: Retrieve a food entry by internal ID.
- `GET /food/gtin/{gtin}`: Retrieve a food entry by GTIN.
- `PATCH /food/id/{id}`: Update an existing food entry (JWT required).
- `GET /foods/name/{name}`: Fuzzy search foods by name.
- `GET /health`: Health check endpoint.
- `GET /openapi.json`: OpenAPI schema.
- `GET /docs`: Swagger UI.

See the [OpenAPI spec](docs/openapi.json) for request and response details.


## Quick Start Setup for Local Development & Execution

### Prerequisites

- **Docker** and **Docker Compose**: **required**
- **Go**: *optional*, needed only if you want to run the API outside Docker
- **sqlc**: *optional*, only needed if modifying database queries
- **goose**: *optional*, only needed if running migrations manually outside Docker

### Running Locally

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Laelapa/PlateOps.git
   cd PlateOps
   ```

2. **Start the development environment:**
   ```bash
   docker compose up
   ```

That's it! The environment includes:
- PostgreSQL with automatic migrations
- Kafka with pre-created topics
- PlateOps API ready to use

Access Points:
- Access the API at: `http://localhost:8081`
- View API docs at: `http://localhost:8081/docs`
- Monitor Kafka at: `http://localhost:8080`

### Modifying the Project

**If you're adding/modifying SQL queries:**
```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Edit queries in internal/queries/
# Then regenerate Go code:
sqlc generate
```

**If you're creating new migrations:**
```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Create a new migration:
goose -dir ./internal/migrations create my_migration sql
```
## For Online Deployment

If you are planning to deploy the application online via the Dockerfile (remember: the Compose file is suggested for local deployment only) consider the following:

### Environment Variables

If you are planning to deploy it online you should set environment variables.
See the `dotenv.example` file for required configuration (documentation in the github wiki pages is a work-in-progress).
Key variables include:

- `ENVIRONMENT` — Sets logging and other functionalities to developer friendly or production configuration
- `DB_URL` — PostgreSQL connection string.
- `POSTGRES_*` — USER, PASSWORD, PORT, DB, basic PostgreSQL env variables.
- `JWT_SECRET` — Secret for signing JWTs (min 16 chars, ideally 32+).
- `SERVICE_NAME` — Used as JWT issuer and Kafka client ID.
- `SERVER_PORT` — Port to bind the API server.
- `STATIC_DIR` — Directory for static file (i.e favicon) serving (default: `static/`).
- `KAFKA_BROKERS` — (optional) Comma-separated Kafka broker URLs.

## Security Considerations

- All passwords are hashed using bcrypt before storage.
- Refresh tokens are invalidated upon logout, expiration, or server-side revocation.
- Refresh tokens are bound to user agent and IP.
- JWT lifetime should be set to something short (i.e. 3 minutes) for production environments. Control of JWT lifetime via environment variable coming soon.
- All endpoints validate user input.
- Database schema uses foreign keys and check constraints for referential integrity.

>[!IMPORTANT]
>The project is currently configured for deployment on [fly.io](https://fly.io) and utilizes platform-specific headers like [`Fly-Client-IP`](https://fly.io/docs/networking/request-headers/#fly-client-ip). For deployments elsewhere this needs to be taken into consideration, as at best they won't be set or they could be spoofed by users when they aren't controlled by the fly environment.

---

## Other Small Notes

- **Error Handling**: All errors are logged; user-facing errors don't broadcast sensitive details.
- **Extensibility**: The codebase is designed for adding new food/inventory/recipe endpoints, synergistic microservices via pub/sub and API, or attaching external integrations such as front-ends to the API.
- **Contributions**: Please open issues, discussions, or PRs to discuss features, bugs, or improvements.

---

*The gopher mascot is a derivative work based on the Go gopher by Renee French (https://reneefrench.blogspot.com/) which carries the Creative Commons 4.0 Attributions license.*
