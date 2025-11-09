package main

import (
	"context"
	"github.com/yasonofriychuk/school-21-workshop/internal/api/handler/current_weather"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/geo_ip"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/weather"
	"github.com/yasonofriychuk/school-21-workshop/internal/usecase/weather_current_temperature_by_ip"
	"github.com/yasonofriychuk/school-21-workshop/pkg/logger"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	log := logger.NewLogger(slog.LevelDebug, "local", os.Stdout)
	httpClient := http.DefaultClient

	geoIpGateway := geo_ip.New(httpClient)
	weatherGateway := weather.New(httpClient)

	currentTemperatureByIpUsecase := weather_current_temperature_by_ip.New(geoIpGateway, weatherGateway)

	apiMux := http.NewServeMux()

	apiMux.Handle("/weather/current", current_weather.New(log, currentTemperatureByIpUsecase))

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	mux.Handle("/api/1/", http.StripPrefix("/api/1", apiMux))

	log.WithContext(ctx).Info("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to start server")
	}
}
