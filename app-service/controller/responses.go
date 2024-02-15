package controller

type AblyTokenRequest struct {
	TTL        int64  `json:"ttl"`
	Capability string `json:"capability"`
	ClientID   string `json:"clientId"`
	Timestamp  int64  `json:"timestamp"`
	KeyName    string `json:"keyName"`
	Nonce      string `json:"nonce"`
	MAC        string `json:"mac"`
}
