FROM node:24-alpine AS frontend

WORKDIR /src
COPY package.json ./
RUN npm install --no-audit --no-fund
COPY . .
RUN npx tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify

FROM golang:1.25-alpine AS builder

WORKDIR /builder
ENV CGO_ENABLED=1
RUN for i in 1 2 3; do \
        apk add --no-cache gcc musl-dev && break; \
        [ "$i" -lt 3 ] || { echo "apk add failed after 3 attempts" >&2; exit 1; }; \
        echo "apk add failed (attempt $i/3), retrying..." >&2; \
        sleep $((i * 5)); \
    done
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@$(go list -f '{{.Version}}' -m github.com/a-h/templ)
COPY . .
RUN templ generate
RUN go build -o bin/main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /builder/bin/main ./main
COPY --from=builder /builder/database/seeds ./database/seeds
COPY --from=builder /builder/static ./static
COPY --from=frontend /src/static/css/style.css ./static/css/style.css
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 8080 8081

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./main", "-d", "postgres"]
