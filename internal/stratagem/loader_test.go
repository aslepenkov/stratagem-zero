package stratagem_test

import (
	"math/rand"
	"testing"

	"stratagem-zero/internal/stratagem"
)

func TestLoadStratagems_Valid(t *testing.T) {
	jsonContent := `[
		{"name": "Reinforce", "sequence": ["UP", "DOWN", "RIGHT", "LEFT", "UP"], "category": "Mission"},
		{"name": "Resupply", "sequence": ["DOWN", "DOWN", "UP", "RIGHT"], "category": "Mission"}
	]`

	list, err := stratagem.LoadStratagems([]byte(jsonContent))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 stratagems, got %d", len(list))
	}

	if list[0].Name != "Reinforce" || len(list[0].Sequence) != 5 {
		t.Errorf("unexpected content in list[0]: %+v", list[0])
	}
}

func TestLoadStratagems_Errors(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "empty array",
			json:    `[]`,
			wantErr: true,
		},
		{
			name:    "empty name",
			json:    `[{"name": "", "sequence": ["UP"]}]`,
			wantErr: true,
		},
		{
			name:    "empty sequence",
			json:    `[{"name": "Test", "sequence": []}]`,
			wantErr: true,
		},
		{
			name:    "invalid direction",
			json:    `[{"name": "Test", "sequence": ["FORWARD"]}]`,
			wantErr: true,
		},
		{
			name:    "duplicate name",
			json:    `[{"name": "Test", "sequence": ["UP"]}, {"name": "Test", "sequence": ["DOWN"]}]`,
			wantErr: true,
		},
		{
			name:    "malformed json",
			json:    `[{name: test}]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := stratagem.LoadStratagems([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadStratagems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRandomSelector_Deterministic(t *testing.T) {
	pool := []stratagem.Stratagem{
		{Name: "S1", Sequence: []stratagem.Direction{stratagem.Up}},
		{Name: "S2", Sequence: []stratagem.Direction{stratagem.Down}},
		{Name: "S3", Sequence: []stratagem.Direction{stratagem.Left}},
	}

	rng1 := rand.New(rand.NewSource(42))
	selector1, err := stratagem.NewRandomSelector(pool, rng1)
	if err != nil {
		t.Fatalf("failed to create selector1: %v", err)
	}

	rng2 := rand.New(rand.NewSource(42))
	selector2, err := stratagem.NewRandomSelector(pool, rng2)
	if err != nil {
		t.Fatalf("failed to create selector2: %v", err)
	}

	for i := 0; i < 10; i++ {
		s1 := selector1.Next()
		s2 := selector2.Next()
		if s1.Name != s2.Name {
			t.Fatalf("iteration %d: expected %s, got %s", i, s1.Name, s2.Name)
		}
	}
}
