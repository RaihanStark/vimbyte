package ui

import (
	"fmt"
	"strings"

	"github.com/RaihanStark/vimbyte/internal/catalog"
	"github.com/RaihanStark/vimbyte/internal/state"
)

const CardWidth = 80
const InnerWidth = CardWidth - 4 // 76 chars between '│ ' and ' │'

// RenderStatusBadge returns an ANSI-styled badge for the progression status.
func RenderStatusBadge(s state.Status) string {
	switch s {
	case state.StatusUnseen:
		return Colorize(BrightBlack, "[ UNSEEN ]")
	case state.StatusLearning:
		return Colorize(BrightCyan, "[ LEARNING ]")
	case state.StatusCompleted:
		return Colorize(BrightGreen, "[ COMPLETED ]")
	case state.StatusMastered:
		return Colorize(BrightYellow+Bold, "[ MASTERED ★ ]")
	default:
		return Colorize(BrightBlack, fmt.Sprintf("[ %s ]", strings.ToUpper(string(s))))
	}
}

// RenderLessonCard renders the full 80-column boxed ANSI learning card.
func RenderLessonCard(lesson *catalog.Lesson, progress state.LessonProgress, trackTotal int, showSolution bool) string {
	var b strings.Builder

	// Top Border
	b.WriteString(Colorize(BrightBlack, "┌"+strings.Repeat("─", CardWidth-2)+"┐") + "\n")

	// Header row
	trackHeader := fmt.Sprintf("TRACK: %s (#%d/%d)", strings.ToUpper(lesson.Track), lesson.Order, trackTotal)
	badge := RenderStatusBadge(progress.Status)
	timeEst := fmt.Sprintf("~%ds", lesson.EstSeconds)
	if progress.RepetitionCount > 0 {
		timeEst = fmt.Sprintf("reps: %d | ~%ds", progress.RepetitionCount, lesson.EstSeconds)
	}

	leftHeader := Colorize(Bold+BrightBlue, trackHeader)
	rightHeader := fmt.Sprintf("%s  %s", Colorize(Dim, timeEst), badge)

	headerContent := formatHeaderLine(leftHeader, rightHeader, InnerWidth)
	b.WriteString(renderCardLine(headerContent) + "\n")

	// Divider
	b.WriteString(renderDivider("") + "\n")

	// Title
	titleText := Colorize(Bold+BrightCyan, fmt.Sprintf("%02d. %s", lesson.Order, lesson.Title))
	b.WriteString(renderCardLine(titleText) + "\n")

	// Summary
	if lesson.Summary != "" {
		summaryText := Colorize(Italic+White, lesson.Summary)
		for _, line := range WrapText(summaryText, InnerWidth) {
			b.WriteString(renderCardLine(line) + "\n")
		}
	}

	// Key bindings section
	if len(lesson.Keys) > 0 {
		b.WriteString(renderDivider("KEYSTROKES") + "\n")
		for _, kb := range lesson.Keys {
			keyCol := Colorize(Bold+BrightGreen, PadRight(kb.Key, 14))
			descCol := Colorize(White, kb.Desc)
			b.WriteString(renderCardLine(fmt.Sprintf("%s %s", keyCol, descCol)) + "\n")
		}
	}

	// Body section
	b.WriteString(renderDivider("CONCEPT") + "\n")
	bodyLines := WrapText(strings.TrimSpace(lesson.Body), InnerWidth)
	for _, line := range bodyLines {
		b.WriteString(renderCardLine(line) + "\n")
	}

	// Drill instructions section
	b.WriteString(renderDivider("DRILL INSTRUCTIONS") + "\n")
	instLines := WrapText(strings.TrimSpace(lesson.Instructions), InnerWidth)
	for _, line := range instLines {
		b.WriteString(renderCardLine(Colorize(Yellow, line)) + "\n")
	}

	// Solution section (if enabled or toggled)
	if showSolution && lesson.Solution != "" {
		b.WriteString(renderDivider("SOLUTION") + "\n")
		solLines := WrapText(strings.TrimSpace(lesson.Solution), InnerWidth)
		for _, line := range solLines {
			b.WriteString(renderCardLine(Colorize(BrightGreen, line)) + "\n")
		}
	}

	// Bottom Border
	b.WriteString(Colorize(BrightBlack, "└"+strings.Repeat("─", CardWidth-2)+"┘") + "\n")

	return b.String()
}

func renderCardLine(content string) string {
	vis := VisibleWidth(content)
	padding := InnerWidth - vis
	if padding < 0 {
		padding = 0
	}
	border := Colorize(BrightBlack, "│")
	return fmt.Sprintf("%s %s%s %s", border, content, strings.Repeat(" ", padding), border)
}

func renderDivider(title string) string {
	border := BrightBlack
	if title == "" {
		return Colorize(border, "├"+strings.Repeat("─", CardWidth-2)+"┤")
	}

	titleSection := fmt.Sprintf("── %s ", title)
	visTitle := VisibleWidth(titleSection)
	remaining := CardWidth - 2 - visTitle
	if remaining < 0 {
		remaining = 0
	}

	return Colorize(border, "├"+titleSection+strings.Repeat("─", remaining)+"┤")
}

func formatHeaderLine(left, right string, totalWidth int) string {
	leftVis := VisibleWidth(left)
	rightVis := VisibleWidth(right)
	spaces := totalWidth - (leftVis + rightVis)
	if spaces < 1 {
		spaces = 1
	}
	return left + strings.Repeat(" ", spaces) + right
}
