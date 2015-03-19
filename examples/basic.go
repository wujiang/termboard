package main

import (
	"github.com/golang/glog"
	"github.com/nsf/termbox-go"
	"github.com/wujiang/termboard"
)

func main() {
	b := termboard.Init(4, 4)
	defer termboard.Close()

	b.Messages = append(b.Messages, termboard.Message{
		Text:        "Title",
		StartPos:    termboard.GPosition{0, -1},
		FG:          termbox.ColorBlue,
		BG:          termbox.ColorDefault,
		AlignCenter: true,
	})
	b.Messages = append(b.Messages, termboard.Message{
		Text: `
Key bindings
- left: C-b, h
- right: C-f, l
- up: C-p, k
- down: C-n, j
- end: C-e, $
- begin: C-a, 0
- pin: Enter, i
`,
		StartPos: termboard.GPosition{0, 5},
		FG:       termbox.ColorDefault,
		BG:       termbox.ColorDefault,
	})

	b.Redraw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			// emacs and arrow key bindings
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				b.MoveCursor(termboard.Left)
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				b.MoveCursor(termboard.Down)
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				b.MoveCursor(termboard.Up)
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				b.MoveCursor(termboard.Right)
			case termbox.KeyEnd, termbox.KeyCtrlE:
				b.MoveCursor(termboard.End)
			case termbox.KeyHome, termbox.KeyCtrlA:
				b.MoveCursor(termboard.Begin)
			case termbox.KeyEnter:
				b.PinCursor(termboard.Tile{'X', termbox.ColorBlue, termbox.ColorWhite})
			}

			// vim key bindings
			switch ev.Ch {
			case 'i':
				b.PinCursor(termboard.Tile{'O', termbox.ColorDefault, termbox.ColorDefault})
			case 'q':
				break mainloop
			case 'h':
				b.MoveCursor(termboard.Left)
			case 'j':
				b.MoveCursor(termboard.Down)
			case 'k':
				b.MoveCursor(termboard.Up)
			case 'l':
				b.MoveCursor(termboard.Right)
			case '$':
				b.MoveCursor(termboard.End)
			case '0':
				b.MoveCursor(termboard.Begin)

			}

		case termbox.EventError:
			glog.Fatal("Termbox EventError")
		}

		b.Redraw()
	}
}
