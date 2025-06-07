# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala dependencias necesarias (git, ca-certificates, etc.)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ⚠️ CGO habilitado por defecto aquí (mejor para compatibilidad)
RUN go build -o main .

# Final stage con Alpine (mínimo pero funcional)
FROM alpine:3.18

# Instala certificados si los necesitas
RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/main"]
