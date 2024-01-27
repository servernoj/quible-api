package jwt

import "errors"

var (
	ErrTokenExpired              = errors.New("token expired")
	ErrTokenInvalidClaims        = errors.New("unable to process token claims")
	ErrTokenInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenInvalidType          = errors.New("invalid token type")
	ErrTokenMissingUserId        = errors.New("unable to extract userId from token")
	ErrTokenMissingTokenId       = errors.New("unable to extract tokenId from token")
	ErrTokenMissingExtraClaims   = errors.New("unable to extract extraClaims from token")
)
