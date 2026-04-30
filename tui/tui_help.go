package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type keyBinding struct {
	key  string
	desc string
}

var helpBindings = map[string][]keyBinding{
	"main": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"Enter", "Open repository"},
		{"q", "Quit"},
		{"h / ?", "Toggle this help"},
	},
	"tree": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"Enter", "Enter directory"},
		{"l", "Open log for selection"},
		{"q", "Back / Quit"},
		{"h / ?", "Toggle this help"},
	},
	"log": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"Enter", "Open revision detail"},
		{"q", "Back"},
		{"h / ?", "Toggle this help"},
	},
	"rev": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"Enter", "Open diff for changed path"},
		{"c", "Open file content (svn cat)"},
		{"q", "Back"},
		{"h / ?", "Toggle this help"},
	},
	"cat": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"q", "Back"},
		{"h / ?", "Toggle this help"},
	},
	"diff": {
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"q", "Back"},
		{"h / ?", "Toggle this help"},
	},
}

func (t *Tui) currentScreenType() string {
	length := len(t.history_screen)
	if length == 0 {
		return "main"
	}
	screen := t.history_screen[length-1][1]
	for prefix := range helpBindings {
		if strings.HasPrefix(screen, prefix) {
			return prefix
		}
	}
	return "main"
}

func (t *Tui) ShowHelp() {
	screenType := t.currentScreenType()
	bindings := helpBindings[screenType]

	title := " Shortcuts — " + screenType + " screen "

	table := tview.NewTable()
	for i, b := range bindings {
		table.SetCell(i, 0,
			tview.NewTableCell("[yellow]"+b.key+"[-]").
				SetAlign(tview.AlignLeft))
		table.SetCell(i, 1,
			tview.NewTableCell(" "+b.desc).
				SetExpansion(1))
	}

	frame := tview.NewFrame(table).
		SetBorders(0, 0, 1, 1, 2, 2)
	frame.SetTitle(title).SetTitleAlign(tview.AlignCenter)
	frame.SetBorder(true).SetBorderColor(tcell.ColorYellow)
	frame.SetTitleColor(tcell.ColorYellow)

	height := len(bindings) + 5
	overlay := tview.NewGrid().
		SetColumns(0, 52, 0).
		SetRows(0, height, 0).
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	overlay.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			t.HideHelp()
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == '?' || event.Rune() == 'h' || event.Rune() == 'H' {
				t.HideHelp()
				return nil
			}
		}
		return event
	})

	t.pages.AddPage("help", overlay, true, true)
	t.app.SetFocus(overlay)
}

func (t *Tui) HideHelp() {
	t.pages.RemovePage("help")
	length := len(t.history_screen)
	if length > 0 {
		entry := t.history_screen[length-1]
		repos := entry[0]
		screen := entry[1]
		if s, ok := t.screen[repos][screen]; ok {
			t.app.SetFocus(s.prim)
		}
	}
}
