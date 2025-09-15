# URLSecure

URLSecure is a fast, secure, and analytics-driven URL shortening service built with Go, MySQL, Redis, and modern web technologies. It provides user authentication, customizable short URLs, real-time analytics, and a responsive design for seamless user experience.

---

## Features

- User Authentication with JWT-based signup and login using username or email.
- URL Shortening with support for custom short codes.
- Real-time Analytics to track clicks, referrers, geographical locations, and device types.
- QR Code Generation for every shortened URL.
- Security features including HTTPS enforcement, encrypted data storage, and protection against common attacks.
- Responsive, mobile-first UI built with Tailwind CSS.
- Link Expiration allowing users to set expiry dates for temporary URLs.
- User Dashboard for personal link management and analytics.

---

## Technology Stack

### Backend

- Go programming language using Gin web framework.
- MySQL relational database.
- Redis for caching and rate limiting.
- JWT for secure authentication.
- Docker and Docker Compose for containerization and orchestration.

### Frontend

- HTML5, CSS3, and modern JavaScript (ES6+).
- Tailwind CSS for styling.
- QRCode.js library for QR code generation.

---

## Prerequisites

- Docker and Docker Compose installed.
- Go version 1.21 or later (for local development).
- MySQL 8.0 or later.
- Redis 7 or later.

---

## Installation and Setup

1. Clone the repository:

    ```
    git clone https://github.com/ConstantineCTF/URLSecure
    cd URLSecure
    ```

2. Configure environment variables:

    Copy the example environment file and update it with your configuration:

    ```
    cp backend/.env.example backend/.env
    # Edit backend/.env with your settings
    ```

3. Start the required services:

    ```
    cd infra
    docker-compose up -d mysql redis
    ```

4. Build and launch the backend service:

    - Using Docker Compose:

      ```
      docker-compose build backend
      docker-compose up -d backend
      ```

    - Or run locally for development:

      ```
      cd backend
      go run cmd/shortener/main.go
      ```

5. Open your browser and navigate to:

    ```
    http://localhost:8080
    ```

---

## Project Structure

```
/backend         Backend Go source code, configuration files, and database migrations
/backend/public  Static frontend assets including HTML, CSS, and JavaScript
/infra           Docker Compose configuration and infrastructure files
```

---

## Usage

- Register and log in to create and manage shortened URLs.
- Access detailed analytics about link usage.
- Generate QR codes for sharing URLs offline.

---

## Contributing

Contributions are welcome. Please submit issues and pull requests on GitHub. Follow standard Go coding practices and Tailwind CSS conventions.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
