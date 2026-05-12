FROM node:20-alpine AS frontend-builder

WORKDIR /app/new-ui

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build || true

FROM golang:1.26-alpine AS builder

WORKDIR /build

# Install dependencies for CGO and SQLite FTS
RUN apk --no-cache add build-base linux-headers sqlite-dev sqlite

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy already-built frontend
COPY --from=frontend-builder /app/new-ui/dist ./new-ui/dist

RUN statik -src=./new-ui/dist -dest=./statik -f || true

# Build with CGO enabled for SQLite FTS5 support
RUN CGO_ENABLED=1 go build -o movielite --tags "fts5" ./cmd/server

# Stage 2: Runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs curl

WORKDIR /app

COPY --from=builder /build/movielite ./movielite
COPY movielite.yaml.tmpl ./movielite.yaml

EXPOSE 8001

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f -u admin:password http://localhost:8001/api/user || exit 1

CMD ["./movielite", "start"]
