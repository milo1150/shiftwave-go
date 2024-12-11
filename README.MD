# NOTE

## DB

### if Docker & Postgres: Failed to bind tcp 0.0.0.0:5432 address already in use

``` bash
sudo lsof -i :5432
```

``` bash
kill -9 <PID>
```

### CLI check db

``` bash
 pgcli postgres://postgres:postgres@127.0.0.1:5432/mydb
 pgcli postgres://postgres:postgres@0.0.0.0:5432/mydb
```

## Docker

### Docker Stop and Remove Everything Related to Docker Compose

```bash
 docker-compose down -v
```

down: Stops and removes containers and networks defined in the docker-compose.yml.
-v: Removes any associated volumes.

#### Remove All Images Built by Docker Compose: To remove images built by Docker Compose, include the --rmi flag

```bash
 docker-compose down --rmi all -v
```

--rmi all: Removes all images built by docker-compose.
-v: Also removes volumes associated with the containers.

## Deploying for Production

### Build the Production Images

```bash
 docker-compose -f docker-compose.production.yml build
```

### Run in Detached Mode

```bash
 docker-compose -f docker-compose.production.yml build
```

### Check Logs

```bash
 docker-compose -f docker-compose.production.yml build
```

### Test run in local pc

```bash
 docker-compose -f docker-compose.production.yml up --build
```

### Folder Structure

 Directory Breakdown
 cmd/: Entry point(s) of your application.
 internal/: Contains the core application code, segregated into:
 handlers/: HTTP handlers managing incoming requests.
 services/: Business logic.
 repositories/: Data access layer interacting with GORM.
 models/: Database models representing your data structures.
 database/: Database connection setup and migrations.
 configs/: Configuration management.
 pkg/, api/, configs/, scripts/, tests/, docs/: As previously described.

```bash
your-service/
├── cmd/
│   └── yourservice/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── root_handler.go
│   │   └── assessment_handler.go
│   ├── services/
│   │   └── assessment_service.go
│   ├── repositories/
│   │   └── assessment_repository.go
│   ├── models/
│   │   └── assessment_model.go
│   ├── types/
│   │   └── assessment_type.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrations/
│   │       └── 001_create_assessments_table.up.sql
│   └── configs/
│       └── config.yaml
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
