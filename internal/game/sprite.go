package game

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

func loadSprites() {
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
}
