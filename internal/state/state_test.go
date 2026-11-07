package state

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestProgressionLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "state.json")

	store, err := OpenPath(storePath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}

	// 1. Initial status should be unseen
	p := store.GetLessonProgress("01-hjkl")
	if p.Status != StatusUnseen {
		t.Fatalf("expected unseen, got %s", p.Status)
	}

	// 2. Drill -> learning
	p, err = store.RecordDrill("01-hjkl")
	if err != nil {
		t.Fatalf("failed to record drill: %v", err)
	}
	if p.Status != StatusLearning || p.DrillCount != 1 {
		t.Fatalf("expected learning with 1 drill, got status=%s, drills=%d", p.Status, p.DrillCount)
	}

	// 3. MarkDone -> completed
	p, err = store.MarkDone("01-hjkl")
	if err != nil {
		t.Fatalf("failed to mark done: %v", err)
	}
	if p.Status != StatusCompleted || p.RepetitionCount != 1 {
		t.Fatalf("expected completed with 1 rep, got status=%s, reps=%d", p.Status, p.RepetitionCount)
	}

	// 4. MarkDone x2 more -> mastered
	_, _ = store.MarkDone("01-hjkl")
	p, err = store.MarkDone("01-hjkl")
	if err != nil {
		t.Fatalf("failed to mark done 3rd time: %v", err)
	}
	if p.Status != StatusMastered || p.RepetitionCount != 3 {
		t.Fatalf("expected mastered with 3 reps, got status=%s, reps=%d", p.Status, p.RepetitionCount)
	}

	// 5. Test reloading from disk
	store2, err := OpenPath(storePath)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	p2 := store2.GetLessonProgress("01-hjkl")
	if p2.Status != StatusMastered || p2.RepetitionCount != 3 {
		t.Fatalf("persisted state mismatch: status=%s, reps=%d", p2.Status, p2.RepetitionCount)
	}

	// 6. ResetTrack
	if err := store2.ResetTrack([]string{"01-hjkl"}); err != nil {
		t.Fatalf("failed reset track: %v", err)
	}
	p3 := store2.GetLessonProgress("01-hjkl")
	if p3.Status != StatusUnseen {
		t.Fatalf("expected unseen after reset, got %s", p3.Status)
	}
}

func TestStoreConcurrency(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "state.json")

	store, err := OpenPath(storePath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			lessonID := "concurrency-test"
			if idx%2 == 0 {
				_, _ = store.RecordDrill(lessonID)
			} else {
				_, _ = store.MarkDone(lessonID)
			}
		}(i)
	}
	wg.Wait()

	p := store.GetLessonProgress("concurrency-test")
	if p.DrillCount == 0 && p.RepetitionCount == 0 {
		t.Fatalf("expected non-zero progress after concurrent operations")
	}
}
