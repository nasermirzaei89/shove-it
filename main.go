package main

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nasermirzaei89/shove-it/internal/game"
	"github.com/pkg/errors"
)

//go:embed assets
var assets embed.FS

func main() {
	game1, err := game.New(assets)
	if err != nil {
		panic(errors.Wrap(err, "error on new game"))
	}

	if err := ebiten.RunGame(game1); err != nil {
		panic(errors.Wrap(err, "error on run game"))
	}
}
