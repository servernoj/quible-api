package service

type UserRegisterDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type UserPatchDTO struct {
	Username *string `json:"username" binding:""`
	Email    *string `json:"email" binding:"omitempty,email"`
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
}

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
