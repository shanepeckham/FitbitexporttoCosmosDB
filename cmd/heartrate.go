package cmd

import (
	"net/http"
	"encoding/json"
	"log"
)

// ActivitiesHeart
type ActivitiesHeart struct {
	DateTime string `json:"dateTime,omitempty"`
	Value *Value `json:"values,omitempty"`
}

type Value struct {
	HeartRateZones []*HeartRateZones `json:"heartrate"`
}

// ActivitiesHeartIntraday
type ActivitiesHeartIntraday struct {
	Dataset []interface{} `json:"dataset,omitempty"`
	DatasetInterval int `json:"datasetInterval,omitempty"`
	DatasetType string `json:"datasetType,omitempty"`
}

// HeartRateZones
type HeartRateZones struct {
	Max int `json:"max,omitempty"`
	Min int `json:"min,omitempty"`
	Name string `json:"name,omitempty"`
}

// HeartRate
type HeartRate struct {
	ActivitiesHeart []*ActivitiesHeart `json:"activities-heart,omitempty"`
	ActivitiesHeartIntraday *ActivitiesHeartIntraday `json:"activities-heart-intraday,omitempty"`
}

func getDailyHeartRateSummary(c *http.Client, date string) HeartRate {

	var response *http.Response
	response, err = c.Get(
			"https://api.fitbit.com/1/user/-/activities/heart/date/" + date + "/1d.json")
	if err != nil {
		log.Fatal(err)
	}

	var heartrate HeartRate
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&heartrate); err != nil {
		log.Fatal(err)
	}
	return heartrate
}