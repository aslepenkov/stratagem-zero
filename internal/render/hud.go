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
			MarginBottom(1)

	stratagemNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(whiteColor).
				Background(lipgloss.Color("#222222")).
				Padding(0, 2).
				MarginBottom(1)

	arrowCompletedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(yellowColor)

	arrowCurrentStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(whiteColor).
				Underline(true)

	arrowUpcomingStyle = lipgloss.NewStyle().
				Foreground(grayColor)

	checkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(greenColor)

	dotStyle = lipgloss.NewStyle().
			Foreground(grayColor)

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

	failBannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(redColor).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(yellowColor).
			Padding(1, 3).
			Width(64).
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

func RenderSequence(seq []stratagem.Direction, inputIndex int) string {
	var topRow []string
	var bottomRow []string

	for i, dir := range seq {
		sym := DirectionSymbol(dir)
		if i < inputIndex {
			topRow = append(topRow, arrowCompletedStyle.Render(sym))
			bottomRow = append(bottomRow, checkStyle.Render("✓"))
		} else if i == inputIndex {
			topRow = append(topRow, arrowCurrentStyle.Render(sym))
			bottomRow = append(bottomRow, dotStyle.Render("·"))
		} else {
			topRow = append(topRow, arrowUpcomingStyle.Render(sym))
			bottomRow = append(bottomRow, dotStyle.Render("·"))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		stringsJoinWithSpaces(topRow, "  "),
		stringsJoinWithSpaces(bottomRow, "  "),
	)
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

	return lipgloss.JoinHorizontal(lipgloss.Center, part1, "     ", part2, "     ", part3)
}

type AnimationState int

const (
	AnimNone AnimationState = iota
	AnimSuccess
	AnimLaunch
	AnimFailure
)

func RenderView(width, height int, stratName string, seq []stratagem.Direction, inputIndex, score, streak int, elapsedMs int64, animState AnimationState, lastCompletedName string) string {
	if width > 0 && height > 0 && (width < 80 || height < 24) {
		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(redColor).Render("Terminal too small.\nResize to at least 80x24."),
		)
	}

	header := titleStyle.Render("⚡ STRATAGEM ZERO ⚡")

	asciiArt := GetGenericStratagemASCII(stratName)

	nameBar := stratagemNameStyle.Render(stratName)

	seqView := RenderSequence(seq, inputIndex)

	var animBanner string
	switch animState {
	case AnimSuccess:
		animBanner = successBannerStyle.Render("✓ STRATAGEM READY")
	case AnimLaunch:
		animBanner = successBannerStyle.Render("🚀 LAUNCHING " + lastCompletedName)
	case AnimFailure:
		animBanner = failBannerStyle.Render("✗ INCORRECT INPUT")
	default:
		animBanner = " "
	}

	hudView := RenderHUD(score, streak, elapsedMs)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		asciiArt,
		"",
		nameBar,
		"",
		seqView,
		"",
		animBanner,
		"",
		hudView,
	)

	boxed := boxStyle.Render(content)

	if width > 0 && height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, boxed)
	}

	return boxed
}
