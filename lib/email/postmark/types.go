package postmark

// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
type Response struct {
	To          string `json:"To"`
	SubmittedAt string `json:"SubmittedAt"`
	MessageID   string `json:"MessageID"`
	ErrorCode   int    `json:"ErrorCode"`
	Message     string `json:"Message"`
}

func (r Response) Error() string {
	return r.Message
}
