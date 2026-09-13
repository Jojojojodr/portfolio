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
| Language  | Go 1.25                                      |
| Web       | Gin                                          |
| ORM       | GORM (SQLite / Postgres)                     |
| Templates | [Templ](https://templ.guide/)                |
| Frontend  | HTMX, [Tailwind CSS v4](https://tailwindcss.com/) |
| DevOps    | Docker, Docker Compose                       |
| Tooling   | [Task](https://taskfile.dev), [Air](https://github.com/air-verse/air) |

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 18+ (for Tailwind CSS build)
- Python 3 (for setup script)

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

## Routes

### Pages

| Method | Path                        | Description          |
| ------ | --------------------------- | -------------------- |
| GET    | `/`                         | Home page            |
| GET    | `/login`                    | Login page           |
| GET    | `/register`                 | Registration page    |
| GET    | `/profile`                  | Own profile (auth)   |
| GET    | `/profile/:id`              | User profile (admin) |
| GET    | `/blog/`                    | Blog listing         |
| GET    | `/blog/post`                | Single blog post     |
| GET    | `/admin/dashboard`          | Admin dashboard      |
| GET    | `/admin/users`              | Admin user list      |
| GET    | `/admin/posts`              | Admin post list      |
| GET    | `/admin/post/create`        | Create blog post     |
| GET    | `/admin/post/edit`          | Edit blog post       |

### API (`/v1`)

| Method | Path                         | Description              |
| ------ | ---------------------------- | ------------------------ |
| GET    | `/v1/health`                 | Health check             |
| GET    | `/v1/users`                  | List users               |
| POST   | `/v1/users`                  | Create user              |
| POST   | `/v1/login`                  | Authenticate (returns JWT) |
| GET    | `/v1/validate`               | Validate JWT token       |
| GET    | `/v1/blog/posts`             | Published blog posts     |
| GET    | `/v1/blog/post`              | Single blog post         |
| POST   | `/v1/blog/create-post`       | Create post (admin)      |

### HTMX / Form (`/handle`)

| Method | Path                           | Description           |
| ------ | ------------------------------ | --------------------- |
| POST   | `/handle/register`             | Register user         |
| POST   | `/handle/profile/update`       | Update profile (auth) |
| POST   | `/handle/auth/login`           | Login form            |
| POST   | `/handle/auth/logout`          | Logout                |
| GET    | `/handle/blog/posts`           | Blog list (HTMX)      |
| GET    | `/handle/blog/post`            | Blog post (HTMX)      |
| GET    | `/handle/blog/comments`        | Load comments         |
| POST   | `/handle/blog/comments/add`    | Add comment           |
| POST   | `/handle/like/post/:postId`    | Toggle post like      |
| POST   | `/handle/like/comment/:id`     | Toggle comment like   |
| POST   | `/handle/admin/post/create`    | Create post (admin)   |
| POST   | `/handle/admin/post/edit`      | Edit post (admin)     |
| POST   | `/handle/admin/preview-title`  | Preview title         |
| POST   | `/handle/admin/preview-markdown` | Preview markdown    |

## License

Licensed under the MIT License — see [LICENSE](LICENSE).
