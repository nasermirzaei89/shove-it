package game

import (
	"embed"

	"github.com/pkg/errors"
)

const defaultFrameRate = 10

type Sprite struct {
	Frames []Frame
	Speed  float64
}

type Frame struct {
	I, J int
}

func NewSprite(frames []Frame, speed float64) *Sprite {
	return &Sprite{
		Frames: frames,
		Speed:  speed,
	}
}

func loadImages(assets embed.FS) (err error) {
	fontImage, err = loadImage(assets, "assets/font.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	playerImage, err = loadImage(assets, "assets/player.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	tileSetImage, err = loadImage(assets, "assets/tileset.png")
	if err != nil {
		return errors.Wrap(err, "error on load image")
	}

	return nil
}

func loadSprites() {
	sprites = map[SpriteName]*Sprite{
		SpriteIdle: NewSprite(
			[]Frame{
				{I: 1, J: 0},
			},
			defaultFrameRate,
		),
		SpriteWalking: NewSprite(
			[]Frame{
				{I: 0, J: 0},
				{I: 1, J: 0},
				{I: 2, J: 0},
				{I: 1, J: 0},
			},
			defaultFrameRate,
		),
		SpritePushing: NewSprite(
			[]Frame{
				{I: 3, J: 0},
				{I: 4, J: 0},
				{I: 5, J: 0},
				{I: 4, J: 0},
			},
			defaultFrameRate,
		),
		SpritePushingIdle: NewSprite(
			[]Frame{
				{I: 4, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground1: NewSprite(
			[]Frame{
				{I: 0, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground2: NewSprite(
			[]Frame{
				{I: 1, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground3: NewSprite(
			[]Frame{
				{I: 2, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground4: NewSprite(
			[]Frame{
				{I: 3, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground5: NewSprite(
			[]Frame{
				{I: 4, J: 0},
			},
			defaultFrameRate,
		),
		SpriteBackground6: NewSprite(
			[]Frame{
				{I: 5, J: 0},
			},
			defaultFrameRate,
		),
		SpriteWall1: NewSprite(
			[]Frame{
				{I: 0, J: 1},
			},
			defaultFrameRate,
		),
		SpriteWall2: NewSprite(
			[]Frame{
				{I: 1, J: 1},
			},
			defaultFrameRate,
		),
		SpriteWall3: NewSprite(
			[]Frame{
				{I: 2, J: 1},
			},
			defaultFrameRate,
		),
		SpriteWall4: NewSprite(
			[]Frame{
				{I: 3, J: 1},
			},
			defaultFrameRate,
		),
		SpriteTile1: NewSprite(
			[]Frame{
				{I: 0, J: 2},
			},
			defaultFrameRate,
		),
		SpriteTile2: NewSprite(
			[]Frame{
				{I: 1, J: 2},
			},
			defaultFrameRate,
		),
		SpriteTile3: NewSprite(
			[]Frame{
				{I: 2, J: 2},
			},
			defaultFrameRate,
		),
		SpriteFlag1: NewSprite(
			[]Frame{
				{I: 0, J: 3},
			},
			defaultFrameRate,
		),
		SpriteFlag2: NewSprite(
			[]Frame{
				{I: 1, J: 3},
			},
			defaultFrameRate,
		),
		SpriteFlag3: NewSprite(
			[]Frame{
				{I: 2, J: 3},
			},
			defaultFrameRate,
		),
		SpriteBox1: NewSprite(
			[]Frame{
				{I: 0, J: 5},
			},
			defaultFrameRate,
		),
		SpriteBox2: NewSprite(
			[]Frame{
				{I: 1, J: 5},
			},
			defaultFrameRate,
		),
		SpriteBox3: NewSprite(
			[]Frame{
				{I: 2, J: 5},
			},
			defaultFrameRate,
		),
		SpriteBox4: NewSprite(
			[]Frame{
				{I: 3, J: 5},
			},
			defaultFrameRate,
		),
		SpriteBox5: NewSprite(
			[]Frame{
				{I: 4, J: 5},
			},
			defaultFrameRate,
		),
		SpriteBoxDone1: NewSprite(
			[]Frame{
				{I: 0, J: 6},
			},
			defaultFrameRate,
		),
		SpriteBoxDone2: NewSprite(
			[]Frame{
				{I: 1, J: 6},
			},
			defaultFrameRate,
		),
		SpriteBoxDone3: NewSprite(
			[]Frame{
				{I: 2, J: 6},
			},
			defaultFrameRate,
		),
		SpriteBoxDone4: NewSprite(
			[]Frame{
				{I: 3, J: 6},
			},
			defaultFrameRate,
		),
		SpriteBoxDone5: NewSprite(
			[]Frame{
				{I: 4, J: 6},
			},
			defaultFrameRate,
		),
	}
}
