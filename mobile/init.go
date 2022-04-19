package mobile

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/nasermirzaei89/shove-it/internal/game"
	"github.com/pkg/errors"
)

//go:generate cp -r ../assets .

//go:embed assets/spritesheet.png
var spriteSheetPNG []byte

//go:embed assets/font.png
var fontPNG []byte

//nolint:gochecknoinits
func init() {
	game1, err := game.New(spriteSheetPNG, fontPNG)
	if err != nil {
		panic(errors.Wrap(err, "error on new game"))
	}

	mobile.SetGame(game1)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
