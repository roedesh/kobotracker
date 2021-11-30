package views

import (
	"cryptokobo/app"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"github.com/leekchan/accounting"
	"github.com/shermp/go-kobo-input/koboin"
)

func renderTrackerScreen(app *app.App, acc accounting.Accounting, lowAcc accounting.Accounting, coinsIndex int) int {
	app.Screen.Clear()
	app.Screen.SetFontSize(175)
	coin := app.Data.Coins[coinsIndex]
	center := float64(app.Screen.State.ScreenHeight) / 2
	app.Screen.GG.DrawStringWrapped(coin.Name, 0, center-350, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	var moneyStr string
	app.Screen.SetFontSize(90)
	if coin.Price < 0.01 {
		moneyStr = lowAcc.FormatMoney(coin.Price)
	} else {
		moneyStr = acc.FormatMoney(coin.Price)
	}
	app.Screen.GG.DrawStringWrapped(moneyStr, 0, center-150, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.DrawChart(coin, 200, center+100, float64(app.Screen.State.ScreenWidth-400), 300)

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
	lowAcc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 10, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}

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

	coinsIndex := renderTrackerScreen(app, acc, lowAcc, 0)
	showNextTicker := time.NewTicker(time.Duration(app.Config.ShowNextInterval) * time.Second)
	updatePricesTicker := time.NewTicker(time.Duration(app.Config.UpdatePriceInterval) * time.Second)

	bus.SubscribeAsync("CHECKINPUT", checkInput, false)
	bus.SubscribeAsync("UPDATE_PRICES", updatePrices, false)
	bus.Publish("CHECKINPUT")

	go func() {
		for {
			select {
			case <-showNextTicker.C:
				coinsIndex = renderTrackerScreen(app, acc, lowAcc, coinsIndex)
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
