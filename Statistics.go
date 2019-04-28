package homie_go

import "time"

type Statistics struct {
	uptime time.Duration
}

func NewStatistics(startTime time.Time) *Statistics {
	return &Statistics{
		uptime: time.Now().Sub(startTime),
	}
}
