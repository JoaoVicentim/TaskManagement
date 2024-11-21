# Usar a imagem base do Golang
FROM golang:1.22-alpine AS builder

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o restante dos arquivos da aplicação
COPY . .

# Compilar a aplicação para ARM
RUN GOOS=linux GOARCH=arm64 go build -o main ./cmd/main.go

# Usar uma imagem Alpine para rodar a aplicação
FROM alpine:latest

# Copiar o binário da aplicação
COPY --from=builder /app/main .

# Expor a porta que a aplicação vai usar
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]