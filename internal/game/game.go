package game

import (
	"embed"
	"fmt"

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
	DrawText(screen, 2, 26, fmt.Sprintf("STEP %d", steps))
	DrawText(screen, 22, 26, fmt.Sprintf("STAGE %s", stages[stageIndex].Name))
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

	loadSprites()

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
