package container

import "github.com/slunak/omgo"

func createWeatherClient() *omgo.Client {
	client, err := omgo.NewClient()
	if err != nil {
		panic(err)
	}

	return &client
}
