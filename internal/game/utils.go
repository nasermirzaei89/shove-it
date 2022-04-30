package game

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/xml"
	"math"
	"path/filepath"
	"sort"
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
	defaultSizeX  = 14
	defaultSizeY  = 10
)

const (
	directionRight = 0
	directionDown  = .5 * math.Pi
	directionLeft  = math.Pi
	directionUp    = 1.5 * math.Pi
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

func (game *Game) scale() float64 {
	scaleX := float64(scaleFactor*defaultSizeX) / float64(game.stages[game.stageIndex].TMX.Width)
	scaleY := float64(scaleFactor*defaultSizeY) / float64(game.stages[game.stageIndex].TMX.Height)

	return math.Min(scaleX, scaleY)
}

func (game *Game) createObjectAt(spriteName SpriteName, i, j int) {
	game.objects = append(game.objects, &object{
		Sprite:    spriteName,
		PositionX: float64(i * tileWidth),
		PositionY: float64(j * tileWidth),
	})
}

func (game *Game) createPlayerAt(i, j int) {
	game.player = &Player{
		PositionX:     float64(i * tileWidth),
		PositionY:     float64(j * tileWidth),
		I:             i,
		J:             j,
		direction:     directionUp,
		animation:     0,
		currentSprite: SpriteIdle,
		idle:          true,
		pushing:       false,
		history:       make([]int, 0),
		boxHistory:    make([]*Box, 0),
	}
}

func (game *Game) createBoxAt(spriteName SpriteName, i, j int) {
	game.boxes = append(game.boxes, &Box{
		PositionX:  float64(i * tileWidth),
		PositionY:  float64(j * tileWidth),
		I:          i,
		J:          j,
		SpriteName: spriteName,
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

func (game *Game) loadStages(assets embed.FS, dir string) error {
	entries, err := assets.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "error on read directory")
	}

	res := make([]Stage, 0)

	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if ext != ".tmx" {
			continue
		}

		file, err := assets.Open(dir + "/" + entry.Name())
		if err != nil {
			return errors.Wrap(err, "error on open file")
		}

		defer func() { _ = file.Close() }()

		var tmx TMX

		err = xml.NewDecoder(file).Decode(&tmx)
		if err != nil {
			return errors.Wrap(err, "error on decode tmx file")
		}

		csvReader := csv.NewReader(bytes.NewBufferString(strings.ReplaceAll(tmx.Data, ",\n", "\n")))

		records, err := csvReader.ReadAll()
		if err != nil {
			return errors.Wrap(err, "error on read all from csv reader")
		}

		data := make([][]int, len(records))
		for j := range records {
			data[j] = make([]int, len(records[j]))

			for i := range records[j] {
				v, err := strconv.Atoi(records[j][i])
				if err != nil {
					return errors.Wrap(err, "error on convert string to int")
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

	sort.SliceStable(res, func(i, j int) bool {
		a, _ := strconv.Atoi(res[i].Name)
		b, _ := strconv.Atoi(res[j].Name)

		return a < b
	})

	game.stages = res

	return nil
}

func (game *Game) nextStage() {
	game.stageIndex++

	if game.stageIndex == len(game.stages) {
		game.stageIndex = 0
	}

	game.startStage()
}

func (game *Game) prevStage() {
	game.stageIndex--

	if game.stageIndex < 0 {
		game.stageIndex = len(game.stages) - 1
	}

	game.startStage()
}

func (game *Game) getThemes() (tileTheme, flagTheme SpriteName) {
	data := game.stages[game.stageIndex].Data

	for j := range data {
		for i := range data[j] {
			switch data[j][i] {
			case ItemTile1:
				tileTheme = SpriteTile1
				flagTheme = SpriteFlag1

				return
			case ItemTile2:
				tileTheme = SpriteTile2
				flagTheme = SpriteFlag2

				return
			case ItemTile3:
				tileTheme = SpriteTile3
				flagTheme = SpriteFlag3

				return
			}
		}
	}

	return
}

func (game *Game) startStage() {
	data := game.stages[game.stageIndex].Data

	game.objects = make([]*object, 0)
	game.boxes = make([]*Box, 0)

	tileTheme, flagTheme := game.getThemes()

	for j := range data {
		for i := range data[j] {
			switch data[j][i] {
			case ItemBackground1:
				game.createObjectAt(SpriteBackground1, i, j)
			case ItemBackground2:
				game.createObjectAt(SpriteBackground2, i, j)
			case ItemBackground3:
				game.createObjectAt(SpriteBackground3, i, j)
			case ItemBackground4:
				game.createObjectAt(SpriteBackground4, i, j)
			case ItemBackground5:
				game.createObjectAt(SpriteBackground5, i, j)
			case ItemWall1:
				game.createObjectAt(SpriteWall1, i, j)
			case ItemWall2:
				game.createObjectAt(SpriteWall2, i, j)
			case ItemWall3:
				game.createObjectAt(SpriteWall3, i, j)
			case ItemWall4:
				game.createObjectAt(SpriteWall4, i, j)
			case ItemTile1:
				game.createObjectAt(SpriteTile1, i, j)
			case ItemTile2:
				game.createObjectAt(SpriteTile2, i, j)
			case ItemTile3:
				game.createObjectAt(SpriteTile3, i, j)
			case ItemTileFlagged1:
				game.createObjectAt(SpriteFlag1, i, j)
			case ItemTileFlagged2:
				game.createObjectAt(SpriteFlag2, i, j)
			case ItemTileFlagged3:
				game.createObjectAt(SpriteFlag3, i, j)
			case ItemPlayer1:
				game.createObjectAt(SpriteTile1, i, j)
				game.createPlayerAt(i, j)
			case ItemPlayer2:
				game.createObjectAt(SpriteTile2, i, j)
				game.createPlayerAt(i, j)
			case ItemPlayer3:
				game.createObjectAt(SpriteTile3, i, j)
				game.createPlayerAt(i, j)
			case ItemBox1:
				game.createObjectAt(tileTheme, i, j)
				game.createBoxAt(SpriteBox1, i, j)
			case ItemBox2:
				game.createObjectAt(tileTheme, i, j)
				game.createBoxAt(SpriteBox2, i, j)
			case ItemBox3:
				game.createObjectAt(tileTheme, i, j)
				game.createBoxAt(SpriteBox3, i, j)
			case ItemBox4:
				game.createObjectAt(tileTheme, i, j)
				game.createBoxAt(SpriteBox4, i, j)
			case ItemBox5:
				game.createObjectAt(tileTheme, i, j)
				game.createBoxAt(SpriteBox5, i, j)
			case ItemBoxDone1:
				game.createObjectAt(flagTheme, i, j)
				game.createBoxAt(SpriteBox1, i, j)
			case ItemBoxDone2:
				game.createObjectAt(flagTheme, i, j)
				game.createBoxAt(SpriteBox2, i, j)
			case ItemBoxDone3:
				game.createObjectAt(flagTheme, i, j)
				game.createBoxAt(SpriteBox3, i, j)
			case ItemBoxDone4:
				game.createObjectAt(flagTheme, i, j)
				game.createBoxAt(SpriteBox4, i, j)
			case ItemBoxDone5:
				game.createObjectAt(flagTheme, i, j)
				game.createBoxAt(SpriteBox5, i, j)
			}
		}
	}

	game.shouldDraw = true
}
