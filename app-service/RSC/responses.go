package RSC

type ScheduleData struct {
	NBA []ScheduleItem `json:"NBA"`
}

type Schedule struct {
	Data ScheduleData `json:"data"`
}
