package game

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/xml"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	player     *Player
	boxes      []*Box
	objects    []*Object
	sprites    map[SpriteName]*Sprite
	stages     []Stage
	stageIndex = 0
)

var (
	fontImage    *ebiten.Image
	playerImage  *ebiten.Image
	tileSetImage *ebiten.Image
)

func Scale() float64 {
	// 336 => 14
	// x(288)   => 12
	// x * y = 336
	scaleX := float64(scaleFactor*14) / float64(stages[stageIndex].TMX.Width)
	// 240 => 10
	// x(288)   => 12
	// x * y = 240
	scaleY := float64(scaleFactor*10) / float64(stages[stageIndex].TMX.Height)

	return math.Min(scaleX, scaleY)
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

func loadImage(assets embed.FS, filename string) (*ebiten.Image, error) {
	b, err := assets.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error on read file from assets")
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "error on new image from file")
	}

	return img, nil
}

func loadStages(assets embed.FS, dir string) ([]Stage, error) {
	entries, err := assets.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "error on read directory")
	}

	res := make([]Stage, 0)

	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if ext != ".tmx" {
			continue
		}

		f, err := assets.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "error on open file")
		}

		defer func() { _ = f.Close() }()

		var tmx TMX

		err = xml.NewDecoder(f).Decode(&tmx)
		if err != nil {
			return nil, errors.Wrap(err, "error on decode tmx file")
		}

		csvReader := csv.NewReader(bytes.NewBufferString(strings.ReplaceAll(tmx.Data, ",\n", "\n")))

		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, errors.Wrap(err, "error on read all from csv reader")
		}

		data := make([][]int, len(records))
		for j := range records {
			data[j] = make([]int, len(records[j]))

			for i := range records[j] {
				v, err := strconv.Atoi(records[j][i])
				if err != nil {
					return nil, errors.Wrap(err, "error on convert string to int")
				}

				data[j][i] = v
			}
		}

		res = append(res, Stage{
			Name: strings.ReplaceAll(entry.Name(), ext, ""),
			Data: data,
			TMX:  tmx,
		})
	}

	return res, nil
}

func nextStage() {
	stageIndex++

	if stageIndex == len(stages) {
		stageIndex = 0
	}

	loadStage()
}

func prevStage() {
	stageIndex--

	if stageIndex < 0 {
		stageIndex = len(stages) - 1
	}

	loadStage()
}

func loadStage() {
	data := stages[stageIndex].Data

	objects = make([]*Object, 0)
	boxes = make([]*Box, 0)

	for j := range data {
		for i := range data[j] {
			switch data[j][i] {
			case ItemBackground1:
				createBackgroundAt(i, j)
			case ItemWall1:
				createWallAt(i, j)
			case ItemTile1:
				createTileAt(i, j)
			case ItemTileFlagged1:
				createFlaggedTileAt(i, j)
			case ItemPlayer1, ItemPlayer2, ItemPlayer3:
				createPlayerAt(i, j)
			case ItemBox1:
				createBoxAt(i, j)
			case ItemBoxDone1:
				createBoxDoneAt(i, j)
			}
		}
	}
}
