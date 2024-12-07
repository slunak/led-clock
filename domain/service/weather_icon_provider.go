package service

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type WeatherType string

const (
	Sun        = "sun"
	Snow       = "snow"
	Moon       = "moon"
	MoonCloud  = "moon-cloud"
	Cloud      = "cloud"
	CloudRain  = "cloud-rain"
	CloudStorm = "cloud-storm"
	CloudSun   = "cloud-sun"
	Wind       = "wind"
)

type WeatherIconProvider struct {
	iconsPath string
}

func NewWeatherIconProvider(iconsPath string) WeatherIconProvider {
	return WeatherIconProvider{
		iconsPath: iconsPath,
	}
}

func (w WeatherIconProvider) Provide(weatherType WeatherType) (*sdl.Surface, error) {
	i, err := img.Load(w.iconsPath + "weather-icon-" + string(weatherType) + ".png")
	if err != nil {
		return nil, err
	}

	return i, nil
}
