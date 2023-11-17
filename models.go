package paddle

import "time"

type TimePeriodInterval string

const (
	TimePeriodIntervalDay   = TimePeriodInterval("day")
	TimePeriodIntervalWeek  = TimePeriodInterval("week")
	TimePeriodIntervalMonth = TimePeriodInterval("month")
	TimePeriodIntervalYear  = TimePeriodInterval("year")
)

type TimeInterval struct {
	Frequency int                `json:"frequency"`
	Interval  TimePeriodInterval `json:"interval"`
}

type Status string

const (
	StatusActive   = Status("active")
	StatusArchived = Status("archived")
)

type TimePeriod struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type MinMax struct {
	Minimum int `json:"minimum"`
	Maximum int `json:"maximum"`
}
