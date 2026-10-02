package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"stratagem-zero/data"
	"stratagem-zero/internal/game"
	"stratagem-zero/internal/input"
	"stratagem-zero/internal/render"
	"stratagem-zero/internal/scoring"
	"stratagem-zero/internal/stratagem"
)

// blinkFrames is how many 33ms ticks each blink phase lasts (~264ms).
const blinkFrames = 8

type tickMsg time.Time
type animTickMsg struct{}
type freezeTickMsg struct{}

type model struct {
	engine            *game.Engine
	width             int
	height            int
	animState         render.AnimationState
	lastCompletedName string
	forceASCII        bool
	freezeSeconds     int
	frame             int
}

func tick() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func animTick(delay time.Duration) tea.Cmd {
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return animTickMsg{}
	})
}

func freezeTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return freezeTickMsg{}
	})
}

func initialModel(engine *game.Engine, forceASCII bool) model {
	return model{
		engine:     engine,
		forceASCII: forceASCII,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.frame++
		return m, tick()

	case animTickMsg:
		if m.animState == render.AnimSuccess {
			m.animState = render.AnimLaunch
			return m, animTick(250 * time.Millisecond)
		} else if m.animState == render.AnimLaunch {
			m.animState = render.AnimNone
		}
		return m, nil

	case freezeTickMsg:
		if m.freezeSeconds > 0 {
			m.freezeSeconds--
			if m.freezeSeconds > 0 {
				return m, freezeTick()
			}
			m.animState = render.AnimNone
		}
		return m, nil

	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

		if m.freezeSeconds > 0 {
			return m, nil
		}

		dir, err := input.ParseKey(k)
		if err != nil {
			return m, nil
		}

		event := m.engine.HandleInput(dir)
		switch event.Type {
		case game.EventRoundSuccess:
			m.animState = render.AnimSuccess
			m.lastCompletedName = event.Stratagem.Name
			return m, animTick(200 * time.Millisecond)

		case game.EventWrongInput:
			m.animState = render.AnimFailure
			m.freezeSeconds = 2
			return m, freezeTick()
		}
	}

	return m, nil
}

func (m model) View() string {
	state := m.engine.State()
	elapsed := time.Since(state.RoundStartTime).Milliseconds()

	return render.RenderView(
		m.width, m.height,
		state.CurrentStratagem.Name,
		state.CurrentStratagem.Sequence,
		state.InputIndex,
		state.Score,
		state.Streak,
		elapsed,
		m.animState,
		m.lastCompletedName,
		m.freezeSeconds > 0,
		(m.frame/blinkFrames)%2 == 0,
	)
}

func main() {
	asciiFlag := flag.Bool("ascii", false, "Force ASCII fallback graphics")
	seedFlag := flag.Int64("seed", 0, "Seed for random stratagem selection")

	flag.Parse()

	var seed int64
	if *seedFlag != 0 {
		seed = *seedFlag
	} else {
		seed = time.Now().UnixNano()
	}

	pool, err := stratagem.LoadStratagems(data.StratagemJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load stratagem data: %v\n", err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(seed))
	selector, err := stratagem.NewRandomSelector(pool, rng)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create selector: %v\n", err)
		os.Exit(1)
	}

	scorer := scoring.NewDefaultScorer()
	engine := game.NewEngine(selector, scorer, time.Now)

	p := tea.NewProgram(
		initialModel(engine, *asciiFlag),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running game: %v\n", err)
		os.Exit(1)
	}
}
