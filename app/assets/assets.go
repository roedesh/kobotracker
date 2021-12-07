package assets

import (
	"bytes"
	_ "embed"
	"image"

	"github.com/nfnt/resize"
)

//go:embed bolt.png
var boltImageBytes []byte

//go:embed signout.png
var signoutImageBytes []byte

var (
	BoltImage    image.Image
	SignOutImage image.Image
)

func init() {
	BoltImage = getImageFromBytes(boltImageBytes, 40, 40)
	SignOutImage = getImageFromBytes(signoutImageBytes, 40, 40)
}

func getImageFromBytes(imageBytes []byte, width uint, height uint) image.Image {
	img, _, _ := image.Decode(bytes.NewReader(imageBytes))
	resizedImage := resize.Resize(width, height, img, resize.Lanczos3)
	return resizedImage
}
