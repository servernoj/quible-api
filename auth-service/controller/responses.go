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
