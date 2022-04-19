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
