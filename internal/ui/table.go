package ui

import (
	"fmt"
	"strings"

	"github.com/RaihanStark/vimbyte/internal/catalog"
	"github.com/RaihanStark/vimbyte/internal/state"
)

// RenderProgressBar generates a visual unicode completion bar: [██████░░░░] 60%
func RenderProgressBar(completed, total int, barWidth int) string {
	if total <= 0 {
		return "[────────────────────]   0%"
	}
	if barWidth <= 0 {
		barWidth = 16
	}

	ratio := float64(completed) / float64(total)
	if ratio > 1.0 {
		ratio = 1.0
	}
	filledChars := int(ratio * float64(barWidth))
	emptyChars := barWidth - filledChars

	barColor := BrightGreen
	if ratio < 0.33 {
		barColor = Yellow
	} else if ratio < 0.75 {
		barColor = Cyan
	}

	filledStr := Colorize(barColor, strings.Repeat("█", filledChars))
	emptyStr := Colorize(BrightBlack, strings.Repeat("░", emptyChars))
	percentStr := fmt.Sprintf("%3.0f%%", ratio*100)

	return fmt.Sprintf("[%s%s] %s", filledStr, emptyStr, percentStr)
}

// RenderTracksTable displays all available lesson tracks with their completion metrics.
func RenderTracksTable(tracks []string, activeTrack string, stats map[string]state.TrackStats) string {
	var b strings.Builder

	b.WriteString(Colorize(Bold+BrightWhite, " VIMBYTE LEARNING TRACKS\n"))
	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))

	header := fmt.Sprintf(" %-2s %-14s %-12s %-24s %-10s",
		"", "TRACK", "COMPLETED", "PROGRESS", "MASTERY")
	b.WriteString(Colorize(Bold+BrightBlue, header) + "\n")
	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))

	for _, track := range tracks {
		st := stats[track]
		marker := "  "
		trackDisplay := track
		if track == activeTrack {
			marker = Colorize(Bold+BrightGreen, "▶ ")
			trackDisplay = Colorize(Bold+BrightWhite, track)
		} else {
			trackDisplay = Colorize(White, track)
		}

		doneCount := st.Completed + st.Mastered
		countStr := fmt.Sprintf("%d/%d", doneCount, st.Total)
		progressBar := RenderProgressBar(doneCount, st.Total, 14)
		masteryStr := fmt.Sprintf("%d★", st.Mastered)
		if st.Mastered > 0 {
			masteryStr = Colorize(BrightYellow+Bold, masteryStr)
		} else {
			masteryStr = Colorize(BrightBlack, masteryStr)
		}

		colTrack := PadRight(trackDisplay, 14)
		colCount := PadRight(countStr, 12)
		colProgress := PadRight(progressBar, 24)
		colMastery := PadRight(masteryStr, 10)

		row := fmt.Sprintf(" %s%s %s %s %s",
			marker,
			colTrack,
			colCount,
			colProgress,
			colMastery,
		)
		b.WriteString(row + "\n")
	}

	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))
	b.WriteString(Colorize(Dim, " Switch track with: vimbyte switch <track>\n"))

	return b.String()
}

// RenderLessonsList formats all lessons within a track as a structured table.
func RenderLessonsList(track string, lessons []*catalog.Lesson, progressFunc func(id string) state.LessonProgress) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s %s\n",
		Colorize(Bold+BrightWhite, " LESSONS IN TRACK:"),
		Colorize(Bold+BrightCyan, strings.ToUpper(track)),
	))
	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))

	header := fmt.Sprintf(" %-3s %-18s %-28s %-14s %-5s %-5s",
		"#", "ID", "TITLE", "STATUS", "REPS", "EST")
	b.WriteString(Colorize(Bold+BrightBlue, header) + "\n")
	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))

	for _, l := range lessons {
		p := progressFunc(l.ID)
		badge := RenderStatusBadge(p.Status)

		title := l.Title
		if VisibleWidth(title) > 27 {
			runes := []rune(title)
			title = string(runes[:25]) + ".."
		}

		reps := fmt.Sprintf("%d", p.RepetitionCount)
		if p.RepetitionCount > 0 {
			reps = Colorize(BrightGreen, reps)
		} else {
			reps = Colorize(BrightBlack, reps)
		}

		est := fmt.Sprintf("%ds", l.EstSeconds)

		colOrder := PadRight(fmt.Sprintf("%-3d", l.Order), 3)
		colID := PadRight(l.ID, 18)
		colTitle := PadRight(title, 28)
		colBadge := PadRight(badge, 14)
		colReps := PadRight(reps, 5)
		colEst := PadRight(est, 5)

		row := fmt.Sprintf(" %s %s %s %s %s %s",
			colOrder,
			colID,
			colTitle,
			colBadge,
			colReps,
			colEst,
		)
		b.WriteString(row + "\n")
	}

	b.WriteString(Colorize(BrightBlack, " "+strings.Repeat("─", CardWidth-2)+"\n"))
	b.WriteString(Colorize(Dim, " Practice any lesson directly with: vimbyte drill <lesson_id>\n"))

	return b.String()
}
