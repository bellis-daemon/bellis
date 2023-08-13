package geo

import (
	"encoding/json"
	"net/http"
	"time"
)

type GeoData struct {
	IP        string  `json:"query"`
	Status    string  `json:"status"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
	Timezone  string  `json:"timezone"`
}

var instance *GeoData

func Geo() (*GeoData, error) {
	if instance != nil {
		return instance, nil
	}
	resp, err := http.Get("http://ip-api.com/json/?fields=57809")
	if err != nil {
		return nil, err
	}
	var g GeoData
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		return nil, err
	}
	instance = &g
	go func() {
		time.Sleep(time.Hour)
		instance = nil
	}()
	return instance, nil
}
