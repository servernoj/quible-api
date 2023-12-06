package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/mail-service/controller"
	"gitlab.com/quible-backend/mail-service/service"
)

const (
	serverTokenEnv  = "POSTMARK_SERVER_TOKEN"
	accountTokenEnv = "POSTMARK_ACCOUNT_TOKEN"
	webPort         = "80"
)

func main() {
	// 从环境变量中读取 Postmark 令牌
	serverToken := os.Getenv(serverTokenEnv)
	if serverToken == "" {
		log.Fatalf("Error: Missing environment variable %s\n", serverTokenEnv)
	}
	accountToken := os.Getenv(accountTokenEnv)
	if accountToken == "" {
		log.Fatalf("Error: Missing environment variable %s\n", accountTokenEnv)
	}

	// 创建 Postmark 客户端
	client := &service.Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      "https://api.postmarkapp.com",
	}

	// 初始化 Gin
	router := gin.Default()

	// 设置路由
	controller.SetupRoutes(router, client)

	// 启动服务器
	log.Printf("Starting mail service on port %s\n", webPort)
	if err := router.Run(":" + webPort); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
