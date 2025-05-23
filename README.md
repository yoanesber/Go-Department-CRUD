# Department API Service – Secure RESTful API with JWT Auth, Redis Caching & Rate Limiting

A secure, scalable, and modular RESTful API service built with **Go**, using the **Gin** web framework and **GORM** ORM. This service manages CRUD operations for the `Department` entity with full **JWT-based authentication** and **Redis caching**. It integrates with **PostgreSQL** as the main database and **Redis** for session/token storage.

---


## ✨ Features

This application provides a **secure**, **token-based authentication system using JWT (JSON Web Tokens)**, **fully integrated with Redis** for optimized token handling, and **PostgreSQL** for persistent storage. Below is a summary of the core features offered:

- Full **JWT authentication** system:
  - `POST /auth/login` — Accepts `username` and `password`, returns:
    - `AccessToken`
    - `RefreshToken`
    - `ExpirationDate`
    - `TokenType`
  - `POST /auth/refresh-token` — Accepts valid `RefreshToken` to generate new `AccessToken`.

- **Token storage in Redis** for faster access:
  - Stored under key format: `access_token:<username>`
  - JSON structure: `{ AccessToken, RefreshToken, ExpirationDate, TokenType }`

- **CRUD API for Department** entity:
  - All routes are protected by JWT Bearer Token via `Authorization` header.


### 🛡️ Security & Middleware

The service is designed with security and extensibility in mind, using several middlewares:

- **Authorization Middleware**:
  - Validates JWT
  - Enforces Role-Based Access Control (RBAC)

- **Context Injection Middleware**:
  - Injects database (PostgreSQL) and Redis connections into the Gin context for downstream handlers

- **Security Headers Middleware**:
  - CORS
  - Request ID
  - Secure HTTP headers (e.g., `X-Frame-Options`, `X-Content-Type-Options`, etc.)

- **Rate Limiter**:
  - Built on `golang.org/x/time/rate`
  - Rate limits based on unique key: `IP + HTTP method + route path`


### 🗄️ Logging

- Uses `github.com/sirupsen/logrus` for structured logging
- Integrates with `gopkg.in/natefinch/lumberjack.v2` for automatic log rotation based on size and age
- Logs are separated by level: **info**, **request**, **warn**, **error**, **fatal**, and **panic**


### 🔐 JWT Key Management

- **RSA key pairs** are used for signing JWTs (instead of symmetric secrets)
- Keys are generated using OpenSSL:
  - `privateKey.pem`, `publicKey.pem` in `/keys`

---

## 🤖 Tech Stack

This project leverages a modern and robust set of technologies to ensure performance, security, and maintainability. Below is an overview of the core tools and libraries used in the development:

| **Component**             | **Description**                                                                         |
|---------------------------|-----------------------------------------------------------------------------------------|
| **Language**              | Go (Golang), a statically typed, compiled language known for concurrency and efficiency |
| **Web Framework**         | Gin, a fast and minimalist HTTP web framework for Go                                    |
| **ORM**                   | GORM, an ORM library for Go supporting SQL and migrations                               |
| **Database**              | PostgreSQL, a powerful open-source relational database system                           |
| **Cache/Session Store**   | Redis, an in-memory data structure store used for caching and session management        |
| **JWT Signing**           | RSA asymmetric keys generated with OpenSSL for secure token signing                     |
| **Logging**               | Logrus for structured logging, combined with Lumberjack for log rotation                |
| **Validation**            | `go-playground/validator.v9` for input validation and data integrity enforcement        |

---

## 🧱 Architecture Overview

This project follows a **modular** and **maintainable** architecture inspired by **Clean Architecture** principles. Each domain feature (e.g., **authentication**, **department management**, **user**, **role**) is organized into self-contained modules with clear separation of concerns.

