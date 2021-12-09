package ui

import (
	"fmt"
	"image"
	"kobotracker/app/utils"

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

	DarkMode bool
	GG       *gg.Context
	State    gofbink.FBInkState
}

func NewScreen() *Screen {
	instance := &Screen{}
	instance.DarkMode = false
	instance.State = gofbink.FBInkState{}

	fbinkOpts := gofbink.FBInkConfig{}
	rOpts := gofbink.RestrictedConfig{}
	instance.fb = gofbink.New(&fbinkOpts, &rOpts)

	instance.fb.Open()
	instance.fb.Init(&fbinkOpts)
	instance.fb.GetState(&fbinkOpts, &instance.State)

	instance.rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(instance.State.ScreenWidth), int(instance.State.ScreenHeight)}})
	instance.GG = gg.NewContextForRGBA(instance.rgba)
	instance.GG.SetRGB(0, 0, 0)
	instance.SetFontSize(80)

	instance.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})

	return instance
}

func (screen *Screen) Clear() {
	screen.GG.SetRGB(1, 1, 1)
	screen.GG.Clear()
	screen.GG.SetRGB(0, 0, 0)
}

func (screen *Screen) Close() {
	screen.fb.Close()
}

func (screen *Screen) DrawChart(pricePoints []float64, min float64, max float64, days int, currency string, x float64, y float64, width float64, height float64) (int, int) {
	var minIndex, maxIndex int

	maxMoneyStr := utils.GetMoneyString(currency, max)
	minMoneyStr := utils.GetMoneyString(currency, min)

	screen.SetFontSize(40)

	screen.GG.DrawStringWrapped(fmt.Sprintf("%d-day chart", days), x, y-45, 0, 0, width, 0, gg.AlignRight)
	screen.GG.DrawStringWrapped(maxMoneyStr, x, y-45, 0, 0, width, 0, gg.AlignLeft)
	screen.GG.DrawStringWrapped(minMoneyStr, x, y+height+15, 0, 0, width, 0, gg.AlignLeft)

	for index, price := range pricePoints {
		stepWidth := width / float64(len(pricePoints))
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

func (screen *Screen) DrawProgressBar(x float64, y float64, width float64, height float64, filled float64) {
	screen.GG.DrawRectangle(x, y, width, height)
	screen.GG.Stroke()

	filled_width := (width / 100) * filled

	screen.GG.DrawRectangle(x, y, filled_width, height)
	screen.GG.Fill()
}

func (screen *Screen) RenderFrame() {
	screen.fb.PrintRBGA(0, 0, screen.rgba, &gofbink.FBInkConfig{IsNightmode: screen.DarkMode})
}

func (screen *Screen) SetFontSize(size float64) {
	screen.GG.LoadFontFace(utils.GetAbsolutePath("assets/font.ttf"), size)
}

func (screen *Screen) TakeScreenshot() {
	screen.GG.SavePNG(utils.GetAbsolutePath("screenshot.png"))
}
