// Package ui provides functions for drawing stuff on the screen.
package ui

import (
	"cryptokobo/app/fs"
	"image"

	"github.com/fogleman/gg"
	"github.com/shermp/go-fbink-v2/gofbink"
)

type PrintConfig struct {
	FontSize int16
	X        int16
	Y        int16
}

type FontConfig struct {
	Size float64
}

type Screen struct {
	fb        *gofbink.FBInk
	ggContext *gg.Context
	ggRGBA    *image.RGBA
	state     gofbink.FBInkState
}

func InitScreen() (screen *Screen) {
	screen = &Screen{}
	screen.state = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{}
	screen.fb = gofbink.New(&fbinkOpts, &rOpts)

	screen.fb.Open()
	screen.fb.Init(&fbinkOpts)
	screen.fb.GetState(&fbinkOpts, &screen.state)

	screen.ggRGBA = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(screen.state.ScreenWidth), int(screen.state.ScreenHeight)}})
	screen.ggContext = gg.NewContextForRGBA(screen.ggRGBA)
	screen.ggContext.SetRGBA(0, 0, 0, 1)
	screen.SetFontSettings(FontConfig{Size: 80})

	return screen
}

func (screen *Screen) SetFontSettings(fontConfig FontConfig) {
	screen.ggContext.LoadFontFace(fs.GetAbsolutePath("assets/font.ttf"), fontConfig.Size)
}

func (screen *Screen) Clear() {
	screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})
}

func (screen *Screen) DrawText(text string, x float64, y float64) {
	screen.ggContext.DrawString(text, x, y)
}

func (screen *Screen) DrawFrame() {
	screen.fb.PrintRBGA(0, 0, screen.ggRGBA, &gofbink.FBInkConfig{})
}

func (screen *Screen) Close() {
	screen.fb.FreeOTfonts()
	screen.fb.Close()
}
