FROM golang:1.26-alpine AS builder

WORKDIR /build

RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o movielite --tags "fts5" ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates && adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /build/movielite ./movielite
COPY movielite.yaml.tmpl ./movielite.yaml

EXPOSE 8001

USER appuser
CMD ["./movielite", "start"]
