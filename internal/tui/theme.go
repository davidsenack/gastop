package tui

import "github.com/gdamore/tcell/v2"

// Theme defines the color scheme for gastop.
// Using a cohesive palette inspired by terminal aesthetics.
type Theme struct {
	// Base colors
	Background     tcell.Color
	Foreground     tcell.Color
	BorderColor    tcell.Color
	TitleColor     tcell.Color
	SelectionBg    tcell.Color
	SelectionFg    tcell.Color

	// Status colors
	Success    tcell.Color
	Warning    tcell.Color
	Error      tcell.Color
	Info       tcell.Color
	Muted      tcell.Color

	// State colors for polecats/beads
	Working    tcell.Color
	Done       tcell.Color
	Idle       tcell.Color
	Stuck      tcell.Color
	InProgress tcell.Color
	Blocked    tcell.Color

	// Accent colors
	Accent1 tcell.Color
	Accent2 tcell.Color
	Accent3 tcell.Color
}

// Tag colors (for tview dynamic colors)
type ThemeTags struct {
	Success    string
	Warning    string
	Error      string
	Info       string
	Muted      string
	Working    string
	Done       string
	Idle       string
	Stuck      string
	InProgress string
	Blocked    string
	Accent1    string
	Accent2    string
	Accent3    string
	Title      string
	Dim        string
}

// DefaultTheme returns the default gastop color scheme.
// A dark theme with cyan/teal accents for a modern terminal look.
func DefaultTheme() *Theme {
	return &Theme{
		// Base - dark background with light text
		Background:  tcell.ColorDefault,
		Foreground:  tcell.NewRGBColor(220, 220, 220), // Light gray
		BorderColor: tcell.NewRGBColor(80, 80, 100),   // Muted blue-gray
		TitleColor:  tcell.NewRGBColor(100, 200, 200), // Teal
		SelectionBg: tcell.NewRGBColor(40, 80, 100),   // Dark teal
		SelectionFg: tcell.ColorWhite,

		// Status
		Success: tcell.NewRGBColor(80, 200, 120),  // Green
		Warning: tcell.NewRGBColor(230, 180, 80),  // Amber
		Error:   tcell.NewRGBColor(220, 80, 80),   // Red
		Info:    tcell.NewRGBColor(80, 160, 220),  // Blue
		Muted:   tcell.NewRGBColor(120, 120, 130), // Gray

		// States
		Working:    tcell.NewRGBColor(100, 180, 255), // Bright blue
		Done:       tcell.NewRGBColor(80, 200, 120),  // Green
		Idle:       tcell.NewRGBColor(120, 120, 130), // Gray
		Stuck:      tcell.NewRGBColor(220, 80, 80),   // Red
		InProgress: tcell.NewRGBColor(230, 180, 80),  // Amber
		Blocked:    tcell.NewRGBColor(180, 100, 180), // Purple

		// Accents
		Accent1: tcell.NewRGBColor(100, 200, 200), // Teal
		Accent2: tcell.NewRGBColor(180, 140, 220), // Purple
		Accent3: tcell.NewRGBColor(255, 180, 100), // Orange
	}
}

// Tags returns tview-compatible color tags for the theme.
func (t *Theme) Tags() *ThemeTags {
	return &ThemeTags{
		Success:    "#50c878",
		Warning:    "#e6b450",
		Error:      "#dc5050",
		Info:       "#50a0dc",
		Muted:      "#787882",
		Working:    "#64b4ff",
		Done:       "#50c878",
		Idle:       "#787882",
		Stuck:      "#dc5050",
		InProgress: "#e6b450",
		Blocked:    "#b464b4",
		Accent1:    "#64c8c8",
		Accent2:    "#b48cdc",
		Accent3:    "#ffb464",
		Title:      "#64c8c8",
		Dim:        "#606070",
	}
}

// Global theme instance
var currentTheme = DefaultTheme()
var currentTags = currentTheme.Tags()

// GetTheme returns the current theme.
func GetTheme() *Theme {
	return currentTheme
}

// GetTags returns the current theme's color tags.
func GetTags() *ThemeTags {
	return currentTags
}
