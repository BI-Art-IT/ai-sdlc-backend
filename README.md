# AirLines Backend

Go backend for the AirLines platform, providing a REST API backed by PostgreSQL.

Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate) using the `pgx/v5` driver. Each table lives in its own numbered migration file under [`migrations/`](./migrations/).

---

## Table of Contents

1. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Configuration](#configuration)
   - [Running Migrations](#running-migrations)
   - [Starting the Server](#starting-the-server)
2. [Project Structure](#project-structure)
3. [Database Design](#database-design)
   - [Overview](#overview)
   - [Entity Relationship Diagram](#entity-relationship-diagram)
   - [Tables](#tables)
   - [Migration Files](#migration-files)

---

## Getting Started

### Prerequisites

| Requirement | Version |
|---|---|
| [Go](https://go.dev/dl/) | ≥ 1.25 |
| [PostgreSQL](https://www.postgresql.org/download/) | ≥ 14 |

### Configuration

Copy the example environment file and fill in your database credentials:

```bash
cp .env.example .env
```

Edit `.env`:

```env
# PostgreSQL connection URL (required)
DATABASE_URL=postgres://postgres:password@localhost:5432/airlines?sslmode=disable

# Path to migration files, relative to the working directory (optional, default: "migrations")
MIGRATIONS_PATH=migrations

# Port the HTTP server listens on (optional, default: "8080")
SERVER_PORT=8080
```

> The application loads `.env` automatically at startup via [godotenv](https://github.com/joho/godotenv).
> You can also export these variables directly into your shell instead of using a file.

### Running Migrations

Migrations can be triggered in two ways:

#### Option 1 — Standalone `migrate` command (recommended for CI/CD and manual runs)

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Roll back the last migration
go run ./cmd/migrate down 1

# Roll back all migrations (full teardown)
go run ./cmd/migrate down
```

Build the binary once and reuse it:

```bash
go build -o bin/migrate ./cmd/migrate

./bin/migrate up
./bin/migrate down 2
```

#### Option 2 — Automatic on server startup

Migrations also run automatically every time the server starts (see `cmd/server/main.go`).
This is convenient for local development; for production, prefer Option 1 so migrations are
decoupled from the application process.

### Starting the Server

```bash
go run ./cmd/server
```

Or build first:

```bash
go build -o bin/server ./cmd/server
./bin/server
```

The server exposes the following endpoint:

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Returns `{"status":"ok"}` when the server is running |

---

## Project Structure

```
.
├── cmd/
│   ├── migrate/        # Standalone migration CLI (go run ./cmd/migrate up|down)
│   │   └── main.go
│   └── server/         # HTTP server entrypoint
│       └── main.go
├── internal/
│   ├── config/         # Environment-based configuration
│   ├── database/       # pgx/v5 connection pool
│   └── migrate/        # golang-migrate Up/Down wrappers
├── migrations/         # SQL migration files (one table per file)
├── .env.example        # Environment variable template
├── go.mod
└── go.sum
```

---

## Database Design

### Overview

The schema is split into two categories:

| Category | Tables | Description |
|---|---|---|
| **Shared** | `airports`, `airlines`, `fare_classes`, `flights` | Reference and operational data shared across multiple platform features |
| **Flight Search–specific** | `fares`, `search_logs` | Pricing data and anonymised search telemetry used only by the Flight Search feature |

All primary keys are `UUID` generated with `gen_random_uuid()` (provided by the `pgcrypto` extension, enabled in migration `000001`).

---

### Entity Relationship Diagram

```
airports ──────────────────────────────────────────────────────┐
  (id)                                                          │
    ▲                                                           │
    │ origin_airport_id / destination_airport_id                │
    │                                                           │
airlines ──── flights ────────────────────────── fares         │
  (id)          (id)         flight_id ──────────► (id)        │
    ▲             │                                  │          │
    │             │                         fare_class_id       │
    └─ airline_id │                                  ▼          │
                  │                           fare_classes      │
                  │                              (id)           │
                  │                                             │
                  └─ destination_airport_id ────────────────────┘

search_logs  (no FK dependencies — standalone telemetry table)
```

---

### Tables

### airports

Stores reference data for all airports supported by the platform.

**Migration:** `000001_create_airports_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `iata_code` | `CHAR(3)` | NO | — | IATA 3-letter airport code (e.g. `LHR`). Unique. |
| `name` | `VARCHAR(255)` | NO | — | Full airport name |
| `city` | `VARCHAR(100)` | NO | — | City in which the airport is located |
| `country` | `VARCHAR(100)` | NO | — | Country name |
| `country_code` | `CHAR(2)` | NO | — | ISO 3166-1 alpha-2 country code |
| `timezone` | `VARCHAR(50)` | NO | — | IANA timezone identifier (e.g. `Europe/London`) |
| `latitude` | `DECIMAL(9,6)` | YES | — | Geographic latitude |
| `longitude` | `DECIMAL(9,6)` | YES | — | Geographic longitude |
| `active` | `BOOLEAN` | NO | `TRUE` | Whether the airport is available for search |
| `created_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record creation timestamp |
| `updated_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record last-updated timestamp |

**Constraints:** `UNIQUE (iata_code)`

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `airports_iata_code_key` | `iata_code` | Unique lookup / FK target |
| `airports_city_idx` | `city` | City-name autocomplete in Flight Search |
| `airports_country_code_idx` | `country_code` | Country-level filtering |

---

### airlines

Stores reference data for airlines operating flights on the platform.

**Migration:** `000002_create_airlines_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `iata_code` | `CHAR(2)` | NO | — | IATA 2-letter airline code (e.g. `BA`). Unique. |
| `name` | `VARCHAR(255)` | NO | — | Full airline name |
| `logo_url` | `VARCHAR(500)` | YES | — | URL to the airline logo image |
| `active` | `BOOLEAN` | NO | `TRUE` | Whether the airline is currently active |
| `created_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record creation timestamp |
| `updated_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record last-updated timestamp |

**Constraints:** `UNIQUE (iata_code)`

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `airlines_iata_code_key` | `iata_code` | Unique lookup / FK target |

---

### fare_classes

Lookup table for the cabin class types offered on the platform. Pre-seeded with the four standard classes.

**Migration:** `000003_create_fare_classes_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `code` | `VARCHAR(20)` | NO | — | Machine-readable class code. Unique. |
| `name` | `VARCHAR(50)` | NO | — | Human-readable display name |
| `sort_order` | `SMALLINT` | NO | — | Display order (ascending) |

**Constraints:** `UNIQUE (code)`

**Seed data:**

| `code` | `name` | `sort_order` |
|---|---|---|
| `economy` | Economy | 1 |
| `premium-economy` | Premium Economy | 2 |
| `business` | Business | 3 |
| `first` | First Class | 4 |

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `fare_classes_code_key` | `code` | Unique lookup / FK target |

---

### flights

Stores individual scheduled flight instances with route, timing, and operational status.

**Migration:** `000004_create_flights_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `flight_number` | `VARCHAR(10)` | NO | — | Flight number including airline prefix (e.g. `BA117`) |
| `airline_id` | `UUID` | NO | — | FK → `airlines.id` |
| `origin_airport_id` | `UUID` | NO | — | FK → `airports.id` (departure airport) |
| `destination_airport_id` | `UUID` | NO | — | FK → `airports.id` (arrival airport) |
| `departure_time` | `TIMESTAMPTZ` | NO | — | Scheduled departure time (UTC) |
| `arrival_time` | `TIMESTAMPTZ` | NO | — | Scheduled arrival time (UTC) |
| `duration_minutes` | `INTEGER` | NO | — | Total scheduled flight duration in minutes. Must be > 0. |
| `stops` | `SMALLINT` | NO | `0` | Number of intermediate stops. Must be >= 0. |
| `aircraft_type` | `VARCHAR(50)` | YES | — | Aircraft model (e.g. `Boeing 777-300ER`) |
| `total_seat_capacity` | `INTEGER` | NO | — | Total seats on this flight. Must be > 0. |
| `status` | `VARCHAR(20)` | NO | `'scheduled'` | Operational status (see allowed values below) |
| `estimated_departure_time` | `TIMESTAMPTZ` | YES | — | Updated EDT when a delay is active |
| `estimated_arrival_time` | `TIMESTAMPTZ` | YES | — | Updated ETA when a delay is active |
| `departure_gate` | `VARCHAR(10)` | YES | — | Assigned departure gate (e.g. `B12`) |
| `departure_terminal` | `VARCHAR(10)` | YES | — | Departure terminal (e.g. `T3`) |
| `arrival_gate` | `VARCHAR(10)` | YES | — | Assigned arrival gate |
| `arrival_terminal` | `VARCHAR(10)` | YES | — | Arrival terminal |
| `created_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record creation timestamp |
| `updated_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record last-updated timestamp |

**Allowed `status` values:** `scheduled`, `boarding`, `departed`, `en_route`, `landed`, `arrived`, `delayed`, `cancelled`

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `flights_route_departure_idx` | `(origin_airport_id, destination_airport_id, departure_time)` | Primary search query index for Flight Search |
| `flights_airline_id_idx` | `airline_id` | Airline-level filtering |
| `flights_flight_number_idx` | `flight_number` | Flight number lookup |
| `flights_departure_time_idx` | `departure_time` | Date-range queries |
| `flights_status_idx` | `status` | Flight Status queries |

---

### fares

Stores pricing information for each flight per cabin class, including taxes and fees.

**Migration:** `000005_create_fares_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `flight_id` | `UUID` | NO | — | FK → `flights.id` |
| `fare_class_id` | `UUID` | NO | — | FK → `fare_classes.id` |
| `base_fare` | `DECIMAL(10,2)` | NO | — | Base fare amount excluding taxes and fees |
| `taxes` | `DECIMAL(10,2)` | NO | — | Tax amount |
| `fees` | `DECIMAL(10,2)` | NO | — | Mandatory fee amount |
| `total_price` | `DECIMAL(10,2)` | NO | — | Computed total (`base_fare + taxes + fees`) |
| `currency` | `CHAR(3)` | NO | — | ISO 4217 currency code (e.g. `GBP`) |
| `available_seats` | `SMALLINT` | NO | — | Seats available at this fare |
| `valid_from` | `TIMESTAMPTZ` | NO | — | Start of fare validity window |
| `valid_until` | `TIMESTAMPTZ` | NO | — | End of fare validity window |
| `created_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record creation timestamp |
| `updated_at` | `TIMESTAMPTZ` | NO | `NOW()` | Record last-updated timestamp |

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `fares_flight_fare_class_idx` | `(flight_id, fare_class_id)` | Primary join path |
| `fares_valid_until_idx` | `valid_until` | Expiry-based queries and cache invalidation |
| `fares_total_price_idx` | `total_price` | Price-ascending sort |

---

### search_logs

Stores anonymised flight search telemetry for analytics and performance monitoring. No personally identifiable information (PII) is stored directly in this table.

> **Retention Policy:** Records are automatically purged after **90 days** to comply with the platform's data minimisation policy.

**Migration:** `000006_create_search_logs_table`

| Column | Type | Nullable | Default | Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | `gen_random_uuid()` | Primary key |
| `user_id_hash` | `CHAR(64)` | YES | — | SHA-256 hash of the authenticated user's ID. `NULL` for anonymous requests. |
| `session_id` | `VARCHAR(64)` | YES | — | Anonymous session identifier |
| `origin_code` | `VARCHAR(3)` | NO | — | Origin airport IATA code |
| `destination_code` | `VARCHAR(3)` | NO | — | Destination airport IATA code |
| `departure_date` | `DATE` | NO | — | Requested departure date |
| `return_date` | `DATE` | YES | — | Requested return date (`NULL` for one-way) |
| `trip_type` | `VARCHAR(20)` | NO | — | `one-way`, `round-trip`, or `multi-city` |
| `adults` | `SMALLINT` | NO | — | Number of adult passengers |
| `children` | `SMALLINT` | NO | — | Number of child passengers |
| `infants` | `SMALLINT` | NO | — | Number of infant passengers |
| `cabin_class` | `VARCHAR(20)` | NO | — | Requested cabin class |
| `flexible_dates` | `BOOLEAN` | NO | `FALSE` | Whether flexible date search was enabled |
| `results_returned` | `INTEGER` | YES | — | Number of results in the response |
| `cache_hit` | `BOOLEAN` | YES | — | Whether the response was served from cache |
| `response_time_ms` | `INTEGER` | YES | — | API response time in milliseconds |
| `created_at` | `TIMESTAMPTZ` | NO | `NOW()` | Log entry timestamp |

**Indexes:**

| Index | Columns | Purpose |
|---|---|---|
| `search_logs_created_at_idx` | `created_at` | Time-range analytics queries |
| `search_logs_route_idx` | `(origin_code, destination_code)` | Route-level analytics |

---

### Migration Files

| # | File | Direction | Description |
|---|---|---|---|
| 000001 | `000001_create_airports_table.up.sql` | up | Creates `airports` table and enables `pgcrypto` |
| 000001 | `000001_create_airports_table.down.sql` | down | Drops `airports` table |
| 000002 | `000002_create_airlines_table.up.sql` | up | Creates `airlines` table |
| 000002 | `000002_create_airlines_table.down.sql` | down | Drops `airlines` table |
| 000003 | `000003_create_fare_classes_table.up.sql` | up | Creates and seeds `fare_classes` table |
| 000003 | `000003_create_fare_classes_table.down.sql` | down | Drops `fare_classes` table |
| 000004 | `000004_create_flights_table.up.sql` | up | Creates `flights` table with FKs and indexes |
| 000004 | `000004_create_flights_table.down.sql` | down | Drops `flights` table |
| 000005 | `000005_create_fares_table.up.sql` | up | Creates `fares` table with FKs and indexes |
| 000005 | `000005_create_fares_table.down.sql` | down | Drops `fares` table |
| 000006 | `000006_create_search_logs_table.up.sql` | up | Creates `search_logs` table |
| 000006 | `000006_create_search_logs_table.down.sql` | down | Drops `search_logs` table |

---

> **Running migrations and starting the application?** See [Getting Started](#getting-started) above.
