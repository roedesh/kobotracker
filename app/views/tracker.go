package views

import (
	"bytes"
	"cryptokobo/app/config"
	"cryptokobo/app/datasource"
	"cryptokobo/app/device"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	_ "embed"
	"image"
	"log"
	"sync"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/montanaflynn/stats"
	"github.com/nfnt/resize"
)

//go:embed bolt.png
var boltImageBytes []byte
var (
	batteryLevel      int
	batteryIsCharging bool
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

func renderTrackerScreen(appConfig *config.AppConfig, coinsDatasource *datasource.CoinsDataSource, screen *ui.Screen, coinsIndex int) int {
	lock.Lock()
	defer lock.Unlock()
	screen.Clear()
	screen.SetFontSize(175)
	coin := coinsDatasource.Coins[coinsIndex]
	center := float64(screen.State.ScreenHeight) / 2
	screen.GG.DrawStringWrapped(coin.Name, 0, center-550, 0, 0, float64(screen.State.ScreenWidth), 1, gg.AlignCenter)

	screen.SetFontSize(40)

	screen.DrawProgressBar(float64(screen.State.ScreenWidth-180), 50, 80, 40, float64(batteryLevel))

	if batteryIsCharging == true {
		bolt, _, _ := image.Decode(bytes.NewReader(boltImageBytes))
		resizedBolt := resize.Resize(40, 40, bolt, resize.Lanczos3)
		screen.GG.DrawImage(resizedBolt, int(screen.State.ScreenWidth-230), 50)
	}

	screen.SetFontSize(100)
	moneyStr := utils.GetMoneyString(appConfig.Fiat, float64(coin.CurrentPrice))
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
	batteryLevel = device.GetBatteryLevel()
	batteryIsCharging = device.GetStatus() == "Charging"
	coinsIndex := 0

	c := make(chan bool)
	defer close(c)
	coinsDatasource.LoadCoinsForIds(appConfig.Ids)

	checkInput := func() {
		_, _, err := touchDevice.GetInput()
		if err == nil {
			c <- true
		}
	}

	checkDeviceChanges := func() {
		newBatteryLevel := device.GetBatteryLevel()
		newBatteryIsCharging := device.GetStatus() == "Charging"

		hasChanges := newBatteryLevel != batteryLevel || newBatteryIsCharging != batteryIsCharging
		batteryIsCharging = newBatteryIsCharging
		batteryLevel = newBatteryLevel
		if hasChanges == true {
			_ = renderTrackerScreen(appConfig, coinsDatasource, screen, coinsIndex)
		}
	}

	updatePrices := func() {
		err := coinsDatasource.UpdatePricesOfCoins(appConfig.Fiat)
		if err != nil {
			log.Println(err.Error())
		}
	}

	showNextCoin := func() {
		coinsIndex = renderTrackerScreen(appConfig, coinsDatasource, screen, coinsIndex)
	}

	updatePrices()
	showNextCoin()

	schedule(checkDeviceChanges, 1*time.Second)
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
