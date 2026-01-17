# NimbleDB Demo Application

A full-stack demo application showcasing [NimbleDB](https://github.com/kelvinwambua/nimbledb) - a custom-built database system with a Go backend API and Svelte frontend.

## Architecture

This demo consists of three containerized services:

- **NimbleDB** - Custom database server running on port 7000
- **Backend API** - Go REST API on port 8080 that interfaces with NimbleDB
- **Frontend** - Svelte web application on port 5173

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed
- [Docker Compose](https://docs.docker.com/compose/install/) installed
- Git

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/kelvinwambua/nimbledb-demo.git
cd nimbledb-demo
```

### 2. Configure Environment Variables

Create a `.env` file in the `backend` directory:

```bash
cd backend
```

Create `backend/.env`:

```env
PORT=8080
APP_ENV=local
DB_ADDR=localhost:7000
JWT_SECRET=your-secret-key-here
FRONTEND_URL="http://localhost:5173"
```

### 3. Build and Run

From the root directory:

```bash
docker-compose up --build
```

This will:

- Build all three services
- Start NimbleDB on port 7000
- Start the backend API on port 8080
- Start the frontend on port 5173

### 4. Access the Application

Open your browser and navigate to:

```
http://localhost:5173
```

The backend API is available at:

```
http://localhost:8080
```

## Development

### Stop the Application

```bash
docker-compose down
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f nimbledb
```

### Rebuild After Changes

```bash
docker-compose up --build
```

### Clean Rebuild (removes volumes)

```bash
docker-compose down -v
docker-compose up --build
```

## Project Structure

```
nimbledb-demo/
├── docker-compose.yml          # Orchestrates all services
├── backend/
│   ├── Dockerfile              # NimbleDB database container
│   ├── Dockerfile.api          # Go backend API container
│   ├── .env                    # Backend environment variables
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # API entry point
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── Dockerfile              # Svelte frontend container
    ├── .dockerignore
    ├── svelte.config.js
    ├── package.json
    └── src/
```

## API Endpoints

The backend exposes the following endpoints (example):

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/posts` - Get all posts
- `POST /api/posts` - Create a post
- `GET /api/posts/:id` - Get specific post
- `PUT /api/posts/:id/edit` - Update post
- `DELETE /api/posts/:id` - Delete post

## Docker Services

### NimbleDB Service

- **Port**: 7000
- **Volume**: `nimbledb-data` (persists database data)
- **Purpose**: Custom database engine

### Backend Service

- **Port**: 8080
- **Dependencies**: NimbleDB
- **Purpose**: REST API for frontend communication

### Frontend Service

- **Port**: 5173 (host) -> 3000 (container)
- **Dependencies**: Backend
- **Purpose**: Web interface

## Links

- [NimbleDB Repository](https://github.com/kelvinwambua/nimbledb)
- [SvelteKit Documentation](https://kit.svelte.dev/)
- [Go Documentation](https://golang.org/doc/)
