package controller

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gitlab.com/quible-backend/lib/models"
)

var APPLICATION_NAME = "Quible"
var ACCESS_TOKEN_DURATION = 4 * time.Hour
var REFRESH_TOKEN_DURATION = 5 * 24 * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type MyClaims struct {
	jwt.StandardClaims
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"isRefresh"`
}

func generateToken(user *models.User, isRefresh bool) (string, error) {

	tokenLifespan := ACCESS_TOKEN_DURATION
	if isRefresh {
		tokenLifespan = REFRESH_TOKEN_DURATION
	}
	standardClaims := jwt.StandardClaims{
		Id:        uuid.New().String(),
		Issuer:    APPLICATION_NAME,
		ExpiresAt: time.Now().Add(tokenLifespan).Unix(),
	}

	claims := MyClaims{
		StandardClaims: standardClaims,
		UserId:         user.ID,
		Email:          user.Email,
		IsRefresh:      isRefresh,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)
	signedToken, err := token.SignedString([]byte(os.Getenv("ENV_JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func verifyJWT(tokenString string, isRefresh bool) (string, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != JWT_SIGNING_METHOD {
				return nil, ErrTokenInvalidSigningMethod
			}
			mapClaims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, ErrTokenInvalidClaims
			}
			if !mapClaims.VerifyExpiresAt(time.Now().Unix(), true) {
				return nil, ErrTokenExpired
			}
			if IsRefresh, ok := mapClaims["isRefresh"].(bool); !ok || IsRefresh != isRefresh {
				return nil, ErrTokenInvalidType
			}
			if _, ok := mapClaims["userId"].(string); !ok {
				return "", ErrTokenMissingUserId
			}
			return []byte(os.Getenv("ENV_JWT_SECRET")), nil
		},
	)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	UserId := claims["userId"].(string)

	return UserId, nil
}
