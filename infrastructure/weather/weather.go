package weather

import (
	"context"
	"github.com/slunak/omgo"
	"led-clock/domain/service"
	"time"
)

const windSpeedWarning float64 = 30

type weather struct {
	client            *omgo.Client
	lat, lon          float64
	temperatureUnit   string
	windspeedUnit     string
	precipitationUnit string
	timezone          string
}

func NewWeather(client *omgo.Client, lat, lon float64, temperatureUnit, windspeedUnit, precipitationUnit, timezone string) service.Weather {
	return weather{
		client:            client,
		lat:               lat,
		lon:               lon,
		temperatureUnit:   temperatureUnit,
		windspeedUnit:     windspeedUnit,
		precipitationUnit: precipitationUnit,
		timezone:          timezone,
	}
}

func (w weather) Current() (service.CurrentWeather, error) {
	var currentWeather service.CurrentWeather
	loc, _ := omgo.NewLocation(w.lat, w.lon)
	opt := omgo.Options{
		TemperatureUnit:   w.temperatureUnit,
		WindspeedUnit:     w.windspeedUnit,
		PrecipitationUnit: w.precipitationUnit,
		Timezone:          w.timezone,
		PastDays:          0,
	}
	response, err := w.client.CurrentWeather(context.TODO(), loc, &opt)
	if err != nil {
		return currentWeather, err
	}

	currentWeather.Temperature = response.Temperature
	currentWeather.Type = w.convertWeatherType(response.WeatherCode, response.WindSpeed, response.IsDay != 0) // response.IsDay != 0 converting int to bool

	return currentWeather, nil
}

func (w weather) Forecast() (service.ForecastWeather, error) {
	var forecastWeather service.ForecastWeather
	var metrics []service.Metric

	loc, _ := omgo.NewLocation(w.lat, w.lon)
	opt := omgo.Options{
		TemperatureUnit:   w.temperatureUnit,
		WindspeedUnit:     w.windspeedUnit,
		PrecipitationUnit: w.precipitationUnit,
		Timezone:          w.timezone,
		PastDays:          0,
		HourlyMetrics:     []string{"temperature_2m", "wind_speed_10m", "wind_gusts_10m", "precipitation", "precipitation_probability", "weather_code", "is_day"},
	}

	response, err := w.client.Forecast(context.TODO(), loc, &opt)
	if err != nil {
		return forecastWeather, err
	}

	location, err := time.LoadLocation(w.timezone)
	if err != nil {
		return forecastWeather, err
	}

	for i, metricTime := range response.HourlyTimes {
		// convert metricTime to local timezone because it is in UTC
		metricTimeLocalized, err := time.ParseInLocation(time.ANSIC, metricTime.Format(time.ANSIC), location)
		if err != nil {
			return forecastWeather, err
		}

		// discard past metrics
		if metricTimeLocalized.Before(time.Now()) {
			continue
		}

		metrics = append(metrics, service.Metric{
			Time:                     metricTimeLocalized,
			Temperature:              response.HourlyMetrics["temperature_2m"][i],
			WindSpeed:                response.HourlyMetrics["wind_speed_10m"][i],
			WindGusts:                response.HourlyMetrics["wind_gusts_10m"][i],
			Precipitation:            response.HourlyMetrics["precipitation"][i],
			PrecipitationProbability: response.HourlyMetrics["precipitation_probability"][i],
			WeatherType:              w.convertWeatherType(response.HourlyMetrics["weather_code"][i], response.HourlyMetrics["wind_speed_10m"][i], response.HourlyMetrics["is_day"][i] != 0),
			IsDay:                    response.HourlyMetrics["is_day"][i] != 0, // converting int to bool
		})
	}

	forecastWeather = service.ForecastWeather{
		HourlyMetrics: metrics,
	}

	return forecastWeather, nil
}

/*
WMO Weather interpretation codes (WW)
Code 	Description
0 	Clear sky
1, 2, 3 	Mainly clear, partly cloudy, and overcast
45, 48 	Fog and depositing rime fog
51, 53, 55 	Drizzle: Light, moderate, and dense intensity
56, 57 	Freezing Drizzle: Light and dense intensity
61, 63, 65 	Rain: Slight, moderate and heavy intensity
66, 67 	Freezing Rain: Light and heavy intensity
71, 73, 75 	Snow fall: Slight, moderate, and heavy intensity
77 	Snow grains
80, 81, 82 	Rain showers: Slight, moderate, and violent
85, 86 	Snow showers slight and heavy
95 * 	Thunderstorm: Slight or moderate
96, 99 * 	Thunderstorm with slight and heavy hail

(*) Thunderstorm forecast with hail is only available in Central Europe
*/
func (w weather) convertWeatherType(WMOCode, windSpeed float64, isDay bool) service.WeatherType {
	switch int(WMOCode) {
	case 0:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return w.sunOrMoonIcon(isDay)
	case 1:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return w.sunOrMoonIcon(isDay)
	case 2:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return w.cloudSunOrMoon(isDay)
	case 3:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return service.Cloud
	case 45:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return service.Cloud
	case 48:
		if windSpeed > windSpeedWarning {
			return service.Wind
		}
		return service.Cloud
	case 51, 53, 55, 56, 57, 61, 63, 65, 66, 67:
		return service.CloudRain
	case 71, 73, 75, 77:
		return service.Snow
	case 80, 81, 82, 85, 86:
		return service.CloudRain
	case 95, 96, 99:
		return service.CloudStorm
	default:
		return service.MoonCloud // Default for unknown codes
	}
}

func (w weather) sunOrMoonIcon(isDay bool) service.WeatherType {
	if isDay {
		return service.Sun
	}
	return service.Moon
}

func (w weather) cloudSunOrMoon(isDay bool) service.WeatherType {
	if isDay {
		return service.CloudSun
	}
	return service.MoonCloud
}
