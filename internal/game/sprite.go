package game

import "image"

type Sprite struct {
	Frames []Frame
	Speed  float64
}

type Frame struct {
	Rect image.Rectangle
}

func NewSprite(frames []Frame, speed float64) *Sprite {
	return &Sprite{
		Frames: frames,
		Speed:  speed,
	}
}
