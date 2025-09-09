package settings

import (
	"time"
)

type HistoricalSourceSettings struct {
	Symbol    string    `json:"symbol" validate:"required"`
	Interval  string    `json:"interval" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

func (d HistoricalSourceSettings) SettingsType() string {
	return "database"
}

var _ Settings = HistoricalSourceSettings{}
