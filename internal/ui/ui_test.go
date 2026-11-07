package ui

import (
	"strings"
	"testing"

	"github.com/RaihanStark/vimbyte/internal/catalog"
	"github.com/RaihanStark/vimbyte/internal/state"
)

func TestVisibleWidthAndWrap(t *testing.T) {
	colored := Colorize(Bold+BrightCyan, "Hello World")
	if VisibleWidth(colored) != 11 {
		t.Errorf("expected visible width 11, got %d", VisibleWidth(colored))
	}

	lines := WrapText("This is a long sentence that should easily wrap across multiple lines when column width is small.", 30)
	if len(lines) < 2 {
		t.Errorf("expected multiple lines, got %d", len(lines))
	}
	for _, l := range lines {
		if VisibleWidth(l) > 30 {
			t.Errorf("line exceeds max column: %s (%d)", l, VisibleWidth(l))
		}
	}
}

func TestRenderLessonCard(t *testing.T) {
	lesson := &catalog.Lesson{
		ID:         "01-hjkl",
		Track:      "basics",
		Order:      1,
		EstSeconds: 45,
		Title:      "Home Row Navigation",
		Summary:    "Navigate with h, j, k, l",
		Keys: []catalog.KeyBinding{
			{Key: "h", Desc: "Left"},
			{Key: "j", Desc: "Down"},
		},
		Body:         "Keep your fingers on the home row.",
		Instructions: "Follow the steps below to practice.",
		Solution:     "Use 4j and 2l.",
	}

	progress := state.LessonProgress{
		Status:          state.StatusLearning,
		DrillCount:      2,
		RepetitionCount: 0,
	}

	card := RenderLessonCard(lesson, progress, 4, true)
	if !strings.Contains(card, "01. Home Row Navigation") {
		t.Errorf("card missing title")
	}
	if !strings.Contains(card, "[ LEARNING ]") {
		t.Errorf("card missing learning status badge")
	}
	if !strings.Contains(card, "KEYSTROKES") {
		t.Errorf("card missing keystrokes section")
	}
	if !strings.Contains(card, "SOLUTION") {
		t.Errorf("card missing solution section")
	}

	// Verify line lengths inside card
	lines := strings.Split(card, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		w := VisibleWidth(line)
		if w != CardWidth {
			t.Errorf("line %d has width %d, expected %d: %q", i+1, w, CardWidth, StripANSI(line))
		}
	}
}

func TestRenderProgressBar(t *testing.T) {
	bar := RenderProgressBar(5, 10, 10)
	if !strings.Contains(bar, "50%") {
		t.Errorf("expected 50%% in progress bar, got %s", bar)
	}
}
