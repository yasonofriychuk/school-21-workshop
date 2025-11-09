// DOCS: https://open-meteo.com/en/docs

package weather

import (
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/weather/methods/get_current_weather_by_coords"
)

type (
	getCurrentWeatherByCoordsMethod struct {
		*get_current_weather_by_coords.Method
	}
)

type Gateway struct {
	getCurrentWeatherByCoordsMethod
}

func New(client gateway.HTTPGetter) *Gateway {
	return &Gateway{
		getCurrentWeatherByCoordsMethod: getCurrentWeatherByCoordsMethod{get_current_weather_by_coords.New(client)},
	}
}
