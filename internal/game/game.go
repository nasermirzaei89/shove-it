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

var shouldDraw bool

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
		startStage()
	}

	return nil
}

const (
	stepX  = 2
	stepY  = 26
	stageX = 22
	stageY = 26
)

func (g *Game) Draw(screen *ebiten.Image) {
	if !shouldDraw {
		ebitenutil.DrawLine(screen, 0, 0, -1, -1, color.Black)

		return
	}

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
	DrawText(screen, stepX, stepY, fmt.Sprintf("STEP %d", steps))
	DrawText(screen, stageX, stageY, fmt.Sprintf("STAGE %s", stages[stageIndex].Name))

	shouldDraw = false
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth * scaleFactor, screenHeight * scaleFactor
}

func New(assets embed.FS) (*Game, error) {
	game1 := Game{}

	err := loadImages(assets)
	if err != nil {
		return nil, errors.Wrap(err, "error on load Images")
	}

	loadSprites()

	err = loadStages(assets, "assets/stages")
	if err != nil {
		return nil, errors.Wrap(err, "error on load stages")
	}

	startStage()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Shove It")
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetWindowSize(screenWidth*scaleFactor, screenHeight*scaleFactor)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)

	return &game1, nil
}
