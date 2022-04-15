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
	walking              bool
	animation            float64
	currentSprite        SpriteName
}

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
	p.currentSprite = sprite
	p.animation = 0
}

func (p *Player) checkLeft() {
	if p.walking || !ebiten.IsKeyPressed(ebiten.KeyLeft) || mapData[p.J][p.I-1] == ItemWall {
		return
	}

	for i := range boxes {
		if boxes[i].I == p.I-1 && boxes[i].J == p.J {
			if mapData[p.J][p.I-2] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I-2 && boxes[j].J == p.J {
					return
				}
			}

			boxes[i].I--
		}
	}

	p.I--
	p.direction = 180
	p.walking = true
	p.SetCurrentSprite("walking")
}

func (p *Player) checkRight() {
	if p.walking || !ebiten.IsKeyPressed(ebiten.KeyRight) || mapData[p.J][p.I+1] == ItemWall {
		return
	}

	for i := range boxes {
		if boxes[i].I == p.I+1 && boxes[i].J == p.J {
			if mapData[p.J][p.I+2] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I+2 && boxes[j].J == p.J {
					return
				}
			}

			boxes[i].I++
		}
	}

	p.I++
	p.direction = 0
	p.walking = true
	p.SetCurrentSprite("walking")
}

func (p *Player) checkUp() {
	if p.walking || !ebiten.IsKeyPressed(ebiten.KeyUp) || mapData[p.J-1][p.I] == ItemWall {
		return
	}

	for i := range boxes {
		if boxes[i].I == p.I && boxes[i].J == p.J-1 {
			if mapData[p.J-2][p.I] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I && boxes[j].J == p.J-2 {
					return
				}
			}

			boxes[i].J--
		}
	}

	p.J--
	p.direction = 270
	p.walking = true
	p.SetCurrentSprite("walking")
}

func (p *Player) checkDown() {
	if p.walking || !ebiten.IsKeyPressed(ebiten.KeyDown) || mapData[p.J+1][p.I] == ItemWall {
		return
	}

	for i := range boxes {
		if boxes[i].I == p.I && boxes[i].J == p.J+1 {
			if mapData[p.J+2][p.I] == ItemWall {
				return
			}

			for j := range boxes {
				if boxes[j].I == p.I && boxes[j].J == p.J+2 {
					return
				}
			}

			boxes[i].J++
		}
	}

	p.J++
	p.direction = 90
	p.walking = true
	p.SetCurrentSprite("walking")
}

func (p *Player) Update() {
	p.checkLeft()
	p.checkRight()
	p.checkUp()
	p.checkDown()

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

	if p.walking && p.PositionX == p.DesiredX() && p.PositionY == p.DesiredY() {
		p.walking = false
		p.SetCurrentSprite("idle")
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

	opts.GeoM.Translate(p.PositionX, p.PositionY)

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
