package infrastructure

import (
	"github.com/labstack/echo/v4"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"led-clock/infrastructure/container"
	"led-clock/infrastructure/route"
)

func Start(container container.Container) error {
	err := initSDL()
	if err != nil {
		return err
	}

	defer sdl.Quit()
	defer ttf.Quit()

	drawer := container.GetDrawer()
	go func() {
		err := drawer.Draw()
		if err != nil {
			panic(err)
		}
	}()

	e := echo.New()
	route.PrepareRoutes(e, container)
	e.Logger.Fatal(e.Start(":8000"))

	return nil
}

func initSDL() error {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	err = ttf.Init()
	if err != nil {
		return err
	}

	return nil
}
