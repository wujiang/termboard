package termboard

import (
	"testing"

	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func mockBoard(tileCountX, tileCountY int) *Board {
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

var (
	board3 *Board = mockBoard(3, 3)
	board4 *Board = mockBoard(4, 4)
)

func TestGetCenter(t *testing.T) {
	assert.Equal(t, board3.GetCenter(), GPosition{X: 1, Y: 1})
	assert.Equal(t, board4.GetCenter(), GPosition{X: 1, Y: 1})
}

func TestToTBPosition(t *testing.T) {
	p := GPosition{X: 1, Y: 1}
	assert.Equal(t, board3.toTBPosition(p), TPosition{X: 0, Y: 0})
	assert.Equal(t, board4.toTBPosition(p), TPosition{X: -5, Y: -2})
}

func TestMoveCursor(t *testing.T) {
	assert.Equal(t, board3.CursorPos, GPosition{X: 0, Y: 0})
	board3.MoveCursor(Up)
	assert.Equal(t, board3.CursorPos, GPosition{X: 0, Y: 0})
	board3.MoveCursor(Down)
	assert.Equal(t, board3.CursorPos, GPosition{X: 0, Y: 1})
	board3.MoveCursor(Up)
	assert.Equal(t, board3.CursorPos, GPosition{X: 0, Y: 0})
}

func TestIsValidPosition(t *testing.T) {
	assert.True(t, board3.isValidPosition(GPosition{0, 0}))
	assert.True(t, board3.isValidPosition(GPosition{0, 1}))
	assert.False(t, board3.isValidPosition(GPosition{0, -1}))
}

func TestPinCursor(t *testing.T) {
	tile := Tile{
		R:  'X',
		FG: termbox.ColorBlue,
		BG: termbox.ColorWhite,
	}
	board3.PinCursor(tile)
	assert.Equal(t, board3.Grid[board3.CursorPos.X][board3.CursorPos.Y],
		tile)
}
