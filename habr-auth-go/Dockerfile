FROM golang:1.24.1 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/main.go

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/myapp .
CMD ["./myapp"]