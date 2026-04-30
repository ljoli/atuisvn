package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) NewTuiCat(repos string, path string, rev string) {
	s := TuiScreen{
		prim: tview.NewGrid(),
	}

	catOutput := t.SvnCat(repos, path, rev)
	lines := strings.Split(catOutput, "\n")

	statusbar := TuiStatusBar(fmt.Sprintf("[%s]cat:%s:%s", repos, path, rev))
	shortcutbar := TuiShortcutBar(" h/? help | j/k move | / search | n/N next/prev | : goto line | q back")
	main := tview.NewTable().SetSelectable(true, false)

	lineNumWidth := len(fmt.Sprintf("%d", len(lines)))
	for i, v := range lines {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		main.SetCell(i, 0,
			tview.NewTableCell(lineNum).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignRight).
				SetSelectable(false))
		main.SetCell(i, 1,
			tview.NewTableCell(tview.Escape(v)).SetExpansion(1))
	}

	searchbar := tview.NewInputField().
		SetLabel(" / ").
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetLabelColor(tcell.ColorAqua)

	matches := []int{}
	currentMatch := -1

	updateMatches := func(query string) {
		matches = []int{}
		currentMatch = -1
		if query == "" {
			return
		}
		q := strings.ToLower(query)
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), q) {
				matches = append(matches, i)
			}
		}
		if len(matches) > 0 {
			currentMatch = 0
			main.Select(matches[0], 1)
		}
	}

	gotoMatch := func(delta int) {
		if len(matches) == 0 {
			return
		}
		currentMatch = (currentMatch + delta + len(matches)) % len(matches)
		main.Select(matches[currentMatch], 1)
	}

	// inputPages: hidden by default, shows searchbar or gotobar on demand
	inputPages := tview.NewPages()
	inputPages.AddPage("empty", tview.NewBox(), true, true)
	inputPages.AddPage("search", searchbar, true, false)

	closeInput := func() {
		inputPages.SwitchToPage("empty")
		t.app.SetFocus(s.prim)
	}

	searchbar.SetChangedFunc(func(text string) {
		updateMatches(text)
	})
	searchbar.SetDoneFunc(func(key tcell.Key) {
		closeInput()
	})

	gotobar := tview.NewInputField().
		SetLabel(" : ").
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetLabelColor(tcell.ColorGreen).
		SetAcceptanceFunc(tview.InputFieldInteger)
	inputPages.AddPage("goto", gotobar, true, false)

	gotobar.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			n, err := strconv.Atoi(strings.TrimSpace(gotobar.GetText()))
			if err == nil && n >= 1 && n <= len(lines) {
				main.Select(n-1, 1)
			}
		}
		gotobar.SetText("")
		closeInput()
	})

	s.prim.
		SetRows(0, 1, 1, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(inputPages, 1, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 2, 0, 1, 3, 0, 0, false).
		AddItem(shortcutbar, 3, 0, 1, 3, 0, 0, false)

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
			case 'n':
				gotoMatch(1)
				return nil
			case 'N':
				gotoMatch(-1)
				return nil
			case '/':
				inputPages.SwitchToPage("search")
				t.app.SetFocus(searchbar)
				return nil
			case ':':
				inputPages.SwitchToPage("goto")
				t.app.SetFocus(gotobar)
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
