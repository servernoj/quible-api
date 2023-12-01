package RSC

type ResponseItem interface {
	ScheduleItem | TeamInfoItem | TeamSeasonStatItem | PlayerInfoItem | PlayerSeasonStatItem | InjuryItem | LiveFeedItem
}

type Response[T ResponseItem] struct {
	Data struct {
		NBA []T `json:"NBA"`
	} `json:"data"`
}
