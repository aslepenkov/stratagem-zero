package input_test

import (
	"testing"

	"stratagem-zero/internal/input"
	"stratagem-zero/internal/stratagem"
)

func TestParseKey(t *testing.T) {
	tests := []struct {
		key     string
		want    stratagem.Direction
		wantErr bool
	}{
		{"w", stratagem.Up, false},
		{"W", stratagem.Up, false},
		{"up", stratagem.Up, false},
		{"s", stratagem.Down, false},
		{"S", stratagem.Down, false},
		{"down", stratagem.Down, false},
		{"a", stratagem.Left, false},
		{"A", stratagem.Left, false},
		{"left", stratagem.Left, false},
		{"d", stratagem.Right, false},
		{"D", stratagem.Right, false},
		{"right", stratagem.Right, false},
		{"x", stratagem.Up, true},
		{"space", stratagem.Up, true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, err := input.ParseKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseKey(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
