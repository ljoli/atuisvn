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
	totalRows := len(lines)

	statusbar := TuiStatusBar(fmt.Sprintf("[%s]cat:%s:%s", repos, path, rev))
	shortcutbar := TuiShortcutBar(" h help | j/k/gg/G | ^d/u ^f/b scroll | / ? search | n/N | { } para | : goto | q back")
	main := tview.NewTable().SetSelectable(true, false)

	lineNumWidth := len(fmt.Sprintf("%d", totalRows))
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

	// ── search state ──────────────────────────────────────────────────────────
	matches := []int{}
	currentMatch := -1
	searchDir := 1 // +1 = forward, -1 = backward

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
		if len(matches) == 0 {
			return
		}
		curRow, _ := main.GetSelection()
		if searchDir == 1 {
			currentMatch = 0
			for i, m := range matches {
				if m >= curRow {
					currentMatch = i
					break
				}
			}
		} else {
			currentMatch = len(matches) - 1
			for i := len(matches) - 1; i >= 0; i-- {
				if matches[i] <= curRow {
					currentMatch = i
					break
				}
			}
		}
		main.Select(matches[currentMatch], 1)
	}

	gotoMatch := func(delta int) {
		if len(matches) == 0 {
			return
		}
		currentMatch = (currentMatch + delta + len(matches)) % len(matches)
		main.Select(matches[currentMatch], 1)
	}

	// ── input overlay (search / goto) ─────────────────────────────────────────
	searchbar := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetLabelColor(tcell.ColorAqua)

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
	searchbar.SetDoneFunc(func(_ tcell.Key) {
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
			if err == nil && n >= 1 && n <= totalRows {
				main.Select(n-1, 1)
			}
		}
		gotobar.SetText("")
		closeInput()
	})

	// ── layout ────────────────────────────────────────────────────────────────
	s.prim.
		SetRows(0, 1, 1, 1).
		SetBorders(false).
		AddItem(main, 0, 0, 1, 3, 0, 0, false).
		AddItem(inputPages, 1, 0, 1, 3, 0, 0, false).
		AddItem(statusbar, 2, 0, 1, 3, 0, 0, false).
		AddItem(shortcutbar, 3, 0, 1, 3, 0, 0, false)

	// ── input capture ─────────────────────────────────────────────────────────
	lastKey := rune(0)

	clamp := func(row, lo, hi int) int {
		if row < lo {
			return lo
		}
		if row > hi {
			return hi
		}
		return row
	}

	const halfPageSize = 15
	const fullPageSize = 30

	s.prim.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := main.GetSelection()

		switch event.Key() {
		// ── arrow keys ──────────────────────────────────────────────────────
		case tcell.KeyUp:
			main.Select(clamp(row-1, 0, totalRows-1), 0)
			return nil
		case tcell.KeyDown:
			main.Select(clamp(row+1, 0, totalRows-1), 0)
			return nil

		// ── Ctrl+D/U/F/B scroll ─────────────────────────────────────────────
		case tcell.KeyCtrlD:
			main.Select(clamp(row+halfPageSize, 0, totalRows-1), 0)
			return nil
		case tcell.KeyCtrlU:
			main.Select(clamp(row-halfPageSize, 0, totalRows-1), 0)
			return nil
		case tcell.KeyCtrlF:
			main.Select(clamp(row+fullPageSize, 0, totalRows-1), 0)
			return nil
		case tcell.KeyCtrlB:
			main.Select(clamp(row-fullPageSize, 0, totalRows-1), 0)
			return nil

		// ── rune keys ───────────────────────────────────────────────────────
		case tcell.KeyRune:
			r := event.Rune()
			defer func() { lastKey = r }()

			switch r {
			case 'j':
				main.Select(clamp(row+1, 0, totalRows-1), 0)
				return nil
			case 'k':
				main.Select(clamp(row-1, 0, totalRows-1), 0)
				return nil

			// gg / G — file navigation
			case 'G':
				main.Select(totalRows-1, 0)
				return nil
			case 'g':
				if lastKey == 'g' {
					main.Select(0, 0)
					lastKey = 0
				}
				return nil

			// { / } — paragraph jumps (empty lines)
			case '}':
				for i := row + 1; i < totalRows; i++ {
					if strings.TrimSpace(lines[i]) == "" {
						main.Select(i, 0)
						break
					}
				}
				return nil
			case '{':
				for i := row - 1; i >= 0; i-- {
					if strings.TrimSpace(lines[i]) == "" {
						main.Select(i, 0)
						break
					}
				}
				return nil

			// / and ? — search
			case '/':
				searchDir = 1
				searchbar.SetLabel(" / ")
				inputPages.SwitchToPage("search")
				t.app.SetFocus(searchbar)
				return nil
			case '?':
				searchDir = -1
				searchbar.SetLabel(" ? ")
				inputPages.SwitchToPage("search")
				t.app.SetFocus(searchbar)
				return nil
			case 'n':
				gotoMatch(searchDir)
				return nil
			case 'N':
				gotoMatch(-searchDir)
				return nil

			// : — goto line
			case ':':
				inputPages.SwitchToPage("goto")
				t.app.SetFocus(gotobar)
				return nil

			case 'q':
				t.BackScreen()
				return nil
			}
		}
		return event
	})
	t.screen[repos]["cat:"+path+":"+rev] = &s
}
