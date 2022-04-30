package game

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pkg/errors"
)

type Game struct {
	fontImage    *ebiten.Image
	playerImage  *ebiten.Image
	tileSetImage *ebiten.Image

	fontCache map[int32]*ebiten.Image

	player     *Player
	boxes      []*Box
	objects    []*object
	sprites    map[SpriteName]*Sprite
	stages     []Stage
	stageIndex int

	shouldDraw bool
}

func (game *Game) Update() error {
	done := true

	for i := range game.boxes {
		game.boxes[i].Update(game)

		if !game.boxes[i].Done(game) {
			done = false
		}
	}

	if game.player != nil {
		game.player.Update(game)
	}

	if done {
		game.nextStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		game.nextStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		game.prevStage()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		game.startStage()
	}

	return nil
}

const (
	stepX  = 2
	stepY  = 26
	stageX = 22
	stageY = 26
)

func (game *Game) Draw(screen *ebiten.Image) {
	if !game.shouldDraw {
		ebitenutil.DrawLine(screen, 0, 0, -1, -1, color.Black)

		return
	}

	for i := range game.objects {
		game.objects[i].Draw(game, screen)
	}

	for i := range game.boxes {
		game.boxes[i].Draw(game, screen)
	}

	steps := 0

	if game.player != nil {
		game.player.Draw(game, screen)
		steps = len(game.player.history)
	}

	// HUD
	game.DrawText(screen, stepX, stepY, fmt.Sprintf("STEP %d", steps))
	game.DrawText(screen, stageX, stageY, fmt.Sprintf("STAGE %s", game.stages[game.stageIndex].Name))

	game.shouldDraw = false
}

func (game *Game) Layout(int, int) (int, int) {
	return screenWidth * scaleFactor, screenHeight * scaleFactor
}

func New(assets embed.FS) (*Game, error) {
	game := Game{
		fontImage:    nil,
		playerImage:  nil,
		tileSetImage: nil,
		fontCache:    make(map[int32]*ebiten.Image),
		player:       nil,
		boxes:        nil,
		objects:      nil,
		sprites:      nil,
		stages:       nil,
		stageIndex:   0,
		shouldDraw:   false,
	}

	err := game.loadImages(assets)
	if err != nil {
		return nil, errors.Wrap(err, "error on load Images")
	}

	game.loadSprites()

	err = game.loadStages(assets, "assets/stages")
	if err != nil {
		return nil, errors.Wrap(err, "error on load stages")
	}

	game.startStage()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetWindowSize(screenWidth*scaleFactor, screenHeight*scaleFactor)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)

	return &game, nil
}