```bash
📁 go-deparment-crud/
├── 📂cert/                                 # Stores self-signed TLS certificates used for local development (e.g., for HTTPS or JWT signing verification)
├── 📂cmd/                                  # Contains the application's entry point.
├── 📂config/
│   └── 📂db/                               # Configuration packages for database connections
│       ├── 📂postgresdb/                   # PostgreSQL connection logic and DSN construction
│       └── 📂redisdb/                      # Redis connection configuration and initialization
├── 📂docker/                               # Docker-related configuration for building and running services
│   ├── 📂app/                              # Contains Dockerfile to build the main Go application image
│   ├── 📂postgres/                         # Contains PostgreSQL container configuration
│   └── 📂redis/                            # Contains Redis container configuration
├── 📂internal/                             # Core domain logic and business use cases, organized by module
│   ├── 📂auth/                             # Authentication logic (login, token generation)
│   ├── 📂dataredis/                        # Handles storing and retrieving data from redis
│   ├── 📂department/                       # Department module
│   ├── 📂refreshtoken/                     # Manages refresh token persistence and validation
│   ├── 📂role/                             # Role management for access control
│   └── 📂user/                             # User module (authentication identity source)
├── 📂keys/                                 # Contains RSA public/private keys used for signing and verifying JWT tokens
├── 📂logs/                                 # Application log files (error, request, info) written and rotated using Logrus + Lumberjack
├── 📂pkg/                                  # Reusable utility and middleware packages shared across modules
│   ├── 📂contextdata/
│   │   ├── 📂dbcontext/                    # Embeds PostgreSQL DB connection into context
│   │   └── 📂metacontext/                  # Provides inject dan extract function of the RequestMeta into/from the context
│   ├── 📂logger/                           # Centralized log initialization and configuration
│   ├── 📂middleware/                       # Request processing middleware
│   │   ├── 📂authorization/                # JWT validation and Role-Based Access Control (RBAC)
│   │   ├── 📂context/                      # Injects DB and Redis connections per request
│   │   ├── 📂headers/                      # Manages request headers like CORS, security, request ID
│   │   ├── 📂logging/                      # Logs incoming requests
│   │   └── 📂ratelimiter/                  # Implements API rate limiting based on IP, path, and method
│   ├── 📂util/                             # General utility functions and helpers
│   │   ├── 📂redisutil/                    # Wrapper utilities for working with Redis data types
│   └── 📂validator/                        # Custom request validation using go-playground/validator.v9
├── 📂routes/                               # Route definitions, groups APIs, and applies middleware per route scope
├── 📂tests/                                # Contains unit or integration tests for business logic
├── .dockerignore
├── .env
├── .gitignore
├── generate-certificate.sh                 # Script to generate self-signed certificates using OpenSSL
├── generate-jwt-key.sh                     # Script to generate RSA key pairs for JWT signing/verification
├── go.mod
├── go.sum
├── import.sql                              # Optional SQL script to import initial database records or mock data
└── Makefile                                # Provides CLI shortcuts to build, run, test, or manage Docker containers and environments
```

---

## 🛠️ Installation & Setup  

Follow the instructions below to get the project up and running in your local development environment. You may run it natively or via Docker depending on your preference.  

### ✅ Prerequisites

Make sure the following tools are installed on your system:

