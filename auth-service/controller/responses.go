package controller

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
type PublicUserRecord struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Image    *string `json:"image"`
}

type AblyTokenRequest struct {
	TTL        int64  `json:"ttl"`
	Capability string `json:"capability"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp"`
	KeyName    string `json:"keyName"`
	Nonce      string `json:"nonce"`
	MAC        string `json:"mac"`
}
