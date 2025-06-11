# Etapa de build
FROM golang:1.24.3 AS builder

WORKDIR /app

# Copia os arquivos Go
COPY go.mod go.sum ./
COPY product/ ./product/

RUN go mod tidy

WORKDIR /app/product/worker

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o worker ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/product/worker/worker .
COPY product/worker/wait-for.sh .
RUN chmod +x wait-for.sh && apk add --no-cache netcat-openbsd

ENTRYPOINT ["/app/wait-for.sh", "/app/worker"]
