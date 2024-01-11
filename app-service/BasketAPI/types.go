package BasketAPI

// Data structure to represent Ably messages for live data publishing
type LiveMessage struct {
	IDs    []uint  `json:"eventIDs"`
	Events []Event `json:"events"`
}

// Data structure to represent a subset of fields from BasketAPI /live response
type LiveData struct {
	Events []Event `json:"events"`
}

type Event struct {
	Tournament     Tournament `json:"tournament"`
	Season         Season     `json:"season"`
	Status         Status     `json:"status"`
	HomeTeam       Team       `json:"homeTeam"`
	AwayTeam       Team       `json:"awayTeam"`
	HomeScore      Score      `json:"homeScore"`
	AwayScore      Score      `json:"awayScore"`
	Time           Time       `json:"time"`
	ID             uint       `json:"id"`
	StartTimestamp uint       `json:"startTimestamp"`
	Slug           string     `json:"slug"`
}

type Score struct {
	Current  uint  `json:"current"`
	Display  uint  `json:"display"`
	Period1  uint  `json:"period1"`
	Period2  *uint `json:"period2,omitempty"`
	Period3  *uint `json:"period3,omitempty"`
	Period4  *uint `json:"period4,omitempty"`
	Overtime *uint `json:"overtime,omitempty"`
}

type Team struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ShortName string  `json:"shortName"`
	NameCode  string  `json:"nameCode"`
	ID        uint    `json:"id"`
	Logo      *string `json:"logoUrl"`
}

type Time struct {
	Played                      uint  `json:"played"`
	PeriodLength                uint  `json:"periodLength"`
	OvertimeLength              uint  `json:"overtimeLength"`
	TotalPeriodCount            uint  `json:"totalPeriodCount"`
	CurrentPeriodStartTimestamp *uint `json:"currentPeriodStartTimestamp,omitempty"`
}

type Status struct {
	Code        uint   `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Season struct {
	Name string `json:"name"`
	Year string `json:"year"`
	ID   uint   `json:"id"`
}

type Tournament struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}
