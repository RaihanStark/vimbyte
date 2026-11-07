package catalog

import (
	"testing"
)

func TestCatalogLoad(t *testing.T) {
	cat, err := Load("")
	if err != nil {
		t.Fatalf("failed to load embedded lessons: %v", err)
	}

	tracks := cat.Tracks()
	if len(tracks) == 0 {
		t.Fatalf("expected tracks, got 0")
	}

	expectedTracks := []string{"basics", "motions", "editing", "visual", "registers"}
	for _, exp := range expectedTracks {
		lessons := cat.List(exp)
		if len(lessons) == 0 {
			t.Errorf("expected lessons in track %s, got 0", exp)
		}
	}

	// Test Get
	lesson, err := cat.Get("01-hjkl")
	if err != nil {
		t.Fatalf("failed to get '01-hjkl': %v", err)
	}
	if lesson.Title == "" || len(lesson.Keys) == 0 || lesson.Drill == "" {
		t.Errorf("unexpected empty fields in lesson: %+v", lesson)
	}

	// Test NextInTrack / PrevInTrack
	next := cat.NextInTrack("basics", "01-hjkl")
	if next == nil || next.ID != "02-modes" {
		t.Errorf("expected 02-modes next, got %+v", next)
	}

	prev := cat.PrevInTrack("basics", "02-modes")
	if prev == nil || prev.ID != "01-hjkl" {
		t.Errorf("expected 01-hjkl prev, got %+v", prev)
	}
}
