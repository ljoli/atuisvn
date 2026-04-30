package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) NewTuiCat(repos string, path string, rev string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}

	catOutput := t.SvnCat(repos, path, rev)

	statusbar := TuiStatusBar(fmt.Sprintf("[%s]cat:%s:%s", repos, path, rev))
	shortcutbar := TuiShortcutBar(" h/? help | j/k move | q back")
	main := tview.NewTable().SetSelectable(true, false)

	for i, v := range strings.Split(catOutput, "\n") {
		main.SetCell(i, 0,
			tview.NewTableCell(tview.Escape(v)).SetExpansion(1))
	}

	s.prim.
		SetRows(0, 1, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 1, 0, 1, 3, 0, 0, false).
		AddItem(shortcutbar, 2, 0, 1, 3, 0, 0, false)

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'k':
				row, _ := main.GetSelection()
				row--
				main.Select(row, 0)
				return nil
			case 'j':
				row, _ := main.GetSelection()
				if row < main.GetRowCount()-1 {
					row++
				}
				main.Select(row, 0)
				return nil
			case 'q':
				t.BackScreen()
				return nil
			}
		case tcell.KeyUp:
			row, _ := main.GetSelection()
			row--
			main.Select(row, 0)
			return nil
		case tcell.KeyDown:
			row, _ := main.GetSelection()
			if row < main.GetRowCount()-1 {
				row++
			}
			main.Select(row, 0)
			return nil
		}
		return event
	})
	t.screen[repos]["cat:"+path+":"+rev] = &s
}
