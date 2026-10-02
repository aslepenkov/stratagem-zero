package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"stratagem-zero/internal/stratagem"
)

var (
	// Direction symbols
	symbolUp    = "▲"
	symbolDown  = "▼"
	symbolLeft  = "◄"
	symbolRight = "►"

	// Colors
	yellowColor = lipgloss.Color("#FFE81F") // Helldivers yellow
	cyanColor   = lipgloss.Color("#00E5FF")
	greenColor  = lipgloss.Color("#00FF66")
	redColor    = lipgloss.Color("#FF3333")
	grayColor   = lipgloss.Color("#555555")
	whiteColor  = lipgloss.Color("#FFFFFF")
	dimColor    = lipgloss.Color("#888888")

	// Lip Gloss Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellowColor).
			MarginBottom(0)

	stratagemNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(whiteColor).
				Padding(0, 1).
				MarginBottom(0)

	hudLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(dimColor)

	hudValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellowColor)

	successBannerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(greenColor).
				Padding(0, 1)

	arrowCompletedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(yellowColor)

	arrowCurrentStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(whiteColor).
				Underline(true)

	arrowUpcomingStyle = lipgloss.NewStyle().
				Foreground(grayColor)

	arrowErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(redColor)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(yellowColor).
			Padding(0, 1).
			Width(58).
			Align(lipgloss.Center)
)

func DirectionSymbol(dir stratagem.Direction) string {
	switch dir {
	case stratagem.Up:
		return symbolUp
	case stratagem.Down:
		return symbolDown
	case stratagem.Left:
		return symbolLeft
	case stratagem.Right:
		return symbolRight
	default:
		return "?"
	}
}

// RenderSequence draws the arrows as single bold glyphs. When errored, every
// arrow is red and blinks: blinkOn=false renders them as blanks so the
// layout stays stable.
func RenderSequence(seq []stratagem.Direction, inputIndex int, errored, blinkOn bool) string {
	arrows := make([]string, 0, len(seq))
	for i, dir := range seq {
		sym := DirectionSymbol(dir)
		style := arrowUpcomingStyle
		switch {
		case errored:
			style = arrowErrorStyle
		case i < inputIndex:
			style = arrowCompletedStyle
		case i == inputIndex:
			style = arrowCurrentStyle
		}
		if errored && !blinkOn {
			sym = " "
		}
		arrows = append(arrows, style.Render(sym))
	}
	return stringsJoinWithSpaces(arrows, "  ")
}

func stringsJoinWithSpaces(items []string, sep string) string {
	res := ""
	for i, item := range items {
		if i > 0 {
			res += sep
		}
		res += item
	}
	return res
}

func RenderHUD(score, streak int, elapsedMs int64) string {
	scoreStr := fmt.Sprintf("%d", score)
	if score > 999 {
		scoreStr = fmt.Sprintf("%d,%03d", score/1000, score%1000)
	}

	timeStr := fmt.Sprintf("%dms", elapsedMs)

	part1 := hudLabelStyle.Render("SCORE ") + hudValueStyle.Render(scoreStr)
	part2 := hudLabelStyle.Render("STREAK ") + hudValueStyle.Render(fmt.Sprintf("%d", streak))
	part3 := hudLabelStyle.Render("TIME ") + hudValueStyle.Render(timeStr)

	return lipgloss.JoinHorizontal(lipgloss.Center, part1, "   ", part2, "   ", part3)
}

type AnimationState int

const (
	AnimNone AnimationState = iota
	AnimSuccess
	AnimLaunch
	AnimFailure
)

func RenderView(width, height int, stratName string, seq []stratagem.Direction, inputIndex, score, streak int, elapsedMs int64, animState AnimationState, lastCompletedName string, frozen, blinkOn bool, upcoming string) string {
	minHeight := 10
	if upcoming != "" {
		minHeight = 8 + lipgloss.Height(upcoming)
	}
	if width > 0 && height > 0 && (width < 60 || height < minHeight) {
		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(redColor).Render(fmt.Sprintf("Terminal too small.\nResize to at least 60x%d.", minHeight)),
		)
	}

	header := titleStyle.Render("⚡ STRATAGEM ZERO ⚡")

	nameBar := stratagemNameStyle.Render(stratName)

	seqView := RenderSequence(seq, inputIndex, frozen, blinkOn)

	var animBanner string
	switch animState {
	case AnimSuccess:
		animBanner = successBannerStyle.Render("✓ STRATAGEM READY")
	default:
		animBanner = " "
	}

	hudView := RenderHUD(score, streak, elapsedMs)

	rows := []string{header}
	if upcoming != "" {
		rows = append(rows, upcoming)
	}
	rows = append(rows, nameBar, seqView, animBanner, hudView)
	content := lipgloss.JoinVertical(lipgloss.Center, rows...)

	box := boxStyle
	if frozen && blinkOn {
		box = box.BorderForeground(redColor)
	}
	boxed := box.Render(content)

	if width > 0 && height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, boxed)
	}

	return boxed
}
