package lib

import (
	"github.com/shermp/go-fbink-v2/gofbink"
)

type PrintConfig struct {
	FontSize int16
	X        int16
	Y        int16
}

type Screen struct {
	fb    *gofbink.FBInk
	state gofbink.FBInkState
}

func InitScreen() (screen *Screen) {
	screen = &Screen{}
	screen.state = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{}
	screen.fb = gofbink.New(&fbinkOpts, &rOpts)

	screen.fb.Open()
	screen.fb.Init(&fbinkOpts)
	screen.fb.AddOTfont("/mnt/onboard/.adds/cryptokobo/inc/font.ttf", gofbink.FntRegular)
	screen.fb.GetState(&fbinkOpts, &screen.state)

	return screen
}

func (screen *Screen) ClearScreen() {
	screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})
}

func (screen *Screen) PrintOT(text string, config PrintConfig) {
	screen.fb.PrintOT(text, &gofbink.FBInkOTConfig{
		Margins: struct {
			Top    int16
			Bottom int16
			Left   int16
			Right  int16
		}{
			Top:  config.Y,
			Left: config.X,
		},
		SizePx: uint16(config.FontSize),
	}, &gofbink.FBInkConfig{})
}

func (screen *Screen) Print(a string) {
	screen.fb.Println(a)
}

func (screen *Screen) Close() {
	screen.fb.FreeOTfonts()
	screen.fb.Close()
}
