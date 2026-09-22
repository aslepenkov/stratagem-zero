package stratagem

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type stratagemDTO struct {
	Name     string   `json:"name"`
	Sequence []string `json:"sequence"`
	Category string   `json:"category,omitempty"`
}

// LoadStratagems parses raw JSON data into a slice of Stratagems with validation.
func LoadStratagems(data []byte) ([]Stratagem, error) {
	var dtos []stratagemDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("failed to parse stratagem JSON: %w", err)
	}

	if len(dtos) == 0 {
		return nil, fmt.Errorf("stratagem list is empty")
	}

	seenNames := make(map[string]bool)
	var stratagems []Stratagem

	for i, dto := range dtos {
		if dto.Name == "" {
			return nil, fmt.Errorf("stratagem at index %d has empty name", i)
		}
		if seenNames[dto.Name] {
			return nil, fmt.Errorf("duplicate stratagem name found: %q", dto.Name)
		}
		seenNames[dto.Name] = true

		if len(dto.Sequence) == 0 {
			return nil, fmt.Errorf("stratagem %q has empty sequence", dto.Name)
		}

		seq := make([]Direction, 0, len(dto.Sequence))
		for _, dirStr := range dto.Sequence {
			dir, err := DirectionFromString(dirStr)
			if err != nil {
				return nil, fmt.Errorf("invalid sequence for %q: %w", dto.Name, err)
			}
			seq = append(seq, dir)
		}

		stratagems = append(stratagems, Stratagem{
			Name:     dto.Name,
			Sequence: seq,
			Category: dto.Category,
		})
	}

	return stratagems, nil
}

// LoadStratagemsFromFile loads stratagems from a file path.
func LoadStratagemsFromFile(filepath string) ([]Stratagem, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stratagems file %q: %w", filepath, err)
	}
	return LoadStratagems(data)
}

// StratagemSelector interface for getting random stratagems.
type StratagemSelector interface {
	Next() Stratagem
}

// RandomSelector selects stratagems uniformly at random.
type RandomSelector struct {
	pool []Stratagem
	rng  *rand.Rand
}

func NewRandomSelector(pool []Stratagem, rng *rand.Rand) (*RandomSelector, error) {
	if len(pool) == 0 {
		return nil, fmt.Errorf("cannot create selector with empty pool")
	}
	if rng == nil {
		return nil, fmt.Errorf("rng source must not be nil")
	}
	return &RandomSelector{
		pool: pool,
		rng:  rng,
	}, nil
}

func (s *RandomSelector) Next() Stratagem {
	idx := s.rng.Intn(len(s.pool))
	return s.pool[idx]
}
