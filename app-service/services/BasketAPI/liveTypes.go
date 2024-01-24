package BasketAPI

type LiveMessage struct {
	IDs    []uint     `json:"eventIDs"`
	Events []LM_Event `json:"events"`
}

// -- LiveMatches (LM) API

type LM_Data struct {
	Events []LM_Event `json:"events"`
}
type LM_Event struct {
	Tournament     LM_Tournament `json:"tournament"`
	Season         LM_Season     `json:"season"`
	Status         LM_Status     `json:"status"`
	HomeTeam       LM_Team       `json:"homeTeam"`
	AwayTeam       LM_Team       `json:"awayTeam"`
	HomeScore      LM_Score      `json:"homeScore"`
	AwayScore      LM_Score      `json:"awayScore"`
	Time           LM_Time       `json:"time"`
	ID             uint          `json:"id"`
	StartTimestamp uint          `json:"startTimestamp"`
	Slug           string        `json:"slug"`
}
type LM_Score struct {
	Current  uint  `json:"current"`
	Display  uint  `json:"display"`
	Period1  uint  `json:"period1"`
	Period2  *uint `json:"period2,omitempty"`
	Period3  *uint `json:"period3,omitempty"`
	Period4  *uint `json:"period4,omitempty"`
	Overtime *uint `json:"overtime,omitempty"`
}
type LM_Team struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ShortName string  `json:"shortName"`
	NameCode  string  `json:"nameCode"`
	ID        uint    `json:"id"`
	Logo      *string `json:"logoUrl"`
}
type LM_Time struct {
	Played                      int64  `json:"played"`
	PeriodLength                int64  `json:"periodLength"`
	OvertimeLength              int64  `json:"overtimeLength"`
	TotalPeriodCount            int64  `json:"totalPeriodCount"`
	CurrentPeriodStartTimestamp *int64 `json:"currentPeriodStartTimestamp,omitempty"`
}
type LM_Status struct {
	Code        uint   `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
type LM_Season struct {
	Name string `json:"name"`
	Year string `json:"year"`
	ID   uint   `json:"id"`
}
type LM_Tournament struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}
