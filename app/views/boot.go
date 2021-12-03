package views

import (
	"cryptokobo/app/config"
	"cryptokobo/app/ui"
	"time"

	"github.com/asaskevich/EventBus"
)

func BootScreen(config *config.AppConfig, bus EventBus.Bus, screen *ui.Screen) {
	screen.Clear()
	screen.GG.DrawString("KoboTracker", 100, 140)
	screen.SetFontSize(42)
	screen.GG.DrawString(config.Version, 100, 220)
	screen.SetFontSize(30)
	screen.GG.DrawString("Created by Ruud Schroën", 100, 350)
	screen.GG.DrawString("Get the latest version @ https://ruud.je/kobotracker", 100, 395)

	if config.SkipCertificateValidation {
		screen.GG.DrawString("Failed to setup SSL certificates!", 100, 475)
	} else {
		screen.GG.DrawString("Successfully setup SSL certificates!", 100, 475)
	}
	screen.GG.DrawString("Loading...", 100, 515)
	screen.RenderFrame()

	time.Sleep(1 * time.Second)

	bus.Publish("ROUTING", "tracker")
}
