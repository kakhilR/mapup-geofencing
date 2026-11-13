# MapUp Geofencing & Real-time Alert System

A full-stack geofencing and real-time vehicle alert system built with Go, PostgreSQL/PostGIS, WebSocket, React, and Leaflet.

## Tech Stack

### Backend
- Go
- PostgreSQL
- PostGIS
- `pgx/v5`
- Gorilla WebSocket
- REST APIs

### Frontend
- React
- Vite
- React Leaflet
- Leaflet
- WebSocket

### Infrastructure
- Docker
- Docker Compose
- Docker Hub
- Render
- Vercel

---

## Project Structure

```text
mapup-geofencing/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── websocket/
├── migrations/
├── frontend/
├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── go.sum
└── SETUP.md
```

---

# 1. Prerequisites

Install:

- Go 1.27+
- Node.js 18+
- npm
- Docker Desktop
- Git

For production deployment:

- Docker Hub account
- Render account
- Vercel account

---

# 2. Local Backend Setup

Clone the repository and enter the project:

```bash
git clone <GITHUB_REPOSITORY_URL>
cd mapup-geofencing
```

Install Go dependencies:

```bash
go mod download
```

---

# 3. Start PostgreSQL + PostGIS

The project uses PostgreSQL with the PostGIS extension.

Start the database:

```bash
docker compose up -d postgres
```

Verify the container:

```bash
docker ps
```

Connect to PostgreSQL:

```bash
docker exec -it mapup-postgres psql \
  -U mapup \
  -d mapup
```

Enable PostGIS:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

Verify:

```sql
SELECT PostGIS_Version();
```

---

# 4. Database Schema

The application requires the following tables:

- `geofences`
- `vehicles`
- `vehicle_locations`
- `vehicle_geofence_states`
- `violations`
- `alert_configs`

PostGIS is used for geofence geometry and point-in-polygon detection.

Geofences use:

```text
SRID 4326
```

and are stored as PostGIS `POLYGON` geometries.

---

# 5. Backend Environment Variables

The backend supports the following environment variables:

```env
DATABASE_URL=postgres://mapup:mapup_password@localhost:5432/mapup
FRONTEND_URL=http://localhost:5173
```

For local development, the backend falls back to the local PostgreSQL connection if `DATABASE_URL` is not provided.

For production, `DATABASE_URL` should point to the deployed PostgreSQL/PostGIS database.

---

# 6. Run Backend Locally

Start the backend:

```bash
go run ./cmd/server
```

The backend runs on:

```text
http://localhost:8080
```

---

# 7. Run Backend with Docker

Build the backend image:

```bash
docker compose build backend
```

Start the backend and database:

```bash
docker compose up
```

The backend is available at:

```text
http://localhost:8080
```

---

# 8. Frontend Setup

Enter the frontend directory:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Create:

```text
frontend/.env
```

For local development:

```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

Start the development server:

```bash
npm run dev
```

The frontend is available at:

```text
http://localhost:5173
```

---

# 9. Production Frontend Environment Variables

For the deployed frontend, configure:

```env
VITE_API_URL=https://mapup-geofencing-backend-latest-gsy9.onrender.com
VITE_WS_URL=wss://mapup-geofencing-backend-latest-gsy9.onrender.com
```

The WebSocket URL uses `wss://` because the deployed backend is accessed over HTTPS.

---

# 10. API Endpoints

## Geofences

### Create Geofence

```http
POST /geofences
```

Example:

```json
{
  "name": "Downtown Delivery Zone",
  "description": "Main delivery area for downtown customers",
  "coordinates": [
    [37.7749, -122.4194],
    [37.7849, -122.4194],
    [37.7849, -122.4094],
    [37.7749, -122.4094],
    [37.7749, -122.4194]
  ],
  "category": "delivery_zone"
}
```

Coordinates are provided as:

```text
[latitude, longitude]
```

---

### List Geofences

```http
GET /geofences
```

---

## Vehicles

### Create Vehicle

```http
POST /vehicles
```

Example:

```json
{
  "vehicle_number": "MAPUP-001",
  "driver_name": "Demo Driver",
  "vehicle_type": "Delivery Van",
  "phone": "9999999999"
}
```

### List Vehicles

```http
GET /vehicles
```

---

## Vehicle Location

### Update Vehicle Location

```http
POST /vehicles/location
```

