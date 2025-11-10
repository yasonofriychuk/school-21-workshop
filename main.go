package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseIp struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ResponseWeather struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
	} `json:"current_weather"`
}

func main() {
	apiMux := http.NewServeMux()
	// web/1/service/weathers?ip=127.0.0.1
	apiMux.HandleFunc("/web/1/service/weathers", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")

		resp, err := http.Get(fmt.Sprintf("https://ipwho.is/%s?lang=ru&output=json", ip))
		if err != nil {

		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {

		}

		var respIp ResponseIp
		if err := json.NewDecoder(resp.Body).Decode(&respIp); err != nil {

		}

		respWeather, err := http.Get(fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
			respIp.Latitude, respIp.Longitude,
		))
		defer respWeather.Body.Close()

		if respWeather.StatusCode != http.StatusOK {

		}

		var respWeatherDao ResponseWeather
		if err := json.NewDecoder(respWeather.Body).Decode(&respWeatherDao); err != nil {

		}

		fmt.Println(respWeatherDao)
	})

	if err := http.ListenAndServe(":8080", apiMux); err != nil {

	}
}
