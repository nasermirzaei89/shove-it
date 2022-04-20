package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	PositionX, PositionY float64
	I, J                 int
	SpriteName           SpriteName
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

	currentSprite := box.SpriteName

	if box.Done() {
		switch {
		case box.SpriteName == SpriteBox1:
			currentSprite = SpriteBoxDone1
		case box.SpriteName == SpriteBox2:
			currentSprite = SpriteBoxDone2
		case box.SpriteName == SpriteBox3:
			currentSprite = SpriteBoxDone3
		case box.SpriteName == SpriteBox4:
			currentSprite = SpriteBoxDone4
		case box.SpriteName == SpriteBox5:
			currentSprite = SpriteBoxDone5
		}
	}

	currentFrame := sprites[currentSprite].Frames[0]

	img := ebiten.NewImageFromImage(tileSetImage.SubImage(image.Rect(currentFrame.I*tileWidth, currentFrame.J*tileWidth, currentFrame.I*tileWidth+tileWidth, currentFrame.J*tileWidth+tileWidth)))

	screen.DrawImage(img, &opts)
}

func (box *Box) Done() bool {
	return stages[stageIndex].IsFlag(box.I, box.J) && box.DesiredX() == box.PositionX && box.DesiredY() == box.PositionY
}

func (box *Box) IsWallOnLeft() bool {
	return stages[stageIndex].IsWall(box.I-1, box.J)
}

func (box *Box) IsWallOnRight() bool {
	return stages[stageIndex].IsWall(box.I+1, box.J)
}

func (box *Box) IsWallOnTop() bool {
	return stages[stageIndex].IsWall(box.I, box.J-1)
}

func (box *Box) IsWallOnBottom() bool {
	return stages[stageIndex].IsWall(box.I, box.J+1)
}
