package weather_current_temperature_by_ip

import (
	"context"
	"fmt"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/weather/methods/get_current_weather_by_coords"
)

type Usecase struct {
	geoIpGateway   geoIpGateway
	weatherGateway weatherGateway
}

func New(geoIpGateway geoIpGateway, weatherGateway weatherGateway) *Usecase {
	return &Usecase{
		geoIpGateway:   geoIpGateway,
		weatherGateway: weatherGateway,
	}
}

func (u *Usecase) GetCurrenTemperatureByIp(ctx context.Context, ip string) (float64, error) {
	coords, err := u.geoIpGateway.GetGeoByIP(ctx, ip)
	if err != nil {
		return 0, fmt.Errorf("u.geoIpGateway.GetGeoByIP: %w", err)
	}

	currentWeather, err := u.weatherGateway.GetCurrentWeatherByCoords(ctx, get_current_weather_by_coords.Coords{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	})
	if err != nil {
		return 0, fmt.Errorf("u.weatherGateway.GetCurrentWeatherByCoords: %w", err)
	}

	return currentWeather.Temperature, nil
}
