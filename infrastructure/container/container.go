package container

import (
	"github.com/slunak/omgo"
	"github.com/veandco/go-sdl2/sdl"
	"led-clock/domain/service"
	"led-clock/infrastructure/weather"

	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type Container interface {
	SetConfigFlags(lat, lon, tz *string) error
	GetDrawer() service.Drawer
	GetCanvasDrawer() service.CanvasDrawer
	GetFontProvider() service.FontProvider
	GetWeatherIconProvider() service.WeatherIconProvider
	GetWiper() service.Wiper
	GetClockDrawer() service.ClockDrawer
	GetWeatherDrawer() service.WeatherDrawer
	GetWeatherForecastDrawer() service.WeatherForecastDrawer
	GetWeather() service.Weather
}

var instance *container

type container struct {
	config        *Config
	weatherClient *omgo.Client
	matrix        *rgbmatrix.Matrix
	canvas        *rgbmatrix.Canvas
	window        *sdl.Window
	windowSurface *sdl.Surface
}

func GetContainer() (Container, error) {
	instance = &container{}
	instance.config = getConfig()
	instance.weatherClient = createWeatherClient()
	instance.matrix = createRGBMatrix(instance.config)
	instance.canvas = createCanvas(instance.matrix)
	instance.window, instance.windowSurface = createSDLWindow(instance.config)

	return instance, nil
}

func (c *container) GetDrawer() service.Drawer {
	return service.NewDrawer(
		c.canvas,
		c.window,
		c.windowSurface,
		c.GetCanvasDrawer(),
		c.GetClockDrawer(),
		c.GetWeatherDrawer(),
		c.GetWeatherForecastDrawer(),
	)
}

func (c *container) SetConfigFlags(lat, lon, tz *string) error {
	latFloat, lonFloat, tzString, err := convertConfigFlags(lat, lon, tz)
	if err != nil {
		return err
	}

	c.config.lat = latFloat
	c.config.lon = lonFloat
	c.config.timezone = tzString

	return nil
}

func (c *container) GetCanvasDrawer() service.CanvasDrawer {
	return service.NewCanvasDrawer()
}

func (c *container) GetFontProvider() service.FontProvider {
	return service.NewFontProvider(c.config.fontsPath)
}

func (c *container) GetWeatherIconProvider() service.WeatherIconProvider {
	return service.NewWeatherIconProvider(c.config.weatherIconsPath)
}

func (c *container) GetWiper() service.Wiper {
	return service.NewWiper()
}

func (c *container) GetClockDrawer() service.ClockDrawer {
	return service.NewClockDrawer(
		c.GetFontProvider(),
		c.GetWiper(),
	)
}

func (c *container) GetWeatherDrawer() service.WeatherDrawer {
	return service.NewWeatherDrawer(
		c.GetWeather(),
		c.GetWeatherIconProvider(),
		c.GetFontProvider(),
		c.GetWiper(),
	)
}

func (c *container) GetWeatherForecastDrawer() service.WeatherForecastDrawer {
	return service.NewWeatherForecastDrawer(
		c.GetWeather(),
		c.GetWeatherIconProvider(),
		c.GetFontProvider(),
		c.GetWiper(),
	)
}

func (c *container) GetWeather() service.Weather {
	return weather.NewWeather(
		c.weatherClient,
		c.config.lat,
		c.config.lon,
		c.config.temperatureUnit,
		c.config.windspeedUnit,
		c.config.precipitationUnit,
		c.config.timezone,
	)
}
