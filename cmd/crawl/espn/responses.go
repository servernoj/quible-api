package espn

type ResponseItem interface {
	TeamItem | PlayerItem | ResponseWithItems
}

type Item struct {
	Ref string `json:"$ref"`
}

type Logo struct {
	Href string   `json:"href"`
	Rel  []string `json:"rel"`
}

type Headshot struct {
	Href string `json:"href"`
	Alt  string `json:"alt"`
}

type ResponseWithItems struct {
	Items []Item `json:"items"`
}

type TeamItem struct {
	Item
	GUID             string `json:"guid"`
	Slug             string `json:"slug"`
	ShortDisplayName string `json:"shortDisplayName"`
	Logos            []Logo `json:"logos"`
}

type PlayerItem struct {
	Item
	GUID        string    `json:"guid"`
	Slug        string    `json:"slug"`
	DisplayName string    `json:"displayName"`
	Headshot    *Headshot `json:"headshot"`
	Team        Item      `json:"team"`
}
