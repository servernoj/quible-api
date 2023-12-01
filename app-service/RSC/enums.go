package RSC

type ScheduleStatus string

const (
	Scheduled  ScheduleStatus = "scheduled"
	Delayed    ScheduleStatus = "delayed"
	Postponed  ScheduleStatus = "postponed"
	Suspended  ScheduleStatus = "suspended"
	Canceled   ScheduleStatus = "canceled"
	Inprogress ScheduleStatus = "inprogress"
	Final      ScheduleStatus = "final"
	Completed  ScheduleStatus = "completed"
)

type ScheduleType string

const (
	Preseason     ScheduleType = "Preseason"
	RegularSeason ScheduleType = "Regular Season"
	Postseason    ScheduleType = "Postseason"
)

type ScheduleEventName string

const (
	FirstRound       ScheduleEventName = "First Round"
	Semifinals       ScheduleEventName = "Semifinals"
	ConferenceFinals ScheduleEventName = "Conference Finals"
	NBAFinals        ScheduleEventName = "NBA Finals"
)

type ScheduleRound int

const (
	One   ScheduleRound = 1
	Two   ScheduleRound = 2
	Three ScheduleRound = 3
	Four  ScheduleRound = 4
)

type PlayerPosition string

const (
	C PlayerPosition = "C"
	F PlayerPosition = "F"
	G PlayerPosition = "G"
	// --
	FC PlayerPosition = "FC"
	GF PlayerPosition = "GF"
	SF PlayerPosition = "SF"
	SG PlayerPosition = "SG"
	PF PlayerPosition = "PF"
	PG PlayerPosition = "PG"
)

type PlayerStatus string

const (
	Act   PlayerStatus = "ACT"
	Inact PlayerStatus = "INACT"
	Ques  PlayerStatus = "QUES"
)
