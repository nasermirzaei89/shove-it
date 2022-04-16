package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	PositionX, PositionY float64
	I, J                 int
	direction            int
	animation            float64
	currentSprite        SpriteName
	idle                 bool
	pushing              bool
	history              []int
	boxHistory           []int
}

const (
	moveLeft = iota
	moveRight
	moveUp
	moveDown
)

func (p Player) DesiredX() float64 {
	return float64(p.I * tileWidth)
}

func (p Player) DesiredY() float64 {
	return float64(p.J * tileWidth)
}

func (p Player) DirectionTheta() float64 {
	return float64(p.direction) * (math.Pi / 180.0)
}

func (p *Player) SetCurrentSprite(sprite SpriteName) {
	if p.currentSprite != sprite {
		p.currentSprite = sprite
		p.animation = 0
	}
}

func (p *Player) checkLeft() {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyLeft) || CurrentRoomData()[p.J][p.I-1] == ItemWall {
		return
	}

	pushing := false

	for i := range boxes {
		if boxes[i].I == p.I-1 && boxes[i].J == p.J {
			if CurrentRoomData()[p.J][p.I-2] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I-2 && boxes[j].J == p.J {
					return
				}
			}

			pushing = true

			boxes[i].I--

			p.boxHistory = append(p.boxHistory, i)
		}
	}

	p.I--
	p.direction = 180
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveLeft)

	if !pushing {
		p.boxHistory = append(p.boxHistory, -1)
	}
}

func (p *Player) checkRight() {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyRight) || CurrentRoomData()[p.J][p.I+1] == ItemWall {
		return
	}

	pushing := false

	for i := range boxes {
		if boxes[i].I == p.I+1 && boxes[i].J == p.J {
			if CurrentRoomData()[p.J][p.I+2] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I+2 && boxes[j].J == p.J {
					return
				}
			}

			pushing = true

			boxes[i].I++

			p.boxHistory = append(p.boxHistory, i)
		}
	}

	p.I++
	p.direction = 0
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveRight)

	if !pushing {
		p.boxHistory = append(p.boxHistory, -1)
	}
}

func (p *Player) checkUp() {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyUp) || CurrentRoomData()[p.J-1][p.I] == ItemWall {
		return
	}

	pushing := false

	for i := range boxes {
		if boxes[i].I == p.I && boxes[i].J == p.J-1 {
			if CurrentRoomData()[p.J-2][p.I] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I && boxes[j].J == p.J-2 {
					return
				}
			}

			pushing = true

			boxes[i].J--

			p.boxHistory = append(p.boxHistory, i)
		}
	}

	p.J--
	p.direction = 270
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveUp)

	if !pushing {
		p.boxHistory = append(p.boxHistory, -1)
	}
}

func (p *Player) checkDown() {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyDown) || CurrentRoomData()[p.J+1][p.I] == ItemWall {
		return
	}

	pushing := false

	for i := range boxes {
		if boxes[i].I == p.I && boxes[i].J == p.J+1 {
			if CurrentRoomData()[p.J+2][p.I] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I && boxes[j].J == p.J+2 {
					return
				}
			}

			pushing = true

			boxes[i].J++

			p.boxHistory = append(p.boxHistory, i)
		}
	}

	p.J++
	p.direction = 90
	p.idle = false
	p.pushing = pushing

	p.history = append(p.history, moveDown)

	if !pushing {
		p.boxHistory = append(p.boxHistory, -1)
	}
}

func (p *Player) checkUndo() {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyBackspace) || len(p.history) == 0 {
		return
	}

	pushing := false

	switch p.history[len(p.history)-1] {
	case moveLeft:
		p.I++
		p.direction = 180

		if p.boxHistory[len(p.history)-1] != -1 {
			pushing = true
			boxes[p.boxHistory[len(p.history)-1]].I++
		}
	case moveRight:
		p.I--
		p.direction = 0

		if p.boxHistory[len(p.history)-1] != -1 {
			pushing = true
			boxes[p.boxHistory[len(p.history)-1]].I--
		}
	case moveUp:
		p.J++
		p.direction = 270

		if p.boxHistory[len(p.history)-1] != -1 {
			pushing = true
			boxes[p.boxHistory[len(p.history)-1]].J++
		}
	case moveDown:
		p.J--
		p.direction = 90

		if p.boxHistory[len(p.history)-1] != -1 {
			pushing = true
			boxes[p.boxHistory[len(p.history)-1]].J--
		}
	}

	p.idle = false
	p.pushing = pushing

	p.history = p.history[:len(p.history)-1]
	p.boxHistory = p.boxHistory[:len(p.boxHistory)-1]
}

func (p *Player) Update() {
	p.checkLeft()
	p.checkRight()
	p.checkUp()
	p.checkDown()
	p.checkUndo()

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

	if !p.idle && p.PositionX == p.DesiredX() && p.PositionY == p.DesiredY() {
		p.idle = true
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	opts.GeoM.Translate(-tileWidth/2, -tileWidth/2)
	opts.GeoM.Rotate(p.DirectionTheta())
	opts.GeoM.Translate(tileWidth/2, tileWidth/2)

	opts.GeoM.Scale(Scale(), Scale())

	opts.GeoM.Translate(p.PositionX*Scale(), p.PositionY*Scale())

	switch {
	case p.idle && !p.pushing:
		p.SetCurrentSprite(SpriteIdle)
	case p.idle && p.pushing:
		p.SetCurrentSprite(SpritePushingIdle)
	case !p.idle && !p.pushing:
		p.SetCurrentSprite(SpriteWalking)
	case !p.idle && p.pushing:
		p.SetCurrentSprite(SpritePushing)
	}

	sprite := sprites[p.currentSprite]

	p.animation += 1.0 / 60.0

	currentFrameIndex := int(p.animation/(1/sprite.Speed)) % len(sprite.Frames)

	currentFrame := sprite.Frames[currentFrameIndex]

	img := ebiten.NewImageFromImage(spriteSheet.SubImage(image.Rect(
		currentFrame.IndexX*tileWidth,
		currentFrame.IndexY*tileWidth,
		(currentFrame.IndexX+1)*tileWidth,
		(currentFrame.IndexY+1)*tileWidth,
	)))

	screen.DrawImage(img, &opts)
}
