package controller

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
