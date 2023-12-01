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

type StatItem struct {
	Fouls                int `json:"fouls"`
	Blocks               int `json:"blocks"`
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

type PlayerStatItem struct {
	StatItem
	Minutes int64 `json:"minutes"`
}
type TeamStatItem struct {
	StatItem
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
}

type TeamSeasonStatItem struct {
	TeamID        int           `json:"team_id"`
	Team          string        `json:"team"`
	RegularSeason *TeamStatItem `json:"regular_season"`
	Postseason    *TeamStatItem `json:"postseason"`
}

type PlayerInfoItem struct {
	PlayerID         int             `json:"player_id"`
	Player           string          `json:"player"`
	TeamID           int             `json:"team_id"`
	Team             string          `json:"team"`
	Number           *int            `json:"number"`
	Status           *PlayerStatus   `json:"status"`
	Position         *PlayerPosition `json:"position"`
	PositionCategory *PlayerPosition `json:"position_category"`
	Height           *string         `json:"height"`
	Weight           *int            `json:"weight"`
	Age              string          `json:"age"`
	College          *string         `json:"college"`
}

type PlayerSeasonStatItem struct {
	PlayerID      int             `json:"player_id"`
	Player        string          `json:"player"`
	Team          string          `json:"team"`
	TeamID        int             `json:"team_id"`
	RegularSeason *PlayerStatItem `json:"regular_season"`
	Postseason    *PlayerStatItem `json:"postseason"`
}

type InjuryItem struct {
	Team     string             `json:"team"`
	TeamID   int                `json:"team_id"`
	Injuries []PlayerInjuryItem `json:"injuries"`
}

// TODO: the response from RSC API reports `player_id` as `string` instead of `int`
// we can fix it on our end by traversing entire slice, but we don't do it now...
type PlayerInjuryItem struct {
	Injury      string `json:"injury"`
	Player      string `json:"player"`
	Returns     string `json:"returns"`
	PlayerID    string `json:"player_id"`
	DateInjured string `json:"date_injured"`
}
