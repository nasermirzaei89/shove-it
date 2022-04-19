package game

import (
	"bytes"
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

	SpriteBackground1 SpriteName = "background1"
	SpriteBackground2 SpriteName = "background2"
	SpriteBackground3 SpriteName = "background3"
	SpriteBackground4 SpriteName = "background4"
	SpriteBackground5 SpriteName = "background5"
	SpriteBackground6 SpriteName = "background6"

	SpriteWall1 SpriteName = "wall1"
	SpriteWall2 SpriteName = "wall2"
	SpriteWall3 SpriteName = "wall3"
	SpriteWall4 SpriteName = "wall4"

	SpriteTile1 SpriteName = "tile1"
	SpriteTile2 SpriteName = "tile2"
	SpriteTile3 SpriteName = "tile3"

	SpriteFlag1 SpriteName = "flag1"
	SpriteFlag2 SpriteName = "flag2"
	SpriteFlag3 SpriteName = "flag3"

	SpriteBox1 SpriteName = "box1"
	SpriteBox2 SpriteName = "box2"
	SpriteBox3 SpriteName = "box3"
	SpriteBox4 SpriteName = "box4"
	SpriteBox5 SpriteName = "box5"

	SpriteBoxDone1 SpriteName = "box-done1"
	SpriteBoxDone2 SpriteName = "box-done2"
	SpriteBoxDone3 SpriteName = "box-done3"
	SpriteBoxDone4 SpriteName = "box-done4"
	SpriteBoxDone5 SpriteName = "box-done5"
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
	objects = append(objects, &Object{
		Sprite:    SpriteBackground1,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
	})
}

func createWallAt(i, j int) {
	objects = append(objects, &Object{
		Sprite:    SpriteWall1,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
	})
}

func createTileAt(i, j int) {
	objects = append(objects, &Object{
		Sprite:    SpriteTile1,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
	})
}

func createFlaggedTileAt(i, j int) {
	objects = append(objects, &Object{
		Sprite:    SpriteFlag1,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
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
				{I: 1, J: 0},
			},
			1,
		),
		SpriteWalking: NewSprite(
			[]Frame{
				{I: 0, J: 0},
				{I: 1, J: 0},
				{I: 2, J: 0},
				{I: 1, J: 0},
			},
			10,
		),
		SpritePushing: NewSprite(
			[]Frame{
				{I: 3, J: 0},
				{I: 4, J: 0},
				{I: 5, J: 0},
				{I: 4, J: 0},
			},
			10,
		),
		SpritePushingIdle: NewSprite(
			[]Frame{
				{I: 4, J: 0},
			},
			1,
		),
		SpriteBackground1: NewSprite(
			[]Frame{
				{I: 0, J: 4},
			},
			1,
		),
		SpriteBackground2: NewSprite(
			[]Frame{
				{I: 1, J: 4},
			},
			1,
		),
		SpriteBackground3: NewSprite(
			[]Frame{
				{I: 2, J: 4},
			},
			1,
		),
		SpriteBackground4: NewSprite(
			[]Frame{
				{I: 3, J: 4},
			},
			1,
		),
		SpriteBackground5: NewSprite(
			[]Frame{
				{I: 4, J: 4},
			},
			1,
		),
		SpriteBackground6: NewSprite(
			[]Frame{
				{I: 5, J: 4},
			},
			1,
		),
		SpriteWall1: NewSprite(
			[]Frame{
				{I: 0, J: 1},
			},
			1,
		),
		SpriteWall2: NewSprite(
			[]Frame{
				{I: 1, J: 1},
			},
			1,
		),
		SpriteWall3: NewSprite(
			[]Frame{
				{I: 2, J: 1},
			},
			1,
		),
		SpriteWall4: NewSprite(
			[]Frame{
				{I: 3, J: 1},
			},
			1,
		),
		SpriteTile1: NewSprite(
			[]Frame{
				{I: 0, J: 5},
			},
			1,
		),
		SpriteTile2: NewSprite(
			[]Frame{
				{I: 1, J: 5},
			},
			1,
		),
		SpriteTile3: NewSprite(
			[]Frame{
				{I: 2, J: 5},
			},
			1,
		),
		SpriteFlag1: NewSprite(
			[]Frame{
				{I: 0, J: 6},
			},
			1,
		),
		SpriteFlag2: NewSprite(
			[]Frame{
				{I: 1, J: 6},
			},
			1,
		),
		SpriteFlag3: NewSprite(
			[]Frame{
				{I: 2, J: 6},
			},
			1,
		),
		SpriteBox1: NewSprite(
			[]Frame{
				{I: 0, J: 2},
			},
			1,
		),
		SpriteBox2: NewSprite(
			[]Frame{
				{I: 1, J: 2},
			},
			1,
		),
		SpriteBox3: NewSprite(
			[]Frame{
				{I: 2, J: 2},
			},
			1,
		),
		SpriteBox4: NewSprite(
			[]Frame{
				{I: 3, J: 2},
			},
			1,
		),
		SpriteBox5: NewSprite(
			[]Frame{
				{I: 4, J: 2},
			},
			1,
		),
		SpriteBoxDone1: NewSprite(
			[]Frame{
				{I: 0, J: 3},
			},
			1,
		),
		SpriteBoxDone2: NewSprite(
			[]Frame{
				{I: 1, J: 3},
			},
			1,
		),
		SpriteBoxDone3: NewSprite(
			[]Frame{
				{I: 2, J: 3},
			},
			1,
		),
		SpriteBoxDone4: NewSprite(
			[]Frame{
				{I: 3, J: 3},
			},
			1,
		),
		SpriteBoxDone5: NewSprite(
			[]Frame{
				{I: 4, J: 3},
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
