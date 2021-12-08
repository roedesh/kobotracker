package views

import (
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/ui"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
)

func BootScreen(appConfig *config.AppConfig, bus EventBus.Bus, screen *ui.Screen, coinsDatasource *datasource.CoinsDataSource) {
	center := float64(screen.State.ScreenHeight) / 2

	screen.Clear()
	screen.SetFontSize(120)
	screen.GG.DrawStringWrapped("Kobotracker", 0, center-250, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)
	screen.SetFontSize(70)
	screen.GG.DrawStringWrapped(appConfig.Version, 0, center-100, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)
	screen.SetFontSize(50)
	screen.GG.DrawStringWrapped("by Ruud Schroën", 0, float64(screen.State.ScreenHeight-100), 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	screen.RenderFrame()

	coinsDatasource.LoadCoinsForIds(appConfig.Ids)
	err := coinsDatasource.UpdatePricesOfCoins(appConfig.Fiat)
	if err != nil {
		panic(err.Error())
	}

	bus.Publish("ROUTING", "tracker")
}
