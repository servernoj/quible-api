package controller

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gitlab.com/quible-backend/lib/models"
)

var APPLICATION_NAME = "Quible"
var ACCESS_TOKEN_DURATION = 30 * time.Second
var REFRESH_TOKEN_DURATION = 5 * 24 * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type MyClaims struct {
	jwt.StandardClaims
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"isRefresh"`
}

type GeneratedToken struct {
	Token string
	ID    string
}

func (gt *GeneratedToken) String() string {
	return gt.Token
}

func generateToken(user *models.User, isRefresh bool) (GeneratedToken, error) {

	tokenLifespan := ACCESS_TOKEN_DURATION
	tokenId := uuid.New().String()
	if isRefresh {
		tokenLifespan = REFRESH_TOKEN_DURATION
	}
	standardClaims := jwt.StandardClaims{
		Id:        tokenId,
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
		return GeneratedToken{}, err
	}

	return GeneratedToken{
		Token: signedToken,
		ID:    tokenId,
	}, nil
}

func verifyJWT(tokenString string, isRefresh bool) (jwt.MapClaims, error) {
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
			if _, ok := mapClaims["jti"].(string); !ok {
				return "", ErrTokenMissingTokenId
			}
			return []byte(os.Getenv("ENV_JWT_SECRET")), nil
		},
	)
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
