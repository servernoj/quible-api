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

type Schedule struct {
	Data ScheduleData `json:"data"`
}
type TeamInfo struct {
	Data TeamInfoData `json:"data"`
}
type TeamSeasonStats struct {
	Data TeamSeasonStatData `json:"data"`
}
