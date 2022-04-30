package game

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
)

const defaultFrameRate = 10

type Sprite struct {
	Images []*ebiten.Image
	Speed  float64
}

type Frame struct {
	I, J int
}

func NewSprite(img *ebiten.Image, frames []image.Rectangle, speed float64) *Sprite {
	images := make([]*ebiten.Image, len(frames))
	for i := range frames {
		images[i], _ = img.SubImage(frames[i]).(*ebiten.Image)
	}

	return &Sprite{
		Images: images,
		Speed:  speed,
	}
}

func (game *Game) loadImages(assets embed.FS) (err error) {
	game.fontImage, err = loadImage(assets, "assets/font.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	game.playerImage, err = loadImage(assets, "assets/player.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	game.tileSetImage, err = loadImage(assets, "assets/tileset.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	return nil
}

func (game *Game) loadSprites() {
	game.sprites = map[SpriteName]*Sprite{
		SpriteIdle: NewSprite(
			game.playerImage,
			[]image.Rectangle{
				image.Rect(tileWidth, 0, tileWidth*2, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteWalking: NewSprite(
			game.playerImage,
			[]image.Rectangle{
				image.Rect(0, 0, tileWidth, tileWidth),
				image.Rect(tileWidth, 0, tileWidth*2, tileWidth),
				image.Rect(tileWidth*2, 0, tileWidth*3, tileWidth),
				image.Rect(tileWidth, 0, tileWidth*2, tileWidth),
			},
			defaultFrameRate,
		),
		SpritePushing: NewSprite(
			game.playerImage,
			[]image.Rectangle{
				image.Rect(tileWidth*3, 0, tileWidth*4, tileWidth),
				image.Rect(tileWidth*4, 0, tileWidth*5, tileWidth),
				image.Rect(tileWidth*5, 0, tileWidth*6, tileWidth),
				image.Rect(tileWidth*4, 0, tileWidth*5, tileWidth),
			},
			defaultFrameRate,
		),
		SpritePushingIdle: NewSprite(
			game.playerImage,
			[]image.Rectangle{
				image.Rect(tileWidth*4, 0, tileWidth*5, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, 0, tileWidth, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, 0, tileWidth*2, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, 0, tileWidth*3, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground4: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*3, 0, tileWidth*4, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground5: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*4, 0, tileWidth*5, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteBackground6: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*5, 0, tileWidth*6, tileWidth),
			},
			defaultFrameRate,
		),
		SpriteWall1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, tileWidth, tileWidth, tileWidth*2),
			},
			defaultFrameRate,
		),
		SpriteWall2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, tileWidth, tileWidth*2, tileWidth*2),
			},
			defaultFrameRate,
		),
		SpriteWall3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, tileWidth, tileWidth*3, tileWidth*2),
			},
			defaultFrameRate,
		),
		SpriteWall4: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*3, tileWidth, tileWidth*4, tileWidth*2),
			},
			defaultFrameRate,
		),
		SpriteTile1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, tileWidth*2, tileWidth, tileWidth*3),
			},
			defaultFrameRate,
		),
		SpriteTile2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, tileWidth*2, tileWidth*2, tileWidth*3),
			},
			defaultFrameRate,
		),
		SpriteTile3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, tileWidth*2, tileWidth*3, tileWidth*3),
			},
			defaultFrameRate,
		),
		SpriteFlag1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, tileWidth*3, tileWidth, tileWidth*4),
			},
			defaultFrameRate,
		),
		SpriteFlag2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, tileWidth*3, tileWidth*2, tileWidth*4),
			},
			defaultFrameRate,
		),
		SpriteFlag3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, tileWidth*3, tileWidth*3, tileWidth*4),
			},
			defaultFrameRate,
		),
		SpriteBox1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, tileWidth*5, tileWidth, tileWidth*6),
			},
			defaultFrameRate,
		),
		SpriteBox2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, tileWidth*5, tileWidth*2, tileWidth*6),
			},
			defaultFrameRate,
		),
		SpriteBox3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, tileWidth*5, tileWidth*3, tileWidth*6),
			},
			defaultFrameRate,
		),
		SpriteBox4: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*3, tileWidth*5, tileWidth*4, tileWidth*6),
			},
			defaultFrameRate,
		),
		SpriteBox5: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*4, tileWidth*5, tileWidth*5, tileWidth*6),
			},
			defaultFrameRate,
		),
		SpriteBoxDone1: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(0, tileWidth*6, tileWidth, tileWidth*7),
			},
			defaultFrameRate,
		),
		SpriteBoxDone2: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth, tileWidth*6, tileWidth*2, tileWidth*7),
			},
			defaultFrameRate,
		),
		SpriteBoxDone3: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*2, tileWidth*6, tileWidth*3, tileWidth*7),
			},
			defaultFrameRate,
		),
		SpriteBoxDone4: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*3, tileWidth*6, tileWidth*4, tileWidth*7),
			},
			defaultFrameRate,
		),
		SpriteBoxDone5: NewSprite(
			game.tileSetImage,
			[]image.Rectangle{
				image.Rect(tileWidth*4, tileWidth*6, tileWidth*5, tileWidth*7),
			},
			defaultFrameRate,
		),
	}
}
