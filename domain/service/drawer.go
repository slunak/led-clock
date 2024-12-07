package service

import (
	"github.com/veandco/go-sdl2/sdl"
	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type Drawer interface {
	Draw() error
}

type drawer struct {
	canvas                *rgbmatrix.Canvas
	window                *sdl.Window
	windowSurface         *sdl.Surface
	canvasDrawer          CanvasDrawer
	clock                 ClockDrawer
	weather               WeatherDrawer
	weatherForecastDrawer WeatherForecastDrawer
}

func NewDrawer(
	canvas *rgbmatrix.Canvas,
	window *sdl.Window,
	windowSurface *sdl.Surface,
	canvasDrawer CanvasDrawer,
	clock ClockDrawer,
	weather WeatherDrawer,
	weatherForecastDrawer WeatherForecastDrawer,
) Drawer {
	return &drawer{
		canvas:                canvas,
		window:                window,
		windowSurface:         windowSurface,
		canvasDrawer:          canvasDrawer,
		clock:                 clock,
		weather:               weather,
		weatherForecastDrawer: weatherForecastDrawer,
	}
}

func (d drawer) Draw() error {
	defer d.window.Destroy()
	defer d.canvas.Close()

	var quit chan struct{}

	// TODO: create a separate struct with all the sizes, for all surfaces
	weatherSurface, err := sdl.CreateRGBSurface(0, 20, 24, 32, 0, 0, 0, 0)
	if err != nil {
		return err
	}
	clockSurface, err := sdl.CreateRGBSurface(0, 44, 24, 32, 0, 0, 0, 0)
	if err != nil {
		return err
	}
	bottomSurface, err := sdl.CreateRGBSurface(0, 64, 40, 32, 0, 0, 0, 0)
	if err != nil {
		return err
	}

	clockChannel := make(chan error)
	go d.clock.Run(clockSurface, clockChannel)
	weatherChannel := make(chan error)
	go d.weather.Run(weatherSurface, weatherChannel)
	bottomChannel := make(chan error)
	go d.weatherForecastDrawer.Run(bottomSurface, bottomChannel)

	for {
		select {
		case err = <-clockChannel:
			if err != nil {
				return err
			}
			err = clockSurface.Blit(nil, d.windowSurface, &sdl.Rect{X: 20, Y: 0, W: 0, H: 0})
			if err != nil {
				return err
			}
			err = d.updateScreen()
			if err != nil {
				return err
			}
		case err = <-weatherChannel:
			if err != nil {
				return err
			}
			err = weatherSurface.Blit(nil, d.windowSurface, &sdl.Rect{X: 0, Y: 0, W: 0, H: 0})
			if err != nil {
				return err
			}
			err = d.updateScreen()
			if err != nil {
				return err
			}
		case err = <-bottomChannel:
			if err != nil {
				return err
			}
			err = bottomSurface.Blit(nil, d.windowSurface, &sdl.Rect{X: 0, Y: 24, W: 0, H: 0})
			if err != nil {
				return err
			}
			err = d.updateScreen()
			if err != nil {
				return err
			}
		case <-quit:
			return nil
		}
	}
}

func (d drawer) updateScreen() error {
	err := d.window.UpdateSurface()
	if err != nil {
		return err
	}

	d.canvasDrawer.Draw(d.windowSurface, d.canvas)
	err = d.canvas.Render()
	if err != nil {
		return err
	}

	return nil
}
