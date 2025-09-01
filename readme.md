# min-notify

内网邮件通知服务 | API 驱动的轻量级通知中心

[![Docker Pulls](https://img.shields.io/docker/pulls/zouzonghua/min-notify)](https://hub.docker.com/r/zouzonghua/min-notify)
[![GitHub release](https://img.shields.io/github/v/release/zouzonghua/min-notify)](https://github.com/zouzonghua/min-notify/releases)
[![License](https://img.shields.io/github/license/zouzonghua/min-notify)](LICENSE)

## 目录

- [特性](#特性)
- [支持场景](#支持场景)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [接口文档](#接口文档)
- [版本发布](#版本发布)
- [安全建议](#安全建议)

## 特性

- 🚀 轻量级 Web API
- 🔗 内网接口
- 📧 跨应用邮件通知
- 🐳 Docker 一键部署

## 支持场景

- 服务器告警
- 系统监控
- 设备通知
- 日志推送

## 技术栈

- Go
- Docker
- SMTP 协议

## 架构设计

1. 内网 HTTP 接口
2. 邮件发送服务
3. 配置管理

## 快速开始

### Docker 方式（推荐）

```bash
# 创建配置目录
mkdir -p data

# 运行容器
docker run -d --name min-notify \
  -p 5001:5001 \
  -v $(pwd)/data:/app/data \
  zouzonghua/min-notify:latest
```

### 源码构建

```bash
# 克隆项目
git clone https://github.com/zouzonghua/min-notify.git
cd min-notify

# 构建镜像
docker build -t min-notify .

# 运行容器
docker run -d --name min-notify -p 5001:5001 min-notify
```

## 配置说明

配置文件位置：`data/config.json`

```json
{
    "smtp_server": "smtp.gmail.com",    // SMTP 服务器地址
    "smtp_port": 587,                   // SMTP 端口
    "smtp_user": "sender@gmail.com",    // 发件人邮箱
    "smtp_pass": "your-password",       // 邮箱授权码
    "to_email": "target@example.com",   // 收件人邮箱
}
```

## 接口文档

### 发送通知

**POST** `/notify`

请求参数：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 通知标题 |
| message | string | 是 | 通知内容 |

响应示例：

```json
{
    "code": 0,
    "message": "success"
}
```

调用示例：

```bash
# curl
curl -X POST http://localhost:5001/notify \
  -H "Content-Type: application/json" \
  -d '{
    "title": "系统告警",
    "message": "磁盘空间不足"
  }'

# python
import requests
requests.post('http://localhost:5001/notify', json={
    'title': '告警',
    'message': '服务异常'
})
```

## 版本发布

### Docker 镜像版本

镜像托管在 Docker Hub，支持以下标签：

- `latest`: 最新版本（main 分支更新时）
- `vx.y.z`: 具体版本号（如 `v1.0.0`）

发布流程：
```bash
# 1. 提交代码到 main 分支会更新 latest 标签
git push origin main

# 2. 发布新版本（会创建对应版本号的镜像标签）
git tag v1.0.0
git push origin v1.0.0
```


> 注意：
> 1. 推送到 main 分支只会更新 `latest` 标签
> 2. 创建并推送 tag 才会生成版本号镜像
> 3. tag 必须以 `v` 开头，如 `v1.0.0`


## 使用示例：
```bash
# 使用最新版本
docker pull zouzonghua/min-notify:latest

# 使用特定版本
docker pull zouzonghua/min-notify:v1.0.0
```
## 安全建议

- 🔒 仅限内网
- 🛡️ 可增加 API Key
- 🔐 HTTPS 加密

## 许可

MIT License
