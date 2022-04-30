package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	PositionX, PositionY float64
	I, J                 int
	direction            float64
	animation            float64
	currentSprite        SpriteName
	idle                 bool
	pushing              bool
	history              []int
	boxHistory           []*Box
}

const (
	moveLeft = iota
	moveRight
	moveUp
	moveDown
)

const fps = 60.0

func (p Player) DesiredX() float64 {
	return float64(p.I * tileWidth)
}

func (p Player) DesiredY() float64 {
	return float64(p.J * tileWidth)
}

func (p *Player) SetCurrentSprite(sprite SpriteName) {
	if p.currentSprite != sprite {
		p.currentSprite = sprite
		p.animation = 0
	}
}

func (p *Player) IsWallAtLeft(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(p.I-1, p.J)
}

func (p *Player) IsWallAtRight(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(p.I+1, p.J)
}

func (p *Player) IsWallAtTop(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(p.I, p.J-1)
}

func (p *Player) IsWallAtBottom(game *Game) bool {
	return game.stages[game.stageIndex].IsWall(p.I, p.J+1)
}

func (p *Player) BoxAtLeft(game *Game) *Box {
	for i := range game.boxes {
		if game.boxes[i].I == p.I-1 && game.boxes[i].J == p.J {
			return game.boxes[i]
		}
	}

	return nil
}

func (p *Player) BoxAtRight(game *Game) *Box {
	for i := range game.boxes {
		if game.boxes[i].I == p.I+1 && game.boxes[i].J == p.J {
			return game.boxes[i]
		}
	}

	return nil
}

func (p *Player) BoxAtTop(game *Game) *Box {
	for i := range game.boxes {
		if game.boxes[i].I == p.I && game.boxes[i].J == p.J-1 {
			return game.boxes[i]
		}
	}

	return nil
}

func (p *Player) BoxAtBottom(game *Game) *Box {
	for i := range game.boxes {
		if game.boxes[i].I == p.I && game.boxes[i].J == p.J+1 {
			return game.boxes[i]
		}
	}

	return nil
}

func (p *Player) checkLeft(game *Game) {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyLeft) || p.IsWallAtLeft(game) {
		return
	}

	pushing := false

	if box := p.BoxAtLeft(game); box != nil {
		if box.IsWallAtLeft(game) {
			return
		}

		if box.IsBoxAtLeft(game) {
			return
		}

		pushing = true

		box.I--

		p.boxHistory = append(p.boxHistory, box)
	}

	p.I--
	p.direction = directionLeft
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveLeft)

	if !pushing {
		p.boxHistory = append(p.boxHistory, nil)
	}
}

func (p *Player) checkRight(game *Game) {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyRight) || p.IsWallAtRight(game) {
		return
	}

	pushing := false

	if box := p.BoxAtRight(game); box != nil {
		if box.IsWallAtRight(game) {
			return
		}

		if box.IsBoxAtRight(game) {
			return
		}

		pushing = true

		box.I++

		p.boxHistory = append(p.boxHistory, box)
	}

	p.I++
	p.direction = directionRight
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveRight)

	if !pushing {
		p.boxHistory = append(p.boxHistory, nil)
	}
}

func (p *Player) checkUp(game *Game) {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyUp) || p.IsWallAtTop(game) {
		return
	}

	pushing := false

	if box := p.BoxAtTop(game); box != nil {
		if box.IsWallAtTop(game) {
			return
		}

		if box.IsBoxAtTop(game) {
			return
		}

		pushing = true

		box.J--

		p.boxHistory = append(p.boxHistory, box)
	}

	p.J--
	p.direction = directionUp
	p.idle = false
	p.pushing = pushing
	p.history = append(p.history, moveUp)

	if !pushing {
		p.boxHistory = append(p.boxHistory, nil)
	}
}

func (p *Player) checkDown(game *Game) {
	if !p.idle || !ebiten.IsKeyPressed(ebiten.KeyDown) || p.IsWallAtBottom(game) {
		return
	}

	pushing := false

	if box := p.BoxAtBottom(game); box != nil {
		if box.IsWallAtBottom(game) {
			return
		}

		if box.IsBoxAtBottom(game) {
			return
		}

		pushing = true

		box.J++

		p.boxHistory = append(p.boxHistory, box)
	}

	p.J++
	p.direction = directionDown
	p.idle = false
	p.pushing = pushing

	p.history = append(p.history, moveDown)

	if !pushing {
		p.boxHistory = append(p.boxHistory, nil)
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
		p.direction = directionLeft

		if p.boxHistory[len(p.history)-1] != nil {
			pushing = true
			p.boxHistory[len(p.history)-1].I++
		}
	case moveRight:
		p.I--
		p.direction = directionRight

		if p.boxHistory[len(p.history)-1] != nil {
			pushing = true
			p.boxHistory[len(p.history)-1].I--
		}
	case moveUp:
		p.J++
		p.direction = directionUp

		if p.boxHistory[len(p.history)-1] != nil {
			pushing = true
			p.boxHistory[len(p.history)-1].J++
		}
	case moveDown:
		p.J--
		p.direction = directionDown

		if p.boxHistory[len(p.history)-1] != nil {
			pushing = true
			p.boxHistory[len(p.history)-1].J--
		}
	}

	p.idle = false
	p.pushing = pushing

	p.history = p.history[:len(p.history)-1]
	p.boxHistory = p.boxHistory[:len(p.boxHistory)-1]
}

func (p *Player) Update(game *Game) {
	p.checkLeft(game)
	p.checkRight(game)
	p.checkUp(game)
	p.checkDown(game)
	p.checkUndo()

	if p.DesiredX() != p.PositionX {
		if math.Signbit(p.DesiredX() - p.PositionX) {
			p.PositionX -= movementSpeed
		} else {
			p.PositionX += movementSpeed
		}

		game.shouldDraw = true
	}

	if p.DesiredY() != p.PositionY {
		if math.Signbit(p.DesiredY() - p.PositionY) {
			p.PositionY -= movementSpeed
		} else {
			p.PositionY += movementSpeed
		}

		game.shouldDraw = true
	}

	if !p.idle && p.PositionX == p.DesiredX() && p.PositionY == p.DesiredY() {
		p.idle = true

		game.shouldDraw = true
	}
}

func (p *Player) Draw(game *Game, screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{
		GeoM:          ebiten.GeoM{},
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	}

	opts.GeoM.Translate(-tileWidth/2, -tileWidth/2)
	opts.GeoM.Rotate(p.direction)
	opts.GeoM.Translate(tileWidth/2, tileWidth/2)

	scale := game.scale()

	opts.GeoM.Scale(scale, scale)

	opts.GeoM.Translate(p.PositionX*scale, p.PositionY*scale)

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

	sprite := game.sprites[p.currentSprite]

	p.animation += 1.0 / fps

	currentFrameIndex := int(p.animation/(1/sprite.Speed)) % len(sprite.Images)

	img := sprite.Images[currentFrameIndex]

	screen.DrawImage(img, &opts)
}
