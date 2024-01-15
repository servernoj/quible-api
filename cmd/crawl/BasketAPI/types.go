package BasketAPI

// -- Season standing API

type Standings struct {
	Standings []Standing `json:"standings"`
}

type Standing struct {
	Type               string        `json:"type"`
	Rows               []StandingRow `json:"rows"`
	UpdatedAtTimestamp uint          `json:"updatedAtTimestamp"`
}

type StandingRow struct {
	Team StandingTeam `json:"team"`
}

type StandingTeam struct {
	ID uint `json:"id"`
}

// -- Team details API

type TeamDetailsData struct {
	Team TeamDetails `json:"team"`
}
type TeamDetails struct {
	Venue      Venue      `json:"venue"`
	ID         uint       `json:"id"`
	TeamColors TeamColors `json:"teamColors"`
}
type TeamColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}
type Venue struct {
	Stadium Stadium `json:"stadium"`
}

type Stadium struct {
	Name     string `json:"name"`
	Capacity uint   `json:"capacity"`
}
