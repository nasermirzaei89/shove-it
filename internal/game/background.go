package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type object struct {
	Sprite               SpriteName
	PositionX, PositionY float64
}

func (obj *object) Draw(game *Game, screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	scale := game.scale()

	opts.GeoM.Scale(scale, scale)

	opts.GeoM.Translate(obj.PositionX*scale, obj.PositionY*scale)

	screen.DrawImage(game.sprites[obj.Sprite].Images[0], &opts)
}
