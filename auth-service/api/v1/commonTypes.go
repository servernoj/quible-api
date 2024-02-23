package v1

type UserSimplified struct {
	ID       string `json:"id" doc:"user ID (UUID)"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

type UserProfile struct {
	ID       string  `json:"id" doc:"user ID (UUID)"`
	FullName string  `json:"full_name"`
	Image    *string `json:"image" doc:"Profile image data URL"`
}

type ImageData struct {
	ContentType   string `json:"contentType"`
	BinaryContent []byte `json:"data"`
}
