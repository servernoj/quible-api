package RSC

type ScheduleItem struct {
	AwayTeam   string             `json:"away_team"`
	HomeTeam   string             `json:"home_team"`
	AwayTeamID int                `json:"away_team_ID"`
	HomeTeamID int                `json:"home_team_ID"`
	GameID     string             `json:"game_ID"`
	GameTime   string             `json:"game_time"`
	SeasonType ScheduleType       `json:"season_type"`
	EventName  *ScheduleEventName `json:"event_name"`
	Round      *ScheduleRound     `json:"round"`
	Season     string             `json:"season"`
	Status     ScheduleStatus     `json:"status"`
	Broadcast  *string            `json:"broadcast"`
}

type TeamInfoItem struct {
	TeamID   int    `json:"team_id"`
	Team     string `json:"team"`
	Abbrv    string `json:"abbrv"`
	Arena    string `json:"arena"`
	Mascot   string `json:"mascot"`
	Conf     string `json:"conf"`
	Location string `json:"location"`
}
