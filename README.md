# Clean Architecture Go CRUD Sample

A complete implementation of Clean Architecture principles in Go for CRUD operations. This project demonstrates how to structure a Go application following Robert C. Martin's Clean Architecture pattern.

## 🏗️ Architecture Overview

Clean Architecture organizes code into concentric layers, with the business rules at the center and technical details at the outer layers. The dependency rule states that dependencies can only point inward.

```
└──────────────────────────────────────────────────────────────────────────────┐
│                        Frameworks & Drivers                 │
│                    (Web, Database, External APIs)           │
├──────────────────────────────────────────────────────────────────────────────┤
│                         Interface Adapters                   │
│                 (Presenters, Controllers, Gateways)          │
├──────────────────────────────────────────────────────────────────────────────┤
│                           Use Cases                          │
│                    (Application Business Rules)              │
├──────────────────────────────────────────────────────────────────────────────┤
│                              Entities                        │
│                     (Enterprise Business Rules)              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── user.go              # Entities and repository interfaces
│   ├── usecase/
│   │   └── user_usecase.go      # Business logic and application rules
│   ├── presenter/
│   │   └── user_presenter.go    # Data formatting and presentation
│   ├── infrastructure/
│   │   └── database/
│   │       ├── connection.go    # Database connection management
│   │       └── user_repository.go # Repository implementation
│   └── delivery/
│       └── http/
│           ├── user_handler.go  # HTTP request handlers
│           └── router.go        # Route configuration
├── schema.sql                   # Database schema
├── go.mod                       # Go module file
└── README.md                    # This file
```

## 🧱 Architecture Layers Explained

### 1. Domain Layer (Entities)
- **Purpose**: Contains the core business entities and rules
- **Contains**: Business models, repository interfaces
- **Dependencies**: None (most inner layer)
- **Example**: `internal/domain/user.go`

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
    Create(user *User) error
    GetByID(id int) (*User, error)
    GetAll() ([]*User, error)
    Update(user *User) error
    Delete(id int) error
}
```

### 2. Use Case Layer (Application Business Rules)
- **Purpose**: Implements application-specific business rules
- **Contains**: Use cases, application services
- **Dependencies**: Domain layer only
- **Example**: `internal/usecase/user_usecase.go`

```go
type UserUsecase struct {
    userRepo  domain.UserRepository
    presenter presenter.UserPresenter
}

