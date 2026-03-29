# 构建阶段
FROM node:20-alpine AS builder
WORKDIR /app
# 复制依赖文件
COPY package.json package-lock.json ./
# 安装依赖
RUN npm install
# 复制所有前端代码
COPY . .
# 打包构建（生成dist目录）
RUN npm run build

# 运行阶段（Nginx部署，无改动）
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
# 复制构建后的dist目录到Nginx
COPY --from=builder /app/dist .
# 复制Nginx配置文件
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
# 等待后端就绪的脚本
COPY ./wait-for-backend.sh /wait-for-backend.sh
RUN chmod +x /wait-for-backend.sh
# 安装 curl 用于就绪探测
RUN apk add --no-cache curl
# 暴露80端口
EXPOSE 80
# 启动脚本：先等待后端健康，再启动 nginx
CMD ["/bin/sh", "-c", "/wait-for-backend.sh && nginx -g 'daemon off;'"]