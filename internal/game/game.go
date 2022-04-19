package game

import (
	"bytes"
	"image"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pkg/errors"
)

const (
	tileWidth     = 24
	movementSpeed = 2
	screenWidth   = 320
	screenHeight  = 224
	scaleFactor   = 3
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
	SpriteFlag        SpriteName = "flag"
	SpriteBox         SpriteName = "box"
	SpriteBoxDone     SpriteName = "box-done"
)

var (
	player      *Player
	boxes       []*Box
	objects     []*Object
	sprites     map[SpriteName]*Sprite
	spriteSheet *ebiten.Image
	fontImage   *ebiten.Image
	stageIndex  = 0
	roomIndex   = 0
)

func Scale() float64 {
	// 336 => 14
	// x(288)   => 12
	// x * y = 336
	scaleX := float64(scaleFactor*14) / float64(CurrentRoomData().Width())
	// 240 => 10
	// x(288)   => 12
	// x * y = 240
	scaleY := float64(scaleFactor*10) / float64(CurrentRoomData().Height())

	return math.Min(scaleX, scaleY)
}

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

	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		g.NextRoom()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		g.PrevRoom()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.LoadRoom()
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
	DrawText(screen, 2, 26, "STEP", false)
	DrawText(screen, 12, 26, strconv.Itoa(len(player.history)), true)

	DrawText(screen, 16, 26, "STAGE", false)
	DrawText(screen, 24, 26, strconv.Itoa(stageIndex+1), true)

	DrawText(screen, 28, 26, "ROOM", false)
	DrawText(screen, 36, 26, strconv.Itoa(roomIndex+1), true)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth * scaleFactor, screenHeight * scaleFactor
}

func CurrentRoomData() Data {
	return stages[stageIndex].Rooms[roomIndex].Data
}

func (g *Game) NextRoom() {
	roomIndex++
	if roomIndex == len(stages[stageIndex].Rooms) {
		roomIndex = 0
		stageIndex++

		if stageIndex == len(stages) {
			stageIndex = 0
		}
	}

	g.LoadRoom()
}

func (g *Game) PrevRoom() {
	roomIndex--

	if roomIndex < 0 {
		stageIndex--

		if stageIndex < 0 {
			stageIndex = len(stages) - 1
		}

		roomIndex = len(stages[stageIndex].Rooms) - 1
	}

	g.LoadRoom()
}

func createBackgroundAt(i, j int) {
	for m := 0; m < 3; m++ {
		for n := 0; n < 3; n++ {
			objects = append(objects, &Object{
				Sprite:    SpriteBackground,
				PositionX: float64(i*tileWidth + m*8),
				PositionY: float64(j*tileWidth + n*8),
			})
		}
	}
}

func createWallAt(i, j int) {
	objects = append(objects, &Object{
		Sprite:    SpriteWall,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
	})
}

func createTileAt(i, j int) {
	for m := 0; m < 3; m++ {
		for n := 0; n < 3; n++ {
			objects = append(objects, &Object{
				Sprite:    SpriteTile,
				PositionX: float64(i*tileWidth + m*8),
				PositionY: float64(j*tileWidth + n*8),
			})
		}
	}
}

func createFlaggedTileAt(i, j int) {
	createTileAt(i, j)

	objects = append(objects, &Object{
		Sprite:    SpriteFlag,
		PositionX: float64(i*tileWidth + 8),
		PositionY: float64(j*tileWidth + 8),
	})
}

func createPlayerAt(i, j int) {
	createTileAt(i, j)

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
		history:       []int{},
		boxHistory:    []int{},
	}
}

func createBoxAt(i, j int) {
	createTileAt(i, j)

	boxes = append(boxes, &Box{
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
		I:         i,
		J:         j,
	})
}

func createBoxDoneAt(i, j int) {
	createFlaggedTileAt(i, j)

	boxes = append(boxes, &Box{
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
		I:         i,
		J:         j,
	})
}

func (g *Game) LoadRoom() {
	data := CurrentRoomData()

	objects = make([]*Object, 0)
	boxes = make([]*Box, 0)

	for j := range data {
		for i := range data[j] {
			switch data[j][i] {
			case ItemBackground:
				createBackgroundAt(i, j)
			case ItemWall:
				createWallAt(i, j)
			case ItemTile:
				createTileAt(i, j)
			case ItemTileFlagged:
				createFlaggedTileAt(i, j)
			case ItemPlayer:
				createPlayerAt(i, j)
			case ItemBox:
				createBoxAt(i, j)
			case ItemBoxDone:
				createBoxDoneAt(i, j)
			}
		}
	}
}

func New(spriteSheetPNG, fontPNG []byte) (*Game, error) {
	game1 := Game{}

	var err error

	spriteSheet, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(spriteSheetPNG))
	if err != nil {
		return nil, errors.Wrap(err, "error on new image from file")
	}

	fontImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(fontPNG))
	if err != nil {
		return nil, errors.Wrap(err, "error on new image from file")
	}

	sprites = map[SpriteName]*Sprite{
		SpriteIdle: NewSprite(
			[]Frame{
				{Rect: image.Rect(24, 88, 48, 112)},
			},
			1,
		),
		SpriteWalking: NewSprite(
			[]Frame{
				{Rect: image.Rect(0, 88, 24, 112)},
				{Rect: image.Rect(24, 88, 48, 112)},
				{Rect: image.Rect(48, 88, 72, 112)},
				{Rect: image.Rect(24, 88, 48, 112)},
			},
			10,
		),
		SpritePushing: NewSprite(
			[]Frame{
				{Rect: image.Rect(72, 88, 96, 112)},
				{Rect: image.Rect(96, 88, 120, 112)},
				{Rect: image.Rect(120, 88, 144, 112)},
				{Rect: image.Rect(96, 88, 120, 112)},
			},
			10,
		),
		SpritePushingIdle: NewSprite(
			[]Frame{
				{Rect: image.Rect(96, 88, 120, 112)},
			},
			1,
		),
		SpriteBackground: NewSprite(
			[]Frame{
				{Rect: image.Rect(0, 16, 8, 24)},
			},
			1,
		),
		SpriteWall: NewSprite(
			[]Frame{
				{Rect: image.Rect(0, 112, 24, 136)},
			},
			1,
		),
		SpriteTile: NewSprite(
			[]Frame{
				{Rect: image.Rect(8, 16, 16, 24)},
			},
			1,
		),
		SpriteFlag: NewSprite(
			[]Frame{
				{Rect: image.Rect(80, 16, 88, 24)},
			},
			1,
		),
		SpriteBox: NewSprite(
			[]Frame{
				{Rect: image.Rect(0, 136, 24, 160)},
			},
			1,
		),
		SpriteBoxDone: NewSprite(
			[]Frame{
				{Rect: image.Rect(0, 160, 24, 184)},
			},
			1,
		),
	}

	game1.LoadRoom()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetWindowSize(screenWidth*scaleFactor, screenHeight*scaleFactor)

	return &game1, nil
}
