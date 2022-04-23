package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	characterWidth = 8
	characterSkip  = 32
)

var fontCache = make(map[int32]*ebiten.Image)

// DrawText renders text on screen.
// x and y are base on 40x28 dimension indexing.
func DrawText(screen *ebiten.Image, x, y int, text string) {
	for i, c := range text {
		if _, ok := fontCache[c]; !ok {
			cx := (int(c) - characterSkip) * characterWidth

			fontCache[c] = ebiten.NewImageFromImage(fontImage.SubImage(image.Rect(cx, 0, cx+characterWidth, characterWidth)))
		}

		opts := new(ebiten.DrawImageOptions)

		opts.GeoM.Scale(scaleFactor, scaleFactor)
		opts.GeoM.Translate(float64(x+i)*characterWidth*scaleFactor, float64(y)*characterWidth*scaleFactor)

		screen.DrawImage(fontCache[c], opts)
	}
}
