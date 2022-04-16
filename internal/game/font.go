package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const characterWidth = 8

var textMap = map[int32]*ebiten.Image{}

func DrawText(screen *ebiten.Image, x, y int, text string, alignRight bool) {
	if alignRight {
		x -= len(text) * characterWidth
	}

	for i, c := range text {
		img := getCharImage(c)

		opts := new(ebiten.DrawImageOptions)

		opts.GeoM.Translate(float64(x+i*characterWidth), float64(y))

		screen.DrawImage(img, opts)
	}
}

func getCharImage(c int32) *ebiten.Image {
	res, ok := textMap[c]
	if !ok {
		i := int(c) % 16
		j := int(c) / 16
		res = ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rect(i*characterWidth, j*characterWidth, i*characterWidth+characterWidth, j*characterWidth+characterWidth)))
		textMap[c] = res
	}

	return res
}
