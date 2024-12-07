package main

import (
	"flag"
	"led-clock/infrastructure"
	"led-clock/infrastructure/container"
)

func main() {
	c, err := container.GetContainer()
	if err != nil {
		panic(err)
	}

	latPointer := flag.String("latitude", "52.3738", "Latitude of the location")
	lonPointer := flag.String("longitude", "4.8910", "Longitude of the location")
	tzPointer := flag.String("timezone", "Europe/Amsterdam", "Timezone of the location")

	flag.Parse()

	err = c.SetConfigFlags(latPointer, lonPointer, tzPointer)
	if err != nil {
		panic(err.Error())
	}

	err = infrastructure.Start(c)
	if err != nil {
		panic(err.Error())
	}
}
