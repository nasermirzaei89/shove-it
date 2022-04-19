package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	PositionX, PositionY float64
	I, J                 int
}

func (box Box) DesiredX() float64 {
	return float64(box.I * tileWidth)
}

func (box Box) DesiredY() float64 {
	return float64(box.J * tileWidth)
}

func (box *Box) Update() {
	if box.DesiredX() != box.PositionX {
		if math.Signbit(box.DesiredX() - box.PositionX) {
			box.PositionX -= movementSpeed
		} else {
			box.PositionX += movementSpeed
		}
	}

	if box.DesiredY() != box.PositionY {
		if math.Signbit(box.DesiredY() - box.PositionY) {
			box.PositionY -= movementSpeed
		} else {
			box.PositionY += movementSpeed
		}
	}
}

func (box *Box) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	opts.GeoM.Scale(Scale(), Scale())

	opts.GeoM.Translate(box.PositionX*Scale(), box.PositionY*Scale())

	currentSprite := SpriteBox
	if box.Done() {
		currentSprite = SpriteBoxDone
	}

	currentFrame := sprites[currentSprite].Frames[0]

	img := ebiten.NewImageFromImage(spriteSheet.SubImage(currentFrame.Rect))

	screen.DrawImage(img, &opts)
}

func (box *Box) Done() bool {
	return (CurrentRoomData()[box.J][box.I] == ItemTileFlagged || CurrentRoomData()[box.J][box.I] == ItemBoxDone) && box.DesiredX() == box.PositionX && box.DesiredY() == box.PositionY
}
