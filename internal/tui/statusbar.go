package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar displays the top status bar.
type StatusBar struct {
	view        *tview.TextView
	refreshTick int
	lastRefresh time.Time
}

// NewStatusBar creates a new status bar.
func NewStatusBar() *StatusBar {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	return &StatusBar{view: view, lastRefresh: time.Now()}
}

// Primitive returns the tview primitive.
func (s *StatusBar) Primitive() tview.Primitive {
	return s.view
}

// Update updates the status bar display.
func (s *StatusBar) Update(townName, rigName, interval string, connected, stale bool, lastError string) {
	s.refreshTick++
	s.lastRefresh = time.Now()

	var status string

	// Connection indicator with spinning refresh indicator
	spinners := []string{"◐", "◓", "◑", "◒"}
	spinner := spinners[s.refreshTick%4]

	if connected {
		if stale {
			status = "[yellow]⚠ Stale[-]"
		} else {
			status = fmt.Sprintf("[green]%s[-]", spinner)
		}
	} else {
		status = "[red]✗[-]"
	}

	// Build the status line
	line := fmt.Sprintf("[::b]gastop[-] | Town: [::b]%s[-]", townName)
	if rigName != "" {
		line += fmt.Sprintf(" | Rig: [::b]%s[-]", rigName)
	}
	line += fmt.Sprintf(" | ↻ %s %s", interval, status)

	if lastError != "" {
		line += fmt.Sprintf(" | [red]%s[-]", truncate(lastError, 30))
	}

	// Add timestamp and help hint
	line += fmt.Sprintf(" | [::d]%s | ? help[-]", time.Now().Format("15:04:05"))

	s.view.SetText(line)
}

// SetBackgroundColor sets the background color.
func (s *StatusBar) SetBackgroundColor(color tcell.Color) {
	s.view.SetBackgroundColor(color)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
