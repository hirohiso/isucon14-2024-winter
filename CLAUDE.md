# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is ISURIDE, a ride-sharing application built for the ISUCON14 (2024 Winter) competition. ISUCON is a performance tuning contest where the goal is to optimize a web application while maintaining correct functionality.

## Tech Stack

- **Language**: Go 1.23
- **Web Framework**: Chi router (v5.1.0)
- **Database**: MySQL with sqlx
- **ID Generation**: ULID
- **Monitoring**: New Relic
- **Frontend**: React (compiled assets in public/assets)

## Build & Run Commands

```bash
# Build the application
cd webapp/go
go build -o isuride .

# Run the application (listens on port 8080)
./isuride

# Initialize database
cd webapp/sql
./init.sh

# Format code
cd webapp/go
go fmt ./...

# Check for compilation errors
cd webapp/go
go build -v ./...

# Run basic Go vet
cd webapp/go
go vet ./...
```

## Environment Variables

```bash
ISUCON_DB_HOST (default: "127.0.0.1")
ISUCON_DB_PORT (default: "3306")
ISUCON_DB_USER (default: "isucon")
ISUCON_DB_PASSWORD (default: "isucon")
ISUCON_DB_NAME (default: "isuride")
NEW_RELIC_APP_NAME
NEW_RELIC_LICENSE_KEY
```

## Architecture Overview

The application implements a three-actor ride-sharing system:

1. **Users** - Request rides through `/api/app/*` endpoints
2. **Chairs** - Autonomous vehicles that fulfill rides through `/api/chair/*` endpoints  
3. **Owners** - Manage fleets of chairs through `/api/owner/*` endpoints

### Key Components

- **app_handlers.go** - User-facing API (registration, ride booking, payments)
- **chair_handlers.go** - Chair operations (location updates, status changes)
- **owner_handlers.go** - Owner management (sales reports, chair registration)
- **internal_handlers.go** - System operations (ride matching algorithm)
- **middlewares.go** - Authentication via cookie-based sessions
- **models.go** - Data structures and database models
- **payment_gateway.go** - Payment processing integration

### Database Schema

Key tables:
- `users` - User accounts with invitation codes
- `chairs` - Vehicle information
- `chair_locations` - Real-time position tracking
- `rides` - Ride requests and assignments
- `ride_statuses` - State machine for ride status
- `coupons` - Discount system
- `payment_tokens` - Payment methods

Current indexes:
- `idx_chairs` on chairs(id)
- `idx_owner_id` on chairs(owner_id)
- `idx_chair_locations_chair_id` on chair_locations(chair_id)
- `idx_ride_statuses_ride_id` on ride_statuses(ride_id)

### Ride State Machine

```
MATCHING → ENROUTE → PICKUP → CARRYING → ARRIVED → COMPLETED
```

### Authentication Pattern

Three middleware functions check cookie-based sessions:
- `appAuthMiddleware` - Validates `app_session` cookie
- `chairAuthMiddleware` - Validates `chair_session` cookie
- `ownerAuthMiddleware` - Validates `owner_session` cookie

All validate against access tokens stored in the database.

### Performance Considerations

Common optimization targets in ISUCON:
1. **Database queries** - Look for N+1 queries, missing indexes, inefficient JOINs
2. **Ride matching** - Currently uses `ORDER BY RAND()` which is inefficient
3. **Caching** - No caching layer implemented
4. **Concurrent requests** - Check for race conditions in status updates
5. **Distance calculations** - Manhattan distance used for fare calculation

### Error Handling

Consistent pattern using `writeError()`:
```go
writeError(w http.ResponseWriter, statusCode int, err error)
```

Check for `sql.ErrNoRows` for not found cases.

### Common Development Tasks

When modifying handlers:
1. Check authentication middleware requirements
2. Use transactions for multi-table updates
3. Follow existing error handling patterns
4. Maintain ride state machine integrity

When adding indexes:
1. Analyze slow queries first
2. Consider write performance impact
3. Update webapp/sql/init.sh for initialization

When optimizing:
1. Profile first (New Relic is integrated)
2. Check for N+1 queries in loops
3. Consider caching frequently accessed data
4. Verify ride status consistency