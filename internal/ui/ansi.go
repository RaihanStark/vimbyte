package ui

import (
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ANSI Color and style constants
const (
	Reset         = "\033[0m"
	Bold          = "\033[1m"
	Dim           = "\033[2m"
	Italic        = "\033[3m"
	Underline     = "\033[4m"
	Black         = "\033[30m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Magenta       = "\033[35m"
	Cyan          = "\033[36m"
	White         = "\033[37m"
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
	BgBlack       = "\033[40m"
	BgBlue        = "\033[44m"
	BgGray        = "\033[100m"

	// Terminal Alternate Screen Buffer and Cursor control sequences
	AltScreenEnter = "\033[?1049h\033[H"
	AltScreenExit  = "\033[?1049l\033[?25h"
	ClearScreen    = "\033[H\033[2J"
	HideCursor     = "\033[?25l"
	ShowCursor     = "\033[?25h"
)

var (
	ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	noColor   = false
)

// EnterAltScreen switches the terminal to the alternate screen buffer (like Vim, less, and other TUIs).
func EnterAltScreen() {
	if IsTerminal() && !IsNoColor() {
		os.Stdout.WriteString(AltScreenEnter + HideCursor)
	}
}

// ExitAltScreen restores the primary screen buffer and cursor.
func ExitAltScreen() {
	if IsTerminal() {
		os.Stdout.WriteString(AltScreenExit)
	}
}

// ClearDisplay clears the entire screen and moves the cursor to row 1, col 1.
func ClearDisplay() {
	if IsTerminal() {
		os.Stdout.WriteString(ClearScreen)
	}
}

// SetNoColor globally configures whether ANSI escapes are emitted.
func SetNoColor(disabled bool) {
	noColor = disabled
}

// IsNoColor returns true if colors are disabled via flag, env var, or non-terminal output.
func IsNoColor() bool {
	if noColor {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	return false
}

// Colorize applies an ANSI sequence if color output is enabled.
func Colorize(code, text string) string {
	if IsNoColor() || code == "" {
		return text
	}
	return code + text + Reset
}

// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// VisibleWidth calculates the visual column width of a string excluding ANSI sequences.
func VisibleWidth(s string) int {
	clean := StripANSI(s)
	return utf8.RuneCountInString(clean)
}

// PadRight pads a string with spaces until its visible width reaches targetWidth.
func PadRight(s string, targetWidth int) string {
	w := VisibleWidth(s)
	if w >= targetWidth {
		return s
	}
	return s + strings.Repeat(" ", targetWidth-w)
}

// PadLeft pads a string on the left until its visible width reaches targetWidth.
func PadLeft(s string, targetWidth int) string {
	w := VisibleWidth(s)
	if w >= targetWidth {
		return s
	}
	return strings.Repeat(" ", targetWidth-w) + s
}

// WrapText wraps text to maxColumn width, preserving existing line breaks.
func WrapText(text string, maxColumn int) []string {
	if maxColumn <= 0 {
		maxColumn = 74
	}

	var result []string
	rawLines := strings.Split(text, "\n")

	for _, rawLine := range rawLines {
		trimmed := strings.TrimRight(rawLine, " \t\r")
		if trimmed == "" {
			result = append(result, "")
			continue
		}

		words := strings.Fields(trimmed)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		// Detect intentional code/indented lines
		leadingSpaces := 0
		for _, r := range rawLine {
			if r == ' ' {
				leadingSpaces++
			} else {
				break
			}
		}

		if leadingSpaces >= 2 {
			// Code block: keep indented line as-is if it fits, or truncate cleanly
			result = append(result, rawLine)
			continue
		}

		var currentLine strings.Builder
		currentWidth := 0

		for _, word := range words {
			wordLen := VisibleWidth(word)
			if currentWidth == 0 {
				currentLine.WriteString(word)
				currentWidth = wordLen
			} else if currentWidth+1+wordLen <= maxColumn {
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
				currentWidth += 1 + wordLen
			} else {
				result = append(result, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				currentWidth = wordLen
			}
		}

		if currentLine.Len() > 0 {
			result = append(result, currentLine.String())
		}
	}

	return result
}
