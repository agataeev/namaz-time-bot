package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiURL = "http://api.aladhan.com/v1/timingsByCity"

// PrayerTimes хранит время намазов
type PrayerTimes struct {
	Fajr    string `json:"Fajr"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

// APIResponse структура для разбора JSON
type APIResponse struct {
	Data struct {
		Timings PrayerTimes `json:"timings"`
	} `json:"data"`
}

// GetPrayerTimes получает время намазов по городу
func GetPrayerTimes(city, country string) (PrayerTimes, error) {
	url := fmt.Sprintf("%s?city=%s&country=%s&method=2", apiURL, city, country)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return PrayerTimes{}, err
	}
	defer resp.Body.Close()

	var result APIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return PrayerTimes{}, err
	}

	return result.Data.Timings, nil
}
