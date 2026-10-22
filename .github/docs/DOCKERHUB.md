# Portfolio

A personal portfolio and blog web application written in Go (Gin, GORM, Templ, Tailwind CSS), with a public site and a separate admin panel.

Source code: [github.com/Jojojojodr/portfolio](https://github.com/Jojojojodr/portfolio)

## Tags

| Tag       | Description                          |
| --------- | ------------------------------------ |
| `latest`  | Latest release                       |
| `X.Y.Z`   | A specific release, e.g. `1.2.3`     |
| `X.Y`     | Latest patch release of a minor line |
| `X`       | Latest release of a major line       |

Images are published for `linux/amd64` and `linux/arm64`.

## Ports and volumes

| Path / Port         | Description                                                         |
| ------------------- | ------------------------------------------------------------------- |
| `8080`              | Public website (`server.port`)                                      |
| `8081`              | Admin panel (`server.admin_port`)                                   |
| `/app/config.yaml`  | Configuration file, **required**                                    |
| `/app/data`         | SQLite database directory, only needed when `database.type: sqlite` |
| `/app/static/img`   | Images uploaded in the admin panel, filled with defaults when empty |

## Configuration

The container is configured entirely through a `config.yaml` file that you mount into the container. Create one next to where you run the container:

```yaml
server:
  host: localhost
  port: 8080
  admin_port: 8081
  jwt_secret: change_me_to_a_long_random_string
database:
  type: sqlite              # "sqlite" or "postgres"
  path: ./data/sqlite.db    # sqlite only, keep it inside ./data so it lands in the volume
  host: postgres            # postgres only, the hostname of the postgres container
  port: 5432
  username: portfolio
  password: change_me
  name: portfolio
  ssl_mode: disable
gin:
  mode: release
  trusted_proxies:
    - 127.0.0.1
```

- `jwt_secret` must be set to a real secret.
- If you change `server.port` or `server.admin_port`, change the published ports to match.
- On first start the database is migrated and seeded with an admin (`admin@example.com` / `adminpass`) and a regular user (`user@example.com` / `userpass`). **Change these passwords right away.**

## Using plain Docker

### SQLite

```bash
docker run -d \
  --name portfolio \
  -p 8080:8080 -p 8081:8081 \
  -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
  -v "$(pwd)/static/img:/app/static/img" \
  -v portfolio_data:/app/data \
  --restart unless-stopped \
  jordyderon/portfolio:latest
```

The site is available at `http://localhost:8080` and the admin panel at `http://localhost:8081`. The database file is stored in the `portfolio_data` volume and uploaded images in `./static/img` on the host.

### PostgreSQL

Set `database.type: "postgres"` and `database.host: "postgres"` in `config.yaml`, then put both containers on the same network:

```bash
docker network create portfolio

docker run -d \
  --name postgres \
  --network portfolio \
  -e POSTGRES_USER=portfolio \
  -e POSTGRES_PASSWORD=change_me \
  -e POSTGRES_DB=portfolio \
  -v portfolio_postgres:/var/lib/postgresql \
  --restart unless-stopped \
  postgres:18

docker run -d \
  --name portfolio \
  --network portfolio \
  -p 8080:8080 -p 8081:8081 \
  -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
  -v "$(pwd)/static/img:/app/static/img" \
  --restart unless-stopped \
  jordyderon/portfolio:latest
```

`database.host` must match the name of the postgres container (`postgres` here), and the username, password and name must match the `POSTGRES_*` variables. Make sure postgres has finished starting before you start the app. If it hasn't, the app exits and `--restart` brings it back once the database is up.

## Using Docker Compose

### SQLite

```yaml
services:
  app:
    image: jordyderon/portfolio:latest
    container_name: portfolio-app
    ports:
      - "8080:8080"
      - "8081:8081"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./static/img:/app/static/img
      - portfolio_data:/app/data
    restart: unless-stopped

volumes:
  portfolio_data:
```

### PostgreSQL

```yaml
services:
  app:
    image: jordyderon/portfolio:latest
    container_name: portfolio-app
    ports:
      - "8080:8080"
      - "8081:8081"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./static/img:/app/static/img
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:18
    container_name: portfolio-postgres
    environment:
      POSTGRES_USER: portfolio
      POSTGRES_PASSWORD: change_me
      POSTGRES_DB: portfolio
    volumes:
      - postgres_data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $${POSTGRES_USER} -d $${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres_data:
```

Set `database.host: "postgres"` (the service name) in `config.yaml`. Then start it:

```bash
docker compose up -d
docker compose logs -f app
```

## Updating

```bash
docker compose pull && docker compose up -d
# or with plain docker
docker pull jordyderon/portfolio:latest && docker rm -f portfolio && docker run ...   # same command as above
```

Your data stays in the volumes, and the config file stays on the host.

## License

MIT, see [LICENSE](https://github.com/Jojojojodr/portfolio/blob/main/LICENSE).
