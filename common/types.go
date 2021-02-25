package common

import (
	"encoding/json"
	"time"
)

// DailyStats represents the stats for a day.
type DailyStats struct {
	FileName   string
	PubDate    time.Time
	FirstDose  int
	SecondDose int
}

// MarshalJSON marshals the data into a format that the API can send to HTTP clients.
func (ds *DailyStats) MarshalJSON() ([]byte, error) {
	data := &struct {
		Date       string `json:"date"`
		FirstDose  int    `json:"first_dose"`
		SecondDose int    `json:"second_dose"`
	}{
		Date:       ds.PubDate.Format("2006-01-02"),
		FirstDose:  ds.FirstDose,
		SecondDose: ds.SecondDose,
	}

	return json.Marshal(data)
}
