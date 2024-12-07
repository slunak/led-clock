package service

import "github.com/veandco/go-sdl2/sdl"

type Wiper struct {
}

func NewWiper() Wiper {
	return Wiper{}
}

func (w Wiper) Wipe(surfaces ...*sdl.Surface) error {
	for _, surface := range surfaces {
		err := surface.FillRect(nil, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
