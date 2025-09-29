# Portfolio Application Release

## Quick Start

1. **Download and extract** the archive for your platform
2. **Copy environment file**: `cp .env.example .env`
3. **Edit .env** with your configuration:
   ```bash
   SECRET_TOKEN=your-super-secret-jwt-key
   PORT=8080
   DB_TYPE=sqlite
   DB_PATH=./database/sqlite.db
   ```
4. **Make executable** (Linux/macOS): `chmod +x portfolio`
5. **Run the application**:
   ```bash
   ./portfolio
   or
   ./portfolio -p 8080 -d sqlite -t your-secret-token
   ```

## What's Included

- ✅ Compiled Go binary
- ✅ Compiled Tailwind CSS
- ✅ All static assets (CSS, images)
- ✅ Database directory structure
- ✅ Example configuration

## Command Line Options

```bash
./portfolio -p <port> -d <database> -t <jwt-secret>
```

- `-p`: Port number (e.g., 8080)
- `-d`: Database type (sqlite or postgres) 
- `-t`: JWT secret token for authentication

## Environment Variables (Alternative)

Instead of command line flags, you can use environment variables:

```bash
export PORT=8080
export DB_TYPE=sqlite  
export DB_PATH=./database/sqlite.db
export SECRET_TOKEN=your-secret-token
```

## First Run

The application will:
1. Create the database and run migrations
2. Seed initial data if database is empty
3. Start the web server on the specified port

Visit `http://localhost:8080` in your browser.

## Troubleshooting

**App won't start?**
- Make sure you have a valid `SECRET_TOKEN` set
- Check that the database directory is writable
- Verify the port isn't already in use

**Static files not loading?**
- Ensure the `static/` directory is in the same location as the binary
- Check file permissions

**Database issues?**
- Make sure the `database/` directory exists and is writable
- For SQLite, ensure the `DB_PATH` points to a valid location