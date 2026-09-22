package stratagem

import "fmt"

// Direction represents a normalized arrow direction.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "UP"
	case Down:
		return "DOWN"
	case Left:
		return "LEFT"
	case Right:
		return "RIGHT"
	default:
		return "UNKNOWN"
	}
}

func DirectionFromString(s string) (Direction, error) {
	switch s {
	case "UP":
		return Up, nil
	case "DOWN":
		return Down, nil
	case "LEFT":
		return Left, nil
	case "RIGHT":
		return Right, nil
	default:
		return Up, fmt.Errorf("invalid direction string: %q", s)
	}
}

// Stratagem represents a stratagem with its name, directional sequence, and optional category.
type Stratagem struct {
	Name     string      `json:"name"`
	Sequence []Direction `json:"sequence"`
	Category string      `json:"category,omitempty"`
}
