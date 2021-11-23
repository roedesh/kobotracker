package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"cryptokobo/ui"
)

var (
	version string
)

func main() {
	screen := ui.InitScreen()

	screen.ClearScreen()
	defer screen.Close()

	screen.DrawText("CryptoKobo", 100, 200)
	screen.SetFontSettings(ui.FontConfig{Size: 40})
	screen.DrawText(fmt.Sprintf("Version: %s", version), 100, 310)
	screen.DrawText("Get the latest version @ https://ruud.je/cryptokobo", 100, 365)

	screen.DrawFrame()

	if _, err := os.Stat("./config.ini"); errors.Is(err, os.ErrNotExist) {
		screen.DrawText("Could not load \"config.ini\"", 100, 460)
		screen.DrawText("Add your \"config.ini\" file to \"/.adds/cryptokobo\"", 100, 515)
		screen.DrawText("Rebooting in 10 seconds...", 100, 570)
	}

	screen.DrawFrame()

	time.Sleep(10 * time.Second)

	return
}
