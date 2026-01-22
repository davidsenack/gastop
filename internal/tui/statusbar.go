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
	theme := GetTheme()
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(theme.Foreground)

	view.SetBackgroundColor(theme.BorderColor)

	return &StatusBar{view: view, lastRefresh: time.Now()}
}

// Primitive returns the tview primitive.
func (s *StatusBar) Primitive() tview.Primitive {
	return s.view
}

// Update updates the status bar display.
func (s *StatusBar) Update(townName, rigName, interval string, connected, stale bool, lastError string) {
	tags := GetTags()
	s.refreshTick++
	s.lastRefresh = time.Now()

	var status string

	// Connection indicator with spinning refresh indicator
	spinners := []string{"◐", "◓", "◑", "◒"}
	spinner := spinners[s.refreshTick%4]

	if connected {
		if stale {
			status = "[" + tags.Warning + "]⚠ Stale[-]"
		} else {
			status = "[" + tags.Success + "]" + spinner + "[-]"
		}
	} else {
		status = "[" + tags.Error + "]✗[-]"
	}

	// Build the status line
	line := fmt.Sprintf("[" + tags.Accent1 + "][::b]gastop[-][-] │ Town: [::b]%s[-]", townName)
	if rigName != "" {
		line += fmt.Sprintf(" │ Rig: [::b]%s[-]", rigName)
	}
	line += fmt.Sprintf(" │ ↻ %s %s", interval, status)

	if lastError != "" {
		line += fmt.Sprintf(" │ [" + tags.Error + "]%s[-]", truncate(lastError, 30))
	}

	// Add timestamp and help hint
	line += fmt.Sprintf(" │ [" + tags.Dim + "]%s[-] │ [" + tags.Muted + "]? help[-]", time.Now().Format("15:04:05"))

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
