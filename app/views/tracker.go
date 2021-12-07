package views

import (
	"cryptokobo/app/assets"
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/device"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"log"
	"sync"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/montanaflynn/stats"
)

var (
	batteryLevel      int
	batteryIsCharging bool
	coinIndex         int
	lock              sync.Mutex
)

func schedule(f func(), interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}

func renderTrackerScreen(appConfig *config.AppConfig, coinsDatasource *datasource.CoinsDataSource, screen *ui.Screen, updateCoinIndex bool) {
	lock.Lock()

	if updateCoinIndex == true {
		nextCoinIndex := coinIndex + 1
		if nextCoinIndex >= len(coinsDatasource.Coins) {
			nextCoinIndex = 0
		}

		coinIndex = nextCoinIndex
	}

	coin := coinsDatasource.Coins[coinIndex]
	center := float64(screen.State.ScreenHeight) / 2

	screen.Clear()

	screen.GG.DrawImage(assets.SignOutImage, 80, 50)
	screen.DrawProgressBar(float64(screen.State.ScreenWidth-180), 50, 80, 40, float64(batteryLevel))
	if batteryIsCharging == true {
		screen.GG.DrawImage(assets.BoltImage, int(screen.State.ScreenWidth-230), 50)
	}

	screen.SetFontSize(175)
	screen.GG.DrawStringWrapped(coin.Name, 0, center-500, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	screen.SetFontSize(100)
	moneyStr := utils.GetMoneyString(appConfig.Fiat, float64(coin.CurrentPrice))
	screen.GG.DrawStringWrapped(moneyStr, 0, center-290, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	min, _ := stats.Min(coin.PricePoints)
	max, _ := stats.Max(coin.PricePoints)

	screen.DrawChart(coin.PricePoints, min, max, appConfig.Fiat, 175, center+50, float64(screen.State.ScreenWidth-400), 425)

	screen.RenderFrame()

	lock.Unlock()
}

func TrackerScreen(appConfig *config.AppConfig, bus EventBus.Bus, screen *ui.Screen, coinsDatasource *datasource.CoinsDataSource) {
	touchDevice := device.GetTouchDevice(int(screen.State.ScreenWidth), int(screen.State.ScreenHeight))
	batteryLevel = device.GetBatteryLevel()
	batteryIsCharging = device.GetStatus() == "Charging"
	coinIndex = -1

	c := make(chan bool)
	defer close(c)

	checkInput := func() {
		for {
			rx, ry, err := touchDevice.GetInput()
			if err == nil && rx <= 150 && ry <= 150 {
				break
			}
		}
		c <- true
	}

	checkDeviceChanges := func() {
		newBatteryLevel := device.GetBatteryLevel()
		newBatteryIsCharging := device.GetStatus() == "Charging"

		hasChanges := newBatteryLevel != batteryLevel || newBatteryIsCharging != batteryIsCharging
		batteryIsCharging = newBatteryIsCharging
		batteryLevel = newBatteryLevel
		if hasChanges == true {
			renderTrackerScreen(appConfig, coinsDatasource, screen, false)
		}
	}

	updatePrices := func() {
		err := coinsDatasource.UpdatePricesOfCoins(appConfig.Fiat)
		if err != nil {
			log.Println(err.Error())
		}
	}

	showNextCoin := func() {
		renderTrackerScreen(appConfig, coinsDatasource, screen, true)
	}

	showNextCoin()

	schedule(checkDeviceChanges, 1500*time.Millisecond)
	schedule(updatePrices, time.Duration(appConfig.UpdatePriceInterval)*time.Second)
	schedule(showNextCoin, time.Duration(appConfig.ShowNextInterval)*time.Second)

	bus.SubscribeAsync("CHECKINPUT", checkInput, false)
	bus.Publish("CHECKINPUT")

	for quit := range c {
		if quit {
			bus.Unsubscribe("CHECKINPUT", checkInput)
			bus.Publish("QUIT")
			break
		}
	}
}
