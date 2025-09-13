# URLSecure

A secure, Go-based URL shortening service with Tailwind-powered frontend.

## Prerequisites

- Docker & Docker Compose  
- Go 1.21+  
- Node.js & NPM  

## Setup

1. Copy environment variables:
cp .env.example .env

2. Start services:
cd infra
docker-compose up -d --build


3. Install frontend dependencies:
cd frontend
npm install


4. Build Tailwind CSS:
npm run build:css


5. Open Adminer at http://localhost:8081 to verify MySQL.

## Project Structure

- **backend**: Go API service  
- **frontend**: Static pages with Tailwind CSS  
- **infra**: Docker Compose stack  

## Next Steps

- Define data models and migrations  
- Implement API handlers and middleware  
- Build frontend pages and integrate with API  