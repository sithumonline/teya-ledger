# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -mod=mod -o /main ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /main /main
EXPOSE 8080
CMD ["/main"]
