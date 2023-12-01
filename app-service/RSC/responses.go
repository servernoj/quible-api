package RSC

type ScheduleData struct {
	NBA []ScheduleItem `json:"NBA"`
}
type TeamInfoData struct {
	NBA []TeamInfoItem `json:"NBA"`
}
type TeamSeasonStatData struct {
	NBA []TeamSeasonStatItem `json:"NBA"`
}
type PlayerInfoData struct {
	NBA []PlayerInfoItem `json:"NBA"`
}
type PlayerSeasonStatData struct {
	NBA []PlayerSeasonStatItem `json:"NBA"`
}
type InjuryData struct {
	NBA []InjuryItem `json:"NBA"`
}

type Schedule struct {
	Data ScheduleData `json:"data"`
}
type TeamInfo struct {
	Data TeamInfoData `json:"data"`
}
type TeamSeasonStats struct {
	Data TeamSeasonStatData `json:"data"`
}

type PlayerInfo struct {
	Data PlayerInfoData `json:"data"`
}

type PlayerSeasonStats struct {
	Data PlayerSeasonStatData `json:"data"`
}

type Injuries struct {
	Data InjuryData `json:"data"`
}
