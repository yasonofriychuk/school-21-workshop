package get_current_weather_by_coords

type Response struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentWeatherUnits  struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature   string `json:"temperature"`
		Windspeed     string `json:"windspeed"`
		Winddirection string `json:"winddirection"`
		IsDay         string `json:"is_day"`
		Weathercode   string `json:"weathercode"`
	} `json:"current_weather_units"`
	CurrentWeather struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection int     `json:"winddirection"`
		IsDay         int     `json:"is_day"`
		Weathercode   int     `json:"weathercode"`
	} `json:"current_weather"`
}
