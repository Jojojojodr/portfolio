# Portfolio Website

A modern, responsive portfolio website built with Go, featuring a blog system, user authentication, and admin panel. This project showcases my skills as a software developer and serves as a platform to share my thoughts through blog posts.

## Features

- **Personal Portfolio**: Professional showcase of skills, experience, and projects
- **Blog System**: Create, edit, and manage blog posts with comments and likes
- **User Authentication**: Secure login/register system with session management
- **Admin Panel**: Comprehensive dashboard for managing users and content
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Dark Theme**: Modern dark color scheme for better user experience
- **HTMX Integration**: Dynamic interactions without JavaScript complexity

## Tech Stack

### Backend
- **Go**: Main programming language
- **Gin**: HTTP web framework
- **SQLite/Postgres**: Database for data persistence
- **Templ**: Go templating engine for type-safe HTML

### Frontend
- **Templ** - Type-safe HTML templating
- **Tailwind CSS**: Utility-first CSS framework
- **HTMX**: Modern interaction library

### DevOps & Tools
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration
- **Task**: Task runner for development workflows
- **Air**: Live reload for Go development

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional)
- Task (optional, for development commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Jojojojodr/portfolio.git
   cd portfolio
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment and download dependencies**
   ```bash
   python ./scripts/setup.py
   ```

4. **Run the dev server**
    ```bash
    task serve
    ```

5. **Access the application**
   - Website: http://localhost:8080
   - Admin Panel: http://localhost:8080/admin/dashboard

## Available Tasks

Use the Task runner for common development tasks:

```bash
task                 # Runs the application using go run
task build           # Build the application
task serve           # Start development server with live reload
task start           # Run the application
task run             # Build and run the application
task docker          # Starts the application in docker
task docker-build    # Rebuild the Docker image and start the application
task docker-stop     # Stops the Docker image
task docker-rebuild  # Rebuild the Docker image and forse recreate
task db              # Migrates the database fresh
task seed            # Seeds to the database
task clean           # Cleans up application
```

## Features Overview

### Portfolio Section
- Professional profile with skills showcase
- Project portfolio with technology tags
- Education history
- Personal hobbies and interests

### Blog System
- Create and edit blog posts
- Comment system with user interactions
- Like/dislike functionality

### User Management
- Secure authentication system
- Admin panel for user management

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [Templ](https://templ.guide/) - Template engine
- [HTMX](https://htmx.org/) - Dynamic HTML
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [GORM](https://gorm.io/) - ORM library
