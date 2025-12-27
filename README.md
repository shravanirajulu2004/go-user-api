# Go User API - Backend Development Task

A production-ready RESTful API built with **Go**, **Fiber**, **PostgreSQL**, and **SQLC** for managing users with **dynamic age calculation**.

 


---

## ✨ Key Features

- ✅ **CRUD operations** for user management
- ✅ **Dynamic age calculation** from date of birth (never stored in database)
- ✅ **Type-safe database queries** with SQLC
- ✅ **Input validation** with go-playground/validator
- ✅ **Structured logging** with Uber Zap
- ✅ **Clean architecture** (handler → service → repository)
- ✅ **Request ID tracking** and duration logging middleware
- ✅ **Unit tests** for core business logic

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Framework** | GoFiber v2.52 |
| **Database** | PostgreSQL 13+ |
| **Query Builder** | SQLC v1.30 |
| **Validation** | go-playground/validator v10 |
| **Logging** | Uber Zap |
| **Architecture** | Clean Architecture Pattern |

---

## 📁 Project Structure

```
go-user-api/
├── cmd/server/          # Application entry point (main.go)
├── config/              # Environment configuration
├── db/
│   ├── migrations/      # SQL migration files
│   ├── queries/         # SQLC query definitions
│   └── sqlc/           # Generated type-safe Go code
├── internal/
│   ├── handler/        # HTTP request handlers
│   ├── service/        # Business logic layer
│   ├── repository/     # Data access layer
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Domain models and DTOs
│   └── logger/         # Logging configuration
├── .env                # Environment variables
├── .env.example        # Example environment file
├── sqlc.yaml          # SQLC configuration
├── go.mod             # Go module definition
├── README.md          # This file
└── reasoning.md       # Design decisions and rationale
```

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21 or higher
- **PostgreSQL** 13 or higher
- **SQLC** CLI tool

### Installation Steps

```bash
# 1. Clone the repository
git clone https://github.com/shravanirajulu2004/go-user-api.git
cd go-user-api

# 2. Install dependencies
go mod download

# 3. Install SQLC (if not already installed)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 4. Generate SQLC code
sqlc generate

# 5. Create database
createdb -U postgres userapi

# 6. Run migration
psql -U postgres -d userapi -f db/migrations/001_create_users_table.sql

# 7. Configure environment
cp .env.example .env
# Edit .env with your PostgreSQL password

# 8. Run the server
go run cmd/server/main.go
```

Server will start at **http://localhost:3000**

---

## 🔌 API Endpoints

### Base URL: `http://localhost:3000`

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| `GET` | `/health` | Health check | - | `{"status":"ok"}` |
| `POST` | `/users` | Create user | `{"name":"Alice","dob":"1990-05-10"}` | User object |
| `GET` | `/users/:id` | Get user by ID | - | User with **calculated age** |
| `GET` | `/users` | List all users | - | Array of users with ages |
| `PUT` | `/users/:id` | Update user | `{"name":"Alice","dob":"1990-05-10"}` | Updated user |
| `DELETE` | `/users/:id` | Delete user | - | HTTP 204 No Content |

---

## 📝 Example Usage

### 1. Create a User

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","dob":"1990-05-10"}'
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice Johnson",
  "dob": "1990-05-10"
}
```

### 2. Get User with Dynamic Age

```bash
curl http://localhost:3000/users/1
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice Johnson",
  "dob": "1990-05-10",
  "age": 34
}
```

**Note:** Age is calculated dynamically based on **today's date**!

### 3. List All Users

```bash
curl http://localhost:3000/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Alice Johnson",
    "dob": "1990-05-10",
    "age": 34
  }
]
```

### 4. Update User

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","dob":"1991-03-15"}'
```

### 5. Delete User

```bash
curl -X DELETE http://localhost:3000/users/1
```

---

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./internal/models/ -v

# Run with coverage
go test ./... -cover
```

**Test Output:**
```
=== RUN   TestCalculateAge
--- PASS: TestCalculateAge (0.00s)
=== RUN   TestCalculateAge_EdgeCases
--- PASS: TestCalculateAge_EdgeCases (0.00s)
PASS
```

---

## ⚙️ Configuration

Create a `.env` file in the project root:

```env
DATABASE_URL=postgresql://postgres:YOUR_PASSWORD@localhost:5432/userapi?sslmode=disable
PORT=3000
ENV=development
```

---

## 🎯 Key Implementation Highlights

### 1. **Dynamic Age Calculation** ⭐

Age is **never stored** in the database. It's calculated in real-time:

```go
func CalculateAge(dob time.Time) int {
    now := time.Now()
    age := now.Year() - dob.Year()
    
    // Subtract 1 if birthday hasn't occurred yet this year
    if now.Month() < dob.Month() || 
       (now.Month() == dob.Month() && now.Day() < dob.Day()) {
        age--
    }
    
    return age
}
```

**Why?** Stored age becomes stale. DOB is the single source of truth.

### 2. **Type-Safe Database Queries with SQLC**

SQLC generates type-safe Go code from SQL:

```sql
-- name: GetUserByID :one
SELECT id, name, dob, created_at, updated_at
FROM users
WHERE id = $1;
```

Generates:
```go
func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error)
```

**Benefits:** Compile-time type safety, no reflection overhead.

### 3. **Clean Architecture**

```
Handler (HTTP) → Service (Business Logic) → Repository (Database)
```

Each layer has a single responsibility, making the code **testable** and **maintainable**.

### 4. **Input Validation**

```go
type CreateUserRequest struct {
    Name string `json:"name" validate:"required"`
    DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}
```

Invalid requests are rejected before reaching business logic.

---

## 🔍 Design Decisions

See **[reasoning.md](reasoning.md)** for detailed explanations of:

- Why SQLC over ORMs (GORM, sqlx)
- Dynamic age calculation approach
- Clean architecture implementation
- Error handling strategy
- Middleware design
- Testing approach
- Challenges faced and solutions

---

## 🐛 Troubleshooting

### "sqlc: command not found"
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### "database does not exist"
```bash
createdb -U postgres userapi
psql -U postgres -d userapi -f db/migrations/001_create_users_table.sql
```

### "port 3000 already in use"
Change `PORT=3001` in `.env` file

### "failed to connect to database"
- Verify PostgreSQL is running: `pg_isready`
- Check credentials in `.env` file

---

## 📚 API Documentation

### Request Validation Rules

- **Name:** Required, non-empty string
- **DOB:** Required, format `YYYY-MM-DD` (e.g., "1990-05-10")

### Error Responses

```json
// Validation error
{
  "error": "Validation failed: Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}

// Not found
{
  "error": "User not found"
}

// Internal server error
{
  "error": "Internal server error"
}
```

---

## 🚀 Future Enhancements

Potential improvements for production deployment:

- [ ] Docker containerization
- [ ] Comprehensive integration tests
- [ ] API documentation with Swagger/OpenAPI
- [ ] Rate limiting
- [ ] Caching layer with Redis
- [ ] JWT authentication
- [ ] Soft deletes
- [ ] Search functionality
- [ ] Pagination metadata

---



