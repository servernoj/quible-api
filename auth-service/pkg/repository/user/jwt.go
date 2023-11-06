package user

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/quible-backend/auth-service/domain"
)

var APPLICATION_NAME = "Quible"
var LOGIN_EXPIRATION_DURATION = time.Duration(24) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("the secret of hogwarts")

type MyClaims struct {
	jwt.StandardClaims
	ID int64 `json:"id"`
	// Username string `json:"username"`
	Email string `json:"email"`
}

func generateToken(user *domain.UserLoginResponse) string {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		ID: user.ID,
		// Username: user.Username,
		Email: user.Email,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return ""
	}

	return signedToken
}

func verifyJWT(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return JWT_SIGNATURE_KEY, nil
	})
	if err != nil {
		return int64(0), "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return int64(0), "", err
	}

	return int64(claims["id"].(float64)), claims["email"].(string), nil
}
