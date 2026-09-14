# Portfolio Release

## Quick Start

1. **Download** the archive for your platform from [Releases](https://github.com/Jojojojodr/portfolio/releases)
2. **Extract** the archive

   ```bash
   tar -xzf portfolio-linux-amd64-*.tar.gz
   cd portfolio
   ```

3. **Edit** `config.yaml` with your values (see [Configuration](#configuration))
4. **Make executable** (Linux/macOS)

   ```bash
   chmod +x portfolio-linux-amd64-*
   ```

5. **Run**

   ```bash
   ./portfolio-linux-amd64-* -p 8080 -d sqlite -t your-secret-token
   ```

## What's Included

```
portfolio-<platform>/
├── portfolio-<platform>-*    # Compiled binary (platform-specific)
├── static/                   # CSS and image assets
├── database/
│   └── seeds/                # Initial seed data (JSON)
└── config.yaml               # Default configuration
```

## CLI Flags

```bash
./portfolio-<platform>-* -p <port> -d <database> -t <jwt-secret>
```

| Flag | Description                                        | Required |
| ---- | -------------------------------------------------- | -------- |
| `-p` | Port number (e.g. `8080`)                          | Yes      |
| `-d` | Database type: `sqlite` or `postgres`              | Yes      |
| `-t` | JWT secret token for authentication                | Yes      |

These can alternatively be set in `config.yaml` (see below). The flags take precedence over the config file.

## Configuration

Edit `config.yaml` before first run:

```yaml
server:
  host: "localhost"
  port: "8080"
  jwt_secret: "your-secure-secret-key"

database:
  # SQLite settings
  type: "sqlite"
  path: "./database/sqlite.db"

  # Postgres settings (only used when type is "postgres")
  # host: "localhost"
  # port: "5432"
  # username: "your_db_user"
  # password: "your_db_password"
  # name: "your_db_name"
  # ssl_mode: "disable"

gin:
  mode: "release"
  trusted_proxies: ["127.0.0.1"]
```

## First Run

On first launch the application will:

1. Connect to the database and run auto-migrations
2. Seed an admin and a regular user if the database is empty

### Default Credentials

| Role  | Email             | Password  |
| ----- | ----------------- | --------- |
| Admin | admin@example.com | adminpass |
| User  | user@example.com  | userpass  |

> **Change these immediately** by updating `database/seeds/users.json` before first run, or create new users via the registration page and delete the seeded accounts via the admin panel.

## Postgres

To use Postgres instead of SQLite, update `config.yaml`:

```yaml
database:
  type: "postgres"
  host: "localhost"
  port: "5432"
  username: "myuser"
  password: "mypassword"
  name: "portfolio"
  ssl_mode: "disable"
```

The database is created automatically on first run; no manual migration step is needed.

## Troubleshooting

**App won't start?**
- Make sure `jwt_secret` is set in `config.yaml` or pass `-t`
- Check that the `database/` directory is writable (SQLite)
- Verify the port isn't already in use

**Static files not loading?**
- Ensure the `static/` directory is in the same directory as the binary
- Check file permissions

**Database errors?**
- Ensure the `database/` directory exists and is writable
- For Postgres, verify the connection details and that the server is running
