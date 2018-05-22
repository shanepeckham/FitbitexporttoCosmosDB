package cmd

import (
	"net/http"
	"encoding/json"
	"log"
)


type ActivitySteps struct {
	Steps []DataPoint `json:"activities-steps"`
}

// Activity holds all the basic information for a single measured activity
type Activity struct {
	ActivityId       uint64  `json:"activityId"`
	ActivityParentId uint64  `json:"activityParentId"`
	Calories         uint64  `json:"calories"`
	Description      string  `json:"description"`
	Duration         uint64  `json:"duration"`
	Distance         float64 `json:"distance"`
	HasStartTime     bool    `json:"hasStartTime"`
	IsFavorite       bool    `json:"isFavorite"`
	LogID            uint64  `json:"logId"`
	Name             string  `json:"name"`
	StartTime        string  `json:"startTime"`
	Steps            uint64  `json:"steps"`
}

// Goal represents all data reached to a given date
type Goal struct {
	CaloriesOut uint64  `json:"caloriesOut"`
	Distance    float64 `json:"distance"`
	Floors      uint64  `json:"floors"`
	Steps       uint64  `json:"steps"`
}

// Distance holds different distances per activity (tracker, total, veryActive, etc.)
type Distance struct {
	Activity string  `json:"activity"`
	Distance float64 `json:"distance"`
}

// Summary holds a summary of all the activities of a given date
type ActivitySummary struct {
	ActivityCalories     uint64      `json:"activityCalories"`
	CaloriesBMR          uint64      `json:"caloriesBMR"`
	CaloriesOut          uint64      `json:"caloriesOut"`
	Distances            []*Distance `json:"distances"`
	Elevation            float64     `json:"elevation"`
	FairlyActiveMinutes  uint64      `json:"fairlyActiveMinutes"`
	Floors               uint64      `json:"floors"`
	LightlyActiveMinutes uint64      `json:"lightlyActiveMinutes"`
	MarginalCalories     uint64      `json:"marginalCalories"`
	SedentaryMinutes     uint64      `json:"sedentaryMinutes"`
	Steps                uint64      `json:"steps"`
	VeryActiveMinutes    uint64      `json:"veryActiveMinutes"`
}

// Activities for a specific given date
type Activities struct {
	Activities []*Activity      `json:"activities"`
	Goals      *Goal            `json:"goals"`
	Summary    *ActivitySummary `json:"summary"`
	DateTime string `json:"dateTime,omitempty"`
}

func getDailyActivitySummary(c *http.Client, date string) Activities {

	var response *http.Response
	response, err = c.Get(
		"https://api.fitbit.com/1/user/-/activities/date/" + date + ".json")
	if err != nil {
		log.Fatal(err)
	}

	var activities Activities
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&activities); err != nil {
		log.Fatal(err)
	}
	activities.DateTime = exportDate
	return activities
}