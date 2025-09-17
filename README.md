# URLSecure

URLSecure is a high-performance, secure URL shortening service built with Go, MySQL, and Redis. It offers user authentication, customizable short URLs, QR code generation, and a responsive user interface for seamless access and management.

***

## Current Status

The URLSecure service is fully deployed and accessible. Users can register, log in, create, manage, and use shortened URLs with real-time analytics and QR code generation. Generated URLs properly redirect to their original destinations.

***

## Features

- JWT-based user authentication allowing signup and login using username or email.  
- Simple and efficient URL shortening with unique short codes.  
- QR Code generation for every shortened URL.  
- Secure design with encrypted data storage and protection against common web threats.  
- Mobile-first responsive UI styled with Tailwind CSS.

***

## Technology Stack

**Backend**  
- Go programming language with Gin web framework for RESTful APIs.  
- MySQL relational database for persistent storage.  
- Redis for caching and rate limiting to improve performance.  
- JSON Web Tokens (JWT) for authentication.  
- Docker and Docker Compose for containerized development and deployment.

**Frontend**  
- Modern HTML5, CSS3, and ES6+ JavaScript.  
- Tailwind CSS for rapid UI development.  
- QRCode.js for client-side QR code rendering.

***

## Prerequisites

- Docker and Docker Compose (recommended for easy setup).  
- Go 1.21+ (required for development and manual run).  
- MySQL 8.0+ and Redis 7+ databases.

***

## Installation and Setup

### Clone Repository

```bash
git clone https://github.com/ConstantineCTF/URLSecure
cd URLSecure
```

### Configure Environment Variables

Create a `.env` file inside the `backend/` directory with the following variables, updated with your own settings:

```ini
APP_ENV=development
HTTP_PORT=8080
DB_HOST=mysql
DB_PORT=3306
DB_USER=shortener
DB_PASS=YourDatabasePassword
DB_NAME=shortener
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=YourJWTSecretKey
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

### Start Infrastructure Services

```bash
cd infra
docker-compose up -d mysql redis
```

### Build and Run Backend

- Using Docker Compose:

```bash
docker-compose build backend
docker-compose up -d backend
```

- Or running manually for development:

```bash
cd backend
go run cmd/shortener/main.go
```

### Access the Application

Open your browser to:

```
http://localhost:8080
```

***

## Project Structure

```
/backend         Backend Go source and migrations
/backend/public  Frontend static assets (HTML, CSS, JS)
/infra           Docker Compose configs and infrastructure setup
```

***

## Usage

- Register and log in to create and manage your short URLs.  
- Generate QR codes for easy offline sharing.

***

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests. Please adhere to Go coding conventions and Tailwind CSS best practices.

***

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
