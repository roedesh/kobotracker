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
	defer cryptokobo.TearDown()

	touchPath := "/dev/input/event1"
	t := koboin.New(touchPath, 1080, 1440)
	if t == nil {
		return
	}
	defer t.Close()

	bus := EventBus.New()

	c := make(chan bool)
	defer close(c)

	bus.SubscribeAsync("QUIT", func() {
		c <- true
		return
	}, false)

	bus.SubscribeAsync("ROUTING", func(routeName string) {
		switch routeName {
		case "boot":
			cryptokobo.LoadConfig()
			views.BootScreen(cryptokobo, bus)
		case "tracker":
			views.TrackerScreen(cryptokobo, bus, t)
		}
	}, false)

	bus.Publish("ROUTING", "boot")

	for quit := range c {
		if quit {
			break
		}
	}
}
