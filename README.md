## Directory Structure

```
fiber-starter/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── db/
│   └── migrations/
│       ├── 000001_init_users.up.sql
│       ├── 000001_init_users.down.sql
├── internal/
│   ├── domain/
│   │   └── entity.go
│   ├── handler/
│   │   └── user_handler.go
│   ├── repository/
│   │   └── user_repository.go
│   ├── service/
│   │   └── user_service.go
├── pkg/
│   └── response/
│       └── response.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```
