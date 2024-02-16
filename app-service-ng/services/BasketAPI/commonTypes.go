package BasketAPI

type Event struct {
	Tournament     Tournament `json:"tournament"`
	Status         Status     `json:"status"`
	HomeTeam       TeamId     `json:"homeTeam"`
	AwayTeam       TeamId     `json:"awayTeam"`
	HomeScore      Score      `json:"homeScore"`
	AwayScore      Score      `json:"awayScore"`
	Time           Time       `json:"time"`
	ID             uint       `json:"id"`
	StartTimestamp int64      `json:"startTimestamp"`
}
type TeamId struct {
	ID uint `json:"id"`
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
type Score struct {
	Current *uint `json:"current,omitempty"`
	Display *uint `json:"display,omitempty"`
}
type Time struct {
	Played                      *int `json:"played,omitempty"`
	PeriodLength                *int `json:"periodLength,omitempty"`
	OvertimeLength              *int `json:"overtimeLength,omitempty"`
	TotalPeriodCount            *int `json:"totalPeriodCount,omitempty"`
	CurrentPeriodStartTimestamp *int `json:"currentPeriodStartTimestamp,omitempty"`
}
type Tournament struct {
	Name string `json:"name"`
}
type Status struct {
	Description string     `json:"description"`
	Type        StatusType `json:"type"`
}

type StatusType string

const (
	StatusType_Finished      StatusType = "finished"
	StatusType_Inprogress    StatusType = "inprogress"
	StatusType_Notstarted    StatusType = "notstarted"
	StatusType_TypePostponed StatusType = "postponed"
)
