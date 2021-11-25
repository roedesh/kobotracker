package main

import (
	"fmt"
	"os"
	"time"

	"cryptokobo/app"
	"cryptokobo/app/cmc"
	"cryptokobo/app/fs"
	"cryptokobo/app/network"
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
	screen.DrawText("Fetching data from CoinMarketCap...", 100, 505)
	screen.DrawFrame()

	cryptocurrencies, err := client.GetIds(app.Tickers)
	if err != nil {
		screen.DrawText("Failed to load cryptocurrency map from CoinMarketCap", 100, 550)
	} else {
		screen.DrawText("Successfully loaded ID map", 100, 550)
	}

	screen.DrawFrame()

	cryptocurrencies, err = client.GetPrices(cryptocurrencies)
	if err != nil {
		screen.DrawText("Failed to load prices from CoinMarketCap", 100, 595)
	} else {
		screen.DrawText("Successfully loaded prices", 100, 595)
	}

	screen.DrawFrame()

	cryptocurrencies, err = client.GetLogos(cryptocurrencies)
	if err != nil {
		screen.DrawText("Failed to load logos from CoinMarketCap", 100, 640)
	} else {
		screen.DrawText("Successfully downloaded logos", 100, 640)
	}

	screen.DrawFrame()

	err = os.MkdirAll(fs.GetAbsolutePath("assets/.downloads"), 0777)
	if err != nil {
		screen.DrawText(err.Error(), 100, 685)
	} else {
		for _, cryptocurrency := range cryptocurrencies {
			if cryptocurrency.Logo != "" {
				imagePath := fs.GetAbsolutePath(fmt.Sprintf("assets/.downloads/logo_%d.png", cryptocurrency.ID))
				if _, err := os.Stat(imagePath); err != nil {
					err = network.DownloadFile(imagePath, cryptocurrency.Logo)
					if err != nil {
						screen.DrawText(err.Error(), 100, 685)
					}
				}
			}
		}
	}

	time.Sleep(1 * time.Second)

	screen.Clear()

	for index, cryptocurrency := range cryptocurrencies {
		if cryptocurrency.Logo != "" {
			imagePath := fs.GetAbsolutePath(fmt.Sprintf("assets/.downloads/logo_%d.png", cryptocurrency.ID))
			screen.DrawImageFile(imagePath, 100, ((index+1)*80)+100)
			screen.DrawText(fmt.Sprintf("%s - €%f", cryptocurrency.Symbol, cryptocurrency.Quote.EUR.Price), 175, float64(((index+1)*80)+100))
		}
	}

	screen.DrawFrame()

	time.Sleep(5 * time.Second)

	return
}
