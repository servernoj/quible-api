package BasketAPI

type GetGameDetailsDTO struct {
	GameId uint `form:"gameId" binding:"required,gt=0"`
}

type MatchDetails struct {
	ID         uint      `json:"id"`
	GameStatus string    `json:"gameStatus"`
	HomeScore  *uint     `json:"homeScore"`
	AwayScore  *uint     `json:"awayScore"`
	Date       string    `json:"date"`
	Event      *MD_Event `json:"-"`
}

type GameDetails struct {
	MatchDetails
	HomeTeam TeamInfoExtended `json:"homeTeam"`
	AwayTeam TeamInfoExtended `json:"awayTeam"`
}

type TeamInfoExtended struct {
	TeamInfo
	Stats   TeamStats      `json:"stats"`
	Players []PlayerEntity `json:"players"`
}

type GameTeamsStats struct {
	HomeTeam TeamStats
	AwayTeam TeamStats
}

type TeamStats struct {
	Rebounds  uint `json:"reb"`
	Assists   uint `json:"ast"`
	Steals    uint `json:"stl"`
	Blocks    uint `json:"blk"`
	Turnovers uint `json:"to"`
	Fouls     uint `json:"fp"`
}

type PlayerEntity struct {
	ID   uint
	Name string
	// Stats PlayerStats
}

type PlayerStats struct {
	min  uint // minutes
	fgm  uint // fieldGoalsMade
	fga  uint // fieldGoalsAttempted
	tpm  uint // threePointsMade
	tpa  uint // threePointsAttempted
	ftm  uint // freeThrowsMade
	fta  uint // freeThrowsAttempted
	oreb uint // offensiveRebounds
	dreb uint // defensiveRebounds
	reb  uint // totalRebounds
	ast  uint // assists
	stl  uint // steals
	blk  uint // blocks
	to   uint // turnovers
	// fp   uint // fouls
}

// -- Match (MD) API

type MD_Data struct {
	Event MD_Event `json:"event"`
}

type MD_Event = MS_Event

const (
	MD_StatusType_Finished      = MS_StatusType_Finished
	MD_StatusType_Inprogress    = MS_StatusType_Inprogress
	MD_StatusType_Notstarted    = MS_StatusType_Notstarted
	MD_StatusType_TypePostponed = MS_StatusType_TypePostponed
)

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
}

type MStat_GroupItemName string

const (
	MStat_GroupItemName_OtherRebounds  MStat_GroupItemName = "Rebounds"
	MStat_GroupItemName_OtherAssists   MStat_GroupItemName = "Assists"
	MStat_GroupItemName_OtherTurnovers MStat_GroupItemName = "Turnovers"
	MStat_GroupItemName_OtherSteals    MStat_GroupItemName = "Steals"
	MStat_GroupItemName_OtherBlocks    MStat_GroupItemName = "Blocks"
	MStat_GroupItemName_OtherFouls     MStat_GroupItemName = "Fouls"
)
