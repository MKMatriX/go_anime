# -------- билдеры ----------
# Основа
FROM golang:1.25-alpine AS builder-base
WORKDIR /app
COPY ./server/go.mod ./server/go.sum .env ./
RUN go mod download
COPY ./server .

# Главный сервис
FROM builder-base AS builder-main
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded-main ./cmd/app/main

# Сервисы
FROM builder-base AS builder-anidb
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded-anidb ./cmd/app/anidb
FROM builder-base AS builder-anilist
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded-anilist ./cmd/app/anilist
FROM builder-base AS builder-animetosho
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded-animetosho ./cmd/app/animetosho
FROM builder-base AS builder-shikimori
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded-shikimori ./cmd/app/shikimori


# -------- сервера ----------
# Основа
FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates

# Главный сервис
FROM runtime AS main
COPY --from=builder-main /builded-main /builded
EXPOSE 8080
CMD ["/builded"]

# Сервисы
FROM runtime AS anidb
COPY --from=builder-anidb /builded-anidb /builded
EXPOSE 8081
CMD ["/builded"]
FROM runtime AS anilist
COPY --from=builder-anilist /builded-anilist /builded
EXPOSE 8082
CMD ["/builded"]
FROM runtime AS animetosho
COPY --from=builder-animetosho /builded-animetosho /builded
EXPOSE 8083
CMD ["/builded"]
FROM runtime AS shikimori
COPY --from=builder-shikimori /builded-shikimori /builded
EXPOSE 8084
CMD ["/builded"]