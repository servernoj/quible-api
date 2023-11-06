package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/quible-backend/auth-service/domain"
)

type controller struct {
	Service Impl
	*domain.Response
}

func NewController(svc Impl) controller {
	return controller{Service: svc}
}

// ShowAccount godoc
//
//	@Summary		Check email
//	@Description	Check if an email address can be used to register.
//	@Description	WARNING: for a public API, it is generally considered as a security concern to allow users to query email validity.
//	@Description	We may want to revisit this API - https://quible.atlassian.net/browse/SPORT-66
//	@Tags			email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserEmailCheckRequest	true	"Email to be checked"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/checkemail/{email} [post]
func (h *controller) CheckEmail(c *gin.Context) {
	var (
		validate   = validator.New()
		checkemail domain.UserEmailCheckRequest
	)

	if err := c.ShouldBindJSON(&checkemail); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&checkemail); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	email := strings.ToLower(checkemail.Email)

	user, err := h.Service.GetUserByEmail(email)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to check email", err)
		return
	}

	if user == nil {
		h.SendSuccess(c, "Email doesn't exist, You can register now", nil)
		return
	}

	h.SendError(c, http.StatusFound, "Email exists, You can't register with this email", nil)
}

func (h *controller) Upload(c *gin.Context) {

	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		h.SendError(c, http.StatusBadRequest, "Image size should be less than 10M", err)
		return
	}

	file, _ := c.FormFile("file")

	// get the root path of the project.
	rootPath, err := os.Getwd()
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed File Upload", err)
		return
	}

	filePath := filepath.Join(rootPath, "images", file.Filename)

	// Upload the file to specific dst.
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed File Upload", err)
		return
	}
	h.SendSuccess(c, "Success file upload", nil)
}

// ShowAccount godoc
//
//	@Summary		Verify email
//	@Description	Try to verify email address by sending out a verification email.
//	@Tags			email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserEmailCheckRequest	true	"Email"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/verify/{email} [post]
func (h *controller) Verify(c *gin.Context) {

	var (
		validate = validator.New()
		verify   domain.UserEmailCheckRequest
	)

	if err := c.ShouldBindJSON(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Generate 4-digit code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(10000)
	codeStr := strconv.Itoa(code)

	// Set up email message
	from := os.Getenv("ENV_EMAIL_ADDRESS")
	password := os.Getenv("ENV_EMAIL_PASSWORD")
	to := verify.Email
	subject := "Your verification code"
	body := "Your verification code is: " + codeStr
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	// Send email
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(message))

	if err != nil {
		h.SendError(c, http.StatusNotAcceptable, "Error sending email", err)
		return
	}

	store := NewInMemoryCodeStore()

	err = store.SaveCode(fmt.Sprintf("%s", to), code)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Error saving code", err)
		return
	}

	fmt.Println("Verification code sent to", to)

	h.SendSuccess(c, "Success sending email", nil)
}

// ShowAccount godoc
//
//	@Summary		Verify email
//	@Description	Try to verify email address by sending out a verification email.
//	@Tags			email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserEmailCheckRequest	true	"Email to verify"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/verify/{email} [post]
func (h *controller) VerifyCode(c *gin.Context) {
	var (
		validate = validator.New()
		verify   domain.UserVerifyCodeRequest
	)

	if err := c.ShouldBindJSON(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	to := verify.Email

	store := NewInMemoryCodeStore()

	code, err := store.GetCode(to)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Verification failed", err)
		return
	}

	if code != verify.Code {
		h.SendError(c, http.StatusNotAcceptable, "Verification code not correct", err)
		return
	}

	fmt.Println("Verification success")
	h.SendSuccess(c, "Verification success", nil)

}

// ShowAccount godoc
//
//	@Summary		Reset password
//	@Description	TODO this API is under development.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserVerifyCodeRequest	true	"TODO"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/resetpassword/{email} [post]
func (h *controller) ResetPassowrd(c *gin.Context) {
	var (
		validate = validator.New()
		verify   domain.UserVerifyCodeRequest
	)

	if err := c.ShouldBindJSON(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	to := verify.Email

	store := NewInMemoryCodeStore()

	code, err := store.GetCode(to)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Reset password failed", err)
		return
	}

	if code != verify.Code {
		h.SendError(c, http.StatusNotAcceptable, "Verification code not correct", err)
		return
	}

	fmt.Println("Password reset success")
	h.SendSuccess(c, "Password reset success", nil)

}

