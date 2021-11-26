package main

import (
	"fmt"
	"log"
	"time"

	"cryptokobo/app"
)

var (
	version string
)

func main() {
	app := app.InitApp(version)
	defer app.TearDown()

	app.Screen.Clear()
	app.Screen.GG.DrawString("CryptoKobo", 100, 140)
	app.Screen.SetFontSize(30)
	app.Screen.GG.DrawString(fmt.Sprintf("Version: %s", app.Version), 100, 310)
	app.Screen.GG.DrawString("Get the latest version @ https://ruud.je/cryptokobo", 100, 355)
	app.Screen.DrawFrame()

	time.Sleep(2 * time.Second)

	app.LoadConfig()

	coins := app.CoinGecko.GetCoinsForIds(app.Ids)
	coins, err := app.CoinGecko.ApplyPricesToCoins(coins)
	if err != nil {
		log.Println(err.Error())
	} else {
		app.Screen.GG.DrawString(coins[0].Name, 100, 600)
	}

	app.Screen.DrawFrame()

	time.Sleep(5 * time.Second)
}
