package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Object struct {
	Sprite               SpriteName
	PositionX, PositionY float64
}

func (obj *Object) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	opts.GeoM.Scale(Scale(), Scale())

	opts.GeoM.Translate(obj.PositionX*Scale(), obj.PositionY*Scale())

	currentFrame := sprites[obj.Sprite].Frames[0]

	img := ebiten.NewImageFromImage(spriteSheet.SubImage(currentFrame.Rect))

	screen.DrawImage(img, &opts)
}
