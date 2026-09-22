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
	output := render.RenderSequence(seq, 1)
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
	output := render.RenderView(60, 15, "Reinforce", []stratagem.Direction{stratagem.Up}, 0, 0, 0, 0, render.AnimNone, "")
	if output == "" {
		t.Errorf("expected output")
	}
	if !contains(output, "Terminal too small") {
		t.Errorf("expected output to contain 'Terminal too small' for 60x15 terminal")
	}
}

func TestRenderView_ValidDimensions(t *testing.T) {
	testCases := []struct {
		width, height int
	}{
		{80, 24},
		{100, 25},
		{66, 20},
	}

	for _, tc := range testCases {
		output := render.RenderView(tc.width, tc.height, "Reinforce", []stratagem.Direction{stratagem.Up}, 0, 0, 0, 0, render.AnimNone, "")
		if contains(output, "Terminal too small") {
			t.Errorf("expected valid rendering without 'Terminal too small' for %dx%d terminal", tc.width, tc.height)
		}
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
