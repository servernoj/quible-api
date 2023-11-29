package RSC

type ScheduleSeasonStatus string

const (
	Scheduled  ScheduleSeasonStatus = "scheduled"
	Delayed    ScheduleSeasonStatus = "delayed"
	Postponed  ScheduleSeasonStatus = "postponed"
	Suspended  ScheduleSeasonStatus = "suspended"
	Canceled   ScheduleSeasonStatus = "canceled"
	Inprogress ScheduleSeasonStatus = "inprogress"
	Final      ScheduleSeasonStatus = "final"
	Completed  ScheduleSeasonStatus = "completed"
)

type ScheduleSeasonType string

const (
	Preseason     ScheduleSeasonType = "Preseason"
	RegularSeason ScheduleSeasonType = "Regular Season"
	Postseason    ScheduleSeasonType = "Postseason"
)

type ScheduleSeasonEventName string

const (
	FirstRound       ScheduleSeasonEventName = "First Round"
	Semifinals       ScheduleSeasonEventName = "Semifinals"
	ConferenceFinals ScheduleSeasonEventName = "Conference Finals"
	NBAFinals        ScheduleSeasonEventName = "NBA Finals"
)

type ScheduleSeasonRound int

const (
	One   ScheduleSeasonRound = 1
	Two   ScheduleSeasonRound = 2
	Three ScheduleSeasonRound = 3
	Four  ScheduleSeasonRound = 4
)
