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
