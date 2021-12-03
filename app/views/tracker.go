package views

import (
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/device"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"log"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/montanaflynn/stats"
)

func renderTrackerScreen(appConfig *config.AppConfig, coinsDatasource *datasource.CoinsDataSource, screen *ui.Screen, coinsIndex int) int {
	screen.Clear()
	screen.SetFontSize(175)
	coin := coinsDatasource.Coins[coinsIndex]
	center := float64(screen.State.ScreenHeight) / 2
	screen.GG.DrawStringWrapped(coin.Name, 0, center-550, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	screen.SetFontSize(40)
	batteryLevel := device.GetBatteryLevel()
	screen.DrawProgressBar(float64(screen.State.ScreenWidth-180), 50, 80, 40, float64(batteryLevel))

	screen.SetFontSize(100)
	moneyStr := utils.GetMoneyString(appConfig.Fiat, float64(coin.Price))
	screen.GG.DrawStringWrapped(moneyStr, 0, center-340, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	min, _ := stats.Min(coin.PricePoints)
	max, _ := stats.Max(coin.PricePoints)

	screen.DrawChart(coin.PricePoints, min, max, appConfig.Fiat, 175, center, float64(screen.State.ScreenWidth-400), 425)

	screen.SetFontSize(40)
	screen.GG.DrawStringWrapped("Touch screen to exit", 0, float64(screen.State.ScreenHeight)-90, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	screen.RenderFrame()

	if coinsIndex+1 == len(coinsDatasource.Coins) {
		return 0
	}

	return coinsIndex + 1
}

func TrackerScreen(appConfig *config.AppConfig, bus EventBus.Bus, screen *ui.Screen, coinsDatasource *datasource.CoinsDataSource) {
	touchDevice := device.GetTouchDevice(int(screen.State.ScreenWidth), int(screen.State.ScreenHeight))

	quit := make(chan struct{})
	coinsDatasource.LoadCoinsForIds(appConfig.Ids)

	checkInput := func() {
		_, _, err := touchDevice.GetInput()
		if err == nil {
			close(quit)
		}
	}

	updatePrices := func() {
		err := coinsDatasource.ApplyPricesToCoins(appConfig.Fiat)
		if err != nil {
			log.Println(err.Error())
		}
	}

	updatePrices()

	coinsIndex := renderTrackerScreen(appConfig, coinsDatasource, screen, 0)
	showNextTicker := time.NewTicker(time.Duration(appConfig.ShowNextInterval) * time.Second)
	updatePricesTicker := time.NewTicker(time.Duration(appConfig.UpdatePriceInterval) * time.Second)

	bus.SubscribeAsync("CHECKINPUT", checkInput, false)
	bus.SubscribeAsync("UPDATE_PRICES", updatePrices, false)
	bus.Publish("CHECKINPUT")

	go func() {
		for {
			select {
			case <-showNextTicker.C:
				coinsIndex = renderTrackerScreen(appConfig, coinsDatasource, screen, coinsIndex)
			case <-updatePricesTicker.C:
				bus.Publish("UPDATE_PRICES")
			case <-quit:
				showNextTicker.Stop()
				updatePricesTicker.Stop()
				bus.Unsubscribe("CHECKINPUT", checkInput)
				bus.Unsubscribe("UPDATE_PRICES", updatePrices)
				bus.Publish("QUIT")
			}
		}
	}()
}
