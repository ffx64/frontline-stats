# Etapa 1: build
FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache tzdata

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    TZ=America/Sao_Paulo

WORKDIR /app

# Copia go.mod e go.sum para cache
COPY go.mod go.sum ./
RUN go mod download

# Copia o resto do código
COPY . .

# Build do binário
RUN go build -o gamestats-backend ./cmd/main.go

# Etapa 2: imagem final minimalista
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata

# Copia o binário do builder
COPY --from=builder /app/gamestats-backend .

# Porta padrão da API
EXPOSE 8080

CMD ["./gamestats-backend"]
