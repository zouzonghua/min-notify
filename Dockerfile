# 使用 golang:1.22-alpine 作为构建阶段的基础镜像
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
# 构建 Go 项目，生成可执行文件 min-notify
RUN go build -o min-notify main.go

# 使用更小的 alpine 镜像作为运行环境
FROM alpine:latest
WORKDIR /app

# 安装时区数据并设置为中国时区
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制可执行文件
COPY --from=builder /app/min-notify .
# 复制静态资源文件夹
COPY static ./static
# 声明配置文件挂载点
VOLUME ["/app/data"]
# 开放 5001 端口
EXPOSE 5001

# 启动前如果没有 config.json，则自动生成一个默认配置文件
ENTRYPOINT ["/bin/sh", "-c", "if [ ! -f /app/data/config.json ]; then echo '{}' > /app/data/config.json; fi && ./min-notify"]
