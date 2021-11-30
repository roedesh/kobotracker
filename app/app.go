// Package app provides functions for initializing and running the application
package app

import (
	"cryptokobo/app/datasource"
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/fogleman/gg"
)

type App struct {
	logFile *os.File

	Data    *datasource.CoinsDataSource
	Config  AppConfig
	Screen  *ui.Screen
	Version string
}

func InitApp(version string) (app *App) {
	app = &App{}
	app.Screen = ui.InitScreen()
	app.Version = version

	logFile, err := os.OpenFile(utils.GetAbsolutePath("debug.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		app.logFile = logFile
		log.SetOutput(app.logFile)
	}

	err = utils.SetupCertificates()
	app.Data = datasource.InitCoinsDataSource(err != nil)

	return app
}

func (app *App) CatchError() {
	if err := recover(); err != nil {
		errorMessage := fmt.Sprintf("%s: %s", err, debug.Stack())
		log.Println(errorMessage)

		app.Screen.Clear()
		app.Screen.SetFontSize(60)
		app.Screen.GG.DrawString("Oops! Something went wrong!", 100, 140)
		app.Screen.SetFontSize(30)
		app.Screen.GG.DrawStringWrapped(fmt.Sprintf("%s", err), 100, 350, 0, 0, float64(app.Screen.State.ScreenWidth-200), 2, gg.AlignLeft)
		app.Screen.GG.DrawString(fmt.Sprintf("Stacktrace was written to \"%s\"", utils.GetAbsolutePath("log.txt")), 100, float64(app.Screen.State.ScreenHeight)-180)
		app.Screen.GG.DrawString("Exiting app in 5 seconds..", 100, float64(app.Screen.State.ScreenHeight)-130)
		app.Screen.SetFontSize(25)

		app.Screen.DrawFrame()

		time.Sleep(5 * time.Second)

		app.Teardown()
		os.Exit(1)
	}
}

func (app *App) Exit() {
	app.CatchError()
	app.Teardown()
}

func (app *App) LoadConfig() {
	app.Config = getConfigFromIniFile()
}

func (app *App) Teardown() {
	app.Screen.Close()
	if app.logFile != nil {
		app.logFile.Close()
	}
}
