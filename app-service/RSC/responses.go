package RSC

type ScheduleSeasonData struct {
	NBA []ScheduleSeasonItem `json:"NBA"`
}

type ScheduleSeason struct {
	Data ScheduleSeasonData `json:"data"`
}
