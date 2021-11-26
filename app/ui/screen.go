package ui

import (
	"cryptokobo/app/utils"
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
	fb   *gofbink.FBInk
	rgba *image.RGBA

	GG    *gg.Context
	State gofbink.FBInkState
}

func InitScreen() (screen *Screen) {
	screen = &Screen{}
	screen.State = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{}
	screen.fb = gofbink.New(&fbinkOpts, &rOpts)

	screen.fb.Open()
	screen.fb.Init(&fbinkOpts)
	screen.fb.GetState(&fbinkOpts, &screen.State)

	screen.rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(screen.State.ScreenWidth), int(screen.State.ScreenHeight)}})
	screen.GG = gg.NewContextForRGBA(screen.rgba)
	screen.GG.SetRGB(0, 0, 0)
	screen.SetFontSize(80)

	screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})

	return screen
}

func (screen *Screen) SetFontSize(size float64) {
	screen.GG.LoadFontFace(utils.GetAbsolutePath("assets/font.ttf"), size)
}

// func (screen *Screen) DrawImageFile(path string, x int, y int) error {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	image, _, err := image.Decode(f)
// 	screen.GG.DrawImage(image, x, y)
// 	return nil
// }

func (screen *Screen) DrawFrame() {
	screen.fb.PrintRBGA(0, 0, screen.rgba, &gofbink.FBInkConfig{})
}

func (screen *Screen) Clear() {
	screen.GG.SetRGB(1, 1, 1)
	screen.GG.Clear()
	screen.GG.SetRGB(0, 0, 0)
}

func (screen *Screen) Close() {
	screen.fb.Close()
}
