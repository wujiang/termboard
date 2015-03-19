package termboard

import (
	"errors"
	"strings"

	"github.com/golang/glog"
	"github.com/nsf/termbox-go"
)

const (
	Up    string = "up"
	Down  string = "down"
	Left  string = "left"
	Right string = "right"
	End   string = "end"
	Begin string = "begin"

	gridSpanX int = 10
	gridSpanY int = 4
)

type TPosition struct {
	X int
	Y int
}

type GPosition struct {
	X int
	Y int
}

type Tile struct {
	R  rune
	FG termbox.Attribute
	BG termbox.Attribute
}

type Board struct {
	TileCountX int
	TileCountY int
	Width      int
	Height     int
	CursorPos  GPosition
	Grid       [][]Tile
	Messages   []Message
}

type Message struct {
	Text        string
	StartPos    GPosition
	AlignCenter bool
	FG          termbox.Attribute
	BG          termbox.Attribute
}

// Get terminal's center
func getTBCenter() TPosition {
	w, h := termbox.Size()
	return TPosition{
		X: w / 2,
		Y: h / 2,
	}
}

// Initialize a board.
func Init(tileCountX, tileCountY int) *Board {
	if err := termbox.Init(); err != nil {
		glog.Fatalln(err)
	}
	grid := make([][]Tile, tileCountX)
	for i := range grid {
		grid[i] = make([]Tile, tileCountY)
	}
	return &Board{
		TileCountX: tileCountX,
		TileCountY: tileCountY,
		Width:      tileCountX * gridSpanX,
		Height:     tileCountY * gridSpanY,
		CursorPos:  GPosition{0, 0},
		Grid:       grid,
	}
}

func Close() {
	termbox.Close()
}

func (b *Board) GetCenter() GPosition {
	return GPosition{
		X: (b.TileCountX - 1) / 2,
		Y: (b.TileCountY - 1) / 2,
	}
}

// Convert grid positions to termbox coordinates
func (b *Board) toTBPosition(p GPosition) TPosition {
	tbCenter := getTBCenter()
	center := b.GetCenter()
	offsetX := 0
	offsetY := 0
	if b.TileCountX%2 == 0 {
		offsetX = gridSpanX / 2
	}
	if b.TileCountY%2 == 0 {
		offsetY = gridSpanY / 2
	}
	x := tbCenter.X - (center.X-p.X)*gridSpanX - offsetX
	y := tbCenter.Y - (center.Y-p.Y)*gridSpanY - offsetY
	return TPosition{x, y}
}

func (b *Board) setTile(p GPosition, t Tile) {
	tbPos := b.toTBPosition(p)
	termbox.SetCell(tbPos.X, tbPos.Y, t.R, t.FG, t.BG)
}

func (b *Board) MoveCursor(direction string) error {
	x := b.CursorPos.X
	y := b.CursorPos.Y
	switch direction {
	case Up:
		y--
	case Down:
		y++
	case Left:
		x--
	case Right:
		x++
	case End:
		x = b.TileCountX - 1
	case Begin:
		x = 0
	}

	if !b.isValidPosition(GPosition{x, y}) {
		return errors.New("Invalid position")
	}
	b.CursorPos.X = x
	b.CursorPos.Y = y
	return nil
}

func (b *Board) isValidPosition(p GPosition) bool {
	return p.X >= 0 && p.X < b.TileCountX && p.Y >= 0 && p.Y < b.TileCountY
}

func (b *Board) SetCursor(p GPosition) {
	b.CursorPos = p
	tbPos := b.toTBPosition(p)
	termbox.SetCursor(tbPos.X, tbPos.Y)
}

func (b *Board) PinCursor(t Tile) {
	b.Grid[b.CursorPos.X][b.CursorPos.Y] = t
}

func (b *Board) fill(x, y, w, h int, t Tile) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, t.R, t.FG, t.BG)
		}
	}
}

func (b *Board) PrintLines(m Message) {
	lines := strings.Split(m.Text, "\n")
	tbCenter := getTBCenter()
	centerX := tbCenter.X
	startPos := b.toTBPosition(m.StartPos)
	x := startPos.X - gridSpanX/2
	y := startPos.Y
	for _, l := range lines {
		if m.AlignCenter {
			x = centerX - len(l)/2
		}
		xstart := x
		for _, c := range l {
			termbox.SetCell(xstart, y, c, m.FG, m.BG)
			xstart++
		}
		y++
	}

}

func (b *Board) drawCells() {
	for x, ts := range b.Grid {
		for y, t := range ts {
			p := GPosition{x, y}
			b.setTile(p, t)
		}
	}
}

func (b *Board) Redraw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	tbCenter := getTBCenter()
	tbLeftXPos := tbCenter.X - b.Width/2
	tbUpYPos := tbCenter.Y - b.Height/2

	// draw the grid
	for yoffset := 0; yoffset <= b.TileCountY; yoffset++ {
		for xoffset := 0; xoffset <= b.TileCountX; xoffset++ {
			xstart := tbLeftXPos + xoffset*gridSpanX
			ystart := tbUpYPos + yoffset*gridSpanY
			// all intersections
			termbox.SetCell(xstart, ystart, '+', termbox.ColorDefault, termbox.ColorDefault)
			if xoffset < b.TileCountX {
				b.fill(xstart+1, ystart, gridSpanX-1, 1, Tile{'-', termbox.ColorDefault, termbox.ColorDefault})
			}
			if yoffset < b.TileCountY {
				b.fill(xstart, ystart+1, 1, gridSpanY-1, Tile{'|', termbox.ColorDefault, termbox.ColorDefault})
			}

		}
	}

	b.SetCursor(b.CursorPos)

	for _, msg := range b.Messages {
		b.PrintLines(msg)
	}

	b.drawCells()
	termbox.Flush()
}
