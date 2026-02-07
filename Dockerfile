FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./server/go.mod ./server/go.sum ./.env ./
RUN go mod download

COPY ./server .
RUN CGO_ENABLED=0 GOOS=linux go build -o /builded ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /builded /builded
# COPY --from=builder /app/migrations /migrations

EXPOSE 8080
CMD ["/builded"]