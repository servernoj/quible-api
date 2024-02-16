package v1

import "github.com/quible-io/quible-api/app-service-ng/services/BasketAPI"

type GetGameDetailsDTO struct {
	GameId uint `form:"gameId" binding:"required,gt=0"`
}

type MatchDetails struct {
	ID         uint             `json:"id"`
	GameStatus string           `json:"gameStatus"`
	HomeScore  *uint            `json:"homeScore"`
	AwayScore  *uint            `json:"awayScore"`
	Date       string           `json:"date"`
	Event      *BasketAPI.Event `json:"-"`
}

type GameDetails struct {
	MatchDetails
	HomeTeam TeamInfoExtended `json:"homeTeam"`
	AwayTeam TeamInfoExtended `json:"awayTeam"`
}

type TeamInfoExtended struct {
	BasketAPI.TeamInfo
	Stats   *TeamStats     `json:"stats"`
	Players []PlayerEntity `json:"players"`
}

type GameTeamsStats struct {
	HomeTeam TeamStats
	AwayTeam TeamStats
}
type GamePlayers struct {
	HomeTeam []PlayerEntity
	AwayTeam []PlayerEntity
}

type TeamStats struct {
	Rebounds           uint `json:"reb"`
	Assists            uint `json:"ast"`
	Steals             uint `json:"stl"`
	Blocks             uint `json:"blk"`
	Turnovers          uint `json:"to"`
	Fouls              uint `json:"fp"`
	FieldGoalsMade     uint `json:"fgm"`
	FieldGoalAttempts  uint `json:"fga"`
	ThreePointsMade    uint `json:"tpm"`
	ThreePointAttempts uint `json:"tpa"`
	FreeThrowsMade     uint `json:"ftm"`
	FreeThrowAttempts  uint `json:"fta"`
}

type PlayerEntity struct {
	ID    uint
	Name  string
	Stats PlayerStats
}

type PlayerStats struct {
	MinutesPlayed      float64 `json:"min"`
	SecondsPlayed      uint    `json:"sec"`
	FieldGoalsMade     uint    `json:"fgm"`
	FieldGoalAttempts  uint    `json:"fga"`
	ThreePointsMade    uint    `json:"tpm"`
	ThreePointAttempts uint    `json:"tpa"`
	FreeThrowsMade     uint    `json:"ftm"`
	FreeThrowAttempts  uint    `json:"fta"`
	OffensiveRebounds  uint    `json:"oreb"`
	DefensiveRebounds  uint    `json:"dreb"`
	Rebounds           uint    `json:"reb"`
	Assists            uint    `json:"ast"`
	Steals             uint    `json:"stl"`
	Blocks             uint    `json:"blk"`
	Turnovers          uint    `json:"to"`
	PersonalFouls      uint    `json:"fp"`
	Points             uint    `json:"pts"`
}

// -- MatchLineups (ML) API

type ML_Data struct {
	Home ML_Team `json:"home"`
	Away ML_Team `json:"away"`
}

type ML_Team struct {
	Players []ML_PlayerElement `json:"players"`
}

type ML_PlayerElement struct {
	Player     ML_Player     `json:"player"`
	Statistics ML_Statistics `json:"statistics"`
}

type ML_Player struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ShortName    string `json:"shortName"`
	Position     string `json:"position"`
	JerseyNumber string `json:"jerseyNumber"`
	ID           uint   `json:"id"`
}

type ML_Statistics struct {
	SecondsPlayed      uint `json:"secondsPlayed"`
	FieldGoalsMade     uint `json:"fieldGoalsMade"`
	FieldGoalAttempts  uint `json:"fieldGoalAttempts"`
	ThreePointsMade    uint `json:"threePointsMade"`
	ThreePointAttempts uint `json:"threePointAttempts"`
	FreeThrowsMade     uint `json:"freeThrowsMade"`
	FreeThrowAttempts  uint `json:"freeThrowAttempts"`
	OffensiveRebounds  uint `json:"offensiveRebounds"`
	DefensiveRebounds  uint `json:"defensiveRebounds"`
	Rebounds           uint `json:"rebounds"`
	Assists            uint `json:"assists"`
	Steals             uint `json:"steals"`
	Blocks             uint `json:"blocks"`
	Turnovers          uint `json:"turnovers"`
	PersonalFouls      uint `json:"personalFouls"`
	Points             uint `json:"points"`
	// -- currently unused
	// TwoPointsMade    uint `json:"twoPointsMade"`
	// TwoPointAttempts uint `json:"twoPointAttempts"`
	// PlusMinus        uint `json:"plusMinus"`
}

// -- Match (MD) API

type MD_Data struct {
	Event BasketAPI.Event `json:"event"`
}

// -- MatchStatistics (MStat) API

type MStat_Data struct {
	Statistics []MStat_StatEntry `json:"statistics"`
}

type MStat_StatEntry struct {
	Period string        `json:"period"`
	Groups []MStat_Group `json:"groups"`
}

type MStat_Group struct {
	GroupName       MStat_GroupName   `json:"groupName"`
	StatisticsItems []MStat_GroupItem `json:"statisticsItems"`
}

type MStat_GroupName string

const (
	MStat_GroupName_Lead    MStat_GroupName = "Lead"
	MStat_GroupName_Other   MStat_GroupName = "Other"
	MStat_GroupName_Scoring MStat_GroupName = "Scoring"
)

type MStat_GroupItem struct {
	Name      MStat_GroupItemName `json:"name"`
	HomeValue uint                `json:"homeValue"`
	AwayValue uint                `json:"awayValue"`
	HomeTotal *uint               `json:"homeTotal,omitempty"`
	AwayTotal *uint               `json:"awayTotal,omitempty"`
}

type MStat_GroupItemName string

const (
	// Other
	MStat_GroupItemName_OtherRebounds  MStat_GroupItemName = "Rebounds"
	MStat_GroupItemName_OtherAssists   MStat_GroupItemName = "Assists"
	MStat_GroupItemName_OtherTurnovers MStat_GroupItemName = "Turnovers"
	MStat_GroupItemName_OtherSteals    MStat_GroupItemName = "Steals"
	MStat_GroupItemName_OtherBlocks    MStat_GroupItemName = "Blocks"
	MStat_GroupItemName_OtherFouls     MStat_GroupItemName = "Fouls"
	// Scoring
	MStat_GroupItemName_ScoringFreeThrows  MStat_GroupItemName = "Free throws"
	MStat_GroupItemName_ScoringTwoPoints   MStat_GroupItemName = "2 pointers"
	MStat_GroupItemName_ScoringThreePoints MStat_GroupItemName = "3 pointers"
	MStat_GroupItemName_ScoringFieldGoals  MStat_GroupItemName = "Field goals"
)
