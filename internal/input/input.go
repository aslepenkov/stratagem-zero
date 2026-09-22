package input

import (
	"fmt"
	"strings"

	"stratagem-zero/internal/stratagem"
)

// ParseKey maps keypress string representations to stratagem.Direction.
func ParseKey(key string) (stratagem.Direction, error) {
	keyLower := strings.ToLower(key)
	switch keyLower {
	case "w", "up":
		return stratagem.Up, nil
	case "s", "down":
		return stratagem.Down, nil
	case "a", "left":
		return stratagem.Left, nil
	case "d", "right":
		return stratagem.Right, nil
	default:
		return stratagem.Up, fmt.Errorf("unsupported key: %q", key)
	}
}
