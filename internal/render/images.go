package render

import (
	"os"
	"strings"
)

type ImageFormat int

const (
	FormatNone ImageFormat = iota
	FormatKitty
	FormatSixel
	FormatASCII
)

type ImageRenderer interface {
	SupportsImages() bool
	Render(imagePath string, width, height int) string
}

func DetectImageSupport(forceASCII bool) ImageFormat {
	if forceASCII {
		return FormatASCII
	}

	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))

	// Kitty graphics detection
	if strings.Contains(term, "kitty") || strings.Contains(termProgram, "kitty") || os.Getenv("KITTY_WINDOW_ID") != "" {
		return FormatKitty
	}

	// Sixel detection heuristic
	if strings.Contains(term, "sixel") || strings.Contains(term, "mlterm") || strings.Contains(term, "xterm") || strings.Contains(termProgram, "wezterm") || strings.Contains(termProgram, "foot") {
		return FormatSixel
	}

	return FormatASCII
}

type ASCIIRenderer struct{}

func (a *ASCIIRenderer) SupportsImages() bool {
	return false
}

func (a *ASCIIRenderer) Render(imagePath string, width, height int) string {
	return ""
}

type NoopRenderer struct{}

func (n *NoopRenderer) SupportsImages() bool {
	return false
}

func (n *NoopRenderer) Render(imagePath string, width, height int) string {
	return ""
}

func GetGenericStratagemASCII(name string) string {
	upper := strings.ToUpper(name)
	if len(upper) > 20 {
		upper = upper[:17] + "..."
	}
	pad := (20 - len(upper)) / 2
	padding := strings.Repeat(" ", pad)

	return "   ┌──────────────────────┐\n" +
		"   │ " + padding + upper + strings.Repeat(" ", 20-len(upper)-pad) + " │\n" +
		"   │          ▲           │\n" +
		"   └──────────────────────┘"
}
