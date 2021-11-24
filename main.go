package main

import (
	"fmt"
	"time"

	"cryptokobo/app"
	"cryptokobo/app/cmc"
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
	screen.DrawText("CryptoKobo", 100, 140)
	screen.SetFontSettings(ui.FontConfig{Size: 30})
	screen.DrawText(fmt.Sprintf("Version: %s", app.Version), 100, 310)
	screen.DrawText("Get the latest version @ https://ruud.je/cryptokobo", 100, 355)

	configErr := app.LoadConfig()
	if configErr != nil {
		screen.DrawText(configErr.Error(), 100, 450)
	} else {
		screen.DrawText("Successfully loaded config", 100, 450)
	}

	screen.DrawFrame()

	client := cmc.InitClient(app.CMCApiKey)

	_, err := client.GetMap()
	if err != nil {
		screen.DrawText("Failed to load cryptocurrency map from CoinMarketCap", 100, 505)
	} else {
		screen.DrawText("Successfully loaded cryptocurrency map", 100, 505)
	}

	screen.DrawFrame()

	time.Sleep(10 * time.Second)

	return
}
