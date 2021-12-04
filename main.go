package main

import (
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/handlers"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"cryptokobo/app/views"
	"log"

	"github.com/asaskevich/EventBus"
)

var (
	version string
)

func main() {
	closeLogger := config.SetupLogger()
	defer closeLogger()

	screen := ui.NewScreen()
	defer handlers.HandlePanic(screen)

	c := make(chan bool)
	defer close(c)

	appConfig := config.NewAppConfigFromFile(utils.GetAbsolutePath("config.ini"))
	appConfig.Version = version
	screen.DarkMode = appConfig.DarkMode

	err := config.SetupSSLCertificates()
	if err != nil {
		// If for whatever reason the SSL certificates cannot be found or setup,
		// disable certificate validation to still allow web requests.
		log.Println(err.Error())
		appConfig.SkipCertificateValidation = true
	}

	bus := EventBus.New()
	coinsDatasource := datasource.NewCoinsDataSource(appConfig.SkipCertificateValidation)

	bus.SubscribeAsync("ROUTING", func(routeName string) {
		defer handlers.HandlePanic(screen)

		switch routeName {
		case "boot":
			views.BootScreen(appConfig, bus, screen)
		case "tracker":
			views.TrackerScreen(appConfig, bus, screen, coinsDatasource)
		}
	}, false)

	bus.SubscribeAsync("QUIT", func() {
		c <- true
	}, false)

	bus.Publish("ROUTING", "boot")

	for quit := range c {
		if quit {
			break
		}
	}
}
