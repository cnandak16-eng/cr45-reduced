# Dockerized Deployment Documentation – cr45-reduced

## Overview

This project is a containerized full-stack timetable management application consisting of:

* React frontend
* Go backend
* PostgreSQL database

The application allows all users to view timetable information while authenticated administrators can submit timetable overrides.

Docker and Docker Compose were used to package, deploy, and run the complete application stack.

---

# Deployment Approach

The deployment was completed in multiple stages:

1. Run the application locally to understand dependencies and behavior.
2. Create a Dockerfile for the Go backend.
3. Create a Dockerfile for the React frontend.
4. Configure Nginx as a reverse proxy for frontend-to-backend communication.
5. Create a Docker Compose configuration to run all services together.
6. Validate application functionality and persistence.

---

# Service Architecture

Frontend (React + Nginx)
|
| HTTP Requests
v
Backend (Go API)
|
| SQL Queries
v
PostgreSQL Database

---

# Dockerfiles

## Backend Dockerfile

The backend Dockerfile uses a multi-stage build.

### Build Stage

* Uses golang:1.26-alpine
* Downloads Go dependencies
* Compiles the backend binary

### Runtime Stage

* Uses alpine:latest
* Copies the compiled binary
* Copies migration files
* Exposes port 8080

Benefits:

* Smaller final image size
* Faster deployment
* Cleaner separation between build and runtime environments

---

## Frontend Dockerfile

The frontend Dockerfile also uses a multi-stage build.

### Build Stage

* Uses node:22-alpine
* Installs dependencies
* Builds React production assets

### Runtime Stage

* Uses nginx:alpine
* Copies generated build files into Nginx
* Uses custom nginx.conf

Benefits:

* Small production image
* Static assets served efficiently
* Reverse proxy support

---

# Docker Compose Configuration

The docker-compose.yml file defines three services:

## postgres

Database service.

Environment Variables:

* POSTGRES_USER=postgres
* POSTGRES_PASSWORD=postgres
* POSTGRES_DB=cr45_reduced

Port Mapping:

5432:5432

Volume:

postgres_data:/var/lib/postgresql/data

This ensures database persistence across container restarts.

---

## backend

Go API service.

Environment Variable:

DATABASE_URL=postgres://postgres:postgres@postgres:5432/cr45_reduced?sslmode=disable

Port Mapping:

8080:8080

The backend automatically executes database migrations during startup.

---

## frontend

React frontend served through Nginx.

Port Mapping:

3000:80

Users access the application through:

http://localhost:3000

---

# Frontend to Backend Communication

The frontend communicates with the backend through Nginx reverse proxying.

Nginx configuration:

location /api/ {
proxy_pass http://backend:8080;
}

This allows browser requests to use:

/api/...

instead of directly calling the backend service.

Advantages:

* Single entry point
* Simplified frontend configuration
* Cleaner deployment architecture

---

# PostgreSQL Configuration

The PostgreSQL database runs in its own container.

Configuration:

User: postgres

Password: postgres

Database: cr45_reduced

Persistence is provided through a Docker volume:

postgres_data

This prevents data loss when containers are restarted.

---

# Database Migrations

The backend automatically executes SQL migrations on startup.

Migration Files:

* 001_init.sql
* 002_seed.sql

These migrations:

* Create required tables
* Seed initial timetable data
* Create the admin account

Default Admin Credentials:

Username: admin

Password: classrep123

---

# Build and Run Instructions

## Start the Application

From the project root:

docker compose up --build -d

---

## Verify Running Containers

docker compose ps

Expected services:

* cr45-postgres
* cr45-backend
* cr45-frontend

---

## Access Application

Frontend:

http://localhost:3000

Backend:

http://localhost:8080

Database:

localhost:5432

---

# Validation Performed

## Timetable Loading

Verified that timetable data loads successfully through the frontend.

Result: PASS

---

## Login Functionality

Verified login using seeded admin account.

Result: PASS

---

## Override Submission

Verified that authenticated users can create timetable overrides.

Result: PASS

---

## Reverse Proxy

Verified that frontend API requests are forwarded through Nginx to the backend.

Result: PASS

---

## Database Persistence

Created timetable overrides and restarted containers.

Data remained available after restart.

Result: PASS

---

# Misconfiguration Case

## Issue

The frontend failed to start correctly due to an incorrect backend hostname in nginx.conf.

Incorrect configuration:

proxy_pass http://cr45-backend:8080;

Docker Compose networking resolves services using service names rather than container names.

The correct configuration was:

proxy_pass http://backend:8080;

---

## Error Observed

Frontend logs:

2026/06/12 11:29:16 [emerg] 1#1: host not found in upstream "cr45-backend" in /etc/nginx/conf.d/default.conf:12

nginx: [emerg] host not found in upstream "cr45-backend" in /etc/nginx/conf.d/default.conf:12

Backend logs:

db ping failed: failed to connect to user=postgres database=cr45_reduced

connect: connection refused

---

## Root Cause

* Nginx could not resolve the backend hostname.
* Backend startup occurred before PostgreSQL was fully ready.

---

## Resolution

1. Updated nginx.conf to use:

proxy_pass http://backend:8080;

2. Restarted the Docker Compose stack.

3. Verified successful communication between all services.

Result: RESOLVED

---

# Common Debugging Commands

View running containers:

docker compose ps

View backend logs:

docker compose logs backend

View frontend logs:

docker compose logs frontend

View database logs:

docker compose logs postgres

Restart services:

docker compose restart

Rebuild all services:

docker compose up --build -d

Stop all services:

docker compose down

---

# Conclusion

The application was successfully containerized using Docker and Docker Compose. The final deployment consists of three interconnected services (frontend, backend, and PostgreSQL) with persistent storage, automatic migrations, reverse proxying through Nginx, and validated application functionality.
