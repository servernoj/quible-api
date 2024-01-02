package userService

type UserRegisterDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required,phone"`
}

type UserPatchDTO struct {
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"omitempty,email"`
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone" binding:"omitempty,phone"`
}

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type UserRefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserRequestNewPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type UserResetPasswordDTO struct {
	Password        string `form:"password" binding:"required,min=6,eqfield=ConfirmPassword"`
	ConfirmPassword string `form:"confirm-password"`
}
