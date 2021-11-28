package views

import (
	"cryptokobo/app"
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/shermp/go-kobo-input/koboin"
)

func TrackerScreen(app *app.App, bus EventBus.Bus, koboin *koboin.TouchDevice) {
	app.Screen.Clear()

	app.Data.LoadCoinsForIds(app.Config.Ids)
	err := app.Data.ApplyPricesToCoins(app.Config.Fiat)
	if err != nil {
		log.Println(err.Error())
	}

	app.Screen.SetFontSize(250)
	symbol_y := float64(app.Screen.State.ScreenHeight)/2 - 250
	app.Screen.GG.DrawStringWrapped(app.Data.Coins[0].Symbol, 0, symbol_y, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.SetFontSize(100)
	app.Screen.GG.DrawStringWrapped(fmt.Sprintf("€%f", app.Data.Coins[0].Price), 0, symbol_y+250, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.DrawFrame()

	time.Sleep(5 * time.Second)

	bus.Publish("QUIT")
}
