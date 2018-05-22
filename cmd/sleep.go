package cmd

import (
	"net/http"
	"encoding/json"
	"log"
)

// Data
type Data struct {
	DateTime string `json:"dateTime,omitempty"`
	Level string `json:"level,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}

// Deep
type Deep struct {
	Count int `json:"count,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
}

// Levels
type Levels struct {
	Data []*Data `json:"data,omitempty"`
	ShortData []*ShortData `json:"shortData,omitempty"`
	Summary *Summary `json:"summary,omitempty"`
}

// Light
type Light struct {
	Count int `json:"count,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
}

// Rem
type Rem struct {
	Count int `json:"count,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
}

// Root
type SleepSummary struct {
	Sleep []*Sleep `json:"sleep,omitempty"`
	Summary *Summary2 `json:"summary,omitempty"`
}

// ShortData
type ShortData struct {
	DateTime string `json:"dateTime,omitempty"`
	Level string `json:"level,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}

// Sleep
type Sleep struct {
	DateOfSleep string `json:"dateOfSleep,omitempty"`
	Duration int `json:"duration,omitempty"`
	Efficiency int `json:"efficiency,omitempty"`
	EndTime string `json:"endTime,omitempty"`
	InfoCode int `json:"infoCode,omitempty"`
	IsMainSleep bool `json:"isMainSleep,omitempty"`
	Levels *Levels `json:"levels,omitempty"`
	LogId int `json:"logId,omitempty"`
	MinutesAfterWakeup int `json:"minutesAfterWakeup,omitempty"`
	MinutesAsleep int `json:"minutesAsleep,omitempty"`
	MinutesAwake int `json:"minutesAwake,omitempty"`
	MinutesToFallAsleep int `json:"minutesToFallAsleep,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	TimeInBed int `json:"timeInBed,omitempty"`
	Type string `json:"type,omitempty"`
}

// Stages
type Stages struct {
	Deep int `json:"deep,omitempty"`
	Light int `json:"light,omitempty"`
	Rem int `json:"rem,omitempty"`
	Wake int `json:"wake,omitempty"`
}

// Summary
type Summary struct {
	Deep *Deep `json:"deep,omitempty"`
	Light *Light `json:"light,omitempty"`
	Rem *Rem `json:"rem,omitempty"`
	Wake *Wake `json:"wake,omitempty"`
}

// Summary2
type Summary2 struct {
	Stages *Stages `json:"stages,omitempty"`
	TotalMinutesAsleep int `json:"totalMinutesAsleep,omitempty"`
	TotalSleepRecords int `json:"totalSleepRecords,omitempty"`
	TotalTimeInBed int `json:"totalTimeInBed,omitempty"`
}

// Wake
type Wake struct {
	Count int `json:"count,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes,omitempty"`
}


func getDailySleepSummary(c *http.Client, date string) SleepSummary {

	var response *http.Response
	response, err = c.Get(
			"https://api.fitbit.com/1.2/user/-/sleep/date/" + date + ".json")
	if err != nil {
		log.Fatal(err)
	}

	var sleep SleepSummary
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&sleep); err != nil {
		log.Fatal(err)
	}
	return sleep
}
