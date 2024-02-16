package BasketAPI

type LiveMessage struct {
	IDs    []uint      `json:"eventIDs"`
	Events []LiveEvent `json:"events"`
}

type LiveEvent struct {
	ID             uint     `json:"id"`
	Status         Status   `json:"status"`
	HomeTeam       TeamInfo `json:"homeTeam"`
	AwayTeam       TeamInfo `json:"awayTeam"`
	HomeScore      Score    `json:"homeScore"`
	AwayScore      Score    `json:"awayScore"`
	Time           Time     `json:"time"`
	StartTimestamp int64    `json:"startTimestamp"`
}

// -- LiveMatches (LM) API

type LM_Data struct {
	Events []Event `json:"events"`
}
