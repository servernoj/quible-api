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
type TeamInfoItemExtended struct {
	TeamInfoItem
	Color string  `json:"color"`
	Logo  *string `json:"logo"`
}

type BaseStatsItem struct {
	Fouls                int `json:"fouls"`
	Blocks               int `json:"blocks"`
	Steals               int `json:"steals"`
	Assists              int `json:"assists"`
	Turnovers            int `json:"turnovers"`
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

type PlayerStatsItem struct {
	BaseStatsItem
	Points      int   `json:"points"`
	GamesPlayed int   `json:"games_played"`
	Minutes     int64 `json:"minutes"`
}
type TeamStatsItem struct {
	BaseStatsItem
	Points      int `json:"points"`
	GamesPlayed int `json:"games_played"`
	Wins        int `json:"wins"`
	Losses      int `json:"losses"`
}

type TeamSeasonStatItem struct {
	TeamID        int            `json:"team_id"`
	Team          string         `json:"team"`
	RegularSeason *TeamStatsItem `json:"regular_season"`
	Postseason    *TeamStatsItem `json:"postseason"`
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
type PlayerInfoItemExtended struct {
	PlayerInfoItem
	Headshot *string `json:"headshot"`
}

type PlayerSeasonStatItem struct {
	PlayerID      int              `json:"player_id"`
	Player        string           `json:"player"`
	Team          string           `json:"team"`
	TeamID        int              `json:"team_id"`
	RegularSeason *PlayerStatsItem `json:"regular_season"`
	Postseason    *PlayerStatsItem `json:"postseason"`
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

// -- Live
type PlayerBox struct {
	AwayTeam map[string]PlayerBoxItem `json:"away_team"`
	HomeTeam map[string]PlayerBoxItem `json:"home_team"`
}

type Current struct {
	Quarter       Quarter `json:"Quarter"`
	TimeRemaining *string `json:"TimeRemaining"`
}

type FullBox struct {
	Current  Current     `json:"current"`
	AwayTeam TeamBoxItem `json:"away_team"`
	HomeTeam TeamBoxItem `json:"home_team"`
}

type PlayerBoxItem struct {
	BaseStatsItem
	Points             int             `json:"points"`
	TwoPointPercentage float64         `json:"two_point_percentage"`
	Player             *string         `json:"player"`
	Status             *PlayerStatus   `json:"status"`
	Minutes            *string         `json:"minutes"`
	Position           *PlayerPosition `json:"position"`
}

type TeamBoxStats struct {
	BaseStatsItem
	TwoPointPercentage float64 `json:"two_point_percentage"`
}

type QuarterScores struct {
	Q1 int  `json:"1"`
	Q2 int  `json:"2"`
	Q3 int  `json:"3"`
	Q4 int  `json:"4"`
	OT *int `json:"OT,omitempty"`
}

type TeamBoxItem struct {
	Abbrv         string        `json:"abbrv"`
	Score         int           `json:"score"`
	Mascot        string        `json:"mascot"`
	Record        string        `json:"record"`
	TeamID        int           `json:"team_id"`
	TeamStats     TeamBoxStats  `json:"team_stats"`
	DivisionName  string        `json:"division_name"`
	QuarterScores QuarterScores `json:"quarter_scores"`
}

type LiveFeedItem struct {
	Round        *any               `json:"round"`
	Sport        Sport              `json:"sport"`
	Season       string             `json:"season"`
	Status       ScheduleStatus     `json:"status"`
	GameID       string             `json:"game_ID"`
	FullBox      FullBox            `json:"full_box"`
	Broadcast    *string            `json:"broadcast"`
	GameTime     string             `json:"game_time"`
	EventName    *ScheduleEventName `json:"event_name"`
	PlayerBox    PlayerBox          `json:"player_box"`
	GameStatus   string             `json:"game_status"`
	SeasonType   ScheduleType       `json:"season_type"`
	GameLocation string             `json:"game_location"`
	AwayTeamName string             `json:"away_team_name"`
	HomeTeamName string             `json:"home_team_name"`
}
