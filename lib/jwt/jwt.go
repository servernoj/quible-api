package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/quible-io/quible-api/lib/models"
)

var APPLICATION_NAME = "Quible"
var DEFAULT_TOKEN_DURATION = 24 * time.Hour
var REFRESH_TOKEN_DURATION = 10 * 24 * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type TokenAction string

const (
	TokenActionAccess                  TokenAction = "Access"
	TokenActionRefresh                 TokenAction = "Refresh"
	TokenActionActivate                TokenAction = "Activate"
	TokenActionPasswordReset           TokenAction = "PasswordReset"
	TokenActionInvitationToPrivateChat TokenAction = "InvitationToPrivateChat"
)

type ExtraClaims = map[string]any

type MyClaims struct {
	jwt.StandardClaims
	UserId      string      `json:"userId"`
	Action      TokenAction `json:"action"`
	ExtraClaims ExtraClaims `json:"extraClaims"`
}

type GeneratedToken struct {
	Token string
	ID    string
}

func (gt *GeneratedToken) String() string {
	return gt.Token
}

func GenerateToken(user *models.User, action TokenAction, extraClaims ExtraClaims) (GeneratedToken, error) {

	tokenId := uuid.New().String()
	var tokenLifespan time.Duration
	switch action {
	case TokenActionRefresh:
		tokenLifespan = REFRESH_TOKEN_DURATION
	default:
		tokenLifespan = DEFAULT_TOKEN_DURATION
	}

	var claims MyClaims = MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenLifespan).Unix(),
		},
		UserId:      user.ID,
		Action:      action,
		ExtraClaims: extraClaims,
	}
	if action == TokenActionAccess || action == TokenActionRefresh {
		claims.StandardClaims.Id = tokenId
		claims.StandardClaims.Issuer = APPLICATION_NAME
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

func VerifyJWT(tokenString string, action TokenAction) (jwt.MapClaims, error) {
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
			if Action, ok := mapClaims["action"].(string); !ok || action != TokenAction(Action) {
				return nil, ErrTokenInvalidType
			}
			if _, ok := mapClaims["userId"].(string); !ok {
				return nil, ErrTokenMissingUserId
			}
			if _, ok := mapClaims["extraClaims"].(ExtraClaims); !ok {
				return nil, ErrTokenMissingExtraClaims
			}
			if action == TokenActionAccess || action == TokenActionRefresh {
				if _, ok := mapClaims["jti"].(string); !ok {
					return nil, ErrTokenMissingTokenId
				}
			}
			return []byte(os.Getenv("ENV_JWT_SECRET")), nil
		},
	)
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
