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
