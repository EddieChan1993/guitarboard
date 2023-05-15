package img

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

//go:embed bass.png
var GuitarBoardImg []byte

var EbitenGuitarBoardImg *ebiten.Image

func InitImg() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(GuitarBoardImg))
	if err != nil {
		log.Fatal(err)
	}
	EbitenGuitarBoardImg = ebiten.NewImageFromImage(img)
}
