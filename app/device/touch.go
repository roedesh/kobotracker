package device

import "github.com/shermp/go-kobo-input/koboin"

func GetTouchDevice(width int, height int) *koboin.TouchDevice {
	touchPath := "/dev/input/event1"
	touchInput := koboin.New(touchPath, width, height)
	if touchInput == nil {
		panic("Could not get touch input")
	}

	return touchInput
}
