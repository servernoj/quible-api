package controller

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/quible-backend/lib/models"
)

// unitest for the jwt.go generateToken test
func TestGenerateToken(t *testing.T) {
	// set the ENV_JWT_SECRET
	os.Setenv("ENV_JWT_SECRET", "your_test_jwt_secret")

	// define the tesycases struct
	testCases := []struct {
		name string
		user *models.User
	}{
		{
			name: "ValidUser",
			user: &models.User{ID: "user1", Email: "user1@example.com"},
		},
		{
			name: "EmptyUser",
			user: &models.User{},
		},
		{
			name: "IDnilUser",
			user: &models.User{ID: "", Email: "user2@example.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := generateToken(tc.user)
			assert.NotEmpty(t, token, "Token should not be empty")

			// The generated token is parsed and validated using jwt.ParseWithClaims.
			//This ensures the token is correctly structured and valid.
			parsedToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("your_test_jwt_secret"), nil
			})
			assert.NoError(t, err, "Token should be valid")

			// Claims Verification
			if claims, ok := parsedToken.Claims.(*MyClaims); ok && parsedToken.Valid {
				expectedClaims := MyClaims{
					StandardClaims: jwt.StandardClaims{
						Issuer:    APPLICATION_NAME,
						ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
					},
					ID:    tc.user.ID,
					Email: tc.user.Email,
				}
				// use reflect.DeepEqual
				assert.True(t, reflect.DeepEqual(expectedClaims, *claims), "Claims should match expected values")
			} else {
				t.Fatal("Failed to parse token claims or token is not valid")
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	// set the ENV_JWT_SECRET
	os.Setenv("ENV_JWT_SECRET", "your_test_jwt_secret")

	// define the testcases
	testCases := []struct {
		name        string
		user        *models.User
		expectError bool
	}{
		{
			name:        "ValidToken",
			user:        &models.User{ID: "user1", Email: "user1@example.com"},
			expectError: false,
		},
		{
			name:        "InvalidToken",
			user:        nil,
			expectError: true,
		},
		{
			name:        "TokenWithoutEmail",
			user:        &models.User{ID: "user2", Email: " "},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var tokenString string
			if tc.user != nil {
				// generate useful token
				tokenString = generateToken(tc.user)
			} else {
				// generate unuseful token
				tokenString = "invalid.token.string"
			}

			// run verifyJWT function
			ID, err := verifyJWT(tokenString)

			if tc.expectError {
				// assert should have faults
				assert.Error(t, err, "Expected an error for invalid token")
			} else {
				// assert should success, have no faults
				assert.NoError(t, err, "No error expected for valid token")
				assert.Equal(t, tc.user.ID, ID, "ID should match the user ID")
			}
		})
	}
}
