package RSC

type ResponseItem interface {
	ScheduleItem | TeamInfoItem | TeamSeasonStatItem | PlayerInfoItem | PlayerSeasonStatItem | InjuryItem
}

type Response[T ResponseItem] struct {
	Data NBA[T] `json:"data"`
}

type NBA[T any] struct {
	NBA []T `json:"NBA"`
}
