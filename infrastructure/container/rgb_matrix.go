package container

import rgbmatrix "github.com/zaggash/go-rpi-rgb-led-matrix"

func createRGBMatrix(c *Config) *rgbmatrix.Matrix {
	dc := &rgbmatrix.DefaultConfig
	dc.Rows = c.Rows
	dc.Cols = c.Cols
	dc.Parallel = c.Parallel
	dc.ChainLength = c.ChainLength
	dc.Brightness = c.Brightness
	dc.GPIOMapping = c.GPIOMapping
	dc.ShowRefreshRate = c.ShowRefresh
	dc.InverseColors = c.InverseColors
	dc.DisableHardwarePulsing = c.DisableHardwarePulsing

	rc := &rgbmatrix.DefaultRtConfig
	rc.GPIOSlowdown = c.GPIOSlowdown

	matrix, err := rgbmatrix.NewRGBLedMatrix(dc, rc)
	if err != nil {
		panic("failed to create RGB matrix")
	}

	return &matrix
}

func createCanvas(matrix *rgbmatrix.Matrix) *rgbmatrix.Canvas {
	return rgbmatrix.NewCanvas(*matrix)
}
