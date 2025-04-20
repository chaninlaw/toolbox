package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	initialCursor     = 1
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	checkboxStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle        = lipgloss.NewStyle().MarginLeft(2)
	destructiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	progressEmpty = subtleStyle.Render(progressEmptyChar)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
