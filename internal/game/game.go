package game

import (
	"time"

	"stratagem-zero/internal/stratagem"
)

type EventType int

const (
	EventCorrectInput EventType = iota
	EventWrongInput
	EventRoundSuccess
	EventRoundStarted
)

type Event struct {
	Type            EventType
	Stratagem       stratagem.Stratagem
	InputIndex      int
	RoundScore      int
	ElapsedDuration time.Duration
}

type GameState struct {
	CurrentStratagem stratagem.Stratagem
	InputIndex       int

	Score  int
	Streak int

	StartedAt     time.Time
	LastSuccessAt time.Time

	TotalAttempts    int
	SuccessfulRounds int
	FailedRounds     int

	RoundStartTime time.Time
}

type TimeProvider func() time.Time

type Engine struct {
	state        GameState
	selector     stratagem.StratagemSelector
	scorer       Scorer
	timeProvider TimeProvider
}

type Scorer interface {
	Score(sequenceLength int, elapsed time.Duration, streak int) int
}

func NewEngine(selector stratagem.StratagemSelector, scorer Scorer, tp TimeProvider) *Engine {
	if tp == nil {
		tp = time.Now
	}
	engine := &Engine{
		selector:     selector,
		scorer:       scorer,
		timeProvider: tp,
	}
	engine.startNewRound()
	return engine
}

func (e *Engine) State() GameState {
	return e.state
}

func (e *Engine) startNewRound() Event {
	e.state.CurrentStratagem = e.selector.Next()
	e.state.InputIndex = 0
	now := e.timeProvider()
	if e.state.StartedAt.IsZero() {
		e.state.StartedAt = now
	}
	e.state.RoundStartTime = now

	return Event{
		Type:      EventRoundStarted,
		Stratagem: e.state.CurrentStratagem,
	}
}

func (e *Engine) HandleInput(dir stratagem.Direction) Event {
	e.state.TotalAttempts++
	expected := e.state.CurrentStratagem.Sequence[e.state.InputIndex]

	if dir == expected {
		e.state.InputIndex++
		if e.state.InputIndex == len(e.state.CurrentStratagem.Sequence) {
			// Round completed successfully
			now := e.timeProvider()
			elapsed := now.Sub(e.state.RoundStartTime)
			roundScore := e.scorer.Score(len(e.state.CurrentStratagem.Sequence), elapsed, e.state.Streak)

			e.state.Score += roundScore
			e.state.Streak++
			e.state.SuccessfulRounds++
			e.state.LastSuccessAt = now

			completedStratagem := e.state.CurrentStratagem
			e.startNewRound()

			return Event{
				Type:            EventRoundSuccess,
				Stratagem:       completedStratagem,
				RoundScore:      roundScore,
				ElapsedDuration: elapsed,
			}
		}

		return Event{
			Type:       EventCorrectInput,
			Stratagem:  e.state.CurrentStratagem,
			InputIndex: e.state.InputIndex,
		}
	}

	// Wrong input
	e.state.FailedRounds++
	e.state.Streak = 0
	failedStratagem := e.state.CurrentStratagem
	e.startNewRound()

	return Event{
		Type:      EventWrongInput,
		Stratagem: failedStratagem,
	}
}
