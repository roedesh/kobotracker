package views

import (
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/ui"

	"github.com/asaskevich/EventBus"
)

func BootScreen(appConfig *config.AppConfig, bus EventBus.Bus, screen *ui.Screen, coinsDatasource *datasource.CoinsDataSource) {
	screen.Clear()
	screen.GG.DrawString("KoboTracker", 100, 140)
	screen.SetFontSize(42)
	screen.GG.DrawString(appConfig.Version, 100, 220)
	screen.SetFontSize(30)
	screen.GG.DrawString("Created by Ruud Schroën", 100, 350)
	screen.GG.DrawString("Get the latest version @ https://ruud.je/kobotracker", 100, 395)

	if appConfig.SkipCertificateValidation {
		screen.GG.DrawString("Failed to setup SSL certificates!", 100, 475)
	} else {
		screen.GG.DrawString("Successfully setup SSL certificates!", 100, 475)
	}
	screen.GG.DrawString("Loading...", 100, 515)
	screen.RenderFrame()

	coinsDatasource.LoadCoinsForIds(appConfig.Ids)
	err := coinsDatasource.UpdatePricesOfCoins(appConfig.Fiat)
	if err != nil {
		panic(err.Error())
	}

	bus.Publish("ROUTING", "tracker")
}