func (u *UserUsecase) CreateUser(name, email string) ([]byte, error) {
    // Business logic validation
    if name == "" || email == "" {
        return u.presenter.PresentError("name and email are required")
    }
    
    // Entity creation
    user := &domain.User{
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Repository interaction
    if err := u.userRepo.Create(user); err != nil {
        return u.presenter.PresentError(err.Error())
    }
    
    return u.presenter.PresentUser(user)
}
```

### 3. Presenter Layer (Interface Adapter)
- **Purpose**: Formats data for presentation and converts between layers
- **Contains**: Presenters, formatters
- **Dependencies**: Domain layer
- **Example**: `internal/presenter/user_presenter.go`

```go
type UserPresenter interface {
    PresentUser(user *domain.User) ([]byte, error)
    PresentUsers(users []*domain.User) ([]byte, error)
    PresentError(message string) ([]byte, error)
}

func (p *HTTPUserPresenter) PresentUser(user *domain.User) ([]byte, error) {
    output := UserOutputData{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
    return json.Marshal(output)
}
```

### 4. Infrastructure Layer (Interface Adapter)
- **Purpose**: Implements technical details and external interfaces
- **Contains**: Database repositories, external API clients
- **Dependencies**: Domain layer (implements interfaces)
- **Example**: `internal/infrastructure/database/user_repository.go`

```go
type MySQLUserRepository struct {
    db *sql.DB
}

func (r *MySQLUserRepository) Create(user *domain.User) error {
    query := "INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)"
    result, err := r.db.Exec(query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
    if err != nil {
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    user.ID = int(id)
    return nil
}
```

### 5. Delivery Layer (Interface Adapter)
- **Purpose**: Handles external requests and coordinates the application
- **Contains**: HTTP handlers, controllers, routers
- **Dependencies**: Use case layer
- **Example**: `internal/delivery/http/user_handler.go`

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Delegate to use case
    data, err := h.userUsecase.CreateUser(req.Name, req.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(data)
}
```

## 🚀 Getting Started

### Option 1: Docker (Recommended)

#### Prerequisites
- Docker and Docker Compose

#### Quick Start

1. Clone the repository:
```bash
git clone <repository-url>
cd clean-arch-sample
```

2. Spin up the Docker environment:
```bash
docker-compose up --build
```

This will:
- Build and start the MySQL database with sample data
- Build and start the Go application
- The API will be available at `http://localhost:8080`

3. Run in background (detached mode):
```bash
docker-compose up -d --build
```

4. Check running containers:
```bash
docker-compose ps
```

5. View logs:
```bash
# View all logs
docker-compose logs

# View API logs only
docker-compose logs api

# View MySQL logs only
docker-compose logs mysql

# Follow logs in real-time
docker-compose logs -f api
```

6. Stop the services:
```bash
docker-compose down
```

7. Remove volumes (clean start):
```bash
docker-compose down -v
```

8. Rebuild and restart:
```bash
docker-compose down
docker-compose up --build
```

#### Docker Commands Reference

```bash
# Start services
docker-compose up

# Start services in background
docker-compose up -d

# Build without starting
docker-compose build

# Build specific service
docker-compose build api

# Stop services
docker-compose stop

# Remove stopped containers
docker-compose rm

# View resource usage
docker-compose top

# Execute command in container
docker-compose exec api sh
docker-compose exec api ./main

# Access MySQL in container
docker-compose exec mysql mysql -u appuser -p cleanarch_sample
```

### Option 2: Local Development

#### Prerequisites
- Go 1.21 or higher
- MySQL database

#### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd clean-arch-sample
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```sql
CREATE DATABASE cleanarch_sample;
USE cleanarch_sample;
```

4. Run the schema:
```bash
mysql -u root -p cleanarch_sample < schema.sql
```

5. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=cleanarch_sample
export PORT=8080
```

6. Run the application:
```bash
go run cmd/api/main.go
```

## 📡 API Endpoints

### Users Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | Get all users |
| POST | `/users` | Create a new user |
| GET | `/user?id={id}` | Get user by ID |
| PUT | `/user?id={id}` | Update user |
| DELETE | `/user?id={id}` | Delete user |

### Request/Response Examples

**Create User:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

**Get All Users:**
```bash
curl http://localhost:8080/users
```

**Get User by ID:**
```bash
curl "http://localhost:8080/user?id=1"
```

**Update User:**
```bash
curl -X PUT "http://localhost:8080/user?id=1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com"}'
```

**Delete User:**
```bash
curl -X DELETE "http://localhost:8080/user?id=1"
```

## 🎯 Benefits of Clean Architecture

1. **Independence from Frameworks**: Business rules don't depend on external libraries
2. **Testability**: Business rules can be tested without UI, database, or external services
3. **Independence from UI**: The UI can change without changing the business rules
4. **Independence from Database**: Business rules are not bound to the database
5. **Independence from External Services**: Business rules don't know about the outside world

## 🔄 Dependency Flow

```
main.go
  ↓
delivery/http (Controllers)
  ↓
usecase (Application Rules)
  ↓
domain (Business Rules)
  ↑
infrastructure/database (Repository Implementation)
```

## 🧪 Testing

The architecture enables easy testing at each layer:

```bash
# Test domain entities
go test ./internal/domain/...

# Test use cases
go test ./internal/usecase/...

# Test presenters
go test ./internal/presenter/...

# Test infrastructure
go test ./internal/infrastructure/...

# Test delivery layer
go test ./internal/delivery/...
```

## 🛠️ Technologies Used

- **Go**: Programming language
- **MySQL**: Database
- **HTTP**: Web framework (standard library)
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration

## 🐳 Docker Configuration

The project includes Docker configuration for easy development and deployment:

### Dockerfile
Multi-stage build for optimized production image:
- **Builder stage**: Compiles Go application
- **Final stage**: Lightweight Alpine Linux image with binary

### Docker Compose Services
- **mysql**: MySQL 8.0 database with persistent storage
- **api**: Go application with health checks and proper dependencies

### Environment Variables
The Docker setup automatically configures:
- `DB_HOST=mysql`
- `DB_PORT=3306`
- `DB_USER=appuser`
- `DB_PASSWORD=apppassword`
- `DB_NAME=cleanarch_sample`
- `PORT=8080`

### Development Workflow
```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Clean up everything
docker-compose down -v --rmi all
```

## 📝 Key Principles

1. **Single Responsibility**: Each component has one reason to change
2. **Dependency Inversion**: High-level modules don't depend on low-level modules
3. **Separation of Concerns**: Each layer handles specific concerns
4. **Testability**: All components can be tested in isolation

## 🤝 Contributing

This is a sample project demonstrating Clean Architecture principles. Feel free to use it as a reference for your own projects.

## 📄 License

This project is for educational purposes.
