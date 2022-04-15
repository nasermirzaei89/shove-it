package game

const (
	ItemBackground = iota
	ItemWall
	ItemTile
	ItemTileFlagged
	ItemPlayer
	ItemBox
	ItemBoxDone
)

var mapData = [][]int{
	{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 2, 1, 1, 1, 1, 0, 0, 0},
	{0, 0, 1, 1, 1, 5, 2, 5, 3, 1, 0, 0, 0},
	{0, 0, 1, 3, 2, 5, 4, 1, 1, 1, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 5, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}
