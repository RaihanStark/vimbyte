# vimbyte Syllabus

Curriculum of all 16 built-in lessons across 5 tracks. Each lesson takes ~60s to read and includes a ~30s practical drill in your editor.

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

---

## Tracks & Lessons

### 1. Basics (`basics`)
Core habits and the modal interaction model.

| ID | Lesson | Keys | Est. |
|---|---|---|---|
| `01-hjkl` | Home Row Navigation | `h`, `j`, `k`, `l` | ~45s |
| `02-modes` | Normal & Insert Modes | `i`, `I`, `<Esc>`, `:wq` | ~50s |
| `03-word-motions` | Word Motions | `w`, `b`, `e`, `W`, `B`, `E` | ~60s |
| `04-simple-edits` | Appending, Lines & Undo | `a`, `A`, `o`, `O`, `u`, `<C-r>` | ~60s |

### 2. Motions (`motions`)
Navigate across lines and blocks without single-character stepping.

| ID | Lesson | Keys | Est. |
|---|---|---|---|
| `05-line-motions` | Line Boundaries | `0`, `^`, `$`, `g_` | ~50s |
| `06-find-char` | Inline Search | `f`, `F`, `t`, `T`, `;`, `,` | ~60s |
| `07-paragraph-motions` | Block Motions | `{`, `}` | ~45s |
| `08-matching-pairs` | Delimiter Jumping | `%` | ~40s |

### 3. Editing Grammar (`editing`)
Combine operator verbs with targets and text objects.

| ID | Lesson | Keys | Est. |
|---|---|---|---|
| `09-operator-motion` | Operator + Motion Grammar | `d`, `c`, `y` (+ motions) | ~60s |
| `10-text-objects` | Text Objects | `iw`, `aw`, `i"`, `i(`, `it` | ~70s |
| `11-replace-change` | Surgical Edits | `r`, `s`, `C`, `D`, `x` | ~50s |
| `12-dot-formula` | The Dot Formula | `.` | ~60s |

### 4. Visual Mode (`visual`)
Selecting regions and editing columns.

| ID | Lesson | Keys | Est. |
|---|---|---|---|
| `13-visual-modes` | Visual Modes | `v`, `V`, `<C-v>` | ~50s |
| `14-visual-block-insert` | Visual Block Insert | `<C-v>` + `I` / `A` | ~60s |

### 5. Registers (`registers`)
Copying, pasting, multiple clipboards, and OS clipboard integration.

| ID | Lesson | Keys | Est. |
|---|---|---|---|
| `15-yank-put` | Yanking & Putting | `yy`, `yw`, `p`, `P` | ~50s |
| `16-named-registers` | Named Registers & Clipboard | `"a` - `"z`, `"+y`, `"+p`, `:reg` | ~60s |
