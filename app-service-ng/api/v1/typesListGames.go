package v1

import "github.com/quible-io/quible-api/app-service-ng/services/BasketAPI"

type Game struct {
	ID         uint               `json:"id"`
	GameStatus string             `json:"gameStatus"`
	HomeTeam   BasketAPI.TeamInfo `json:"homeTeam"`
	AwayTeam   BasketAPI.TeamInfo `json:"awayTeam"`
	HomeScore  *uint              `json:"homeScore"`
	AwayScore  *uint              `json:"awayScore"`
	Date       string             `json:"date"`
}

// -- MatchSchedules (MS) API
type MS_Data struct {
	Events []BasketAPI.Event `json:"events"`
}
