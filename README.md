# vimbyte

A terminal tool for learning Vim through bite-sized lessons (~60s reads) and short hands-on drills (~30s practice in your real editor).

No streaks, no accounts, and no calendar lockouts.

## Install

```bash
go install github.com/RaihanStark/vimbyte@latest
```

Or build from source:

```bash
git clone https://github.com/RaihanStark/vimbyte.git
cd vimbyte
go build
```

## Quick Start

Run `vimbyte` with no arguments to open the next lesson:

```bash
vimbyte
```

Use single keys to interact:

- `d` - Open the exercise in Vim (uses `$EDITOR`, `$VISUAL`, or `vim`)
- `Enter` - Mark the lesson done
- `s` - Toggle solution and hints
- `n` / `p` - Next or previous lesson
- `q` - Quit

## Commands

```bash
vimbyte                  # Open the next lesson card
vimbyte drill [id]       # Jump straight into an editor drill
vimbyte done [id]        # Mark a lesson as completed
vimbyte review           # Review completed lessons to build mastery
vimbyte tracks           # View available tracks and progress
vimbyte switch <track>   # Switch track (basics, motions, editing, visual, registers)
vimbyte list [track]     # List all lessons in a track
vimbyte reset [track]    # Reset progress for a track (or --all)
```

## Syllabus

```text
vimbyte
├── basics
│   ├── 01-hjkl                # Home row navigation (h, j, k, l)
│   ├── 02-modes               # Normal & Insert modes (i, I, <Esc>)
│   ├── 03-word-motions        # Word jumps (w, b, e, W, B, E)
│   └── 04-simple-edits        # Append, new lines & undo (a, A, o, O, u, <C-r>)
├── motions
│   ├── 05-line-motions        # Line boundaries (0, ^, $, g_)
│   ├── 06-find-char           # Inline character search (f, F, t, T, ;, ,)
│   ├── 07-paragraph-motions   # Block jumping ({ and })
│   └── 08-matching-pairs      # Bracket & delimiter matching (%)
├── editing
│   ├── 09-operator-motion     # Verb + noun grammar (d, c, y + motions)
│   ├── 10-text-objects        # Inside/around targets (iw, aw, i", i()
│   ├── 11-replace-change      # Quick edits (r, s, C, D, x)
│   └── 12-dot-formula         # Repeating changes (.)
├── visual
│   ├── 13-visual-modes        # Character, line & block visual (v, V, <C-v>)
│   └── 14-visual-block-insert # Column editing (<C-v> + I / A)
└── registers
    ├── 15-yank-put            # Copying and pasting (yy, yw, p, P)
    └── 16-named-registers     # Named clips & OS clipboard ("a-"z, "+y, "+p)
```

## Flags

- `--clean`: Run Vim without plugins or configs (`vim -u NONE -c "set number nocompatible"`)
- `--editor <cmd>`: Use a specific editor (e.g. `--editor nvim`)
- `--dir <path>`: Load extra custom lessons from a folder of YAML files
- `--no-color`: Disable ANSI colors

## Progress

Progress is saved locally at `~/.config/vimbyte/state.json` (or `$XDG_CONFIG_HOME/vimbyte/state.json`).

Lessons move through four stages:
- **unseen**: Not started yet
- **learning**: Opened or drilled
- **completed**: Marked done
- **mastered**: Practiced 3 or more times

## License

MIT
