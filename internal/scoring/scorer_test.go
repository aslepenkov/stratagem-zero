package scoring_test

import (
	"testing"
	"time"

	"stratagem-zero/internal/scoring"
)

func TestDefaultScorer(t *testing.T) {
	scorer := scoring.NewDefaultScorer()

	t.Run("basic scoring", func(t *testing.T) {
		// base = 5 * 100 = 500
		// speed_bonus = max(0, 1000 - 200 * 2) = 600
		// sum = 1100
		// streak_multiplier = 1.0 + 0 * 0.1 = 1.0
		// score = 1100
		score := scorer.Score(5, 200*time.Millisecond, 0)
		if score != 1100 {
			t.Errorf("expected 1100, got %d", score)
		}
	})

	t.Run("faster completion produces higher score", func(t *testing.T) {
		fastScore := scorer.Score(5, 100*time.Millisecond, 0)
		slowScore := scorer.Score(5, 400*time.Millisecond, 0)

		if fastScore <= slowScore {
			t.Errorf("expected fastScore (%d) > slowScore (%d)", fastScore, slowScore)
		}
	})

	t.Run("longer sequence produces higher base score", func(t *testing.T) {
		longSeqScore := scorer.Score(8, 200*time.Millisecond, 0)
		shortSeqScore := scorer.Score(4, 200*time.Millisecond, 0)

		if longSeqScore <= shortSeqScore {
			t.Errorf("expected longSeqScore (%d) > shortSeqScore (%d)", longSeqScore, shortSeqScore)
		}
	})

	t.Run("streak increases multiplier", func(t *testing.T) {
		noStreak := scorer.Score(5, 200*time.Millisecond, 0)
		withStreak := scorer.Score(5, 200*time.Millisecond, 5)

		if withStreak <= noStreak {
			t.Errorf("expected withStreak (%d) > noStreak (%d)", withStreak, noStreak)
		}
	})

	t.Run("streak multiplier caps at 10", func(t *testing.T) {
		streak10 := scorer.Score(5, 200*time.Millisecond, 10)
		streak15 := scorer.Score(5, 200*time.Millisecond, 15)

		if streak10 != streak15 {
			t.Errorf("expected streak10 (%d) == streak15 (%d)", streak10, streak15)
		}
	})

	t.Run("score never negative", func(t *testing.T) {
		score := scorer.Score(0, 10000*time.Millisecond, 0)
		if score < 0 {
			t.Errorf("expected score >= 0, got %d", score)
		}
	})
}
