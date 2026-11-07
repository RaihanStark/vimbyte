package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/RaihanStark/vimbyte/internal/catalog"
)

// Runner manages the temporary sandbox and launches the user's Vim editor.
type Runner struct {
	EditorOverride string
	Clean          bool
}

// DrillResult contains telemetry and outcomes from an interactive drill session.
type DrillResult struct {
	Duration time.Duration
	Modified bool
}

// New creates a new drill runner instance.
func New(editorOverride string, clean bool) *Runner {
	return &Runner{
		EditorOverride: editorOverride,
		Clean:          clean,
	}
}

// ResolveEditor identifies the target editor binary using environment variables and path lookups.
func (r *Runner) ResolveEditor() (string, error) {
	if r.EditorOverride != "" {
		path, err := exec.LookPath(r.EditorOverride)
		if err == nil {
			return path, nil
		}
		return r.EditorOverride, nil
	}

	if env := os.Getenv("EDITOR"); strings.TrimSpace(env) != "" {
		fields := strings.Fields(env)
		if len(fields) > 0 {
			if path, err := exec.LookPath(fields[0]); err == nil {
				return path, nil
			}
			return fields[0], nil
		}
	}

	if visual := os.Getenv("VISUAL"); strings.TrimSpace(visual) != "" {
		fields := strings.Fields(visual)
		if len(fields) > 0 {
			if path, err := exec.LookPath(fields[0]); err == nil {
				return path, nil
			}
			return fields[0], nil
		}
	}

	// Standard fallbacks in order of preference
	for _, candidate := range []string{"vim", "nvim", "vi"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	return "vim", nil
}

// Run prepares the temporary buffer, executes the drill in the editor, and reliably cleans up.
func (r *Runner) Run(lesson *catalog.Lesson) (*DrillResult, error) {
	editorPath, err := r.ResolveEditor()
	if err != nil {
		return nil, fmt.Errorf("resolving editor: %w", err)
	}

	// 1. Build drill buffer content
	bufferContent := formatDrillBuffer(lesson)

	// 2. Create temporary file
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("vimbyte-drill-%s-*.txt", sanitizeID(lesson.ID)))
	if err != nil {
		return nil, fmt.Errorf("creating temporary drill buffer: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Guarantee cleanup regardless of how the function or process terminates
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.WriteString(bufferContent); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("writing drill content: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("closing drill buffer: %w", err)
	}

	// Snapshot original bytes to check if user made edits
	origBytes, _ := os.ReadFile(tmpPath)

	// 3. Assemble command line arguments
	editorBase := strings.ToLower(filepath.Base(editorPath))
	var args []string

	if r.Clean {
		if strings.Contains(editorBase, "nvim") {
			args = []string{"--clean", "-c", "set number", tmpPath}
		} else {
			// Vim / Vi clean sandbox
			args = []string{"-u", "NONE", "-c", "set number nocompatible", tmpPath}
		}
	} else {
		args = []string{tmpPath}
	}

	cmd := exec.Command(editorPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	runErr := cmd.Run()
	duration := time.Since(startTime)

	if runErr != nil {
		return nil, fmt.Errorf("editor session terminated with error: %w", runErr)
	}

	// 4. Check if buffer was modified
	currentBytes, _ := os.ReadFile(tmpPath)
	modified := !bytes.Equal(origBytes, currentBytes)

	return &DrillResult{
		Duration: duration,
		Modified: modified,
	}, nil
}

func formatDrillBuffer(lesson *catalog.Lesson) string {
	var keysSummary []string
	for _, kb := range lesson.Keys {
		if kb.Desc != "" {
			keysSummary = append(keysSummary, fmt.Sprintf("  %-14s :: %s", kb.Key, kb.Desc))
		} else {
			keysSummary = append(keysSummary, fmt.Sprintf("  %s", kb.Key))
		}
	}

	var buf strings.Builder
	buf.WriteString("\" ==============================================================================\n")
	buf.WriteString(fmt.Sprintf("\" VIMBYTE DRILL :: %s\n", lesson.Title))
	buf.WriteString(fmt.Sprintf("\" Track: %-10s | Lesson: #%-3d | Estimated: ~%ds\n", lesson.Track, lesson.Order, lesson.EstSeconds))
	buf.WriteString("\" ------------------------------------------------------------------------------\n")
	buf.WriteString("\" INSTRUCTIONS:\n")
	for _, line := range strings.Split(strings.TrimSpace(lesson.Instructions), "\n") {
		buf.WriteString(fmt.Sprintf("\"   %s\n", line))
	}
	buf.WriteString("\"\n")
	if len(keysSummary) > 0 {
		buf.WriteString("\" KEYSTROKES / CHEAT SHEET:\n")
		for _, k := range keysSummary {
			buf.WriteString(fmt.Sprintf("\" %s\n", k))
		}
		buf.WriteString("\"\n")
	}
	buf.WriteString("\" PRACTICE TARGET: Edit or navigate on the text below.\n")
	buf.WriteString("\" Finish: Save & quit (:wq) or exit without saving (:q!).\n")
	buf.WriteString("\" ==============================================================================\n\n")

	buf.WriteString(strings.TrimSpace(lesson.Drill))
	buf.WriteString("\n")

	return buf.String()
}

func sanitizeID(id string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, id)
}
