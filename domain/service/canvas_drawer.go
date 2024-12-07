package service

import (
	"github.com/veandco/go-sdl2/sdl"
	rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"
)

type CanvasDrawer struct {
}

func NewCanvasDrawer() CanvasDrawer {
	return CanvasDrawer{}
}

func (cd CanvasDrawer) Draw(surface *sdl.Surface, canvas *rgbmatrix.Canvas) {
	for xx := 0; xx < 64; xx++ {
		for yy := 0; yy < 64; yy++ {
			pixel := surface.At(xx, yy)
			canvas.Set(xx, yy, pixel)
		}
	}
}
