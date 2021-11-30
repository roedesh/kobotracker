package ui

import (
	"cryptokobo/app/datasource"
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

func (screen *Screen) Clear() {
	screen.GG.SetRGB(1, 1, 1)
	screen.GG.Clear()
	screen.GG.SetRGB(0, 0, 0)
}

func (screen *Screen) Close() {
	screen.fb.Close()
}

func (screen *Screen) DrawFrame() {
	screen.fb.PrintRBGA(0, 0, screen.rgba, &gofbink.FBInkConfig{})
}

func (screen *Screen) DrawChart(coin datasource.Coin, minMoneyStr string, maxMoneyStr string, x float64, y float64, width float64, height float64) (int, int) {
	min, max := coin.GetBaselinePrices()

	var minIndex, maxIndex int

	screen.SetFontSize(30)

	screen.GG.DrawStringWrapped(maxMoneyStr, x, y-35, 0, 0, width, 0, gg.AlignRight)
	screen.GG.DrawStringWrapped(minMoneyStr, x, y+height+15, 0, 0, width, 0, gg.AlignLeft)

	for index, price := range coin.PricePoints {
		stepWidth := width / float64(len(coin.PricePoints))
		newX := float64(x + (stepWidth * float64(index)))
		percentage := (price - min) / (max - min)
		newY := y + (height - (height * percentage))

		if index == 0 {
			screen.GG.MoveTo(newX, newY)
		} else {
			screen.GG.LineTo(newX, newY)
		}
	}

	screen.GG.SetLineWidth(2.25)
	screen.GG.Stroke()

	screen.GG.SetRGBA(0, 0, 0, 0.2)
	screen.GG.MoveTo(x, y)
	screen.GG.LineTo(x+width, y)
	screen.GG.MoveTo(x, y+height)
	screen.GG.LineTo(x+width, y+height)
	screen.GG.Stroke()

	screen.GG.SetRGB(0, 0, 0)

	return minIndex, maxIndex
}

func (screen *Screen) SetFontSize(size float64) {
	screen.GG.LoadFontFace(utils.GetAbsolutePath("assets/font.ttf"), size)
}
