# 使用 golang:1.22-alpine 作为构建阶段的基础镜像
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
# 构建 Go 项目，生成可执行文件 min-notify
RUN go build -o min-notify main.go

# 使用更小的 alpine 镜像作为运行环境
FROM alpine:latest
WORKDIR /app

# 安装时区数据（时区由环境变量 TZ 控制，不硬编码）
RUN apk add --no-cache tzdata

# 从构建阶段复制可执行文件
COPY --from=builder /app/min-notify .
# 复制静态资源文件夹
COPY static ./static
# 声明配置文件挂载点
VOLUME ["/app/data"]
# 开放 5001 端口
EXPOSE 5001

# 启动脚本：根据 TZ 环境变量设置时区，并初始化配置文件
ENTRYPOINT ["/bin/sh", "-c", "\
    if [ -n \"$TZ\" ]; then \
        ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
        echo $TZ > /etc/timezone; \
    fi && \
    if [ ! -f /app/data/config.json ]; then \
        echo '{}' > /app/data/config.json; \
    fi && \
    ./min-notify"]
