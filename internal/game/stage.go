package game

import "encoding/xml"

const (
	ItemBackground1 = iota + 1
	ItemBackground2
	ItemBackground3
	ItemBackground4
	ItemBackground5
	ItemWall1 = iota + 4
	ItemWall2
	ItemWall3
	ItemWall4
	ItemTile1 = iota + 8
	ItemTile2
	ItemTile3
	ItemTileFlagged1 = iota + 13
	ItemTileFlagged2
	ItemTileFlagged3
	ItemPlayer1 = iota + 18
	ItemPlayer2
	ItemPlayer3
	ItemBox1 = iota + 23
	ItemBox2
	ItemBox3
	ItemBox4
	ItemBox5
	ItemBoxDone1 = iota + 26
	ItemBoxDone2
	ItemBoxDone3
	ItemBoxDone4
	ItemBoxDone5
)

type Stage struct {
	Name string
	Data [][]int
	TMX  TMX
}

type TMX struct {
	XMLName xml.Name `xml:"map"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Data    string   `xml:"layer>data"`
}

func (stg Stage) Width() int {
	return stg.TMX.Width * tileWidth
}

func (stg Stage) Height() int {
	return stg.TMX.Height * tileWidth
}

func (stg Stage) ValueAt(i, j int) int {
	return stg.Data[j][i]
}

func (stg Stage) IsWall(i, j int) bool {
	switch stg.Data[j][i] {
	case ItemWall1, ItemWall2, ItemWall3, ItemWall4:
		return true
	default:
		return false
	}
}
