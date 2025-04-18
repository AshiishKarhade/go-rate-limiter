# Go Rate Limiter in API Gateway with Redis and bcrypt

This is a simple API gateway implemented in Go that supports user registration, login, and rate limiting using Redis. Passwords are hashed using **bcrypt** for security. This project is intended for learning purposes and demonstrates the use of Redis and bcrypt in building a basic user authentication system.

## Features:
- **User Registration**: Allows users to register with a username and password.
- **User Login**: Allows users to log in by verifying their password against a bcrypt hash.
- **Rate Limiting**: Implements basic rate limiting for users using Redis (Token Bucket Algorithm).
- **Redis for Data Storage**: Stores user data and rate limiting information in Redis.

## Prerequisites

- Go (1.15+)
- Redis (installed and running locally)

## Setup and Installation

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/go-api-gateway-with-redis.git
cd go-api-gateway-with-redis
```

### 2. Install Dependencies

Install the required dependencies:
```bash
go mod tidy
```

### 3. Start Redis

Ensure you have Redis running locally. If Redis isn't installed yet, follow the instructions from Redis Installation Guide.

To start Redis locally, use the following command:
```bash
redis-server
```

### 4. Run the API Gateway

```bash
go run main.go
```

## API Endpoints

### 1. **Register User**
- **URL**: `/register`
- **Method**: `POST`
- **Description**: Registers a new user by accepting a `username` and `password`. The password is hashed using bcrypt before storing it in Redis.
- **Request Body**:
    - Content-Type: `application/json`
    - Example:

      ```json
      {
          "username": "john_doe",
          "password": "securepassword"
      }
      ```

- **Response**:
    - **Success (201 Created)**:
      ```json
      {
        "message": "User registered successfully"
      }
      ```
    - **Error (400 Bad Request)**:
        - Missing or invalid `username` or `password`.
      ```json
      {
        "error": "Invalid user data"
      }
      ```

### 2. **Login User**
- **URL**: `/login`
- **Method**: `POST`
- **Description**: Allows a user to log in by verifying their `username` and `password`. The password is compared against the hashed password stored in Redis.
- **Request Body**:
    - Content-Type: `application/json`
    - Example:

      ```json
      {
          "username": "john_doe",
          "password": "securepassword"
      }
      ```

- **Response**:
    - **Success (200 OK)**:
      ```json
      {
        "message": "User logged in successfully"
      }
      ```
    - **Error (401 Unauthorized)**:
        - Invalid username or incorrect password.
      ```json
      {
        "error": "Incorrect password"
      }
      ```
        - User not found.
      ```json
      {
        "error": "User not found"
      }
      ```

### 3. **Rate Limiting**
- **URL**: `/api/v1/token` (or any other protected endpoint you may add)
- **Method**: `GET`
- **Description**: A dummy endpoint that implements rate limiting for each user. Each user is allowed a maximum of 2 requests per minute.
- **Request**:
    - Requires the `User-ID` header to identify the user. You should send a valid `User-ID` as part of the header when making a request.
    - Example `curl` command to send the request:
      ```bash
      curl -X GET http://localhost:8080/api/v1/token -H "User-ID: john_doe"
      ```

- **Response**:
    - **Success (200 OK)**:
      ```json
      {
        "message": "Token generated successfully"
      }
      ```
    - **Error (429 Too Many Requests)**:
        - The user has exceeded the allowed request limit (2 requests per minute).
      ```json
      {
        "error": "Rate limit exceeded"
      }
      ```
        - Response headers will also include a `Retry-After` header indicating when the user can try again (usually in seconds).
          Example:
          ```
          Retry-After: 60
          ```
