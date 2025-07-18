# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is ARK Server Commander - a Docker-based ARK Survival Evolved server management system. The project consists of a Go backend (Gin + GORM + SQLite) that manages Docker containers and a Nuxt.js frontend for the web interface.

## Architecture

- **Backend**: Go API server in `/server` directory
  - Uses Gin web framework with GORM ORM
  - SQLite database for server configurations
  - Docker SDK for container management
  - JWT authentication with initialization flow
  - Swagger API documentation at `/swagger/index.html`

- **Frontend**: Nuxt.js SPA in `/ui` directory  
  - Vue 3 with TypeScript
  - Nuxt UI component library
  - Pinia for state management
  - i18n support (English/Chinese)
  - Static site generation for production

- **Docker Integration**: Each ARK server runs in isolated containers
  - Uses `tbro98/ase-server:latest` base image
  - Automatic volume management for game data
  - Container lifecycle management (start/stop/restart)

## Development Commands

### Frontend (ui/ directory)
- `pnpm install` - Install dependencies
- `pnpm dev` - Start development server (port 3000)
- `pnpm build` - Build for production
- `pnpm generate` - Generate static site
- `pnpm lint` - Run ESLint
- `pnpm lint:fix` - Fix ESLint issues

### Backend (server/ directory)
- `go mod tidy` - Clean up dependencies
- `go run main.go` - Start development server (port 8080)
- `go build -o ark-commander main.go` - Build binary
- `swag init` - Generate Swagger documentation

### Docker
- `docker-compose up -d` - Start the full application stack
- The system requires Docker socket access for container management

## Key Components

### Backend Structure
- `main.go` - Entry point with Docker environment checks
- `routes/routes.go` - API routing and CORS configuration
- `controllers/` - HTTP handlers for auth and server management
- `service/docker_manager/` - Docker container operations
- `models/` - Database models (User, Server, ServerArgs)
- `middleware/auth.go` - JWT authentication middleware

### Frontend Structure  
- `pages/` - Route components (login, servers, init)
- `components/servers/` - Server management UI components
- `stores/` - Pinia stores for auth and server state
- `utils/` - ARK server parameter definitions and config parsers

## Important Configuration

- Backend serves on port 8080, frontend dev on port 3000
- Frontend proxies `/api` requests to backend during development
- Production builds frontend as static files served by backend
- System requires privileged Docker access for container management
- Initial setup requires admin user creation via `/init` endpoint

## Docker Images

The system manages ARK server containers using:
- Base image: `tbro98/ase-server:latest` 
- Automatic image pulling and version checking
- Manual image update workflow (see ToDo.md)

## Authentication Flow

1. Check initialization status via `/api/auth/check-init`
2. If not initialized, create admin user via `/api/auth/init`  
3. Login via `/api/auth/login` to get JWT token
4. Include token in Authorization header for protected endpoints