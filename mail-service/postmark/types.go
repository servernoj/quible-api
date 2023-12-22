package postmark

// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
type PostmarkResponse struct {
	To          string `json:"To"`
	SubmittedAt string `json:"SubmittedAt"`
	MessageID   string `json:"MessageID"`
	ErrorCode   int    `json:"ErrorCode"`
	Message     string `json:"Message"`
}

func (r *PostmarkResponse) String() string {
	return r.Message
}

type PostmarkPayload interface {
	EmailDTO
}

type RequestParams[T PostmarkPayload] struct {
	Method  string
	Path    string
	Payload *T
}
