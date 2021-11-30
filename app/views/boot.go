package views

import (
	"cryptokobo/app"
	"fmt"
	"time"

	"github.com/asaskevich/EventBus"
)

func BootScreen(app *app.App, bus EventBus.Bus) {
	app.Screen.Clear()
	app.Screen.GG.DrawString(fmt.Sprintf("CryptoKobo %s", app.Version), 100, 140)
	app.Screen.SetFontSize(30)
	app.Screen.GG.DrawString("Created by Ruud Schroën", 100, 310)
	app.Screen.GG.DrawString("Get the latest version @ https://ruud.je/kobotracker", 100, 355)

	if app.Data.Insecure {
		app.Screen.GG.DrawString("Failed to setup SSL certificates!", 100, 445)
	} else {
		app.Screen.GG.DrawString("Successfully setup SSL certificates!", 100, 445)
	}

	app.Screen.DrawFrame()

	time.Sleep(2 * time.Second)

	bus.Publish("ROUTING", "tracker")
}
