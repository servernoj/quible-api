package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.com/quible-backend/lib/models"
)

var APPLICATION_NAME = "Quible"
var LOGIN_EXPIRATION_DURATION = time.Duration(24) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type MyClaims struct {
	jwt.StandardClaims
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func generateToken(user *models.User) string {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		ID:    user.ID,
		Email: user.Email,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)
	signedToken, err := token.SignedString([]byte(os.Getenv("ENV_JWT_SECRET")))
	if err != nil {
		log.Printf("unable to generate token: %q", err)
		return ""
	}

	return signedToken
}

func verifyJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(os.Getenv("ENV_JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, err
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unable to extract ID from token")
	}
	return int(id), nil
}
