package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	characterWidth = 8
	characterSkip  = 32
)

// DrawText renders text on screen.
// x and y are base on 40x28 dimension indexing.
func (game *Game) DrawText(screen *ebiten.Image, posX, posY int, text string) {
	for i, c := range text {
		if _, ok := game.fontCache[c]; !ok {
			cx := (int(c) - characterSkip) * characterWidth

			game.fontCache[c] = ebiten.NewImageFromImage(game.fontImage.SubImage(image.Rect(cx, 0, cx+characterWidth, characterWidth)))
		}

		opts := new(ebiten.DrawImageOptions)

		opts.GeoM.Scale(scaleFactor, scaleFactor)
		opts.GeoM.Translate(float64(posX+i)*characterWidth*scaleFactor, float64(posY)*characterWidth*scaleFactor)

		screen.DrawImage(game.fontCache[c], opts)
	}
}
