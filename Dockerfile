FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    make \
    g++ \
    gcc \
    musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o bin/web ./cmd/web

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install PostgreSQL client
RUN apk add --no-cache \
    postgresql-client \
    ca-certificates

# Copy binary and schema
COPY --from=builder /app/bin/web /app/bin/
COPY --from=builder /app/schema.sql /app/
COPY --from=builder /app/tls /app/tls

# Create entrypoint script
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'cd /app' >> /app/entrypoint.sh && \
    echo 'max_retries=30' >> /app/entrypoint.sh && \
    echo 'retry_count=0' >> /app/entrypoint.sh && \
    echo 'while ! pg_isready -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do' >> /app/entrypoint.sh && \
    echo '    echo "Waiting for PostgreSQL... (Attempt $retry_count of $max_retries)"' >> /app/entrypoint.sh && \
    echo '    retry_count=$((retry_count+1))' >> /app/entrypoint.sh && \
    echo '    if [ $retry_count -ge $max_retries ]; then' >> /app/entrypoint.sh && \
    echo '        echo "Failed to connect to PostgreSQL after $max_retries attempts"' >> /app/entrypoint.sh && \
    echo '        exit 1' >> /app/entrypoint.sh && \
    echo '    fi' >> /app/entrypoint.sh && \
    echo '    sleep 2' >> /app/entrypoint.sh && \
    echo 'done' >> /app/entrypoint.sh && \
    echo 'echo "PostgreSQL is ready!"' >> /app/entrypoint.sh && \
    echo 'DSN="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST/$POSTGRES_DB?sslmode=disable"' >> /app/entrypoint.sh && \
    echo 'exec bin/web --dsn "$DSN" "$@"' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]