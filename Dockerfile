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

# Install MySQL client
RUN apk add --no-cache \
    mysql-client \
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
    echo 'while ! mysql -h$MYSQL_HOST -u$MYSQL_USER -p$MYSQL_PASSWORD -e "SELECT 1" >/dev/null 2>&1; do' >> /app/entrypoint.sh && \
    echo '    echo "Waiting for MySQL... (Attempt $retry_count of $max_retries)"' >> /app/entrypoint.sh && \
    echo '    retry_count=$((retry_count+1))' >> /app/entrypoint.sh && \
    echo '    if [ $retry_count -ge $max_retries ]; then' >> /app/entrypoint.sh && \
    echo '        echo "Failed to connect to MySQL after $max_retries attempts"' >> /app/entrypoint.sh && \
    echo '        exit 1' >> /app/entrypoint.sh && \
    echo '    fi' >> /app/entrypoint.sh && \
    echo '    sleep 2' >> /app/entrypoint.sh && \
    echo 'done' >> /app/entrypoint.sh && \
    echo 'echo "MySQL is ready!"' >> /app/entrypoint.sh && \
    echo 'DSN="$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST)/$MYSQL_DATABASE?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"' >> /app/entrypoint.sh && \
    echo 'exec bin/web --dsn "$DSN" "$@"' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]