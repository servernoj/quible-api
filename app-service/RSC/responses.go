package RSC

type ScheduleData struct {
	NBA []ScheduleItem `json:"NBA"`
}
type TeamInfoData struct {
	NBA []TeamInfoItem `json:"NBA"`
}

type Schedule struct {
	Data ScheduleData `json:"data"`
}
type TeamInfo struct {
	Data TeamInfoData `json:"data"`
}