// ShowAccount godoc
//
//	@Summary		Validate email
//	@Description	TODO is it the same as validate? Try to verify email address by sending out a verification email.
//	@Tags			email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserEmailCheckRequest	true	"Email"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/verify/{email} [post]
func (h *controller) R(c *gin.Context) {

	var (
		validate = validator.New()
		verify   domain.UserEmailCheckRequest
	)

	if err := c.ShouldBindJSON(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&verify); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Generate 4-digit code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(10000)
	codeStr := strconv.Itoa(code)

	// Set up email message
	from := os.Getenv("ENV_EMAIL_ADDRESS")
	password := os.Getenv("ENV_EMAIL_PASSWORD")
	to := verify.Email
	subject := "Your verification code"
	body := "Your verification code is: " + codeStr
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)

	// Send email
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(message))

	if err != nil {
		h.SendError(c, http.StatusNotAcceptable, "Error sending email", err)
		return
	}

	store := NewInMemoryCodeStore()

	err = store.SaveCode(fmt.Sprintf("%s", to), code)

	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Error saving code", err)
		return
	}

	fmt.Println("Verification code sent to", to)

	h.SendSuccess(c, "Success sending email", nil)
}

// ShowAccount godoc
//
//	@Summary		Register
//	@Description	Register a new user.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserRegisterRequest	true	"User registration information"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/register/{email} [post]
func (h *controller) Register(c *gin.Context) {
	var (
		validate = validator.New()
		user     domain.UserRegisterRequest
	)
	if err := c.ShouldBindJSON(&user); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&user); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if user, _ := h.Service.GetUserByEmail(user.Email); user != nil {
		h.SendError(c, http.StatusBadRequest, "Email exists, You can't register with this email", nil)
		return
	}

	if user, err := h.Service.GetByUsername(user.Username); user != nil {
		h.SendError(c, http.StatusBadRequest, "Username exists, You can't register with this username", err)
		return
	}

	id, err := h.Service.Create(user)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	h.SendSuccess(c, "Success add new user", id)
}

// ShowAccount godoc
//
//	@Summary		Login
//	@Description	User login.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserLoginRequest	true	"User login information"
//	@Success		200		{object}	domain.Response
//	@Failure		400		{object}	domain.Response
//	@Router			/login/{email} [post]
func (h *controller) Login(c *gin.Context) {
	var (
		validate = validator.New()
		req      domain.UserLoginRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.Service.GetLoginCredential(req.Email)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	if user == nil {
		h.SendError(c, http.StatusInternalServerError, "Email not found", err)
		return
	}

	if err := h.Service.ValidatePassword(user.HashedPassword, req.Password); err != nil {
		h.SendError(c, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	h.SendSuccess(c, "Sign in success", generateToken(user))
}

// ShowAccount godoc
//
//	@Summary		Get user
//	@Description	Get the profile of the user currently logged in.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.Response
//	@Router			/user [get]
func (h *controller) Get(c *gin.Context) {
	user := h.currentUser(c)

	h.SendSuccess(c, "Success get user", user)

	// TODO: SendError if no user is logged in.
}

// ShowAccount godoc
//
//	@Summary		Update user
//	@Description	Update the user that is currently logged in.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.UserUpdateRequest	true	"User update information"
//	@Success		200		{object}	domain.Response
//	@Success		400		{object}	domain.Response
//	@Router			/user [put]
func (h *controller) Update(c *gin.Context) {
	var (
		validate = validator.New()
		u        domain.UserUpdateRequest
		user     = h.currentUser(c)
	)

	if err := c.ShouldBindJSON(&u); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := validate.Struct(&u); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userid, err := h.Service.Update(user.ID, u)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	h.SendSuccess(c, "Success update user", userid)
}

// ShowAccount godoc
//
//	@Summary		Delete user
//	@Description	Delete the user that is currently logged in.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.Response
//	@Success		400	{object}	domain.Response
//	@Router			/user [delete]
func (h *controller) Delete(c *gin.Context) {
	var (
		user = h.currentUser(c)
	)

	u, err := h.Service.Delete(user.ID)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	h.SendSuccess(c, "Success delete user", u)
}
