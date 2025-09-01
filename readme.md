# min-notify

内网邮件通知服务 | API 驱动的轻量级通知中心

## 特性
- 🚀 轻量级 Web API
- 🔗 内网接口
- 📧 跨应用邮件通知
- 🐳 Docker 一键部署

## 架构设计
1. 内网 HTTP 接口
2. 邮件发送服务
3. 配置管理

## 快速部署
构建镜像
```
docker build -t min-notify .
```


运行容器
```
docker run -d --name min-notify \
  -p 5001:5001 \
  -v $(pwd)/config.json:/app/config.json \
  min-notify
```

配置文件 (config.json)
```
{
    "smtp_server": "smtp.gmail.com",
    "smtp_port": 587,
    "smtp_user": "sender@gmail.com",
    "smtp_pass": "授权码",
    "to_email": "target@example.com"
}
```

调用示例
```
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

## 支持场景
- 服务器告警
- 系统监控
- 设备通知
- 日志推送

## 技术栈
- Go
- Docker
- SMTP 协议

## 安全建议
- 🔒 仅限内网
- 🛡️ 可增加 API Key
- 🔐 HTTPS 加密

## 许可
MIT License
