# 使用官方 Golang 镜像作为基础镜像
# 构建阶段和运行阶段分离，使用多阶段构建优化镜像大小

# 构建阶段
FROM golang:1.24-alpine AS builder
# 工作目录设为根目录（对应宿主机的geekedu-project）
WORKDIR /app
# 复制根目录的go.mod和go.sum（关键！）
COPY go.mod go.sum ./
# 下载依赖
RUN go mod download
# 复制所有代码（根目录下的common、proto、logic-server等）
COPY . .
# 编译logic-server（指定输出到logic-server目录）
RUN CGO_ENABLED=0 GOOS=linux go build -o logic-server/logic-server logic-server/main.go

# 运行阶段
FROM alpine:3.18
WORKDIR /app
# 复制编译后的二进制文件
COPY --from=builder /app/logic-server/logic-server .
# 暴露端口
EXPOSE 50051
# 启动
CMD ["./logic-server"]