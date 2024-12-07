package service

import "time"

type Weather interface {
	Current() (CurrentWeather, error)
	Forecast() (ForecastWeather, error)
}

type CurrentWeather struct {
	Temperature float64
	Type        WeatherType
}

type ForecastWeather struct {
	HourlyMetrics []Metric
}

type Metric struct {
	Time                     time.Time
	Temperature              float64
	WindSpeed                float64
	WindGusts                float64
	Precipitation            float64
	PrecipitationProbability float64
	WeatherType              WeatherType
	IsDay                    bool
}
