// Package app provides functions for initializing and running the application
package app

import (
	"cryptokobo/app/network"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"gopkg.in/ini.v1"
)

type App struct {
	logFile *os.File

	CMCApiKey string
	CoinGecko *network.CoinGeckoClient
	Ids       []string
	Screen    *ui.Screen
	Version   string
}

func InitApp(version string) (app *App) {
	app = &App{}
	app.Screen = ui.InitScreen()
	app.CoinGecko = network.InitCoinGecko()
	app.Version = version

	logFile, err := os.OpenFile(utils.GetAbsolutePath("log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		app.logFile = logFile
		log.SetOutput(app.logFile)
	}

	return app
}

func (app *App) LoadConfig() {
	config, err := ini.Load(utils.GetAbsolutePath("config.ini"))
	if err != nil {
		panic(fmt.Sprintf("Could not load \"%s\".", utils.GetAbsolutePath("config.ini")))
	}

	ids := config.Section("").Key("ids").String()
	app.Ids = strings.Fields(ids)
	if len(app.Ids) == 0 {
		panic("No CoinGecko ids set. Add \"ids\" to your \"config.ini\".")
	}
}

func (app *App) TearDown() {
	if err := recover(); err != nil {
		error_message := fmt.Sprintf("%s: %s", err, debug.Stack())
		log.Println(error_message)

		app.Screen.Clear()
		app.Screen.SetFontSize(60)
		app.Screen.GG.DrawString("Oops! Something went wrong!", 100, 140)
		app.Screen.SetFontSize(30)
		app.Screen.GG.DrawString(fmt.Sprintf("Stacktrace was written to \"%s\"", utils.GetAbsolutePath("log.txt")), 100, 220)
		app.Screen.GG.DrawString("Exiting app in 5 seconds..", 100, 270)
		app.Screen.SetFontSize(25)
		app.Screen.GG.DrawStringWrapped(error_message, 100, 350, 0, 0, float64(app.Screen.State.ScreenWidth-100), 2, gg.AlignLeft)
		app.Screen.DrawFrame()

		time.Sleep(5 * time.Second)
	}

	app.Screen.Close()
	if app.logFile != nil {
		app.logFile.Close()
	}
}
