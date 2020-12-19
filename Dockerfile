FROM node:lts-alpine AS nodebuilder
WORKDIR /go/src/ms/movielite
COPY . .
WORKDIR /go/src/ms/movielite/movieui
RUN npm install
RUN npm run build



FROM golang:alpine AS builder
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
#ENV GOBIN=/usr/local/bin
RUN go get github.com/rakyll/statik

WORKDIR /go/src/ms/movielite
COPY . .
COPY movielite.yaml /go/src/ms/movielite
COPY --from=nodebuilder /go/src/ms/movielite/movieui/dist /go/src/ms/movielite/movieui/dist
RUN /go/bin/statik  -src=/go/src/ms/movielite/movieui/dist
RUN CGO_ENABLED=1 GOOS=linux  go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" --tags "fts5" -o movielite cmd/server/main.go





FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/ms/movielite/movielite .
#COPY --from=nodebuilder /movieui/dist /dist
COPY Docker/movielite.yaml .
CMD ["./movielite", "start"]
