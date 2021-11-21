package screener

import (
	"github.com/shermp/go-fbink-v2/gofbink"
)

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

	screen.fb.GetState(&fbinkOpts, &screen.state)

	return
}

func (screen *Screen) ClearScreen() {
	screen.fb.ClearScreen(&gofbink.FBInkConfig{IsFlashing: true})
}

func (screen *Screen) Print(a ...interface{}) {
	screen.fb.Println(a)
}

func (screen *Screen) Close() {
	screen.fb.Close()
}
