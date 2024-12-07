package service

import "github.com/veandco/go-sdl2/ttf"

type FontType string

const (
	CLOCK           = "clock-12"
	WEATHER_CURRENT = "clock-12"
	DATE            = "date-9"
	WEEKDAY         = "weekday-6"
)

type FontProvider struct {
	fontsPath string
}

func NewFontProvider(fontsPath string) FontProvider {
	return FontProvider{
		fontsPath: fontsPath,
	}
}

func (f FontProvider) Provide(fontType FontType) (*ttf.Font, error) {
	font, err := ttf.OpenFont(f.fontsPath+string(fontType)+".bdf", 0)
	if err != nil {
		return nil, err
	}

	return font, nil
}