Example:

```json
{
  "vehicle_id": "<VEHICLE_ID>",
  "latitude": 37.7799,
  "longitude": -122.4144
}
```

### Get Latest Location

```http
GET /vehicles/location/{vehicle_id}
```

---

## Alert Configuration

### Configure Alert

```http
POST /alerts/configure
```

Example:

```json
{
  "vehicle_id": "<VEHICLE_ID>",
  "geofence_id": "<GEOFENCE_ID>",
  "event_type": "ENTRY",
  "enabled": true
}
```

Supported event types:

```text
ENTRY
EXIT
```

### List Alert Configurations

```http
GET /alerts
```

---

## Violations

### List Violation History

```http
GET /violations/history
```

---

# 11. WebSocket

The real-time alert WebSocket endpoint is:

```text
/ws/alerts
```

Local:

```text
ws://localhost:8080/ws/alerts
```

Production:

```text
wss://mapup-geofencing-backend-latest-gsy9.onrender.com/
```

When a configured geofence transition occurs, the backend broadcasts the violation through the WebSocket.

Example alert:

```json
{
  "id": "<VIOLATION_ID>",
  "vehicle_id": "<VEHICLE_ID>",
  "geofence_id": "<GEOFENCE_ID>",
  "event_type": "ENTRY",
  "message": "Vehicle MAPUP-001 ENTRY geofence Downtown Delivery Zone",
  "created_at": "<TIMESTAMP>"
}
```

---

# 12. Geofence Detection Flow

Vehicle location updates follow this flow:

```text
Vehicle Location
       |
       v
Save Location
       |
       v
PostGIS Point-in-Polygon
       |
       v
Find Current Geofences
       |
       v
Compare Previous State
       |
       +---- ENTRY
       |
       +---- EXIT
       |
       v
Check Alert Configuration
       |
       v
Create Violation
       |
       v
Broadcast via WebSocket
       |
       v
Frontend Live Alerts
```

The application maintains the previous vehicle/geofence state in:

```text
vehicle_geofence_states
```

This allows the system to distinguish between:

- vehicle entering a geofence
- vehicle remaining inside a geofence
- vehicle exiting a geofence

---

# 13. API Response Timing

All API responses include:

```json
{
  "time_ns": 1234567
}
```

`time_ns` represents the API processing time measured using Go's `time` package.

---

# 14. Docker Hub

The backend Docker image is published at:

```text
kar1998/mapup-geofencing-backend:latest
```

The image is built for:

```text
linux/amd64
```

This allows deployment on the target cloud environment.

---

# 15. Production Deployment

## Backend

The backend Docker image is deployed using Render.

Required environment variables:

```env
DATABASE_URL=<>
FRONTEND_URL=https://mapu-plum.vercel.app
```

The Go server listens on:

```text
8080
```

## Database

Production PostgreSQL is hosted on Render with PostGIS enabled.

Required extension:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

## Frontend

The React/Vite frontend is deployed using Vercel.

Production environment variables:

```env
VITE_API_URL=https://mapup-geofencing-backend-latest-gsy9.onrender.com
VITE_WS_URL=wss://mapup-geofencing-backend-latest-gsy9.onrender.com
```

---

# 16. Demo Flow

To demonstrate the application:

1. Open the deployed frontend.
2. Verify the WebSocket status shows `Connected`.
3. Verify the geofence is displayed on the map.
4. Verify the vehicle marker is displayed.
5. Configure an `ENTRY` alert for a vehicle/geofence.
6. Send a vehicle location outside the geofence.
7. Send another location inside the geofence.
8. The backend detects the `ENTRY` transition.
9. A violation is created.
10. The violation is broadcast through WebSocket.
11. The frontend displays the real-time alert.

Example outside location:

```json
{
  "vehicle_id": "<VEHICLE_ID>",
  "latitude": 37.7900,
  "longitude": -122.4144
}
```

Example inside location:

```json
{
  "vehicle_id": "<VEHICLE_ID>",
  "latitude": 37.7799,
  "longitude": -122.4144
}
```

---

# 17. Production URLs

Frontend:

```text
https://mapu-plum.vercel.app
```

Backend:

```text
https://mapup-geofencing-backend-latest-gsy9.onrender.com
```

Docker Hub:

```text
kar1998/mapup-geofencing-backend:latest
```