package handlers

import (
	"cryptokobo/app/ui"
	"cryptokobo/app/utils"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/fogleman/gg"
)

func HandlePanic(screen *ui.Screen) {
	if err := recover(); err != nil {
		errorMessage := fmt.Sprintf("%s: %s", err, debug.Stack())
		log.Println(errorMessage)

		screen.Clear()
		screen.SetFontSize(60)
		screen.GG.DrawString("Oops! Something went wrong!", 100, 140)
		screen.SetFontSize(30)
		screen.GG.DrawStringWrapped(fmt.Sprintf("%s", err), 100, 350, 0, 0, float64(screen.State.ScreenWidth-200), 2, gg.AlignLeft)
		screen.GG.DrawString(fmt.Sprintf("Stacktrace was written to \"%s\"", utils.GetAbsolutePath("log.txt")), 100, float64(screen.State.ScreenHeight)-180)
		screen.GG.DrawString("Exiting app in 5 seconds..", 100, float64(screen.State.ScreenHeight)-130)
		screen.SetFontSize(25)

		screen.DrawFrame()

		time.Sleep(5 * time.Second)

		screen.Close()
		os.Exit(1)
	}
}
