package container

import "strconv"

type Config struct {
	RGBMatrixConfig
	Resource
	Weather
}

type RGBMatrixConfig struct {
	Rows                   int    // "led-rows" number of rows supported
	Cols                   int    // "led-cols" number of columns supported
	Parallel               int    // "led-parallel" number of daisy-chained panels
	ChainLength            int    // "led-chain" number of displays daisy-chained
	Brightness             int    // "brightness" brightness (0-100)
	GPIOMapping            string // "led-gpio-mapping" Name of GPIO mapping used.
	ShowRefresh            bool   // "led-show-refresh" Show refresh rate.
	InverseColors          bool   // "led-inverse" Switch if your matrix has inverse colors on.
	DisableHardwarePulsing bool   // "led-no-hardware-pulse" Don't use hardware pin-pulse generation.
	GPIOSlowdown           int    // "led-gpio-slowdown" Slow down writing to GPIO.
}

type Resource struct {
	fontsPath        string
	weatherIconsPath string
}

type Weather struct {
	lat, lon          float64
	temperatureUnit   string
	windspeedUnit     string
	precipitationUnit string
	timezone          string
}

func getConfig() *Config {
	return &Config{
		RGBMatrixConfig: RGBMatrixConfig{
			Rows:                   64,
			Cols:                   64,
			Parallel:               1,
			ChainLength:            1,
			Brightness:             70,
			GPIOMapping:            "regular",
			ShowRefresh:            false,
			InverseColors:          false,
			DisableHardwarePulsing: false,
		},
		Resource: Resource{
			fontsPath:        "resources/fonts/",
			weatherIconsPath: "resources/weather/",
		},
		Weather: Weather{
			lat:               52.3738,
			lon:               4.8910,
			temperatureUnit:   "celsius",
			windspeedUnit:     "kmh",
			precipitationUnit: "mm",
			timezone:          "Europe/Amsterdam",
		},
	}
}

func convertConfigFlags(lat, lon, tz *string) (latFloat, lonFloat float64, tzString string, err error) {
	latFloat, err = strconv.ParseFloat(*lat, 64)
	if err != nil {
		return
	}

	lonFloat, err = strconv.ParseFloat(*lon, 64)
	if err != nil {
		return
	}

	tzString = *tz

	return
}
