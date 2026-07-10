# Etapa 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cachea dependencias.
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main main.go

# Etapa 2: Imagen final mínima
FROM alpine:3.20

WORKDIR /root/

# Certificados para llamadas HTTPS salientes.
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
