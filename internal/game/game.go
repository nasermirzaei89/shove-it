package game

import (
	"bytes"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/pkg/errors"
)

const (
	tileWidth     = 24
	movementSpeed = 2
)

type SpriteName string

const (
	SpriteIdle        SpriteName = "idle"
	SpriteWalking     SpriteName = "walking"
	SpritePushing     SpriteName = "pushing"
	SpritePushingIdle SpriteName = "pushing-idle"
	SpriteBackground  SpriteName = "background"
	SpriteWall        SpriteName = "wall"
	SpriteTile        SpriteName = "tile"
	SpriteTileFlagged SpriteName = "tile-flagged"
	SpriteBox         SpriteName = "box"
	SpriteBoxDone     SpriteName = "box-done"
)

var (
	player             *Player
	boxes              []*Box
	objects            []*Object
	sprites            map[SpriteName]*Sprite
	spriteSheet        *ebiten.Image
	steps              = 0
	stageIndex         = 0
	roomIndex          = 0
	keyPageUpPressed   = false
	keyPageDownPressed = false
)

type Game struct{}

func (g *Game) Update() error {
	done := true

	for i := range boxes {
		boxes[i].Update()

		if !boxes[i].Done() {
			done = false
		}
	}

	player.Update()

	if done {
		g.NextRoom()
	}

	if ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		if !keyPageUpPressed {
			g.NextRoom()

			keyPageUpPressed = true
		}
	} else {
		keyPageUpPressed = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		if !keyPageDownPressed {
			g.PrevRoom()

			keyPageDownPressed = true
		}
	} else {
		keyPageDownPressed = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range objects {
		objects[i].Draw(screen)
	}

	for i := range boxes {
		boxes[i].Draw(screen)
	}

	player.Draw(screen)

	// HUD
	hudY := len(CurrentRoomData())*tileWidth - 2*characterWidth

	DrawText(screen, 2*characterWidth, hudY, "STEP", false)
	DrawText(screen, 12*characterWidth, hudY, strconv.Itoa(steps), true)

	DrawText(screen, 16*characterWidth, hudY, "STAGE", false)
	DrawText(screen, 24*characterWidth, hudY, strconv.Itoa(stageIndex+1), true)

	DrawText(screen, 28*characterWidth, hudY, "ROOM", false)
	DrawText(screen, 36*characterWidth, hudY, strconv.Itoa(roomIndex+1), true)
}

func (g *Game) Layout(int, int) (int, int) {
	return len(CurrentRoomData()[0]) * tileWidth, len(CurrentRoomData()) * tileWidth
}

func CurrentRoomData() [][]int {
	return stages[stageIndex].Rooms[roomIndex].Data
}

func (g *Game) NextRoom() {
	roomIndex++
	if roomIndex == len(stages[stageIndex].Rooms) {
		roomIndex = 0
		stageIndex++

		if stageIndex == len(stages) {
			panic(errors.New("game finished")) // TODO
		}
	}

	g.LoadRoom()
}

func (g *Game) PrevRoom() {
	roomIndex--

	if roomIndex < 0 {
		stageIndex--

		if stageIndex < 0 {
			panic(errors.New("game finished")) // TODO
		}

		roomIndex = len(stages[stageIndex].Rooms) - 1
	}

	g.LoadRoom()
}

func (g *Game) LoadRoom() {
	data := CurrentRoomData()

	objects = make([]*Object, 0)
	boxes = make([]*Box, 0)
	steps = 0

	for j := range data {
		for i := range data[j] {
			switch data[j][i] {
			case ItemBackground:
				objects = append(objects, &Object{
					Sprite:    SpriteBackground,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})
			case ItemWall:
				objects = append(objects, &Object{
					Sprite:    SpriteWall,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})
			case ItemTile:
				objects = append(objects, &Object{
					Sprite:    SpriteTile,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})
			case ItemTileFlagged:
				objects = append(objects, &Object{
					Sprite:    SpriteTileFlagged,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})
			case ItemPlayer:
				objects = append(objects, &Object{
					Sprite:    SpriteTile,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})
				player = &Player{
					PositionX:     float64(i * tileWidth),
					PositionY:     float64(j * tileWidth),
					I:             i,
					J:             j,
					direction:     270,
					animation:     0,
					currentSprite: SpriteIdle,
					idle:          true,
					pushing:       false,
				}
			case ItemBox:
				objects = append(objects, &Object{
					Sprite:    SpriteTile,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})

				boxes = append(boxes, &Box{
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
					I:         i,
					J:         j,
				})
			case ItemBoxDone:
				objects = append(objects, &Object{
					Sprite:    SpriteTileFlagged,
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
				})

				boxes = append(boxes, &Box{
					PositionX: float64(i * tileWidth),
					PositionY: float64(j * tileWidth),
					I:         i,
					J:         j,
				})
			}
		}
	}
}

func New(spriteSheetPNG []byte) (*Game, error) {
	game1 := Game{}

	var err error

	spriteSheet, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(spriteSheetPNG))
	if err != nil {
		return nil, errors.Wrap(err, "error on new image from file")
	}

	sprites = map[SpriteName]*Sprite{
		SpriteIdle: NewSprite(
			[]Frame{
				{IndexX: 1, IndexY: 3},
			},
			1,
		),
		SpriteWalking: NewSprite(
			[]Frame{
				{IndexX: 0, IndexY: 3},
				{IndexX: 1, IndexY: 3},
				{IndexX: 2, IndexY: 3},
				{IndexX: 1, IndexY: 3},
			},
			10,
		),
		SpritePushing: NewSprite(
			[]Frame{
				{IndexX: 0, IndexY: 4},
				{IndexX: 1, IndexY: 4},
				{IndexX: 2, IndexY: 4},
				{IndexX: 1, IndexY: 4},
			},
			10,
		),
		SpritePushingIdle: NewSprite(
			[]Frame{
				{IndexX: 1, IndexY: 4},
			},
			1,
		),
		SpriteBackground: NewSprite(
			[]Frame{
				{IndexX: 0, IndexY: 2},
			},
			1,
		),
		SpriteWall: NewSprite(
			[]Frame{
				{IndexX: 1, IndexY: 2},
			},
			1,
		),
		SpriteTile: NewSprite(
			[]Frame{
				{IndexX: 2, IndexY: 2},
			},
			1,
		),
		SpriteTileFlagged: NewSprite(
			[]Frame{
				{IndexX: 3, IndexY: 2},
			},
			1,
		),
		SpriteBox: NewSprite(
			[]Frame{
				{IndexX: 3, IndexY: 3},
			},
			1,
		),
		SpriteBoxDone: NewSprite(
			[]Frame{
				{IndexX: 3, IndexY: 4},
			},
			1,
		),
	}

	game1.LoadRoom()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)

	return &game1, nil
}
