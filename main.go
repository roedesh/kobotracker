package main

import (
	"cryptokobo/app"
	"cryptokobo/app/views"

	"github.com/asaskevich/EventBus"
	"github.com/shermp/go-kobo-input/koboin"
)

var (
	version string
)

func main() {
	cryptokobo := app.InitApp(version)
	defer cryptokobo.Exit()

	touchPath := "/dev/input/event1"
	touchInput := koboin.New(touchPath, 1080, 1440)
	if touchInput == nil {
		return
	}
	defer touchInput.Close()

	bus := EventBus.New()

	c := make(chan bool)
	defer close(c)

	bus.SubscribeAsync("QUIT", func() {
		c <- true
	}, false)

	bus.SubscribeAsync("ROUTING", func(routeName string) {
		defer cryptokobo.CatchError()

		switch routeName {
		case "boot":
			cryptokobo.LoadConfig()
			views.BootScreen(cryptokobo, bus)
		case "tracker":
			views.TrackerScreen(cryptokobo, bus, touchInput)
		}
	}, false)

	bus.Publish("ROUTING", "boot")

	for quit := range c {
		if quit {
			break
		}
	}
}
