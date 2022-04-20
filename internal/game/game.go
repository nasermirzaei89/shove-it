package game

import (
	"embed"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pkg/errors"
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

	if player != nil {
		player.Update()
	}

	if done {
		nextStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		nextStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		prevStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		loadStage()
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

	steps := 0

	if player != nil {
		player.Draw(screen)
		steps = len(player.history)
	}

	// HUD
	DrawText(screen, 2, 26, "STEP", false)
	DrawText(screen, 12, 26, strconv.Itoa(steps), true)

	DrawText(screen, 22, 26, "STAGE", false)
	DrawText(screen, 38, 26, stages[stageIndex].Name, true)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth * scaleFactor, screenHeight * scaleFactor
}

func New(assets embed.FS) (*Game, error) {
	game1 := Game{}

	var err error

	fontImage, err = loadImage(assets, "assets/font.png")
	if err != nil {
		return nil, errors.Wrap(err, "error on load image")
	}

	playerImage, err = loadImage(assets, "assets/player.png")
	if err != nil {
		return nil, errors.Wrap(err, "error on load image")
	}

	tileSetImage, err = loadImage(assets, "assets/tileset.png")
	if err != nil {
		return nil, errors.Wrap(err, "error on load image")
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
				{I: 0, J: 0},
			},
			1,
		),
		SpriteBackground2: NewSprite(
			[]Frame{
				{I: 1, J: 0},
			},
			1,
		),
		SpriteBackground3: NewSprite(
			[]Frame{
				{I: 2, J: 0},
			},
			1,
		),
		SpriteBackground4: NewSprite(
			[]Frame{
				{I: 3, J: 0},
			},
			1,
		),
		SpriteBackground5: NewSprite(
			[]Frame{
				{I: 4, J: 0},
			},
			1,
		),
		SpriteBackground6: NewSprite(
			[]Frame{
				{I: 5, J: 0},
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
				{I: 0, J: 2},
			},
			1,
		),
		SpriteTile2: NewSprite(
			[]Frame{
				{I: 1, J: 2},
			},
			1,
		),
		SpriteTile3: NewSprite(
			[]Frame{
				{I: 2, J: 2},
			},
			1,
		),
		SpriteFlag1: NewSprite(
			[]Frame{
				{I: 0, J: 3},
			},
			1,
		),
		SpriteFlag2: NewSprite(
			[]Frame{
				{I: 1, J: 3},
			},
			1,
		),
		SpriteFlag3: NewSprite(
			[]Frame{
				{I: 2, J: 3},
			},
			1,
		),
		SpriteBox1: NewSprite(
			[]Frame{
				{I: 0, J: 5},
			},
			1,
		),
		SpriteBox2: NewSprite(
			[]Frame{
				{I: 1, J: 5},
			},
			1,
		),
		SpriteBox3: NewSprite(
			[]Frame{
				{I: 2, J: 5},
			},
			1,
		),
		SpriteBox4: NewSprite(
			[]Frame{
				{I: 3, J: 5},
			},
			1,
		),
		SpriteBox5: NewSprite(
			[]Frame{
				{I: 4, J: 5},
			},
			1,
		),
		SpriteBoxDone1: NewSprite(
			[]Frame{
				{I: 0, J: 6},
			},
			1,
		),
		SpriteBoxDone2: NewSprite(
			[]Frame{
				{I: 1, J: 6},
			},
			1,
		),
		SpriteBoxDone3: NewSprite(
			[]Frame{
				{I: 2, J: 6},
			},
			1,
		),
		SpriteBoxDone4: NewSprite(
			[]Frame{
				{I: 3, J: 6},
			},
			1,
		),
		SpriteBoxDone5: NewSprite(
			[]Frame{
				{I: 4, J: 6},
			},
			1,
		),
	}

	stages, err = loadStages(assets, "assets/stages")
	if err != nil {
		return nil, errors.Wrap(err, "error on load stages")
	}

	loadStage()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetWindowSize(screenWidth*scaleFactor, screenHeight*scaleFactor)

	return &game1, nil
}
