FROM node:lts-alpine AS nodebuilder
WORKDIR /go/src/github.com/ms900ft/movielite
COPY . .
WORKDIR /go/src/github.com/ms900ft/movielite/movieui
RUN npm install
RUN npm run build



FROM golang:alpine AS builder
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
#ENV GOBIN=/usr/local/bin
RUN go get github.com/rakyll/statik

WORKDIR /go/src/github.com/ms900ft/movielite
COPY . .
COPY docker/movielite.yaml /go/src/github.com/ms900ft/movielite
COPY --from=nodebuilder /go/src/github.com/ms900ft/movielite/movieui/dist /go/src/github.com/ms900ft/movielite/movieui/dist
RUN /go/bin/statik  -src=/go/src/github.com/ms900ft/movielite/movieui/dist
RUN CGO_ENABLED=1 GOOS=linux  go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" --tags "fts5" -o movielite cmd/server/main.go





FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/ms900ft/movielite/movielite .
#COPY --from=nodebuilder /movieui/dist /dist
COPY movielite.yaml.tmpl movielite.yaml
COPY example/movielite.db example/
CMD ["./movielite", "start"]
