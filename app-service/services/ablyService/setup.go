package ablyService

import (
	"os"

	"github.com/ably/ably-go/ably"
)

var (
	ablyRealTime *ably.Realtime
)

func Setup() error {
	client, err := ably.NewRealtime(
		ably.WithKey(os.Getenv("ENV_ABLY_KEY")),
		ably.WithClientID("backend"),
	)
	if err != nil {
		return err
	}
	ablyRealTime = client
	return nil
}

func GetAbly() *ably.Realtime {
	return ablyRealTime
}

func CreateTokenRequest(params *ably.TokenParams, opts ...ably.AuthOption) (*ably.TokenRequest, error) {
	return ablyRealTime.Auth.CreateTokenRequest(params, opts...)
}
