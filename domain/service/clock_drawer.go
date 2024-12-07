package service

import (
	"github.com/veandco/go-sdl2/sdl"
	"image/color"
	"strings"
	"time"
)

type ClockDrawer struct {
	fontProvider FontProvider
	wiper        Wiper
}

func NewClockDrawer(fontProvider FontProvider, wiper Wiper) ClockDrawer {
	return ClockDrawer{
		fontProvider: fontProvider,
		wiper:        wiper,
	}
}

func (c ClockDrawer) Run(clockSurface *sdl.Surface, clockChannel chan error) {
	err := c.update(clockSurface)
	if err != nil {
		clockChannel <- err
		return
	}

	clockChannel <- err

	// wait until next minute
	time.Sleep(time.Second * time.Duration(60-time.Now().Second()))

	for {
		select {
		case <-time.Tick(time.Minute):
			err = c.update(clockSurface)
			if err != nil {
				clockChannel <- err
				return
			}
			clockChannel <- nil
		}
	}
}

func (c ClockDrawer) update(clockSurface *sdl.Surface) error {
	clockFont, err := c.fontProvider.Provide(CLOCK)
	if err != nil {
		return err
	}
	defer clockFont.Close()

	dateFont, err := c.fontProvider.Provide(DATE)
	if err != nil {
		return err
	}
	defer dateFont.Close()

	weekdayFont, err := c.fontProvider.Provide(WEEKDAY)
	if err != nil {
		return err
	}
	defer weekdayFont.Close()

	clock, err := clockFont.RenderUTF8Blended(time.Now().Format("15:04"), sdl.Color{R: 233, G: 234, B: 236, A: 255})
	if err != nil {
		return err
	}
	defer clock.Free()

	date, err := dateFont.RenderUTF8Blended(strings.ToUpper(time.Now().Format("_2.01")), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return err
	}
	defer date.Free()

	weekday, err := weekdayFont.RenderUTF8Blended(strings.ToUpper(time.Now().Format("Mon")), sdl.Color{R: 128, G: 128, B: 128, A: 255})
	if err != nil {
		return err
	}
	defer weekday.Free()

	err = c.wiper.Wipe(clockSurface)
	if err != nil {
		return err
	}

	err = clock.Blit(nil, clockSurface, &sdl.Rect{X: clockSurface.W - clock.W - 1, Y: 1, W: 0, H: 0})
	if err != nil {
		return err
	}

	err = date.Blit(nil, clockSurface, &sdl.Rect{X: clockSurface.W - date.W - 1, Y: 13, W: 0, H: 0})
	if err != nil {
		return err
	}

	err = weekday.Blit(nil, clockSurface, &sdl.Rect{X: clockSurface.W - date.W - weekday.W - 1, Y: 15, W: 0, H: 0})
	if err != nil {
		return err
	}

	//c.drawDays(clockSurface)

	return nil
}

func (c ClockDrawer) drawDays(clockSurface *sdl.Surface) {
	days := [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	day := time.Now().Format("Monday")

	x := 14
	y := 19

	notToday := color.RGBA{
		R: 205,
		G: 209,
		B: 228,
		A: 255,
	}

	today := color.RGBA{
		R: 255,
		G: 76,
		B: 48,
		A: 255,
	}

	for i := 0; i < len(days); i++ {
		if days[i] == day {
			clockSurface.Set(x+i, y, today)
			continue
		}
		clockSurface.Set(x+i, y, notToday)
	}
}
