package game

import (
	"image"

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

	opts.GeoM.Translate(obj.PositionX, obj.PositionY)

	currentFrame := sprites[obj.Sprite].Frames[0]

	img := ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rect(
		currentFrame.IndexX*tileWidth,
		currentFrame.IndexY*tileWidth,
		(currentFrame.IndexX+1)*tileWidth,
		(currentFrame.IndexY+1)*tileWidth,
	)))

	screen.DrawImage(img, &opts)
}
