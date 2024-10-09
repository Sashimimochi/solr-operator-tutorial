# Build stage
FROM golang:1.20-alpine AS builder

# ソースコードをコピーしてビルド
WORKDIR /app
COPY . .

RUN go mod init go-web-app && go mod tidy
RUN go build -o /go-web-app

# 実行ステージ
FROM alpine:3.20

WORKDIR /

COPY --from=builder /go-web-app /go-web-app

EXPOSE 8080

CMD ["/go-web-app"]
