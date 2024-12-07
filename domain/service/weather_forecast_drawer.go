package service

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"image/color"
	"time"
)

type WeatherForecastDrawer struct {
	weather             Weather
	weatherIconProvider WeatherIconProvider
	fontProvider        FontProvider
	wiper               Wiper
}

func NewWeatherForecastDrawer(
	weather Weather, weatherIconProvider WeatherIconProvider,
	fontProvider FontProvider, wiper Wiper,
) WeatherForecastDrawer {
	return WeatherForecastDrawer{
		weather:             weather,
		weatherIconProvider: weatherIconProvider,
		wiper:               wiper,
		fontProvider:        fontProvider,
	}
}

func (w WeatherForecastDrawer) Run(bottomSurface *sdl.Surface, bottomChannel chan error) {
	forecastWeather, err := w.weather.Forecast()
	if err != nil {
		bottomChannel <- err
		return
	}

	err = w.update(bottomSurface, forecastWeather)
	if err != nil {
		bottomChannel <- err
		return
	}

	bottomChannel <- err

	for {
		select {
		case <-time.Tick(time.Minute * 5):
			forecastWeather, err = w.weather.Forecast()
			if err != nil {
				// TODO: log error, no need to stop the whole application, log locally withing the app
				continue
			}
			err = w.update(bottomSurface, forecastWeather)
			if err != nil {
				bottomChannel <- err
				return
			}
			bottomChannel <- nil
		}
	}
}

func (w WeatherForecastDrawer) update(bottomSurface *sdl.Surface, forecastWeather ForecastWeather) error {
	err := w.wiper.Wipe(bottomSurface)
	if err != nil {
		return err
	}

	var x int32 = 0

	for i := 0; i < 4; i++ {
		weather := forecastWeather.HourlyMetrics[i]

		weatherIcon, err := w.weatherIconProvider.Provide(weather.WeatherType)
		if err != nil {
			return err
		}

		font, err := w.fontProvider.Provide(WEEKDAY)
		if err != nil {
			return err
		}

		err = weatherIcon.Blit(nil, bottomSurface, &sdl.Rect{X: x, Y: 0, W: 0, H: 0})
		if err != nil {
			return err
		}

		background, err := w.drawBackground(weather)
		if err != nil {
			return err
		}

		err = background.Blit(nil, bottomSurface, &sdl.Rect{X: x, Y: 11, W: 0, H: 0})
		if err != nil {
			return err
		}

		x = x + 16

		weatherIcon.Free()
		font.Close()
	}

	return nil
}

func (w WeatherForecastDrawer) drawBackground(weather Metric) (*sdl.Surface, error) {
	var width, height int32 = 16, 29

	background, err := sdl.CreateRGBSurface(0, width, height, 32, 0, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// precipitation probability only applies when precipitation is 0.1 or more
	switch weather.Precipitation {
	case 0.0:
		err = w.drawWithTemperature(background, weather.Temperature)
		if err != nil {
			return nil, err
		}
	default:
		err = w.drawWithPrecipitation(background, weather.Precipitation, weather.PrecipitationProbability)
		if err != nil {
			return nil, err
		}
	}

	font, err := w.fontProvider.Provide(WEEKDAY)
	if err != nil {
		return nil, err
	}

	hour, err := font.RenderUTF8Blended(fmt.Sprintf("%d", weather.Time.Hour()), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return nil, err
	}

	err = hour.Blit(nil, background, &sdl.Rect{X: background.W/2 - hour.W/2, Y: background.H - hour.H, W: 0, H: 0})
	if err != nil {
		return nil, err
	}

	font.Close()
	hour.Free()

	return background, nil
}

func (w WeatherForecastDrawer) drawWithPrecipitation(background *sdl.Surface, precipitation, precipitationProbability float64) error {
	height, width := background.H, background.W

	rainColor := color.RGBA{
		R: 0,
		G: 124,
		B: 255,
		A: 255,
	}

	precipitationVisualStart := int(height - 8)
	precipitationVisualEnd := precipitationVisualStart - int(precipitationProbability)/10

	for x := 0; x < int(width); x++ {
		for y := precipitationVisualStart; y > precipitationVisualEnd; y-- {
			background.Set(x, y, rainColor)
		}
	}

	font, err := w.fontProvider.Provide(WEEKDAY)
	if err != nil {
		return err
	}

	precipitationAmount, err := font.RenderUTF8Blended(fmt.Sprintf("%.1f", precipitation), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return err
	}

	err = precipitationAmount.Blit(nil, background, &sdl.Rect{X: background.W/2 - precipitationAmount.W/2, Y: 0, W: 0, H: 0})
	if err != nil {
		return err
	}

	font.Close()
	precipitationAmount.Free()

	return nil
}

func (w WeatherForecastDrawer) drawWithTemperature(background *sdl.Surface, temperature float64) error {
	font, err := w.fontProvider.Provide(WEEKDAY)
	if err != nil {
		return err
	}

	temperatureText, err := font.RenderUTF8Blended(fmt.Sprintf("%.0f°", temperature), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return err
	}

	err = temperatureText.Blit(nil, background, &sdl.Rect{X: background.W/2 - temperatureText.W/2, Y: 0, W: 0, H: 0})
	if err != nil {
		return err
	}

	font.Close()
	temperatureText.Free()

	return nil
}
