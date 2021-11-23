// Package app provides functions for initializing and running the application
package app

import (
	"cryptokobo/app/fs"
	"cryptokobo/app/ui"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type App struct {
	coinMarketCapApiKey string
	logFile             *os.File
	tickers             []string

	Screen  *ui.Screen
	Version string
}

func InitApp(version string) (app *App) {
	app = &App{}
	app.Screen = ui.InitScreen()
	app.Version = version

	return app
}

func (app *App) SetupLogger() {
	logFile, err := os.OpenFile(fs.GetAbsolutePath("log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	app.logFile = logFile
	log.SetOutput(app.logFile)
}

func (app *App) LoadConfig() error {
	config, err := ini.Load(fs.GetAbsolutePath("config.ini"))
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load \"%s\".", fs.GetAbsolutePath("config.ini")))
	}

	app.coinMarketCapApiKey = config.Section("").Key("cmc_api_key").String()
	if app.coinMarketCapApiKey == "" {
		return errors.New("CoinMarketCap API key not set. Add \"cmc_api_key\" to your \"config.ini\".")
	}

	tickers := config.Section("").Key("tickers").String()
	app.tickers = strings.Fields(tickers)

	if len(app.tickers) == 0 {
		return errors.New("No tickers set. Add \"tickers\" to your \"config.ini\".")
	}

	return nil
}

func (app *App) TearDown() {
	app.Screen.Close()
	if app.logFile != nil {
		app.logFile.Close()
	}
}
