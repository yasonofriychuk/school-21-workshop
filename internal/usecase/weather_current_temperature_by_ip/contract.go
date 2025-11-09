package weather_current_temperature_by_ip

import (
	"context"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/geo_ip/methods/get_geo_by_ip"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/weather/methods/get_current_weather_by_coords"
)

type geoIpGateway interface {
	GetGeoByIP(ctx context.Context, ipAdr string) (get_geo_by_ip.Coords, error)
}

type weatherGateway interface {
	GetCurrentWeatherByCoords(_ context.Context, coords get_current_weather_by_coords.Coords) (get_current_weather_by_coords.CurrentWeather, error)
}
