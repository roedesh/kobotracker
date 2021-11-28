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

func renderCoin(app *app.App, acc accounting.Accounting, lowAcc accounting.Accounting, index int) int {
	app.Screen.Clear()
	app.Screen.SetFontSize(250)
	symbol_y := float64(app.Screen.State.ScreenHeight)/2 - 250
	app.Screen.GG.DrawStringWrapped(app.Data.Coins[index].Symbol, 0, symbol_y, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.SetFontSize(100)
	var moneyStr string

	if app.Data.Coins[index].Price < 0.01 {
		moneyStr = lowAcc.FormatMoney(app.Data.Coins[index].Price)
	} else {
		moneyStr = acc.FormatMoney(app.Data.Coins[index].Price)
	}

	app.Screen.GG.DrawStringWrapped(moneyStr, 0, symbol_y+250, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.SetFontSize(40)
	app.Screen.GG.DrawStringWrapped("Touch screen to exit", 0, float64(app.Screen.State.ScreenHeight)-90, 0, 0, float64(app.Screen.State.ScreenWidth), 1, gg.AlignCenter)

	app.Screen.DrawFrame()

	if index+1 == len(app.Data.Coins) {
		return 0
	}

	return index + 1
}

func TrackerScreen(app *app.App, bus EventBus.Bus, koboin *koboin.TouchDevice) {
	app.Screen.Clear()

	app.Data.LoadCoinsForIds(app.Config.Ids)
	err := app.Data.ApplyPricesToCoins(app.Config.Fiat)
	if err != nil {
		log.Println(err.Error())
	}

	localeInfo := accounting.LocaleInfo[strings.ToUpper(app.Config.Fiat)]
	acc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 2, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}
	lowAcc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 10, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}

	showNextTicker := time.NewTicker(time.Duration(app.Config.ShowNextInterval) * time.Second)
	updatePricesTicker := time.NewTicker(time.Duration(app.Config.UpdatePriceInterval) * time.Second)
	quit := make(chan struct{})

	checkInput := func() {
		_, _, err := koboin.GetInput()
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

	coinIndex := renderCoin(app, acc, lowAcc, 0)

	bus.SubscribeAsync("CHECKINPUT", checkInput, false)
	bus.SubscribeAsync("UPDATE_PRICES", updatePrices, false)

	bus.Publish("CHECKINPUT")

	updatePrices()

	go func() {
		for {
			select {
			case <-showNextTicker.C:
				coinIndex = renderCoin(app, acc, lowAcc, coinIndex)
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
