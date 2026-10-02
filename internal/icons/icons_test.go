package icons_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"stratagem-zero/assets"
	"stratagem-zero/data"
	"stratagem-zero/internal/icons"
	"stratagem-zero/internal/stratagem"
)

func TestEveryStratagemIconRenders(t *testing.T) {
	pool, err := stratagem.LoadStratagems(data.StratagemJSON)
	if err != nil {
		t.Fatal(err)
	}
	set, err := icons.NewSet(assets.IconsFS, "icons", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range pool {
		if s.Icon == "" {
			t.Errorf("%s has no icon", s.Name)
			continue
		}
		strip := set.Strip([]string{s.Icon})
		if strip == "" || !strings.Contains(strip, "▀") {
			t.Errorf("%s: icon %s was not rendered", s.Name, s.Icon)
		}
	}
}

func TestStripLayout(t *testing.T) {
	set, err := icons.NewSet(assets.IconsFS, "icons", 10)
	if err != nil {
		t.Fatal(err)
	}
	strip := set.Strip([]string{
		"https://x/Meltagun_Stratagem_Icon_Background.svg",
		"https://x/Reinforce_Stratagem_Icon_Background.svg",
		"https://x/unknown.svg",
	})
	if w := lipgloss.Width(strip); w != 3*10+2+2 {
		t.Errorf("width = %d, want 34", w)
	}
	if h := lipgloss.Height(strip); h != 5+2 {
		t.Errorf("height = %d, want 7", h)
	}
}

func TestNilSetIsEmpty(t *testing.T) {
	var set *icons.Set
	if set.Strip([]string{"a.svg"}) != "" {
		t.Error("nil set should render nothing")
	}
}
