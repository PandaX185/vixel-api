# Vixel API

Vixel API is a Go-based REST API for image management and processing. It allows users to register, authenticate, upload images, manage their image collections, and apply transformations to images.

## Features

- User registration and authentication with JWT tokens
- Image upload and storage using MinIO object storage
- Image retrieval and management
- Image transformation capabilities
- PostgreSQL database for data persistence

## Technologies Used

- **Go** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Database
- **MinIO** - Object storage
- **JWT** - Authentication
- **GORM** - ORM for database operations

## API Endpoints

### Authentication

- `POST /api/v1/users` - Register a new user
- `POST /api/v1/users/login` - Login user

### Images

- `POST /api/v1/images` - Upload an image (requires authentication)
- `GET /api/v1/images/:id` - Get image details (requires authentication)
- `GET /api/v1/users/:user_id/images` - List user's images (requires authentication)
- `DELETE /api/v1/images/:id` - Delete an image (requires authentication)

### Processing

- `POST /api/v1/images/:id/transform` - Transform an image

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose
- PostgreSQL (or use Docker)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd vixel-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables. Create a `.env` file in the root directory with the following variables:
   ```
   PORT=:8080
   DATABASE_URL=postgres://user:password@localhost:5432/vixel?sslmode=disable
   JWT_SECRET=your-jwt-secret
   MINIO_ENDPOINT=localhost:9000
   MINIO_ACCESS_KEY=minioadmin
   MINIO_SECRET_KEY=minioadmin
   MINIO_BUCKET=vixel-bucket
   ```

4. Start MinIO using Docker Compose:
   ```bash
   docker-compose up -d
   ```

5. Run the application:
   ```bash
   go run app/main.go
   ```

The API will be available at `http://localhost:8080`.

## Usage

### Authentication

First, register a new user:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

Login to get a JWT token:

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

Use the returned token in subsequent requests by including it in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Uploading Images

Upload an image file:

```bash
curl -X POST http://localhost:8080/api/v1/images \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "file=@/path/to/your/image.jpg" \
  -F "alt_text=Description of the image"
```

### Transforming Images

Apply transformations to an uploaded image:

```bash
curl -X POST http://localhost:8080/api/v1/images/1/transform \
  -H "Content-Type: application/json" \
  -d '{"transformations": ["resize", "grayscale"]}'
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.</content>