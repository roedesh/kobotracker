package main

import (
	"fmt"
	"time"

	"cryptokobo/app"
	"cryptokobo/app/ui"
)

var (
	version string
)

func main() {
	app := app.InitApp(version)
	defer app.TearDown()

	screen := app.Screen

	screen.Clear()
	screen.DrawText("CryptoKobo", 100, 200)
	screen.SetFontSettings(ui.FontConfig{Size: 40})
	screen.DrawText(fmt.Sprintf("Version: %s", app.Version), 100, 310)
	screen.DrawText("Get the latest version @ https://ruud.je/cryptokobo", 100, 365)

	configErr := app.LoadConfig()
	if configErr != nil {
		screen.DrawText(configErr.Error(), 100, 460)
	} else {
		screen.DrawText("Successfully loaded config", 100, 460)
	}

	screen.DrawFrame()

	time.Sleep(10 * time.Second)

	return
}
