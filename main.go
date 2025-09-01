package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
)

type Config struct {
	SMTPServer string `json:"smtp_server"`
	SMTPPort   int    `json:"smtp_port"`
	SMTPUser   string `json:"smtp_user"`
	SMTPPass   string `json:"smtp_pass"`
	ToEmail    string `json:"to_email"`
}

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

var (
	config     Config
	configLock sync.RWMutex
	configFile = "config.json"
)

// ---------------- 配置管理 ----------------

// 加载配置
func loadConfig() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config = Config{
			SMTPServer: "smtp.qq.com",
			SMTPPort:   465, // 默认使用 SSL 端口
			SMTPUser:   "",
			SMTPPass:   "",
			ToEmail:    "",
		}
		saveConfig()
	} else {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatalf("读取配置失败: %v", err)
		}
		_ = json.Unmarshal(data, &config)
	}
}

// 保存配置
func saveConfig() {
	data, _ := json.MarshalIndent(config, "", "  ")
	_ = ioutil.WriteFile(configFile, data, 0644)
}

// ---------------- 邮件发送 ----------------

func sendEmail(subject, body string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	addr := fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)

	log.Printf("准备发送邮件: server=%s port=%d user=%s to=%s",
		config.SMTPServer, config.SMTPPort, config.SMTPUser, config.ToEmail)

	// 构造邮件头
	header := make(map[string]string)
	header["From"] = config.SMTPUser
	header["To"] = config.ToEmail
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	// 建立 TLS 连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.SMTPServer,
	})
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, config.SMTPServer)
	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}
	defer client.Quit()

	// 登录认证
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, config.SMTPServer)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(config.SMTPUser); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(config.ToEmail); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("写入数据失败: %v", err)
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("发送内容失败: %v", err)
	}
	_ = w.Close()

	log.Println("邮件发送成功")
	return nil
}

// ---------------- HTTP API ----------------

// 提供邮件通知接口
func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var notif Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if notif.Title == "" {
		notif.Title = "通知"
	}
	if notif.Message == "" {
		notif.Message = "无内容"
	}

	if err := sendEmail(notif.Title, notif.Message); err != nil {
		log.Println("邮件发送失败:", err)
		http.Error(w, "发送邮件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

// 提供配置 API
func configHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		configLock.RLock()
		defer configLock.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(config)

	case http.MethodPost:
		var newConfig Config
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, "Invalid config", http.StatusBadRequest)
			return
		}
		configLock.Lock()
		config = newConfig
		configLock.Unlock()
		saveConfig()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "saved"}`))
	}
}

// ---------------- Main ----------------

func main() {
	loadConfig()

	// 静态文件 (配置页面)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/config", configHandler)
	http.HandleFunc("/notify", notifyHandler)

	fmt.Println("服务启动: http://0.0.0.0:5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
