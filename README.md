# Stratagem Zero

<img width="790" height="343" alt="image" src="https://github.com/user-attachments/assets/c122ecb0-c672-4d2d-8516-3d93f1a3abbb" />


A fast, responsive, arcade-like terminal game inspired by the stratagem input mechanic from **Helldivers 2**.

## Features

- **Fast & Responsive:** Instant key-handling using WASD or Arrow keys.
- **105 Stratagems:** Complete current stratagem pool embedded in the binary.
- **Continuous Gameplay:** Immediate transitions on success/failure without blocking menus or result screens.
- **Scoring & Streaks:** Combo multipliers and speed bonuses rewarded for rapid inputs.
- **Rich Visuals:** Helldivers-inspired TUI built with Bubble Tea and Lip Gloss.
- **Fallback Support:** Works over SSH, small terminal displays (with resize warnings), and without image support.

## Requirements

- Linux / macOS / BSD terminal environment (Linux primary target)
- **Go 1.21+** (only required if building from source)

## Installation

Install using curl (installs into `~/.local/bin`):

```bash
curl -sSL https://raw.githubusercontent.com/aslepenkov/stratagem-zero/main/install.sh | sh
```

Or build locally:

```bash
make build
```

Or install via Go:

```bash
go install ./cmd/stratagem-zero
```

## Usage

Simply run:

```bash
stratagem-zero
```

### CLI Flags

- `--ascii`: Force ASCII fallback graphics mode.
- `--seed <number>`: Set a specific random seed for deterministic stratagem selection.

```bash
stratagem-zero --ascii
```

## Controls

| Key | Direction |
| --- | --- |
| `W` / `↑` | UP |
| `S` / `↓` | DOWN |
| `A` / `←` | LEFT |
| `D` / `→` | RIGHT |
| `q` / `Esc` / `Ctrl+C` | Quit Game |

## Scoring System

- **Base Score:** `100 * sequence_length`
- **Speed Bonus:** `max(0, 1000 - elapsed_ms * 2)`
- **Streak Multiplier:** `1.0 + min(streak, 10) * 0.1`
- **Round Score:** `(base + speed_bonus) * streak_multiplier`

Failing a sequence immediately resets your active streak to 0.

## License

[MIT](LICENSE)
