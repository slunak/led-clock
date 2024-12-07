package service

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type WeatherDrawer struct {
	weather             Weather
	weatherIconProvider WeatherIconProvider
	fontProvider        FontProvider
	wiper               Wiper
}

func NewWeatherDrawer(weather Weather, weatherIconProvider WeatherIconProvider, fontProvider FontProvider, wiper Wiper) WeatherDrawer {
	return WeatherDrawer{
		weather:             weather,
		weatherIconProvider: weatherIconProvider,
		wiper:               wiper,
		fontProvider:        fontProvider,
	}
}

func (w WeatherDrawer) Run(weatherSurface *sdl.Surface, weatherChannel chan error) {
	currentWeather, err := w.weather.Current()
	if err != nil {
		weatherChannel <- err
		return
	}

	err = w.update(weatherSurface, currentWeather.Type, currentWeather.Temperature)
	if err != nil {
		weatherChannel <- err
		return
	}

	weatherChannel <- err

	for {
		select {
		case <-time.Tick(time.Minute * 5):
			currentWeather, err = w.weather.Current()
			if err != nil {
				// TODO: log error, no need to stop the whole application
				continue
			}
			err = w.update(weatherSurface, currentWeather.Type, currentWeather.Temperature)
			if err != nil {
				weatherChannel <- err
				return
			}
			weatherChannel <- nil
		}
	}
}

func (w WeatherDrawer) update(weatherSurface *sdl.Surface, weather WeatherType, temperature float64) error {
	// TODO: add temperature to the weather icon
	weatherIcon, err := w.weatherIconProvider.Provide(weather)
	if err != nil {
		return err
	}
	defer weatherIcon.Free()

	font, err := w.fontProvider.Provide(DATE)
	if err != nil {
		return err
	}
	defer font.Close()

	text, err := font.RenderUTF8Blended(fmt.Sprintf("%.0f°", temperature), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return err
	}
	defer text.Free()

	err = w.wiper.Wipe(weatherSurface)
	if err != nil {
		return err
	}

	err = weatherIcon.Blit(nil, weatherSurface, &sdl.Rect{X: weatherSurface.W/2 - weatherIcon.W/2, Y: 1, W: 0, H: 0})
	if err != nil {
		return err
	}

	err = text.Blit(nil, weatherSurface, &sdl.Rect{X: weatherSurface.W/2 - text.W/2, Y: 13, W: 0, H: 0})
	if err != nil {
		return err
	}

	return nil
}
