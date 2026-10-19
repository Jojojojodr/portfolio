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

- Go 1.25+
- Node.js 18+ (for Tailwind CSS build)
- Python (for setup script)

### Setup

```bash
git clone https://github.com/Jojojojodr/portfolio.git
cd portfolio
go mod download
python3 ./scripts/setup.py   # installs templ, air, task; runs npm i; copies config.yaml; creates database/sqlite.db
```

### Run

```bash
task serve   # live reload — watches CSS + Go code (via Air)
```

The server starts at `http://localhost:8080`.

On first run the application:
1. Connects to the database and runs auto-migrations
2. Seeds an admin and a regular user if the database is empty

### Docker

`docker compose up -d` builds the app and starts it together with a PostgreSQL container. No `.env` file is needed: the app is configured entirely by `config.yaml`, which is mounted into the container as a volume. `database.type` in that file selects the database:

- **sqlite**: keep `path: "./data/sqlite.db"`. The database file is stored in `./data` on the host, which is mounted at `/app/data`.
- **postgres**: set `database.host: "postgres"` (the compose service name, not `localhost`). Keep `username` / `password` / `name` matching the `POSTGRES_*` values in `docker-compose.yaml` (`your_db_username` / `your_db_password` / `your_db_name` by default).

The published ports `8080` / `8081` must match `server.port` / `server.admin_port` in the config.

### Default Credentials

| Role  | Email              | Password   |
| ----- | ------------------ | ---------- |
| Admin | admin@example.com  | adminpass  |
| User  | user@example.com   | userpass   |

## Available Tasks

```bash
task                # Run the application (templ generate + go run)
task serve          # Start dev server with live reload (Air + Tailwind watcher)
task build          # Build the binary to bin/
task run            # Build then run
task docker         # docker-compose up -d
task docker-build   # Rebuild Docker image and start
task docker-stop    # docker-compose down
task docker-rebuild # Full rebuild (stop + build + recreate)
task clean          # Remove bin/ and tmp/
```

## License

Licensed under the MIT License — see [LICENSE](LICENSE).
