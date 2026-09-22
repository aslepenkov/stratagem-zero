package scoring

import (
	"time"
)

// Scorer calculates the score for completing a stratagem sequence.
type Scorer interface {
	Score(sequenceLength int, elapsed time.Duration, streak int) int
}

// DefaultScorer implements the standard scoring logic.
type DefaultScorer struct{}

func NewDefaultScorer() *DefaultScorer {
	return &DefaultScorer{}
}

func (s *DefaultScorer) Score(sequenceLength int, elapsed time.Duration, streak int) int {
	base := sequenceLength * 100

	elapsedMs := int(elapsed.Milliseconds())
	speedBonus := 1000 - elapsedMs*2
	if speedBonus < 0 {
		speedBonus = 0
	}

	effectiveStreak := streak
	if effectiveStreak > 10 {
		effectiveStreak = 10
	}
	if effectiveStreak < 0 {
		effectiveStreak = 0
	}

	streakMultiplier := 1.0 + float64(effectiveStreak)*0.1

	roundScore := int(float64(base+speedBonus) * streakMultiplier)
	if roundScore < 0 {
		roundScore = 0
	}

	return roundScore
}
