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

func (box *Box) Update(game *Game) {
	if box.DesiredX() != box.PositionX {
		if math.Signbit(box.DesiredX() - box.PositionX) {
			box.PositionX -= movementSpeed
		} else {
			box.PositionX += movementSpeed
		}

		game.shouldDraw = true
	}

	if box.DesiredY() != box.PositionY {
		if math.Signbit(box.DesiredY() - box.PositionY) {
			box.PositionY -= movementSpeed
		} else {
			box.PositionY += movementSpeed
		}

		game.shouldDraw = true
	}
}

func (box *Box) Draw(game *Game, screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	scale := game.scale()

	opts.GeoM.Scale(scale, scale)

	opts.GeoM.Translate(box.PositionX*scale, box.PositionY*scale)

	currentSprite := box.SpriteName

	if box.Done(game) {
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

	screen.DrawImage(game.sprites[currentSprite].Images[0], &opts)
}

func (box *Box) Done(game *Game) bool {
	return game.stages[game.stageIndex].IsFlag(box.I, box.J) && box.DesiredX() == box.PositionX && box.DesiredY() == box.PositionY
}

func (box *Box) IsWallAtLeft(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(box.I-1, box.J)
}

func (box *Box) IsWallAtRight(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(box.I+1, box.J)
}

func (box *Box) IsWallAtTop(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(box.I, box.J-1)
}

func (box *Box) IsWallAtBottom(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(box.I, box.J+1)
}

func (box *Box) IsBoxAtLeft(game *Game) bool {
	for j := range game.boxes {
		if game.boxes[j].I == box.I-1 && game.boxes[j].J == box.J {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtRight(game *Game) bool {
	for j := range game.boxes {
		if game.boxes[j].I == box.I+1 && game.boxes[j].J == box.J {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtTop(game *Game) bool {
	for j := range game.boxes {
		if game.boxes[j].I == box.I && game.boxes[j].J == box.J-1 {
			return true
		}
	}

	return false
}

func (box *Box) IsBoxAtBottom(game *Game) bool {
	for j := range game.boxes {
		if game.boxes[j].I == box.I && game.boxes[j].J == box.J+1 {
			return true
		}
	}

	return false
}
