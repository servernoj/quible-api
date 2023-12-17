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
	})
	tokenParams := &ably.TokenParams{
		ClientID:   userId,
		Capability: string(capabilities),
	}
	return ablyRealTime.Auth.CreateTokenRequest(tokenParams)
}
