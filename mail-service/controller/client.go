package controller

import (
	"net/http"
	"os"
	"time"

	"gitlab.com/quible-backend/mail-service/service"
)

func NewClient() *service.Client {
	serverToken := os.Getenv("ENV_SERVER_TOKEN")
	accountToken := os.Getenv("ENV_ACCOUNT_TOKEN")

	return &service.Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      "https://api.postmarkapp.com",
	}
}
