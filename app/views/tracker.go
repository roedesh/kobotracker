package views

import (
	"cryptokobo/app"
	"cryptokobo/app/device"
	"cryptokobo/app/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/leekchan/accounting"
	"github.com/shermp/go-kobo-input/koboin"
)

func renderTrackerScreen(app *app.App, acc accounting.Accounting, middleAcc accounting.Accounting, lowAcc accounting.Accounting, coinsIndex int) int {
	app.Screen.Clear()
	app.Screen.SetFontSize(175)
	coin := app.Data.Coins[coinsIndex]
	center := float64(app.Screen.State.ScreenHeight) / 2
	app.Screen.GG.DrawStringWrapped(coin.Name, 0, center-500, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.SetFontSize(40)
	batteryLevel := device.GetBatteryLevel()
	app.Screen.GG.DrawStringWrapped(fmt.Sprintf("%d%%", batteryLevel), 100, 50, 0, 0, float64(app.Screen.State.ScreenWidth-200), 0, gg.AlignRight)

	min, max := coin.GetBaselinePrices()
	app.Screen.SetFontSize(90)

	moneyStr := utils.GetMoneyString(app.Config.Fiat, float64(coin.Price))
	minMoneyStr := utils.GetMoneyString(app.Config.Fiat, min)
	maxMoneyStr := utils.GetMoneyString(app.Config.Fiat, max)

	app.Screen.GG.DrawStringWrapped(moneyStr, 0, center-300, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.DrawChart(coin, minMoneyStr, maxMoneyStr, 200, center-50, float64(app.Screen.State.ScreenWidth-400), 500)

	app.Screen.SetFontSize(40)
	app.Screen.GG.DrawStringWrapped("Touch screen to exit", 0, float64(app.Screen.State.ScreenHeight)-90, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.DrawFrame()

	if coinsIndex+1 == len(app.Data.Coins) {
		return 0
	}

	return coinsIndex + 1
}

func TrackerScreen(app *app.App, bus EventBus.Bus) {
	touchPath := "/dev/input/event1"
	touchInput := koboin.New(touchPath, int(app.Screen.State.ScreenWidth), int(app.Screen.State.ScreenHeight))
	if touchInput == nil {
		panic("Could not get touch input")
	}

	app.Screen.Clear()

	app.Data.LoadCoinsForIds(app.Config.Ids)

	localeInfo := accounting.LocaleInfo[strings.ToUpper(app.Config.Fiat)]
	acc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 2, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}
	middleAcc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 4, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}
	lowAcc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 8, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}

	quit := make(chan struct{})

	checkInput := func() {
		_, _, err := touchInput.GetInput()
		if err == nil {
			close(quit)
		}
	}

	updatePrices := func() {
		err := app.Data.ApplyPricesToCoins(app.Config.Fiat)
		if err != nil {
			log.Println(err.Error())
		}
	}

	updatePrices()

	coinsIndex := renderTrackerScreen(app, acc, middleAcc, lowAcc, 0)
	showNextTicker := time.NewTicker(time.Duration(app.Config.ShowNextInterval) * time.Second)
	updatePricesTicker := time.NewTicker(time.Duration(app.Config.UpdatePriceInterval) * time.Second)

	bus.SubscribeAsync("CHECKINPUT", checkInput, false)
	bus.SubscribeAsync("UPDATE_PRICES", updatePrices, false)
	bus.Publish("CHECKINPUT")

	go func() {
		for {
			select {
			case <-showNextTicker.C:
				coinsIndex = renderTrackerScreen(app, acc, middleAcc, lowAcc, coinsIndex)
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
