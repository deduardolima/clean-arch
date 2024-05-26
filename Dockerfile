FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/google/wire/cmd/wire@latest
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

WORKDIR /app/cmd/ordersystem

RUN /go/bin/wire

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM debian:buster-slim as app

WORKDIR /root/

COPY --from=builder /app/cmd/ordersystem/main .
COPY --from=builder /app/.env ./
COPY --from=builder /app/sql/migrations /root/migrations
COPY --from=builder /app/Makefile /root/
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

RUN apt-get update && apt-get install -y default-mysql-client make

RUN chown -R root:root /root/migrations
RUN chmod -R 755 /root/migrations

EXPOSE 8000
EXPOSE 8080
EXPOSE 50051

CMD ["./main"]
