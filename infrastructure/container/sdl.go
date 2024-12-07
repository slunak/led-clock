package container

import "github.com/veandco/go-sdl2/sdl"

func createSDLWindow(c *Config) (*sdl.Window, *sdl.Surface) {
	window, err := sdl.CreateWindow("Loading images", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(c.Cols), int32(c.Rows), sdl.WINDOW_FOREIGN)
	if err != nil {
		panic(err)
	}

	windowSurface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	return window, windowSurface
}
