# Guide

## Example .env

```bash
# Echo app
JWT=gummy
BASE_URL=http://localhost:8080
APP_URL=http://localhost:4321
APP_ENV=development
LOCAL_TIMEZONE=Asia/Bangkok
OPENAI_API_KEY=youropenaitoken
ADMIN_PASSWORD=hello
REDIS_PASSWORD=hello

# Postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=mydb
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5436/mydb
```

- Shell db

```bash
pgcli postgres://postgres:postgres@127.0.0.1:5436/mydb
```
