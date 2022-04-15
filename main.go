package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nasermirzaei89/shove-it/internal/game"
	"github.com/pkg/errors"
)

func main() {
	game1, err := game.New()
	if err != nil {
		panic(errors.Wrap(err, "error on new game"))
	}

	if err := ebiten.RunGame(game1); err != nil {
		panic(errors.Wrap(err, "error on run game"))
	}
}
