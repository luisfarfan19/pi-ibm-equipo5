# Build stage
FROM golang:1.24.1 AS builder

WORKDIR /app

# Copy go.mod and go.sum before other files to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage
FROM gcr.io/distroless/static-debian11

WORKDIR /

# Copy the binary
COPY --from=builder /app/main .

# Run the binary
ENTRYPOINT ["/main"]
