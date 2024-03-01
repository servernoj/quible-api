package v1

import (
	"regexp"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/auth-service/services/userService"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
)

// -- Authorization header containing Bearer access token. Injects `UserId` into `input` struct
type AuthorizationHeaderResolver struct {
	Authorization string `header:"authorization"`
	UserId        string
}

func (f *AuthorizationHeaderResolver) Resolve(ctx huma.Context) (errs []error) {
	re, _ := regexp.Compile(`\s+`)
	headerParts := re.Split(f.Authorization, -1)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "invalid format of the authorization header",
			Location: "header.authorization",
			Value:    f.Authorization,
		})
		return
	}
	token := headerParts[1]
	tokenClaims, err := jwt.VerifyJWT(token, jwt.TokenActionAccess)
	if err != nil {
		errs = append(errs, &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "header.authorization.bearer",
			Value:    token,
		})
		return
	}
	UserId := tokenClaims["userId"].(string)
	if exists, _ := models.UserExistsG(ctx.Context(), UserId); !exists {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "referenced user not found",
			Location: "header.authorization.bearer",
			Value:    token,
		})
		return
	}
	f.UserId = UserId
	return
}

// -- Password in request body. Injects `HashedPassword` field into `input` struct
type PasswordResolver struct {
	Password       string `json:"password" doc:"at least 6 characters long"`
	hashedPassword string
}

func (f *PasswordResolver) Resolve(ctx huma.Context) (errs []error) {
	hasSufficientLength := len(f.Password) >= 6
	if !hasSufficientLength {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "password has insufficient complexity",
			Location: "body.password",
			Value:    f.Password,
		})
	}
	if len(errs) == 0 {
		us := userService.UserService{}
		hash, err := us.HashPassword(f.Password)
		if err != nil {
			errs = append(errs, &huma.ErrorDetail{
				Message:  err.Error(),
				Location: "body.password",
				Value:    f.Password,
			})
		}
		f.hashedPassword = hash
	}
	return
}
