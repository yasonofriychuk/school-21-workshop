package current_weather

import "context"

type usecase interface {
	GetCurrenTemperatureByIp(ctx context.Context, ip string) (float64, error)
}
