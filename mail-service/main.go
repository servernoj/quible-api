package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/quible-backend/mail-service/service"
)

func main() {

	// create Postmark client
	client := service.Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		ServerToken:  "your-server-token",           // your Postmark token
		AccountToken: "your-account-token",          // your Postmark account
		BaseURL:      "https://api.postmarkapp.com", // Postmark API URL
	}

	// define the email to be sent
	email := service.Email{
		From:       "no-reply@example.com",
		To:         "tito@example.com",
		Subject:    "Reset your password",
		HTMLBody:   "<p>Your password reset link is here.</p>",
		TextBody:   "Your password reset link is here.",
		Tag:        "password-reset",
		TrackOpens: true,
	}

	// send the email
	response, err := client.SendEmail(context.Background(), email)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}

	// print the response
	fmt.Printf("Email sent! Response: %+v\n", response)
}
