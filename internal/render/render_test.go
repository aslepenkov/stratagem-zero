package render_test

import (
	"testing"

	"stratagem-zero/internal/render"
	"stratagem-zero/internal/stratagem"
)

func TestDirectionSymbol(t *testing.T) {
	if render.DirectionSymbol(stratagem.Up) != "▲" {
		t.Errorf("expected ▲ for Up")
	}
	if render.DirectionSymbol(stratagem.Down) != "▼" {
		t.Errorf("expected ▼ for Down")
	}
	if render.DirectionSymbol(stratagem.Left) != "◄" {
		t.Errorf("expected ◄ for Left")
	}
	if render.DirectionSymbol(stratagem.Right) != "►" {
		t.Errorf("expected ► for Right")
	}
}

func TestRenderSequence(t *testing.T) {
	seq := []stratagem.Direction{stratagem.Up, stratagem.Down, stratagem.Right}
	output := render.RenderSequence(seq, 1, false, false)
	if output == "" {
		t.Errorf("expected non-empty output")
	}
}

func TestRenderHUD(t *testing.T) {
	output := render.RenderHUD(12450, 7, 842)
	if output == "" {
		t.Errorf("expected non-empty HUD output")
	}
}

func TestRenderView_TerminalTooSmall(t *testing.T) {
	testCases := []struct {
		width, height int
	}{
		{59, 10},
		{60, 9},
	}

	for _, tc := range testCases {
		output := render.RenderView(tc.width, tc.height, "Reinforce", []stratagem.Direction{stratagem.Up}, 0, 0, 0, 0, render.AnimNone, "", false, false)
		if output == "" {
			t.Errorf("expected output")
		}
		if !contains(output, "Terminal too small") {
			t.Errorf("expected output to contain 'Terminal too small' for %dx%d terminal", tc.width, tc.height)
		}
	}
}

func TestRenderView_ValidDimensions(t *testing.T) {
	testCases := []struct {
		width, height int
	}{
		{80, 24},
		{100, 25},
		{70, 20},
		{60, 10},
	}

	for _, tc := range testCases {
		output := render.RenderView(tc.width, tc.height, "Reinforce", []stratagem.Direction{stratagem.Up}, 0, 0, 0, 0, render.AnimNone, "", false, false)
		if contains(output, "Terminal too small") {
			t.Errorf("expected valid rendering without 'Terminal too small' for %dx%d terminal", tc.width, tc.height)
		}
	}
}

func TestRenderView_FrozenBlinksArrows(t *testing.T) {
	seq := []stratagem.Direction{stratagem.Up}
	on := render.RenderView(80, 24, "Reinforce", seq, 0, 0, 0, 0, render.AnimFailure, "", true, true)
	off := render.RenderView(80, 24, "Reinforce", seq, 0, 0, 0, 0, render.AnimFailure, "", true, false)
	if !contains(on, render.DirectionSymbol(stratagem.Up)) {
		t.Errorf("expected arrow visible during blink-on phase")
	}
	if contains(off, render.DirectionSymbol(stratagem.Up)) {
		t.Errorf("expected arrow hidden during blink-off phase")
	}
	if contains(on, "INCORRECT INPUT") || contains(on, "Lockout") {
		t.Errorf("expected no overlay or banner during freeze")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || searchSubstring(s, substr))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
