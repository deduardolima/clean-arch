# Etapa de construção
FROM golang:1.20-buster as builder

WORKDIR /app

# Copiar go.mod e go.sum e baixar as dependências
COPY go.mod go.sum ./
RUN go mod download

# Instalar Wire e migrate
RUN go install github.com/google/wire/cmd/wire@latest
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copiar o código-fonte
COPY . .

WORKDIR /app/cmd/ordersystem

# Gerar código com Wire
RUN /go/bin/wire

# Construir o executável incluindo wire_gen.go
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa final
FROM debian:buster-slim as app

WORKDIR /root/

COPY --from=builder /app/cmd/ordersystem/main .
COPY --from=builder /app/.env ./
COPY --from=builder /app/sql/migrations /root/migrations
COPY --from=builder /app/Makefile /root/
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Instalar cliente MySQL e make
RUN apt-get update && apt-get install -y default-mysql-client make

# Ajustar permissões do diretório de migrações
RUN chown -R root:root /root/migrations
RUN chmod -R 755 /root/migrations

# Expor portas
EXPOSE 8000
EXPOSE 8080
EXPOSE 50051

CMD ["./main"]
