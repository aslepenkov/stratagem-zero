// Package icons rasterizes stratagem SVG icons into small truecolor
// half-block ("▀") terminal art.
package icons

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var (
	styleBlockRe = regexp.MustCompile(`(?s)<style[^>]*>(.*?)</style>`)
	cssRuleRe    = regexp.MustCompile(`\.([A-Za-z0-9_-]+)\s*\{([^}]*)\}`)
	classAttrRe  = regexp.MustCompile(`class="([^"]*)"`)
	clipStyleRe  = regexp.MustCompile(`clip-path\s*:\s*url\([^)]*\)\s*;?`)
	rootSizeRe   = regexp.MustCompile(`(<svg[^>]*?)\s(?:width|height)="[^"]*"`)
	clipElemRe   = regexp.MustCompile(`(?s)<clipPath\b.*?</clipPath>`)
	transformRe  = regexp.MustCompile(`transform="[^"]*"`)
	scaleOneRe   = regexp.MustCompile(`scale\(\s*([-+0-9.eE]+)\s*\)`)
	closeParenRe = regexp.MustCompile(`\)([A-Za-z])`)
	clipAttrRe   = regexp.MustCompile(`\sclip-path\s*=\s*"url\([^)]*\)"`)
)

// inlineCSS rewrites class="x y" attributes into style="..." using the rules
// from the SVG's <style> block, because oksvg does not support CSS classes.
func inlineCSS(svg []byte) []byte {
	// The clip paths in these icons only clip to the canvas, and oksvg hides
	// any group that references one, so drop them.
	// Unit-suffixed root sizes (e.g. "2.66667in") confuse oksvg; the viewBox is enough.
	for rootSizeRe.Match(svg) {
		svg = rootSizeRe.ReplaceAll(svg, []byte("$1"))
	}
	svg = clipStyleRe.ReplaceAll(svg, nil)
	// oksvg needs whitespace between chained transforms: "translate(1 2)scale(3)".
	svg = transformRe.ReplaceAllFunc(svg, func(t []byte) []byte {
		t = closeParenRe.ReplaceAll(t, []byte(") $1"))
		// A single-argument scale(s) leaves the y scale at 0 in oksvg.
		return scaleOneRe.ReplaceAll(t, []byte("scale($1 $1)"))
	})
	svg = clipAttrRe.ReplaceAll(svg, nil)
	svg = clipElemRe.ReplaceAll(svg, nil)
	rules := map[string]string{}
	for _, block := range styleBlockRe.FindAllSubmatch(svg, -1) {
		text := strings.NewReplacer("<![CDATA[", "", "]]>", "").Replace(string(block[1]))
		for _, m := range cssRuleRe.FindAllStringSubmatch(text, -1) {
			rules[m[1]] += strings.TrimSpace(m[2]) + ";"
		}
	}
	svg = styleBlockRe.ReplaceAll(svg, nil)
	return classAttrRe.ReplaceAllFunc(svg, func(attr []byte) []byte {
		classes := classAttrRe.FindSubmatch(attr)[1]
		var style strings.Builder
		for _, c := range strings.Fields(string(classes)) {
			style.WriteString(rules[c])
		}
		return []byte(`style="` + style.String() + `"`)
	})
}

// Rasterize draws the SVG into a size x size RGBA image.
func Rasterize(svg []byte, size int) (*image.RGBA, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(inlineCSS(svg)), oksvg.IgnoreErrorMode)
	if err != nil {
		return nil, fmt.Errorf("parse svg: %w", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)
	icon.SetTarget(0, 0, float64(size), float64(size))
	icon.Draw(rasterx.NewDasher(size, size, rasterx.NewScannerGV(size, size, img, img.Bounds())), 1)
	return img, nil
}

// HalfBlocks renders a square image as len(rows)=size/2 lines of "▀" cells,
// each cell encoding two vertically stacked pixels.
func HalfBlocks(img *image.RGBA) []string {
	b := img.Bounds()
	var lines []string
	for y := b.Min.Y; y+1 < b.Max.Y; y += 2 {
		var sb strings.Builder
		for x := b.Min.X; x < b.Max.X; x++ {
			top := img.RGBAAt(x, y)
			bot := img.RGBAAt(x, y+1)
			fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", top.R, top.G, top.B, bot.R, bot.G, bot.B)
		}
		sb.WriteString("\x1b[0m")
		lines = append(lines, sb.String())
	}
	return lines
}

// Set holds pre-rendered icons keyed by SVG file name.
type Set struct {
	art map[string]string
}

// NewSet rasterizes every SVG in fsys/dir to a size x size/2 cell block.
// It does all the work up front so rendering a frame is a map lookup, and it
// silences oksvg's logging, which would otherwise corrupt the TUI.
func NewSet(fsys fs.FS, dir string, size int) (*Set, error) {
	prev := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(prev)

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read icon dir: %w", err)
	}
	set := &Set{art: make(map[string]string, len(entries))}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".svg") {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}
		img, err := Rasterize(data, size)
		if err != nil {
			continue // a bad icon just won't be shown
		}
		set.art[e.Name()] = strings.Join(HalfBlocks(img), "\n")
	}
	return set, nil
}

var highlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFE81F"))

// Strip lays the icons for the given icon URLs/paths side by side, separated
// by one column, and marks the first one with yellow corner brackets. The
// brackets sit in a one-cell margin around the strip, so the result is two
// columns wider and two rows taller than the icons alone. Stratagems without
// a known icon are left as blank space.
func (s *Set) Strip(icons []string) string {
	if s == nil || len(icons) == 0 {
		return ""
	}
	var blank string
	for _, a := range s.art {
		w, h := lipgloss.Width(strings.SplitN(a, "\n", 2)[0]), strings.Count(a, "\n")+1
		blank = strings.TrimRight(strings.Repeat(strings.Repeat(" ", w)+"\n", h), "\n")
		break
	}
	iconW := lipgloss.Width(strings.SplitN(blank, "\n", 2)[0])

	parts := make([]string, 0, 2*len(icons))
	for i, icon := range icons {
		if i > 0 {
			parts = append(parts, " ")
		}
		art, ok := s.art[path.Base(icon)]
		if !ok {
			art = blank
		}
		parts = append(parts, art)
	}
	body := lipgloss.JoinHorizontal(lipgloss.Top, parts...)

	// Corners hug the first icon: one cell left of it and one cell right of it
	// (the gap column before the second icon).
	total := max(lipgloss.Width(strings.SplitN(body, "\n", 2)[0]), iconW) + 2
	corners := func(left, right string) string {
		return highlightStyle.Render(left) + strings.Repeat(" ", iconW) + highlightStyle.Render(right) +
			strings.Repeat(" ", total-iconW-2)
	}
	lines := []string{corners("┏", "┓")}
	for _, l := range strings.Split(body, "\n") {
		lines = append(lines, " "+l+" ")
	}
	lines = append(lines, corners("┗", "┛"))
	return strings.Join(lines, "\n")
}
