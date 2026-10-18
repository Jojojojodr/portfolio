#!/usr/bin/env sh
set -eu

q() {
    printf "'%s'" "$(printf '%s' "$1" | sed "s/'/''/g")"
}

if [ -z "${JWT_SECRET:-}" ]; then
    echo "error: JWT_SECRET is not set." >&2
    echo "Set it in .env (JWT_SECRET=...) or pass -e JWT_SECRET=<value>." >&2
    exit 1
fi

if [ "$JWT_SECRET" = "your_jwt_secret_key" ]; then
    echo "warning: JWT_SECRET is the placeholder value, set a real secret before deploying." >&2
fi

DB_TYPE="${DB_TYPE:-postgres}"
DB_PATH="${DB_PATH:-./data/sqlite.db}"

mkdir -p "$(dirname "$DB_PATH")"

cat > config.yaml <<EOF
server:
  host: $(q "${SERVER_HOST:-localhost}")
  port: $(q "${SERVER_PORT:-8080}")
  admin_port: $(q "${ADMIN_PORT:-8081}")
  jwt_secret: $(q "$JWT_SECRET")
database:
  type: $(q "$DB_TYPE")
  path: $(q "$DB_PATH")
  host: $(q "${DB_HOST:-postgres}")
  port: $(q "${DB_PORT:-5432}")
  username: $(q "${DB_USER:-your_db_username}")
  password: $(q "${DB_PASSWORD:-your_db_password}")
  name: $(q "${DB_NAME:-your_db_name}")
  ssl_mode: $(q "${DB_SSLMODE:-disable}")
gin:
  mode: $(q "${GIN_MODE:-release}")
  trusted_proxies: [$(q "${GIN_TRUSTED_PROXIES:-127.0.0.1}")]
EOF

exec ./main "$@"
