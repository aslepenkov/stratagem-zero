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
				Background(lipgloss.Color("#222222")).
				Padding(0, 1).
				MarginBottom(0)

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

// DirectionBigSymbol returns a 2-line representation of the symbol for x2 size
func DirectionBigSymbol(dir stratagem.Direction) (string, string) {
	sym := DirectionSymbol(dir)
	// Build a 2x2 multi-line string for 2x font size visual appearance
	return sym + sym, sym + sym
}

func RenderSequence(seq []stratagem.Direction, inputIndex int) string {
	var arrowRow1 []string
	var arrowRow2 []string
	var statusRow []string

	for i, dir := range seq {
		line1, line2 := DirectionBigSymbol(dir)
		if i < inputIndex {
			arrowRow1 = append(arrowRow1, arrowCompletedStyle.Render(line1))
			arrowRow2 = append(arrowRow2, arrowCompletedStyle.Render(line2))
			statusRow = append(statusRow, checkStyle.Render("✓ "))
		} else if i == inputIndex {
			arrowRow1 = append(arrowRow1, arrowCurrentStyle.Render(line1))
			arrowRow2 = append(arrowRow2, arrowCurrentStyle.Render(line2))
			statusRow = append(statusRow, dotStyle.Render("· "))
		} else {
			arrowRow1 = append(arrowRow1, arrowUpcomingStyle.Render(line1))
			arrowRow2 = append(arrowRow2, arrowUpcomingStyle.Render(line2))
			statusRow = append(statusRow, dotStyle.Render("· "))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		stringsJoinWithSpaces(arrowRow1, " "),
		stringsJoinWithSpaces(arrowRow2, " "),
		stringsJoinWithSpaces(statusRow, " "),
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

	return lipgloss.JoinHorizontal(lipgloss.Center, part1, "   ", part2, "   ", part3)
}

type AnimationState int

const (
	AnimNone AnimationState = iota
	AnimSuccess
	AnimLaunch
	AnimFailure
)

func RenderView(width, height int, stratName string, seq []stratagem.Direction, inputIndex, score, streak int, elapsedMs int64, animState AnimationState, lastCompletedName string, freezeSeconds int) string {
	if width > 0 && height > 0 && (width < 60 || height < 10) {
		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(redColor).Render("Terminal too small.\nResize to at least 60x10."),
		)
	}

	header := titleStyle.Render("⚡ STRATAGEM ZERO ⚡")

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
		nameBar,
		seqView,
		animBanner,
		hudView,
	)

	boxed := boxStyle.Render(content)

	if freezeSeconds > 0 || animState == AnimFailure {
		sec := freezeSeconds
		if sec <= 0 {
			sec = 3
		}
		overlayMsg := fmt.Sprintf(" Lockout: %ds ", sec)
		overlayBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(redColor).
			Background(lipgloss.Color("#1A0000")).
			Foreground(whiteColor).
			Bold(true).
			Padding(0, 1).
			Render(failBannerStyle.Render("✗ INCORRECT INPUT") + "\n" + overlayMsg)

		boxed = lipgloss.Place(
			58, 10,
			lipgloss.Center, lipgloss.Center,
			overlayBox,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(grayColor),
		)
	}

	if width > 0 && height > 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, boxed)
	}

	return boxed
}
