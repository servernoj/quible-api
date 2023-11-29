package RSC

type ScheduleSeasonItem struct {
	AwayTeam   string                   `json:"away_team"`
	HomeTeam   string                   `json:"home_team"`
	AwayTeamID int                      `json:"away_team_ID"`
	HomeTeamID int                      `json:"home_team_ID"`
	GameID     string                   `json:"game_ID"`
	GameTime   string                   `json:"game_time"`
	SeasonType ScheduleSeasonType       `json:"season_type"`
	EventName  *ScheduleSeasonEventName `json:"event_name"`
	Round      *ScheduleSeasonRound     `json:"round"`
	Season     string                   `json:"season"`
	Status     ScheduleSeasonStatus     `json:"status"`
	Broadcast  *string                  `json:"broadcast"`
}
