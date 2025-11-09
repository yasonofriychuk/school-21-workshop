package get_current_weather_by_coords

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/weather/utils"
	"net/http"
)

type Method struct {
	client  gateway.HTTPGetter
	baseURL string
}

func New(client gateway.HTTPGetter) *Method {
	return &Method{
		client:  client,
		baseURL: utils.BaseURL,
	}
}

func (m *Method) GetCurrentWeatherByCoords(_ context.Context, coords Coords) (CurrentWeather, error) {
	url := fmt.Sprintf("%s/v1/forecast?latitude=%f&longitude=%f&current_weather=true", m.baseURL, coords.Latitude, coords.Longitude)
	resp, err := m.client.Get(url)
	if err != nil {
		return CurrentWeather{}, fmt.Errorf("m.client.Get: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CurrentWeather{}, fmt.Errorf("status code not 200: code: %d", resp.StatusCode)
	}

	var body Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return CurrentWeather{}, fmt.Errorf("json.NewDecoder: Decode: %w", err)
	}

	return CurrentWeather{
		Temperature: body.CurrentWeather.Temperature,
	}, nil
}
