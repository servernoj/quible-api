package BasketAPI

type Game struct {
	ID         uint     `json:"id"`
	GameStatus string   `json:"gameStatus"`
	HomeTeam   TeamInfo `json:"homeTeam"`
	AwayTeam   TeamInfo `json:"awayTeam"`
	HomeScore  *uint    `json:"homeScore"`
	AwayScore  *uint    `json:"awayScore"`
	Date       string   `json:"date"`
}

type GetGamesDTO struct {
	Date               string `form:"date" binding:"required,datetime=2006-01-02"`
	LocalTimeZoneShift *int   `form:"localTimeZoneShift" binding:"omitempty,number,lte=0"`
}

// -- MatchSchedules (MS) API

type MS_Data struct {
	Events []MS_Event `json:"events"`
}
type MS_Event struct {
	Tournament     MS_Tournament `json:"tournament"`
	Status         MS_Status     `json:"status"`
	HomeTeam       TeamId        `json:"homeTeam"`
	AwayTeam       TeamId        `json:"awayTeam"`
	HomeScore      MS_Score      `json:"homeScore"`
	AwayScore      MS_Score      `json:"awayScore"`
	Time           MS_Time       `json:"time"`
	ID             uint          `json:"id"`
	StartTimestamp int64         `json:"startTimestamp"`
}
type MS_Tournament struct {
	Name string `json:"name"`
}
type MS_Score struct {
	Current *uint `json:"current,omitempty"`
	Display *uint `json:"display,omitempty"`
}
type MS_Time struct {
	Played                      *int `json:"played,omitempty"`
	PeriodLength                *int `json:"periodLength,omitempty"`
	OvertimeLength              *int `json:"overtimeLength,omitempty"`
	TotalPeriodCount            *int `json:"totalPeriodCount,omitempty"`
	CurrentPeriodStartTimestamp *int `json:"currentPeriodStartTimestamp,omitempty"`
}
type MS_Status struct {
	Description string        `json:"description"`
	Type        MS_StatusType `json:"type"`
}
type MS_StatusType string

const (
	MS_StatusType_Finished      MS_StatusType = "finished"
	MS_StatusType_Inprogress    MS_StatusType = "inprogress"
	MS_StatusType_Notstarted    MS_StatusType = "notstarted"
	MS_StatusType_TypePostponed MS_StatusType = "postponed"
)
