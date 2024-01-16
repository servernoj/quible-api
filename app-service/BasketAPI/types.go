package BasketAPI

// Data structure to represent a subset of fields from BasketAPI /live response
type LiveData struct {
	Events []LiveEvent `json:"events"`
}

type LiveEvent struct {
	Tournament     LiveTournament `json:"tournament"`
	Season         LiveSeason     `json:"season"`
	Status         LiveStatus     `json:"status"`
	HomeTeam       LiveTeam       `json:"homeTeam"`
	AwayTeam       LiveTeam       `json:"awayTeam"`
	HomeScore      LiveScore      `json:"homeScore"`
	AwayScore      LiveScore      `json:"awayScore"`
	Time           LiveTime       `json:"time"`
	ID             uint           `json:"id"`
	StartTimestamp uint           `json:"startTimestamp"`
	Slug           string         `json:"slug"`
}

type LiveScore struct {
	Current  uint  `json:"current"`
	Display  uint  `json:"display"`
	Period1  uint  `json:"period1"`
	Period2  *uint `json:"period2,omitempty"`
	Period3  *uint `json:"period3,omitempty"`
	Period4  *uint `json:"period4,omitempty"`
	Overtime *uint `json:"overtime,omitempty"`
}

type LiveTeam struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ShortName string  `json:"shortName"`
	NameCode  string  `json:"nameCode"`
	ID        uint    `json:"id"`
	Logo      *string `json:"logoUrl"`
}

type LiveTime struct {
	Played                      uint  `json:"played"`
	PeriodLength                uint  `json:"periodLength"`
	OvertimeLength              uint  `json:"overtimeLength"`
	TotalPeriodCount            uint  `json:"totalPeriodCount"`
	CurrentPeriodStartTimestamp *uint `json:"currentPeriodStartTimestamp,omitempty"`
}

type LiveStatus struct {
	Code        uint   `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type LiveSeason struct {
	Name string `json:"name"`
	Year string `json:"year"`
	ID   uint   `json:"id"`
}

type LiveTournament struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

// -- Match schedules
type MatchScheduleData struct {
	Events []MatchScheduleEvent `json:"events"`
}
type MatchScheduleEvent struct {
	Tournament     EventTournament     `json:"tournament"`
	Status         MatchScheduleStatus `json:"status"`
	HomeTeam       MatchScheduleTeam   `json:"homeTeam"`
	AwayTeam       MatchScheduleTeam   `json:"awayTeam"`
	HomeScore      MatchScheduleScore  `json:"homeScore"`
	AwayScore      MatchScheduleScore  `json:"awayScore"`
	Time           MatchScheduleTime   `json:"time"`
	ID             uint                `json:"id"`
	StartTimestamp int64               `json:"startTimestamp"`
}
type EventTournament struct {
	Name string `json:"name"`
}
type MatchScheduleScore struct {
	Current *uint `json:"current,omitempty"`
	Display *uint `json:"display,omitempty"`
}

type MatchScheduleTeam struct {
	ID uint `json:"id"`
}

type MatchScheduleTime struct {
	Played                      *int `json:"played,omitempty"`
	PeriodLength                *int `json:"periodLength,omitempty"`
	OvertimeLength              *int `json:"overtimeLength,omitempty"`
	TotalPeriodCount            *int `json:"totalPeriodCount,omitempty"`
	CurrentPeriodStartTimestamp *int `json:"currentPeriodStartTimestamp,omitempty"`
}

type MatchScheduleStatus struct {
	Description string                  `json:"description"`
	Type        MatchScheduleStatusType `json:"type"`
}

type MatchScheduleStatusType string

const (
	Finished      MatchScheduleStatusType = "finished"
	Inprogress    MatchScheduleStatusType = "inprogress"
	Notstarted    MatchScheduleStatusType = "notstarted"
	TypePostponed MatchScheduleStatusType = "postponed"
)

// -- Our reported types

// Data structure to represent Ably messages for live data publishing
type LiveMessage struct {
	IDs    []uint      `json:"eventIDs"`
	Events []LiveEvent `json:"events"`
}

// Data structure to represent a single entry of GetGames() response
type Game struct {
	ID         uint     `json:"id"`
	GameStatus string   `json:"gameStatus"`
	HomeTeam   TeamInfo `json:"homeTeam"`
	AwayTeam   TeamInfo `json:"awayTeam"`
	HomeScore  *uint    `json:"homeScore"`
	AwayScore  *uint    `json:"awayScore"`
	Date       string   `json:"date"`
}

type TeamInfo struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	ShortName      string  `json:"shortName"`
	Abbr           string  `json:"abbr"`
	ArenaName      string  `json:"arenaName"`
	ArenaSize      int     `json:"arenaSize"`
	Color          string  `json:"color"`
	SecondaryColor string  `json:"secondaryColor"`
	Logo           *string `json:"logo"`
}

type GetGamesDTO struct {
	Date               string `form:"date" binding:"required,datetime=2006-01-02"`
	LocalTimeZoneShift *int   `form:"localTimeZoneShift" binding:"omitempty,number,lte=0"`
}
