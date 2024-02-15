package jwt

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/stretchr/testify/assert"
)

// unit test for the jwt.go generateToken test
func TestGenerateToken(t *testing.T) {
	// set the ENV_JWT_SECRET
	os.Setenv("ENV_JWT_SECRET", "your_test_jwt_secret")

	// define the test cases struct
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
			token, _ := GenerateToken(tc.user, TokenActionAccess, nil)
			assert.NotEmpty(t, token.Token, "Token should not be empty")

			// The generated token is parsed and validated using jwt.ParseWithClaims.
			//This ensures the token is correctly structured and valid.
			parsedToken, err := jwt.ParseWithClaims(token.Token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("your_test_jwt_secret"), nil
			})
			assert.NoError(t, err, "Token should be valid")

			// Claims Verification
			if claims, ok := parsedToken.Claims.(*MyClaims); ok && parsedToken.Valid {
				assert.Equal(t, claims.UserId, tc.user.ID)
			} else {
				t.Fatal("Failed to parse token claims or token is not valid")
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	// set the ENV_JWT_SECRET
	os.Setenv("ENV_JWT_SECRET", "your_test_jwt_secret")

	// define the test cases
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
			var generatedToken GeneratedToken
			if tc.user != nil {
				// generate useful token
				generatedToken, _ = GenerateToken(tc.user, TokenActionAccess, nil)
			} else {
				// generate unuseful token
				generatedToken = GeneratedToken{
					Token: "invalid.token.string",
					ID:    "invalid.id",
				}
			}
			// run verifyJWT function
			claims, err := VerifyJWT(generatedToken.Token, TokenActionAccess)

			if tc.expectError {
				// assert should have faults
				assert.Error(t, err, "Expected an error for invalid token")
			} else {
				// assert should success, have no faults
				assert.NoError(t, err, "No error expected for valid token")
				userId := claims["userId"].(string)
				assert.Contains(t, tc.user.ID, userId, "userId from claims should match the user ID")
			}
		})
	}
}
