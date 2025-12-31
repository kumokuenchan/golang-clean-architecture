# Clean Architecture Go CRUD Sample

A complete implementation of Clean Architecture principles in Go for CRUD operations. This project demonstrates how to structure a Go application following Robert C. Martin's Clean Architecture pattern with modern tooling including OpenAPI specification, Docker containerization, and jinzhu/copier for object mapping.

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
│   ├── dto/
│   │   └── user_dto.go          # Data transfer objects (requests/responses)
│   ├── mapper/
│   │   └── user_mapper.go       # Object mapping using jinzhu/copier
│   ├── presenter/
│   │   └── user_presenter.go    # Data formatting and presentation
│   ├── api/
│   │   ├── server.go            # HTTP API server implementation
│   │   └── errors.go            # Common API error definitions
│   ├── infrastructure/
│   │   └── database/
│   │       ├── connection.go    # Database connection management
│   │       └── user_repository.go # Repository implementation
│   └── delivery/
│       └── http/
│           └── router.go        # Route configuration
├── docker/
│   └── mysql/
│       └── init.sql             # Database initialization script
├── Dockerfile                   # Multi-stage Docker build configuration
├── docker-compose.yml           # Container orchestration
├── oapi.yaml                    # OpenAPI 3.0 specification
├── oapi-codegen.yaml            # Code generation configuration
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
- **Purpose**: Implements application-specific business rules and orchestration
- **Contains**: Use cases, business logic, validation, error handling
- **Dependencies**: Domain, DTO, Mapper, and Presenter layers
- **Example**: `internal/usecase/user_usecase.go`

```go
func (u *UserUsecase) CreateUser(name, email string) ([]byte, error) {
    if name == "" || email == "" {
        return u.presenter.PresentError("name and email are required")
    }
    
    // Use copier mapper to map request to domain
    req := &dto.CreateUserRequest{Name: name, Email: email}
    domainUser := mapper.MapCreateRequestToDomain(req)
    
    if err := u.userRepo.Create(domainUser); err != nil {
        return u.presenter.PresentError(err.Error())
    }
    
    return u.presenter.PresentUser(domainUser)
}
```

### 3. DTO Layer (Data Transfer Objects)
- **Purpose**: Defines request/response structures for API layer
- **Contains**: Request/response types, validation tags
- **Dependencies**: None
- **Example**: `internal/dto/user_dto.go`

```go
type CreateUserRequest struct {
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type UserResponse struct {
    ID        int    `json:"id" example:"1"`
    Name      string `json:"name" example:"John Doe"`
    Email     string `json:"email" example:"john@example.com"`
    CreatedAt string `json:"created_at" example:"2023-01-01 12:00:00"`
    UpdatedAt string `json:"updated_at" example:"2023-01-01 12:00:00"`
}
```

### 4. Mapper Layer (Object Mapping)
- **Purpose**: Handles object transformation between layers using jinzhu/copier
- **Contains**: Mapping functions, data transformation logic
- **Dependencies**: Domain and DTO layers
- **Example**: `internal/mapper/user_mapper.go`

```go
func MapCreateRequestToDomain(req *dto.CreateUserRequest) *domain.User {
    user := &domain.User{}
    copier.Copy(user, req)
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    return user
}

func MapDomainToResponse(user *domain.User) *dto.UserResponse {
    response := &dto.UserResponse{}
    copier.Copy(response, user)
    response.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
    response.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
    return response
}
```

### 5. Presenter Layer (Interface Adapter)
- **Purpose**: Formats data for presentation and converts between layers
- **Contains**: Presenters, formatters
- **Dependencies**: Domain and DTO layers
- **Example**: `internal/presenter/user_presenter.go`

```go
type UserPresenter interface {
    PresentUser(user *domain.User) ([]byte, error)
    PresentUsers(users []*domain.User) ([]byte, error)
    PresentError(message string) ([]byte, error)
}

func (p *HTTPUserPresenter) PresentUser(user *domain.User) ([]byte, error) {
    response := mapper.MapDomainToResponse(user)
    return json.Marshal(response)
}
```

