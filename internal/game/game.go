package game

import (
	"fmt"
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
	player      *Player
	boxes       []*Box
	objects     []*Object
	sprites     map[SpriteName]*Sprite
	spriteSheet *ebiten.Image
	done        bool
	steps       int
	stage       = 1
	room        = 1
)

type Game struct{}

func (g *Game) Update() error {
	check := true

	for i := range boxes {
		boxes[i].Update()

		if !boxes[i].Done() {
			check = false
		}
	}

	player.Update()

	done = check

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

	DrawText(screen, 2*characterWidth, 200, "STEP", false)
	DrawText(screen, 12*characterWidth, 200, strconv.Itoa(steps), true)

	DrawText(screen, 16*characterWidth, 200, "STAGE", false)
	DrawText(screen, 24*characterWidth, 200, strconv.Itoa(stage), true)

	DrawText(screen, 28*characterWidth, 200, "ROOM", false)
	DrawText(screen, 36*characterWidth, 200, strconv.Itoa(room), true)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%t", done))
}

func (g *Game) Layout(int, int) (int, int) {
	return len(mapData[0]) * tileWidth, len(mapData) * tileWidth
}

func New() (*Game, error) {
	game1 := Game{}

	var err error

	spriteSheet, _, err = ebitenutil.NewImageFromFile("./spritesheet.png")
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

	for j := range mapData {
		for i := range mapData[j] {
			switch mapData[j][i] {
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
					direction:     0,
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

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)

	return &game1, nil
}
