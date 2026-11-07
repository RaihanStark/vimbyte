package runner

import (
	"strings"
	"testing"

	"github.com/RaihanStark/vimbyte/internal/catalog"
)

func TestFormatDrillBuffer(t *testing.T) {
	lesson := &catalog.Lesson{
		ID:         "01-test",
		Track:      "basics",
		Order:      1,
		EstSeconds: 30,
		Title:      "Test Title",
		Keys: []catalog.KeyBinding{
			{Key: "h", Desc: "Left"},
		},
		Instructions: "Move left",
		Drill:        "Sample text here",
	}

	content := formatDrillBuffer(lesson)
	if !strings.Contains(content, "VIMBYTE DRILL :: Test Title") {
		t.Errorf("buffer missing title")
	}
	if !strings.Contains(content, "Sample text here") {
		t.Errorf("buffer missing drill snippet")
	}
	if !strings.Contains(content, "Left") {
		t.Errorf("buffer missing key description")
	}
}

func TestRunnerExecution(t *testing.T) {
	runner := New("true", false)
	lesson := &catalog.Lesson{
		ID:           "01-test",
		Track:        "basics",
		Order:        1,
		EstSeconds:   30,
		Title:        "Test Title",
		Instructions: "Practice",
		Drill:        "Sample text",
	}

	res, err := runner.Run(lesson)
	if err != nil {
		t.Fatalf("runner failed: %v", err)
	}
	if res == nil {
		t.Fatalf("expected result, got nil")
	}
}
