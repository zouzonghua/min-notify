# min-notify

内网邮件通知服务 | API 驱动的轻量级通知中心

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

## 快速部署

构建镜像

```bash
docker build -t min-notify .
docker run -d --name min-notify -p 5001:5001 min-notify
```

运行容器

```bash
docker run -d --name min-notify -p 5001:5001 zouzonghua/min-notify:latest
```

## 版本发布

### Docker 镜像版本

镜像托管在 Docker Hub，支持以下标签：

- `latest`: 最新版本
- `x.y.z`: 具体版本号（如 `1.0.0`）
- `x.y`: 主次版本号（如 `1.0`）


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

## 调用示例

```bash
# Curl 测试
curl -X POST http://localhost:5001/notify \
  -H "Content-Type: application/json" \
  -d '{
    "title":"系统告警",
    "message":"磁盘空间不足"
  }'

# Python 测试
import requests
requests.post('http://localhost:5001/notify', json={
    'title': '告警',
    'message': '服务异常'
})
```

## 安全建议

- 🔒 仅限内网
- 🛡️ 可增加 API Key
- 🔐 HTTPS 加密

## 许可

MIT License
