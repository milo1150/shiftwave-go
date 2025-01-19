## Basic CLI

- [Click](https://github.com/milo1150/cli-cmds)

## Example .env

```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=mydb
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/mydb
JWT=helloworld
BASE_URL=http://localhost:8080
APP_ENV=development
LOCAL_TIMEZONE=Asia/Bangkok
```

## Directory Breakdown

- internal/: Contains the core application code, segregated into

- database/: Database connection setup and migrations

- handlers/: HTTP handlers managing incoming requests

- repositories/: Data access layer interacting with GORM

- services/: Business logic

- setup/: JSON data loader

- models/: Database models representing your data structures

- configs/: Configuration management

- cmd/: Entry point(s) of your application

- pkg/, api/, configs/, scripts/, tests/, docs/: As previously described

```bash
your-service/
├── cmd/
│   └── yourservice/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── auth.go          // General authentication logic
│   │   ├── password.go      // Specific password-related functions
│   │   └── jwt.go           // JWT-related logic
│   ├── handler/
│   │   ├── root_handler.go
│   │   └── review_handler.go
│   ├── dto/
│   │   └── review_dto.go
│   ├── services/
│   │   └── review_service.go
│   ├── repository/
│   │   └── review_repository.go
│   ├── resources/
│   │   └── master_data.json
│   ├── model/
│   │   └── review_model.go
│   ├── types/
│   │   └── review_type.go
│   ├── database/
│   │   ├── connection.go
│   ├── utils/
│   │   └── validation.go
│   └── configs/
│   │   └── config.yaml
│   └── middleware.go
├── pkg/
│   └── (shared libraries)
├── api/
│   └── (API definitions)
├── configs/
│   └── config.yaml
├── scripts/
│   └── build.sh
├── tests/
│   └── (optional additional tests)
├── docs/
│   └── README.md
├── go.mod
└── go.sum
```
