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

type TeamStatItem struct {
	WINS                 int `json:"wins"`
	Fouls                int `json:"fouls"`
	Blocks               int `json:"blocks"`
	Losses               int `json:"losses"`
	Points               int `json:"points"`
	Steals               int `json:"steals"`
	Assists              int `json:"assists"`
	Turnovers            int `json:"turnovers"`
	GamesPlayed          int `json:"games_played"`
	TotalRebounds        int `json:"total_rebounds"`
	TwoPointsMade        int `json:"two_points_made"`
	FieldGoalsMade       int `json:"field_goals_made"`
	FreeThrowsMade       int `json:"free_throws_made"`
	ThreePointsMade      int `json:"three_points_made"`
	DefensiveRebounds    int `json:"defensive_rebounds"`
	OffensiveRebounds    int `json:"offensive_rebounds"`
	TwoPointsAttempted   int `json:"two_points_attempted"`
	FieldGoalsAttempted  int `json:"field_goals_attempted"`
	FreeThrowsAttempted  int `json:"free_throws_attempted"`
	ThreePointsAttempted int `json:"three_points_attempted"`
}

type TeamSeasonStatItem struct {
	TeamID        int64         `json:"team_id"`
	Team          string        `json:"team"`
	RegularSeason *TeamStatItem `json:"regular_season"`
	Postseason    *TeamStatItem `json:"postseason"`
}
