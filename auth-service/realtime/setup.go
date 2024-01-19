package realtime

import (
	"encoding/json"
	"os"

	"github.com/ably/ably-go/ably"
)

var (
	ablyRealTime *ably.Realtime
)

func Setup() error {
	client, err := ably.NewRealtime(
		ably.WithKey(os.Getenv("ENV_ABLY_KEY")),
	)
	if err != nil {
		return err
	}
	ablyRealTime = client
	return nil
}

func GetToken(userId string) (*ably.TokenRequest, error) {
	capabilities, _ := json.Marshal(&map[string][]string{
		"chat:*": {"*"},
		"live:*": {"subscribe", "history"},
	})
	tokenParams := &ably.TokenParams{
		ClientID:   userId,
		Capability: string(capabilities),
	}
	return ablyRealTime.Auth.CreateTokenRequest(tokenParams)
}

func CreateTokenRequest(params *ably.TokenParams, opts ...ably.AuthOption) (*ably.TokenRequest, error) {
	return ablyRealTime.Auth.CreateTokenRequest(params, opts...)
}
