package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StateData models the JSON schema persisted to disk.
type StateData struct {
	Version     int                       `json:"version"`
	ActiveTrack string                    `json:"active_track"`
	Progress    map[string]LessonProgress `json:"progress"`
}

// Store provides an atomic, thread-safe JSON repository for user progress.
type Store struct {
	mu       sync.RWMutex
	filePath string
	data     StateData
}

// DefaultStatePath resolves $XDG_CONFIG_HOME/vimbyte/state.json or ~/.config/vimbyte/state.json.
func DefaultStatePath() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "vimbyte", "state.json")
	}

	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".config", "vimbyte", "state.json")
	}

	return filepath.Join(".", ".vimbyte", "state.json")
}

// Open loads existing state from the default path or initializes a fresh store.
func Open() (*Store, error) {
	return OpenPath(DefaultStatePath())
}

// OpenPath opens or initializes a store at a specified path.
func OpenPath(path string) (*Store, error) {
	s := &Store{
		filePath: path,
		data: StateData{
			Version:     1,
			ActiveTrack: "basics",
			Progress:    make(map[string]LessonProgress),
		},
	}

	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// First run: save empty default state to verify write permissions
			return s.saveLocked()
		}
		return fmt.Errorf("reading state file %s: %w", s.filePath, err)
	}

	if len(data) == 0 {
		return nil
	}

	var loaded StateData
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("parsing state JSON: %w", err)
	}

	if loaded.Progress == nil {
		loaded.Progress = make(map[string]LessonProgress)
	}
	if loaded.ActiveTrack == "" {
		loaded.ActiveTrack = "basics"
	}

	s.data = loaded
	return nil
}

// saveLocked writes state atomically via a temporary file and rename.
func (s *Store) saveLocked() error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating state directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state data: %w", err)
	}

	// Create temp file in same directory to guarantee atomic rename across filesystems
	tmpFile, err := os.CreateTemp(dir, "state-*.json.tmp")
	if err != nil {
		return fmt.Errorf("creating temporary state file: %w", err)
	}
	tmpName := tmpFile.Name()

	// Ensure cleanup on write/sync failure
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("writing state data: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("syncing state file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp state file: %w", err)
	}

	if err := os.Rename(tmpName, s.filePath); err != nil {
		return fmt.Errorf("replacing state file: %w", err)
	}

	cleanup = false
	return nil
}

// GetActiveTrack returns the currently selected lesson track.
func (s *Store) GetActiveTrack() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.data.ActiveTrack == "" {
		return "basics"
	}
	return s.data.ActiveTrack
}

// SetActiveTrack sets and persists the current active track.
func (s *Store) SetActiveTrack(track string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.ActiveTrack = track
	return s.saveLocked()
}

// GetLessonProgress returns progression details for the requested lesson ID.
func (s *Store) GetLessonProgress(lessonID string) LessonProgress {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data.Progress[lessonID]
	if !ok {
		return LessonProgress{
			Status:          StatusUnseen,
			DrillCount:      0,
			RepetitionCount: 0,
		}
	}
	return p
}

// RecordDrill tracks when a user enters the interactive Vim buffer for a lesson.
func (s *Store) RecordDrill(lessonID string) (LessonProgress, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	p, ok := s.data.Progress[lessonID]
	if !ok {
		p = LessonProgress{
			Status:          StatusLearning,
			DrillCount:      1,
			RepetitionCount: 0,
			FirstSeenAt:     &now,
			LastPracticedAt: &now,
		}
	} else {
		p.DrillCount++
		p.LastPracticedAt = &now
		if p.FirstSeenAt == nil {
			p.FirstSeenAt = &now
		}
		if p.Status == StatusUnseen {
			p.Status = StatusLearning
		}
	}

	s.data.Progress[lessonID] = p
	return p, s.saveLocked()
}

// MarkDone advances a lesson to completed or mastered status.
func (s *Store) MarkDone(lessonID string) (LessonProgress, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	p, ok := s.data.Progress[lessonID]
	if !ok {
		p = LessonProgress{
			Status:          StatusCompleted,
			DrillCount:      1,
			RepetitionCount: 1,
			FirstSeenAt:     &now,
			LastPracticedAt: &now,
			LastCompletedAt: &now,
		}
	} else {
		p.RepetitionCount++
		p.LastPracticedAt = &now
		p.LastCompletedAt = &now
		if p.FirstSeenAt == nil {
			p.FirstSeenAt = &now
		}
		if p.RepetitionCount >= MasteryThreshold {
			p.Status = StatusMastered
		} else {
			p.Status = StatusCompleted
		}
	}

	s.data.Progress[lessonID] = p
	return p, s.saveLocked()
}

// ResetTrack clears user progress records for the given lesson IDs.
func (s *Store) ResetTrack(lessonIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range lessonIDs {
		delete(s.data.Progress, id)
	}
	return s.saveLocked()
}

// ResetAll completely clears all tracked progress.
func (s *Store) ResetAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.Progress = make(map[string]LessonProgress)
	return s.saveLocked()
}

// GetFilePath returns the underlying file location on disk.
func (s *Store) GetFilePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filePath
}
