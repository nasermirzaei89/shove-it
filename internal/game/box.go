package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	PositionX, PositionY float64
	I, J                 int
}

func (p Box) DesiredX() float64 {
	return float64(p.I * tileWidth)
}

func (p Box) DesiredY() float64 {
	return float64(p.J * tileWidth)
}

func (p *Box) Update() {
	if p.DesiredX() != p.PositionX {
		if math.Signbit(p.DesiredX() - p.PositionX) {
			p.PositionX -= movementSpeed
		} else {
			p.PositionX += movementSpeed
		}
	}

	if p.DesiredY() != p.PositionY {
		if math.Signbit(p.DesiredY() - p.PositionY) {
			p.PositionY -= movementSpeed
		} else {
			p.PositionY += movementSpeed
		}
	}
}

func (p *Box) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	opts.GeoM.Translate(p.PositionX, p.PositionY)

	currentSprite := SpriteBox
	if mapData[p.J][p.I] == ItemTileFlagged {
		currentSprite = SpriteBoxDone
	}

	currentFrame := sprites[currentSprite].Frames[0]

	img := ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rect(
		currentFrame.IndexX*tileWidth,
		currentFrame.IndexY*tileWidth,
		(currentFrame.IndexX+1)*tileWidth,
		(currentFrame.IndexY+1)*tileWidth,
	)))

	screen.DrawImage(img, &opts)
}
