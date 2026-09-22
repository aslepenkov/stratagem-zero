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
	output := render.RenderView(60, 20, "Reinforce", []stratagem.Direction{stratagem.Up}, 0, 0, 0, 0, render.AnimNone, "")
	if !testing.Short() && output == "" {
		t.Errorf("expected output")
	}
}