### 6. API Layer (Interface Adapter)
- **Purpose**: Implements HTTP handlers with DRY principles and helper methods
- **Contains**: HTTP handlers, response utilities, error definitions
- **Dependencies**: Use case and DTO layers
- **Example**: `internal/api/server.go`

```go
// Helper methods eliminate code duplication
func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, data []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(data)
}

func (s *Server) extractAndValidateID(r *http.Request) (int, error) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        return 0, ErrIDRequired
    }
    return strconv.Atoi(idStr)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    if err := s.decodeJSONBody(r, &req); err != nil {
        s.writeErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }
    
    data, err := s.userUsecase.CreateUser(req.Name, req.Email)
    if err != nil {
        s.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    s.writeJSONResponse(w, http.StatusCreated, data)
}
```

### 7. Infrastructure Layer (Interface Adapter)
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

### 8. Delivery Layer (Interface Adapter)
- **Purpose**: Handles external requests and coordinates the application
- **Contains**: HTTP routers, middleware
- **Dependencies**: API layer
- **Example**: `internal/delivery/http/router.go`

```go
func (r *Router) SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
        switch req.Method {
        case http.MethodGet:
            r.server.GetAllUsers(w, req)
        case http.MethodPost:
            r.server.CreateUser(w, req)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    return mux
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

## 🔄 Object Mapping with jinzhu/copier

This project uses **jinzhu/copier** for automatic object mapping between layers:

### Key Features
- **Automatic field copying**: No manual assignment needed
- **Type safety**: Compile-time checking of field mappings
- **Performance optimized**: Efficient reflection-based copying
- **Flexible mapping**: Handles nested structures and custom transformations

### Mapping Examples
```go
// Request to Domain
req := &dto.CreateUserRequest{Name: "John", Email: "john@example.com"}
domainUser := mapper.MapCreateRequestToDomain(req)

// Domain to Response
response := mapper.MapDomainToResponse(domainUser)

// List mapping
usersResponse := mapper.MapDomainListToResponse(domainUsers)
```

## 🎯 DRY Principles in API Layer

The API layer follows **Don't Repeat Yourself (DRY)** principles to eliminate code duplication:

### Helper Methods
- **`writeJSONResponse()`**: Centralized JSON response writing
- **`writeErrorResponse()`**: Centralized error response handling
- **`extractAndValidateID()`**: Reusable ID extraction and validation
- **`decodeJSONBody()`**: Reusable JSON decoding

### Error Management
- **`internal/api/errors.go`**: Common error definitions
- **Consistent error handling**: All endpoints use the same error patterns
- **Maintainable**: Changes to response format only need updates in one place

### Benefits
- **Reduced code duplication**: From ~200 lines to ~120 lines per handler
- **Consistent responses**: All endpoints use the same response patterns
- **Easier maintenance**: Changes propagate automatically across all handlers
- **Better testability**: Helper methods can be unit tested independently

## 📋 OpenAPI Specification

The API is defined using OpenAPI 3.0 specification in `oapi.yaml`:

### Features
- **Complete API documentation**: All endpoints documented with examples
- **Type definitions**: Request/response schemas with validation
- **Code generation ready**: Configured for oapi-codegen tool
- **Standard compliance**: Follows OpenAPI 3.0 standards

### Key Endpoints
- `POST /users` - Create new user
- `GET /users` - List all users  
- `GET /user?id={id}` - Get user by ID
- `PUT /user?id={id}` - Update user
- `DELETE /user?id={id}` - Delete user

### Code Generation
```bash
# Generate Go code from OpenAPI spec
oapi-codegen -config oapi-codegen.yaml oapi.yaml
```

### Option 2: Local Development

#### Prerequisites
- Go 1.21 or higher
- MySQL database

#### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd golang-clean-architecture
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
- **jinzhu/copier**: Automatic object mapping between structs
- **OpenAPI 3.0**: API specification and documentation
- **oapi-codegen**: Code generation from OpenAPI specs

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
5. **Specification-First Development**: API defined in OpenAPI before implementation
6. **Automated Object Mapping**: Uses copier to reduce boilerplate and errors

## 🤝 Contributing

This is a sample project demonstrating Clean Architecture principles. Feel free to use it as a reference for your own projects.

## 📄 License

This project is for educational purposes.
