package state

import (
	"time"
)

// Status represents the progression phase of a lesson.
type Status string

const (
	StatusUnseen    Status = "unseen"
	StatusLearning  Status = "learning"
	StatusCompleted Status = "completed"
	StatusMastered  Status = "mastered"
)

// MasteryThreshold defines how many completed repetitions elevate a lesson to mastered.
const MasteryThreshold = 3

// LessonProgress stores user practice metrics and progression milestones for one lesson.
type LessonProgress struct {
	Status          Status     `json:"status"`
	DrillCount      int        `json:"drill_count"`
	RepetitionCount int        `json:"repetition_count"`
	FirstSeenAt     *time.Time `json:"first_seen_at,omitempty"`
	LastPracticedAt *time.Time `json:"last_practiced_at,omitempty"`
	LastCompletedAt *time.Time `json:"last_completed_at,omitempty"`
}

// TrackStats summarizes learning metrics across an entire track.
type TrackStats struct {
	Track            string
	Total            int
	Unseen           int
	Learning         int
	Completed        int
	Mastered         int
	PercentCompleted float64
}

// IsFinished returns true if the lesson has reached completed or mastered state.
func (p LessonProgress) IsFinished() bool {
	return p.Status == StatusCompleted || p.Status == StatusMastered
}
