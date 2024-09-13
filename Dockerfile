FROM golang:1.22.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api ./cmd/api

FROM debian:bullseye-slim

ENV SERVER_ADDRESS=0.0.0.0:8080
ENV POSTGRES_CONN=postgres://postgres:1234@db:5432/tender?sslmode=disable
ENV POSTGRES_JDBC_URL=jdbc:postgresql://db:5432/tender
ENV POSTGRES_USERNAME=postgres
ENV POSTGRES_PASSWORD=1234
ENV POSTGRES_HOST=db
ENV POSTGRES_PORT=5432
ENV POSTGRES_DATABASE=PostgreSQL
ENV ENV=prod

WORKDIR /app

COPY --from=builder /app/api .

COPY internal/storage/migrations /app/internal/storage/migrations

CMD ["/app/api"]