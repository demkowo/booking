# Booking Service

## Overview
The **Booking Service** is a Golang-based API that manages reservations and room availability for a booking system. It is built with `Gin Gonic` for the HTTP router and `PostgreSQL` for the database. The service allows users to add, update, and delete reservations, as well as check room availability.

## Features
- **Reservation Management**: Add, update, retrieve, and delete reservations.
- **Room Management**: Add, update, and retrieve room information.
- **Room Availability**: Check room availability for a given date range.
- **Soft Deletion**: Reservations are soft-deleted to preserve booking history.
- **Transaction Management**: Ensures data consistency in reservation operations.

## Directory Structure
```
booking/
|-- app/          # Initialize application componenst and routes
│-- models/       # Contains data models
│-- repositories/ # Data access layer (PostgreSQL implementation)
│   ├── postgres/
│   │   ├── repository.go # Booking repository implementation
│-- services/     # Business logic layer
│   ├── booking.go # Service layer for booking logic
│-- handlers/     # HTTP handlers for API endpoints
│-- main.go       # Service entry point
```

## API Endpoints
| Method | Endpoint | Description |
|--------|----------------------------------------|------------------------------|
| `POST` | `/api/v1/reservations/add` | Create a new reservation |
| `DELETE` | `/api/v1/reservations/:reservation_id` | Delete a reservation |
| `GET`  | `/api/v1/reservations/` | Retrieve all reservations |
| `GET`  | `/api/v1/reservations/find/:room_id` | Retrieve reservations for a specific room |
| `GET`  | `/api/v1/reservations/:reservation_id` | Get reservation by ID |
| `PUT`  | `/api/v1/reservations/:reservation_id` | Update reservation details |
| `POST` | `/api/v1/rooms/add` | Add a new room |
| `GET`  | `/api/v1/rooms/` | Retrieve all rooms |
| `POST` | `/api/v1/rooms/find-available` | Find available rooms for a date range |
| `GET`  | `/api/v1/rooms/:room_id` | Get room details by ID |
| `POST` | `/api/v1/rooms/:room_id/availability-check` | Check room availability by ID |
| `PUT`  | `/api/v1/rooms/:room_id` | Update room details |

## Database Schema
The service interacts with the following tables:

### `reservations`
```sql
CREATE TABLE reservations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    room_id UUID NOT NULL,
    status INT NOT NULL DEFAULT 0,
    created TIMESTAMPTZ DEFAULT NOW(),
    updated TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);
```

### `rooms`
```sql
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    created TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ NOT NULL
);
```

## Usage

### Create a Reservation
```sh
curl -X POST http://localhost:8080/api/v1/reservations/add -H "Content-Type: application/json" -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "start_date": "2025-02-15T12:00:00Z",
    "end_date": "2025-02-20T12:00:00Z",
    "room_id": "456e7890-b12c-34d5-e678-910111213141"
}'
```

### Get a Reservation by ID
```sh
curl -X GET http://localhost:8080/api/v1/reservations/{reservation_id}
```

### Check Room Availability
```sh
curl -X POST http://localhost:8080/api/v1/rooms/find-available -H "Content-Type: application/json" -d '{
    "start_date": "2025-02-15T12:00:00Z",
    "end_date": "2025-02-20T12:00:00Z"
}'
```

### Update a Reservation
```sh
curl -X PUT http://localhost:8080/api/v1/reservations/{reservation_id} -H "Content-Type: application/json" -d '{
    "start_date": "2025-03-01T12:00:00Z",
    "end_date": "2025-03-05T12:00:00Z",
    "status": 1
}'
```

## Transactions & Error Handling
- All **write operations** (`Add`, `Update`, `Delete`) use transactions to ensure atomicity.
- **Soft deletion** is implemented for reservations to prevent accidental data loss.
- Errors are handled gracefully, returning appropriate HTTP status codes.

## Development Setup
### Prerequisites
- Golang (>=1.18)
- PostgreSQL
- `Gin Gonic`

### Install Dependencies
```sh
go mod tidy
```

### Run Service
```sh
go run main.go
```
