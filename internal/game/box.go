package game

import (
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

		shouldDraw = true
	}

	if box.DesiredY() != box.PositionY {
		if math.Signbit(box.DesiredY() - box.PositionY) {
			box.PositionY -= movementSpeed
		} else {
			box.PositionY += movementSpeed
		}

		shouldDraw = true
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

	screen.DrawImage(sprites[currentSprite].Images[0], &opts)
}

func (box *Box) Done() bool {
	return stages[stageIndex].IsFlag(box.I, box.J) && box.DesiredX() == box.PositionX && box.DesiredY() == box.PositionY
}

func (box *Box) IsWallAtLeft() bool {
	return stages[stageIndex].IsWall(box.I-1, box.J)
}

func (box *Box) IsWallAtRight() bool {
	return stages[stageIndex].IsWall(box.I+1, box.J)
}

func (box *Box) IsWallAtTop() bool {
	return stages[stageIndex].IsWall(box.I, box.J-1)
}

func (box *Box) IsWallAtBottom() bool {
	return stages[stageIndex].IsWall(box.I, box.J+1)
}

func (box *Box) IsBoxAtLeft() bool {
	for j := range boxes {
		if boxes[j].I == box.I-1 && boxes[j].J == box.J {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtRight() bool {
	for j := range boxes {
		if boxes[j].I == box.I+1 && boxes[j].J == box.J {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtTop() bool {
	for j := range boxes {
		if boxes[j].I == box.I && boxes[j].J == box.J-1 {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtBottom() bool {
	for j := range boxes {
		if boxes[j].I == box.I && boxes[j].J == box.J+1 {
			return true
		}
	}

	return false
}
