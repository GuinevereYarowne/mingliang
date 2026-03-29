# 构建阶段
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 编译web-server
RUN CGO_ENABLED=0 GOOS=linux go build -o web-server/web-server web-server/main.go

# 运行阶段
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/web-server/web-server .
EXPOSE 8080
CMD ["./web-server"]