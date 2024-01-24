package RSC

type ResponseItem interface {
	ScheduleItem | TeamInfoItemExtended | TeamSeasonStatItem | PlayerInfoItemExtended | PlayerSeasonStatItem | InjuryItem | LiveFeedItem
}

type Response[T ResponseItem] struct {
	Data struct {
		NBA []T `json:"NBA"`
	} `json:"data"`
}
