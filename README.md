# Portfolio Website

A modern, responsive portfolio website built with Go, featuring a blog system, user authentication, and admin panel.

## Features

- **Personal Portfolio**: Professional showcase of skills, experience, and projects
- **Blog System**: Create, edit, and manage blog posts with comments and likes
- **User Authentication**: Secure login/register system with JWT session management
- **Admin Panel**: Dashboard for managing users and content
- **HTMX Integration**: Dynamic interactions without JavaScript complexity
- **Responsive Design**: Mobile-first design with Tailwind CSS

## Tech Stack

| Layer     | Technology                                   |
| --------- | -------------------------------------------- |
| Language  | [Go 1.25](https://go.dev)                    |
| Web       | [Gin](https://gin-gonic.com)                 |
| ORM       | [GORM](https://gorm.io/index.html) (SQLite / Postgres) |
| Templates | [Templ](https://templ.guide/)                |
| Frontend  | [HTMX](https://htmx.org), [Tailwind CSS v4](https://tailwindcss.com/) |
| DevOps    | [Docker](https://docker.com), [Docker Compose](https://docs.docker.com/compose) |
| Tooling   | [Task](https://taskfile.dev), [Air](https://github.com/air-verse/air) |

## Getting Started

### Prerequisites

- Go 1.25+ with a C compiler (`gcc`); the SQLite driver uses CGO
- Node.js 20+ (for the Tailwind CSS build)
- Python 3 (for the setup script)
- Docker with the Compose plugin (only for running in Docker)

### Setup

```bash
git clone https://github.com/Jojojojodr/portfolio.git
cd portfolio
go mod download
python3 ./scripts/setup.py
```

The setup script:

1. Installs the Go tools `templ`, `air` and `task` with `go install` (make sure `$(go env GOPATH)/bin` is in your `PATH`)
2. Runs `npm i` to install Tailwind CSS
3. Copies `config.example.yaml` to `config.yaml` if it doesn't exist yet

### Configuration

The app reads `config.yaml` from the directory it is started in. There are no environment variables.

```yaml
server:
  host: localhost
  port: 8080            # public website
  admin_port: 8081      # admin panel
  # public_url: https://portfolio.example.com  # optional, see "Running behind a reverse proxy"
  # admin_url: https://admin.example.com       # optional
  jwt_secret: your_jwt_secret_key

database:
  type: sqlite          # "sqlite" or "postgres"
  path: ./data/sqlite.db  # sqlite only
  host: postgres        # postgres only
  port: 5432
  username: your_db_username
  password: your_db_password
  name: your_db_name
  ssl_mode: disable

gin:
  mode: debug           # "debug" or "release"
  trusted_proxies:
    - 127.0.0.1
```

- Set `jwt_secret` to a long random string for anything other than local development.
- `database.host: postgres` only resolves inside Docker Compose. For a PostgreSQL server outside Docker use its real hostname, e.g. `localhost`.

#### Running behind a reverse proxy

When the website and the admin panel are served on their own domains (for example with nginx, Caddy or Traefik), set `public_url` and `admin_url` so the links between the two sites (Blog, Profile, Admin Dashboard) point to those domains:

```yaml
server:
  port: 8080
  admin_port: 8081
  public_url: https://portfolio.example.com
  admin_url: https://admin.example.com
```

- Point the website domain to `server.port` and the admin domain to `server.admin_port`.
- Leave both empty for local development. The links then use the current host with `port` / `admin_port`, e.g. `http://localhost:8081`.
- The app tells the two sites apart by the port the request came in on, so the proxy has to forward each domain to its own port.

On startup the application connects to the database, runs the auto-migrations, and seeds it from `database/seeds/*.json` if it is empty.

## Running Locally

### Development

```bash
task serve   # live reload: Tailwind watcher + Air (rebuilds on .go / .templ changes)
task         # one-off run: templ generate + go run, no CSS build or reload
```

Both run from the project root, so they use `./config.yaml`, `./static`, `./database/seeds` and (with the default config) `./data/sqlite.db`.

- Website: `http://localhost:8080`
- Admin panel: `http://localhost:8081`

`task` does not build the CSS. Run `task css:build-css` first if `static/css/style.css` is missing or out of date.

### Building a Binary

```bash
task build   # build into bin/
task run     # build, then start the binary from inside bin/
```

`task build` generates the templates, builds the minified CSS and compiles the binary (`app`, or `app.exe` on Windows). It then sets up `bin/` as a self-contained folder:

```
bin/
├── app                  # the binary
├── config.yaml          # copied from config.example.yaml
├── database/seeds/      # copied from database/seeds
└── static/              # copied from static
```

`static`, `database/seeds` and `config.yaml` are only copied when they don't exist in `bin/` yet, so changes you make to `bin/config.yaml` survive rebuilds. This also means CSS, JS and image changes don't reach `bin/static` on a rebuild. Run `task clean` to start over with fresh copies.

`task run` starts the binary from inside `bin/`, so it uses `bin/config.yaml`. With the default SQLite path it needs a `bin/data/` directory, which gets automaticly created when the app runs.

## Running with Docker

The image is built in three stages: Node builds the CSS, Go compiles the binary (with CGO for SQLite), and a small Alpine image runs it. The container is configured only by `config.yaml`, which is mounted from the host.

### From Source with Docker Compose

1. Make sure `config.yaml` exists in the project root (`cp config.example.yaml config.yaml`). If it's missing, Docker creates an empty **directory** with that name and the app won't start.
2. Choose a database in `config.yaml` (see below).
3. Build and start:

```bash
task docker:build   # docker compose --profile postgres up -d --build
```

| Command                                        | What it does                                   |
| ---------------------------------------------- | ---------------------------------------------- |
| `docker compose up -d --build`                 | Build and start only the app (SQLite)          |
| `docker compose --profile postgres up -d --build` | Build and start the app and PostgreSQL      |
| `docker compose --profile postgres down`       | Stop and remove both containers                |
| `docker compose logs -f app`                   | Follow the app logs                            |

The `postgres` service is behind a [Compose profile](https://docs.docker.com/compose/how-tos/profiles/), so it only starts when `--profile postgres` is passed. The `docker:*` tasks always pass it.

#### SQLite

```yaml
database:
  type: sqlite
  path: ./data/sqlite.db
```

Keep the path inside `./data`. That directory is mounted at `/app/data`, so the database file is stored in `./data` on the host and survives rebuilds.

#### PostgreSQL

```yaml
database:
  type: postgres
  host: postgres          # the compose service name, not localhost
  port: 5432
  username: your_db_username
  password: your_db_password
  name: your_db_name
```

`username`, `password` and `name` must match `POSTGRES_USER`, `POSTGRES_PASSWORD` and `POSTGRES_DB` in `docker-compose.yaml`. The app waits for the database healthcheck before it starts. The data is stored in the `postgres_data` volume; `docker compose down -v` deletes it.

#### Volumes and Ports

| Host                | Container          | Purpose                                                     |
| ------------------- | ------------------ | ----------------------------------------------------------- |
| `./config.yaml`     | `/app/config.yaml` | Configuration, mounted read-only                            |
| `./data/`           | `/app/data/`       | SQLite database file                                        |
| `./static/img`      | `/app/static/img`  | Images uploaded in the admin panel                          |
| `8080`              | `8080`             | Public website, must match `server.port`                    |
| `8081`              | `8081`             | Admin panel, must match `server.admin_port`                 |

Uploaded images in `./static/img` survive rebuilds. If that folder is empty when the container starts, the default images from the image are copied into it.

### Prebuilt Image

Release images for `linux/amd64` and `linux/arm64` are published on Docker Hub as `jordyderon/portfolio`. See [.github/docs/DOCKERHUB.md](.github/docs/DOCKERHUB.md) for `docker run` and Compose examples that use the prebuilt image.

## Available Tasks

Run `task --list-all` to see every task.

| Task                 | Description                                                                 |
| -------------------- | --------------------------------------------------------------------------- |
| `task`               | Generate templates and run the app with `go run` from the project root      |
| `task serve`         | Dev server with live reload (Tailwind watcher + Air)                        |
| `task build`         | Generate templates, build CSS, compile the binary and set up `bin/`         |
| `task run`           | `task build`, then start the binary from inside `bin/`                      |
| `task clean`         | Remove `bin/` and `tmp/`                                                    |
| `task css:build-css` | Build minified CSS to `static/css/style.css`                                |
| `task css:watch-css` | Rebuild CSS on every change                                                 |
| `task dev:air`       | Run Air only (uses `.air_linux.toml` or `.air_windows.toml`)                |
| `task docker:up`     | `docker compose --profile postgres up -d`                                   |
| `task docker:build`  | `docker compose --profile postgres up -d --build` (rebuild the image first) |
| `task docker:down`   | `docker compose --profile postgres down`                                    |

## Default Credentials

The seed data creates these accounts. **Change the passwords right away** on any deployment that isn't local.

| Role  | Email              | Password   |
| ----- | ------------------ | ---------- |
| Admin | admin@example.com  | adminpass  |
| User  | user@example.com   | userpass   |

## License

Licensed under the MIT License — see [LICENSE](LICENSE).
