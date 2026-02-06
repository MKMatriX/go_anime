FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go_test_task ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go_test_task /go_test_task
COPY --from=builder /app/migrations /migrations

EXPOSE 8080
CMD ["/go_test_task"]