| **Tool**                                                      | **Description**                           |
|---------------------------------------------------------------|-------------------------------------------|
| [Go](https://go.dev/dl/)                                      | Go programming language (v1.20+)          |
| [Make](https://www.gnu.org/software/make/)                    | Build automation tool (`make`)            |
| [Redis](https://redis.io/)                                    | In-memory data store                      |
| [PostgreSQL](https://www.postgresql.org/)                     | Relational database system (v14+)         |
| [Docker](https://www.docker.com/)                             | Containerization platform (optional)      |
| [dotenv](https://github.com/motdotla/dotenv) (optional)       | To load `.env` files in local development |

### ⚙️ Configure `.env` File  

Set up your **database**, **Redis**, and **JWT configuration** in `.env` file. Create a `.env` file at the project root directory:  

```properties
# Application configuration
ENV=DEVELOPMENT
API_VERSION=1.0
PORT=1000
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASS=app@123
DB_NAME=department
DB_SSL=disable
DB_TIMEZONE=Asia/Jakarta
DB_MIGRATE=TRUE
DB_SEED=TRUE
DB_SEED_FILE=import.sql
# Set to INFO for development and staging, SILENT for production
DB_LOG=SILENT

# Redis configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USER=default
REDIS_PASS=your_redis_password
REDIS_DB=0
# 1 hour
ACCESS_TOKEN_TTL_MINUTES=60

# JWT configuration
JWT_SECRET=your_jwt_secret_key
# 2 days
JWT_EXPIRATION_HOUR=48
JWT_ISSUER=your_jwt_issuer
JWT_AUDIENCE=your_jwt_audience
# 30 days
JWT_REFRESH_TOKEN_EXPIRATION_HOUR=720
JWT_PRIVATE_KEY_PATH=./keys/privateKey.pem
JWT_PUBLIC_KEY_PATH=./keys/publicKey.pem
# RS256 or HS256
JWT_ALGORITHM=RS256
# Bearer or JWT
TOKEN_TYPE=Bearer
```

- **🔐 Notes**:  
  - `IS_SSL=TRUE`: Enable this if you want your app to run over `HTTPS`. Make sure to run `generate-certificate.sh` to generate **self-signed certificates** and place them in the `./cert/` directory (e.g., `mycert.key`, `mycert.cer`).
  - `JWT_ALGORITHM=RS256`: Set this if you're using **asymmetric JWT signing**. Be sure to run `generate-jwt-key.sh` to generate **RSA key pairs** and place `privateKey.pem` and `publicKey.pem` in the `./keys/` directory.
  - Make sure your paths (`./cert/`, `./keys/`) exist and are accessible by the application during runtime.
  - `DB_TIMEZONE=Asia/Jakarta`: Adjust this value to your local timezone (e.g., `America/New_York`, etc.).
  - `DB_MIGRATE=TRUE`: Set to `TRUE` to automatically run `GORM` migrations for all entity definitions on app startup.
  - `DB_SEED=TRUE` & `DB_SEED_FILE=import.sql`: Use these settings if you want to insert predefined data into the database using the SQL file provided.
  - `DB_USER=appuser`, `DB_PASS=app@123`: It's strongly recommended to create a dedicated database user instead of using the default postgres superuser.

### 🔑 Generate RSA Key for JWT (If Using `RS256`)  

If you are using `JWT_ALGORITHM=RS256`, generate the **RSA key** pair for **JWT signing** by running this file:  
```bash
./generate-jwt-key.sh
```

- **Notes**:  
  - On **Linux/macOS**: Run the script directly
  - On **Windows**: Use **WSL** to execute the `.sh` script

This will generate:
  - `./keys/privateKey.pem`
  - `./keys/publicKey.pem`


These files will be referenced by your `.env`:
```properties
JWT_PRIVATE_KEY_PATH=./keys/privateKey.pem
JWT_PUBLIC_KEY_PATH=./keys/publicKey.pem
JWT_ALGORITHM=RS256
```

### 🔐 Generate Certificate for HTTPS (Optional)  

If `IS_SSL=TRUE` in your `.env`, generate the certificate files by running this file:  
```bash
./generate-certificate.sh
```

- **Notes**:  
  - On **Linux/macOS**: Run the script directly
  - On **Windows**: Use **WSL** to execute the `.sh` script

This will generate:
  - `./cert/mycert.key`
  - `./cert/mycert.cer`


Ensure these are correctly referenced in your `.env`:
```properties
IS_SSL=TRUE
SSL_KEYS=./cert/mycert.key
SSL_CERT=./cert/mycert.cer
```

### 👤 Create Dedicated PostgreSQL User (Recommended)

For security reasons, it's recommended to avoid using the default postgres superuser. Use the following SQL script to create a dedicated user (`appuser`) and assign permissions:

```sql
-- Create appuser and set permissions
CREATE USER appuser WITH PASSWORD 'app@123';

GRANT CONNECT ON DATABASE department TO appuser;
GRANT CREATE ON SCHEMA public TO appuser;

GRANT USAGE ON SCHEMA public TO appuser;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appuser;
```

Update your `.env` accordingly:
```properties
DB_USER=appuser
DB_PASS=app@123
```

---


## 🚀 Running the Application  

This section provides step-by-step instructions to run the application either **locally** or via **Docker containers**.

- **Notes**:  
  - All commands are defined in the `Makefile`.
  - To run using `make`, ensure that `make` is installed on your system.
  - To run the application in containers, make sure `Docker` is installed and running.

### 🧪 Run Unit Tests

```bash
make test
```

### 🔧 Run Locally (Non-containerized)

Ensure Redis and PostgreSQL are running locally, then:

```bash
make run
```

### 🐳 Run Using Docker

To build and run all services (Redis, PostgreSQL, Go app):

```bash
make start-all
```

To stop and remove all containers:

```bash
make stop-all
```

- **Notes**:  
  - Before running the application inside Docker, make sure to update your environment variables `.env`
    - Change `DB_HOST=localhost` to `DB_HOST=postgres-server`.
    - Change `REDIS_HOST=localhost` to `REDIS_HOST=redis-server`.

### Application is Running

Now your application is accessible at:
```bash
http://localhost:1000
```

or 

```bash
https://localhost:1000 (if SSL is enabled)
```

---

## 🧪 Testing Scenarios  

### 🔐 Login API

**Endpoint**: `POST https://localhost:1000/auth/login`

#### ✅ Scenario 1: Successful Login

**Request**:

```json
{
  "username": "admin",
  "password": "P@ssw0rd"
}
```

**Response**:

```json
{
  "message": "Login successful",
  "error": null,
  "path": "/auth/login",
  "status": 200,
  "data": {
    "accessToken": "<JWT>",
    "refreshToken": "<UUID>",
    "expirationDate": "2025-05-25T12:58:00Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T12:58:00Z"
}
```

#### ❌ Scenario 2: Invalid Credentials

**Request with invalid user**:
```json
{
  "username": "invalid_user",
  "password": "P@ssw0rd"
}
```

**Response**:
```json
{
  "message": "Failed to login",
  "error": "user with the given username not found",
  "path": "/auth/login",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:18:23Z"
}
```

**Request with invalid password**:
```json
{
  "username": "admin",
  "password": "invalid_password"
}
```

**Response**:
```json
{
    "message": "Failed to login",
    "error": "invalid password",
    "path": "/auth/login",
    "status": 401,
    "data": null,
    "timestamp": "2025-05-23T15:51:39.288150079Z"
}
```

#### 🚫 Scenario 3: Disabled User

Precondition:
```sql
UPDATE users SET is_enabled = false WHERE id = 2;
```

**Request**:
```json
{
  "username": "userone",
  "password": "P@ssw0rd"
}
```

**Response**:
```json
{
  "message": "Failed to login",
  "error": "user is not enabled",
  "path": "/auth/login",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:19:24Z"
}
```

#### Scenario 4: Rate Limit Exceeded on Login

Precondition:
  - The rate limiter is configured as:
    - **rate.Limit**: rate.Every(30 * time.Second)
    - **burst**: 1
    - **expireAfter**: 5 * time.Minute
  - **Artinya**: allow `1 request` every `30 seconds`, with a burst capacity of `1`, within a `5-minute` window

**Request**: repeated quickly using valid credentials

```json
{
    "username": "admin",
    "password": "P@ssw0rd"
}
```

  - Steps:
    - Send the request once → receive access token.
    - Send the same request again shortly after (before 30 seconds pass).

**Response will be**:
```json
{
    "message": "Rate limit exceeded",
    "error": "You have exceeded the rate limit. Please try again later.",
    "path": "/auth/login",
    "status": 429,
    "data": null,
    "timestamp": "2025-05-23T16:01:30.407871957Z"
}
```


### 🔄 Refresh Token API

**Endpoint**: `POST https://localhost:1000/auth/refresh-token`

#### ✅ Scenario 1: Successful Refresh Token

**Request**:
```json
{
  "refreshToken": "<valid_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Token refreshed successfully",
  "error": null,
  "path": "/auth/refresh-token",
  "status": 200,
  "data": {
    "accessToken": "<JWT>",
    "refreshToken": "<new_UUID>",
    "expirationDate": "2025-05-25T15:23:51Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T15:23:51Z"
}
```

#### ❌ Scenario 2: Invalid Refresh Token

**Request**:
```json
{
  "refreshToken": "<invalid_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Failed to refresh token",
  "error": "record not found",
  "path": "/auth/refresh-token",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:24:47Z"
}
```

#### 🔁 Scenario 3: Expired Refresh Token (Auto Regenerate)

**Request**:
```json
{
  "refreshToken": "<expired_refresh_token>"
}
```

**Response**:
```json
{
  "message": "Token refreshed successfully",
  "error": null,
  "path": "/auth/refresh-token",
  "status": 200,
  "data": {
    "accessToken": "<new_JWT>",
    "refreshToken": "<new_UUID>",
    "expirationDate": "2025-05-25T15:29:02Z",
    "tokenType": "Bearer"
  },
  "timestamp": "2025-05-23T15:29:02Z"
}
```

### 🏢 Department API

**Endpoint**: `https://localhost:1000/api/v1/departments`

#### ✅ Scenario 1: Successful Retrieval (All Departments)

**Request Header**: `Authorization: Bearer <valid_token>`

**Response**:
```json
{
  "message": "All Departments retrieved successfully",
  "error": null,
  "path": "/api/v1/departments",
  "status": 200,
  "data": [
    {
      "id": "d001",
      "deptName": "Marketing",
      "active": true,
      "createdBy": 1,
      "createdAt": "2025-05-23T15:40:37Z",
      "updatedBy": 1,
      "updatedAt": "2025-05-23T15:40:37Z"
    },
    ...
  ],
  "timestamp": "2025-05-23T15:44:41Z"
}
```

#### ❌ Scenario 2: Expired Access Token

**Request Header**: `Authorization: Bearer <expired_token>`

**Response**:
```json
{
  "message": "Invalid token",
  "error": "token has invalid claims: token is expired",
  "path": "/api/v1/departments",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:45:50Z"
}
```

#### ❌ Scenario 3: Invalid Access Token

**Request Header**: `Authorization: Bearer <invalid_token>`

**Response**:
```json
{
  "message": "Invalid token",
  "error": "token signature is invalid: crypto/rsa: verification error",
  "path": "/api/v1/departments",
  "status": 401,
  "data": null,
  "timestamp": "2025-05-23T15:46:37Z"
}
```

### 💾 Redis Token Retrieval API Testing Scenario

**Endpoint**: `GET https://localhost:1000/api/v1/dataredis/json/{:key}`
**Description**: Retrieve an access token and related information stored in Redis using a specific Redis key.

#### ✅ Scenario 1: Successfully Retrieve Access Token from Redis

**Redis Key**: `access_token:admin`

**Request**: `GET https://localhost:1000/api/v1/dataredis/json/access_token:admin`

**Response**:
```json
{
    "message": "JSON value retrieved successfully",
    "error": null,
    "path": "/api/v1/dataredis/json/access_token:admin",
    "status": 200,
    "data": {
        "accessToken": "<JWT>",
        "expirationDate": "2025-05-25T16:01:29Z",
        "refreshToken": "<UUID>",
        "tokenType": "Bearer"
    },
    "timestamp": "2025-05-23T16:40:54.770426202Z"
}
```
