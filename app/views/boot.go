package views

import (
	"cryptokobo/app"
	"fmt"
	"time"

	"github.com/asaskevich/EventBus"
)

func BootScreen(app *app.App, bus EventBus.Bus) {
	app.Screen.Clear()
	app.Screen.GG.DrawString("CryptoKobo", 100, 140)
	app.Screen.SetFontSize(30)
	app.Screen.GG.DrawString(fmt.Sprintf("Version: %s", app.Version), 100, 310)
	app.Screen.GG.DrawString("Get the latest version @ https://ruud.je/cryptokobo", 100, 355)
	app.Screen.DrawFrame()

	time.Sleep(2 * time.Second)

	bus.Publish("ROUTING", "tracker")
}
