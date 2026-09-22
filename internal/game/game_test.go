package game_test

import (
	"testing"
	"time"

	"stratagem-zero/internal/game"
	"stratagem-zero/internal/scoring"
	"stratagem-zero/internal/stratagem"
)

type mockSelector struct {
	stratagems []stratagem.Stratagem
	index      int
}

func (m *mockSelector) Next() stratagem.Stratagem {
	s := m.stratagems[m.index%len(m.stratagems)]
	m.index++
	return s
}

func TestGameEngine_SequenceProgression(t *testing.T) {
	s1 := stratagem.Stratagem{
		Name:     "TestStratagem",
		Sequence: []stratagem.Direction{stratagem.Up, stratagem.Down, stratagem.Left},
	}
	s2 := stratagem.Stratagem{
		Name:     "NextStratagem",
		Sequence: []stratagem.Direction{stratagem.Right, stratagem.Up},
	}

	selector := &mockSelector{stratagems: []stratagem.Stratagem{s1, s2}}
	scorer := scoring.NewDefaultScorer()

	currentTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	tp := func() time.Time {
		return currentTime
	}

	engine := game.NewEngine(selector, scorer, tp)

	state := engine.State()
	if state.CurrentStratagem.Name != "TestStratagem" {
		t.Fatalf("expected initial stratagem TestStratagem, got %s", state.CurrentStratagem.Name)
	}
	if state.InputIndex != 0 {
		t.Fatalf("expected initial InputIndex 0, got %d", state.InputIndex)
	}

	// 1. Correct first input (Up)
	event := engine.HandleInput(stratagem.Up)
	if event.Type != game.EventCorrectInput {
		t.Errorf("expected EventCorrectInput, got %v", event.Type)
	}
	if engine.State().InputIndex != 1 {
		t.Errorf("expected InputIndex 1, got %d", engine.State().InputIndex)
	}

	// 2. Correct second input (Down)
	event = engine.HandleInput(stratagem.Down)
	if event.Type != game.EventCorrectInput {
		t.Errorf("expected EventCorrectInput, got %v", event.Type)
	}
	if engine.State().InputIndex != 2 {
		t.Errorf("expected InputIndex 2, got %d", engine.State().InputIndex)
	}

	// Advance mock time by 500ms
	currentTime = currentTime.Add(500 * time.Millisecond)

	// 3. Correct final input (Left) -> Round Success
	event = engine.HandleInput(stratagem.Left)
	if event.Type != game.EventRoundSuccess {
		t.Errorf("expected EventRoundSuccess, got %v", event.Type)
	}
	if engine.State().Score <= 0 {
		t.Errorf("expected Score > 0, got %d", engine.State().Score)
	}
	if engine.State().Streak != 1 {
		t.Errorf("expected Streak 1, got %d", engine.State().Streak)
	}
	if engine.State().CurrentStratagem.Name != "NextStratagem" {
		t.Errorf("expected new stratagem NextStratagem, got %s", engine.State().CurrentStratagem.Name)
	}
	if engine.State().InputIndex != 0 {
		t.Errorf("expected reset InputIndex 0, got %d", engine.State().InputIndex)
	}
}

func TestGameEngine_WrongInputFailsRound(t *testing.T) {
	s1 := stratagem.Stratagem{
		Name:     "TestStratagem",
		Sequence: []stratagem.Direction{stratagem.Up, stratagem.Down},
	}
	s2 := stratagem.Stratagem{
		Name:     "NextStratagem",
		Sequence: []stratagem.Direction{stratagem.Right},
	}

	selector := &mockSelector{stratagems: []stratagem.Stratagem{s1, s2}}
	scorer := scoring.NewDefaultScorer()
	engine := game.NewEngine(selector, scorer, time.Now)

	// Advance 1 key correctly
	_ = engine.HandleInput(stratagem.Up)

	// Wrong key on second step
	event := engine.HandleInput(stratagem.Left)
	if event.Type != game.EventWrongInput {
		t.Errorf("expected EventWrongInput, got %v", event.Type)
	}

	state := engine.State()
	if state.Streak != 0 {
		t.Errorf("expected streak reset to 0, got %d", state.Streak)
	}
	if state.FailedRounds != 1 {
		t.Errorf("expected 1 failed round, got %d", state.FailedRounds)
	}
	if state.CurrentStratagem.Name != "NextStratagem" {
		t.Errorf("expected new stratagem NextStratagem, got %s", state.CurrentStratagem.Name)
	}
	if state.InputIndex != 0 {
		t.Errorf("expected InputIndex 0 for new round, got %d", state.InputIndex)
	}
}
