# Etapa de build
FROM golang:1.24.3 AS builder

WORKDIR /app

# Copia os arquivos Go
COPY go.mod go.sum ./
RUN go mod download

COPY product/worker/. .

# Compila o binário de forma estática e otimizada
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Etapa de runtime
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /app/app .

# Executa como usuário não-root (UID 65532)
USER nonroot:nonroot

ENTRYPOINT ["/app"]
