# vimbyte

> **Bite-sized Vim micro-learning for your terminal.**  
> ~60s conceptual reads, ~30s hands-on drills. Zero streaks, zero calendar-based lockouts, zero runtime bloat.

`vimbyte` is a self-contained, terminal-based CLI tool designed to help developers build deep muscle memory for Vim. Built entirely in Go with zero external runtime dependencies.

---

## Features

- **Micro-Learning Design**: Every lesson is calibrated for ~60s reading and a ~30s focused drill in a real Vim sandbox.
- **No Calendar Lockouts / Streaks**: Learn at your own pace without artificial streaks, notifications, or daily gates.
- **Progression State Machine**:
  - `[ UNSEEN ]`: Lessons awaiting discovery.
  - `[ LEARNING ]`: Lessons practiced in the editor buffer.
  - `[ COMPLETED ]`: Lessons marked done.
  - `[ MASTERED ★ ]`: Lessons reinforced with 3+ practice repetitions.
- **80-Column Boxed ANSI Card Interface**: Formatted with box borders, status headers, aligned shortcut tables, and single-key action prompts.
- **Vim Sandbox Runner**:
  - Resolves `$EDITOR`, `$VISUAL`, `vim`, `nvim`, or `vi`.
  - Automatically loads drill instructions and starting snippets into a temporary buffer via `os.CreateTemp`.
  - Reliable automatic cleanup on exit.
  - Supports `--clean` mode (`vim -u NONE -c "set number nocompatible"`).
- **Atomic, Thread-Safe JSON Store**: Saves user state to `$XDG_CONFIG_HOME/vimbyte/state.json` (or `~/.config/vimbyte/state.json`).
- **Embedded Catalog + External Extensions**: Ships with 16 lessons embedded via `//go:embed` across 5 tracks (`basics`, `motions`, `editing`, `visual`, `registers`), with optional external directory loading via `--dir`.

---

## Installation

```bash
# Clone the repository
git clone https://github.com/RaihanStark/vimbyte.git
cd vimbyte

# Build binary
go build -o vimbyte .

# Optional: Install to your PATH
go install
```

---

## CLI Usage

```
vimbyte [command] [options]
```

### Commands

| Command | Description |
|---|---|
| `vimbyte` | Launch the interactive learning card for the next recommended lesson |
| `vimbyte drill [lesson_id]` | Open the drill directly in your Vim editor sandbox |
| `vimbyte done [lesson_id]` | Mark a lesson completed and advance progression |
| `vimbyte review` | Spaced review mode: drill completed lessons to build mastery |
| `vimbyte tracks` | View all tracks, completion percentages, and mastery counts |
| `vimbyte switch <track>` | Switch the active track (`basics`, `motions`, `editing`, `visual`, `registers`) |
| `vimbyte list [track]` | Tabular view of all lessons with statuses, repetitions, and times |
| `vimbyte reset [track\|--all]`| Reset user progress for a track or globally |

### Global Flags

- `--clean`: Run Vim in a pristine configuration sandbox (`vim -u NONE -c "set number nocompatible"`).
- `--dir <path>`: Load additional or custom lesson YAML files from an external directory.
- `--editor <cmd>`: Override editor binary (e.g. `--editor=nvim`).
- `--no-color`: Disable all ANSI color sequences.
- `-v, --version`: Display application version.
- `-h, --help`: Display help and examples.

---

## Interactive Shortcuts

When viewing a lesson card in interactive mode:
- `[d]`: Launch the drill in Vim.
- `[Enter]`: Mark the lesson done and advance to the next lesson.
- `[s]`: Toggle solution text and keystroke explanations.
- `[n]`: Move to the next lesson in the track.
- `[p]`: Move to the previous lesson in the track.
- `[q]` / `[Esc]` / `[Ctrl-C]`: Exit.

---

## Project Structure

```
├── main.go                     # Entry point & CLI subcommand routing
├── internal/
│   ├── catalog/                # Embedded YAML lesson loader & domain models
│   │   ├── catalog.go
│   │   ├── model.go
│   │   └── lessons/            # Embedded YAML lessons across 5 tracks
│   │       ├── 01_h_j_k_l.yaml
│   │       ├── 02_modes.yaml
│   │       ├── ...
│   │       └── 16_named_registers.yaml
│   ├── state/                  # Atomic JSON state store & progression engine
│   │   ├── state.go
│   │   └── progression.go
│   ├── runner/                 # Sandboxed editor launcher & temp buffer engine
│   │   └── runner.go
│   └── ui/                     # 80-column ANSI card display, tables, and input
│       ├── ansi.go
│       ├── card.go
│       ├── table.go
│       └── input.go
└── README.md
```

---

## License

MIT
