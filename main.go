package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/RaihanStark/vimbyte/internal/catalog"
	"github.com/RaihanStark/vimbyte/internal/runner"
	"github.com/RaihanStark/vimbyte/internal/state"
	"github.com/RaihanStark/vimbyte/internal/ui"
)

const version = "1.0.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s %s\n", ui.Colorize(ui.Bold+ui.BrightRed, "Error:"), err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var (
		flagClean   bool
		flagDir     string
		flagEditor  string
		flagNoColor bool
		flagHelp    bool
		flagVersion bool
	)

	fs := flag.NewFlagSet("vimbyte", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = printUsage

	fs.BoolVar(&flagClean, "clean", false, "Run Vim in pristine sandbox (vim -u NONE -c \"set number nocompatible\")")
	fs.StringVar(&flagDir, "dir", "", "Path to additional external lesson YAML files")
	fs.StringVar(&flagEditor, "editor", "", "Override editor binary (defaults to $EDITOR, $VISUAL, or vim)")
	fs.BoolVar(&flagNoColor, "no-color", false, "Disable ANSI color styling")
	fs.BoolVar(&flagHelp, "h", false, "Show help message")
	fs.BoolVar(&flagHelp, "help", false, "Show help message")
	fs.BoolVar(&flagVersion, "v", false, "Show version information")
	fs.BoolVar(&flagVersion, "version", false, "Show version information")

	// Pre-extract subcommand and intermixed flags
	var filteredFlags []string
	var subArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			filteredFlags = append(filteredFlags, arg)
			// Handle flags that accept values
			if (arg == "--dir" || arg == "-dir" || arg == "--editor" || arg == "-editor") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				filteredFlags = append(filteredFlags, args[i+1])
				i++
			}
		} else {
			subArgs = append(subArgs, arg)
		}
	}

	if err := fs.Parse(filteredFlags); err != nil {
		return err
	}

	if flagNoColor {
		ui.SetNoColor(true)
	}

	if flagVersion {
		fmt.Printf("vimbyte v%s - micro-learning for terminal Vim\n", version)
		return nil
	}

	if flagHelp && len(subArgs) == 0 {
		printUsage()
		return nil
	}

	// 1. Initialize State Store
	store, err := state.Open()
	if err != nil {
		return fmt.Errorf("opening state store: %w", err)
	}

	// 2. Load Catalog
	cat, err := catalog.Load(flagDir)
	if err != nil {
		return fmt.Errorf("loading lesson catalog: %w", err)
	}

	// 3. Initialize Drill Runner
	drillRunner := runner.New(flagEditor, flagClean)

	// 4. Dispatch subcommand
	subcommand := "next"
	var cmdParams []string
	if len(subArgs) > 0 {
		subcommand = subArgs[0]
		cmdParams = subArgs[1:]
	}

	switch subcommand {
	case "next":
		return handleInteractiveCard(cat, store, drillRunner, "")
	case "drill":
		targetID := ""
		if len(cmdParams) > 0 {
			targetID = cmdParams[0]
		}
		return handleDrill(cat, store, drillRunner, targetID)
	case "done":
		targetID := ""
		if len(cmdParams) > 0 {
			targetID = cmdParams[0]
		}
		return handleDone(cat, store, targetID)
	case "review":
		return handleReview(cat, store, drillRunner)
	case "tracks":
		return handleTracks(cat, store)
	case "switch":
		if len(cmdParams) == 0 {
			return fmt.Errorf("missing track name: usage: vimbyte switch <track>")
		}
		return handleSwitch(cat, store, cmdParams[0])
	case "list":
		track := store.GetActiveTrack()
		if len(cmdParams) > 0 {
			track = cmdParams[0]
		}
		return handleList(cat, store, track)
	case "reset":
		allFlag := false
		targetTrack := ""
		for _, p := range cmdParams {
			if p == "--all" {
				allFlag = true
			} else {
				targetTrack = p
			}
		}
		return handleReset(cat, store, targetTrack, allFlag)
	default:
		// Check if user passed a lesson ID directly (e.g. vimbyte 01-hjkl)
		if _, err := cat.Get(subcommand); err == nil {
			return handleInteractiveCard(cat, store, drillRunner, subcommand)
		}
		printUsage()
		return fmt.Errorf("unknown command '%s'", subcommand)
	}
}

func printUsage() {
	fmt.Printf(`%s - Terminal-based micro-learning tool for Vim

%s
  vimbyte [command] [options]

%s
  (default)            Open interactive lesson card for next unfinished lesson
  drill  [lesson_id]   Launch drill directly in Vim sandbox
  done   [lesson_id]   Mark a lesson completed and advance progression
  review               Review completed lessons to build mastery repetitions
  tracks               Display all tracks and completion progress
  switch <track>       Switch the active lesson track
  list   [track]       List all lessons in current or specified track
  reset  [track|--all] Reset learning progress for a track or globally

%s
  --clean              Run Vim in pristine sandbox (vim -u NONE -c "set number nocompatible")
  --dir <path>         Path to additional external lesson YAML files
  --editor <cmd>       Override editor binary (defaults to $EDITOR, $VISUAL, or vim)
  --no-color           Disable ANSI color styling
  -v, --version        Show version information
  -h, --help           Show this help message

%s
  $ vimbyte                  # Practice next lesson interactively
  $ vimbyte drill 01-hjkl    # Practice home-row navigation in Vim
  $ vimbyte done 01-hjkl     # Mark lesson completed
  $ vimbyte tracks           # View progress across all tracks
  $ vimbyte switch motions   # Switch to the motions track
`,
		ui.Colorize(ui.Bold+ui.BrightWhite, "vimbyte"),
		ui.Colorize(ui.Bold+ui.BrightBlue, "USAGE:"),
		ui.Colorize(ui.Bold+ui.BrightBlue, "COMMANDS:"),
		ui.Colorize(ui.Bold+ui.BrightBlue, "FLAGS:"),
		ui.Colorize(ui.Bold+ui.BrightBlue, "EXAMPLES:"),
	)
}

func handleInteractiveCard(cat *catalog.Catalog, store *state.Store, r *runner.Runner, lessonID string) error {
	activeTrack := store.GetActiveTrack()
	trackLessons := cat.List(activeTrack)
	if len(trackLessons) == 0 {
		// Fallback to first available track if active track is empty
		allTracks := cat.Tracks()
		if len(allTracks) == 0 {
			return fmt.Errorf("no lessons found in catalog")
		}
		activeTrack = allTracks[0]
		_ = store.SetActiveTrack(activeTrack)
		trackLessons = cat.List(activeTrack)
	}

	var currentLesson *catalog.Lesson
	if lessonID != "" {
		l, err := cat.Get(lessonID)
		if err != nil {
			return err
		}
		currentLesson = l
	} else {
		// Find first unfinished lesson in active track
		for _, l := range trackLessons {
			p := store.GetLessonProgress(l.ID)
			if !p.IsFinished() {
				currentLesson = l
				break
			}
		}
		if currentLesson == nil && len(trackLessons) > 0 {
			// All completed: pick the one with lowest repetitions to review
			currentLesson = trackLessons[0]
			minReps := store.GetLessonProgress(currentLesson.ID).RepetitionCount
			for _, l := range trackLessons {
				reps := store.GetLessonProgress(l.ID).RepetitionCount
				if reps < minReps {
					currentLesson = l
					minReps = reps
				}
			}
		}
	}

	showSolution := false

	// Interactive Loop
	for {
		progress := store.GetLessonProgress(currentLesson.ID)
		trackList := cat.List(currentLesson.Track)
		cardOutput := ui.RenderLessonCard(currentLesson, progress, len(trackList), showSolution)

		fmt.Print(cardOutput)

		if !ui.IsTerminal() {
			// Non-interactive environment: print card and finish
			return nil
		}

		fmt.Print(ui.RenderPromptBar(showSolution))

		action, err := ui.ReadSingleKey()
		if err != nil {
			return nil
		}

		switch action {
		case ui.ActionDrill:
			_, _ = store.RecordDrill(currentLesson.ID)
			res, runErr := r.Run(currentLesson)
			if runErr != nil {
				fmt.Printf("%s %v\n", ui.Colorize(ui.Bold+ui.BrightRed, "Drill error:"), runErr)
			} else {
				modNotice := ""
				if res.Modified {
					modNotice = " (buffer modified)"
				}
				fmt.Printf("\n%s Drill session ended (%0.1fs)%s.\n",
					ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
					res.Duration.Seconds(),
					modNotice,
				)
			}
		case ui.ActionDone:
			newProg, err := store.MarkDone(currentLesson.ID)
			if err != nil {
				return fmt.Errorf("marking lesson done: %w", err)
			}
			badge := ui.RenderStatusBadge(newProg.Status)
			fmt.Printf("\n%s Marked '%s' as done! %s (repetitions: %d)\n\n",
				ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
				currentLesson.ID,
				badge,
				newProg.RepetitionCount,
			)

			// Advance to next lesson in track if available
			next := cat.NextInTrack(currentLesson.Track, currentLesson.ID)
			if next != nil {
				currentLesson = next
				showSolution = false
			}
		case ui.ActionNext:
			next := cat.NextInTrack(currentLesson.Track, currentLesson.ID)
			if next != nil {
				currentLesson = next
				showSolution = false
			} else {
				fmt.Println(ui.Colorize(ui.Dim, " End of track reached."))
			}
		case ui.ActionPrev:
			prev := cat.PrevInTrack(currentLesson.Track, currentLesson.ID)
			if prev != nil {
				currentLesson = prev
				showSolution = false
			}
		case ui.ActionToggleSolution:
			showSolution = !showSolution
		case ui.ActionQuit:
			fmt.Println()
			return nil
		default:
			// Ignore unrecognized single keystroke
		}
	}
}

func handleDrill(cat *catalog.Catalog, store *state.Store, r *runner.Runner, lessonID string) error {
	var targetLesson *catalog.Lesson
	if lessonID != "" {
		l, err := cat.Get(lessonID)
		if err != nil {
			return err
		}
		targetLesson = l
	} else {
		// Pick next unfinished lesson in active track
		activeTrack := store.GetActiveTrack()
		for _, l := range cat.List(activeTrack) {
			p := store.GetLessonProgress(l.ID)
			if !p.IsFinished() {
				targetLesson = l
				break
			}
		}
		if targetLesson == nil {
			lessons := cat.List(activeTrack)
			if len(lessons) > 0 {
				targetLesson = lessons[0]
			} else {
				return fmt.Errorf("no lessons found in track '%s'", activeTrack)
			}
		}
	}

	_, _ = store.RecordDrill(targetLesson.ID)
	res, err := r.Run(targetLesson)
	if err != nil {
		return err
	}

	modNotice := ""
	if res.Modified {
		modNotice = " (buffer modified)"
	}
	fmt.Printf("\n%s Completed drill for '%s' in %0.1fs%s\n",
		ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
		targetLesson.Title,
		res.Duration.Seconds(),
		modNotice,
	)

	// Single key prompt to mark done
	if ui.IsTerminal() {
		fmt.Printf("Mark '%s' as done? %s: ", targetLesson.ID, ui.Colorize(ui.Bold, "[Enter/y = yes, n = no]"))
		action, _ := ui.ReadSingleKey()
		if action == ui.ActionDone || action == ui.ActionDrill { // Enter / y
			newProg, err := store.MarkDone(targetLesson.ID)
			if err != nil {
				return err
			}
			badge := ui.RenderStatusBadge(newProg.Status)
			fmt.Printf("\n%s Lesson completed! %s (repetitions: %d)\n",
				ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
				badge,
				newProg.RepetitionCount,
			)
		} else {
			fmt.Println()
		}
	}

	return nil
}

func handleDone(cat *catalog.Catalog, store *state.Store, lessonID string) error {
	if lessonID == "" {
		// Find first unfinished lesson in active track
		activeTrack := store.GetActiveTrack()
		for _, l := range cat.List(activeTrack) {
			p := store.GetLessonProgress(l.ID)
			if !p.IsFinished() {
				lessonID = l.ID
				break
			}
		}
		if lessonID == "" {
			return fmt.Errorf("all lessons in track '%s' already marked done; specify an ID: vimbyte done <id>", activeTrack)
		}
	}

	lesson, err := cat.Get(lessonID)
	if err != nil {
		return err
	}

	newProg, err := store.MarkDone(lesson.ID)
	if err != nil {
		return fmt.Errorf("marking lesson done: %w", err)
	}

	badge := ui.RenderStatusBadge(newProg.Status)
	fmt.Printf("%s Marked '%s' (%s) as completed!\n   Status: %s | Repetitions: %d\n",
		ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
		lesson.ID,
		lesson.Title,
		badge,
		newProg.RepetitionCount,
	)
	return nil
}

func handleReview(cat *catalog.Catalog, store *state.Store, r *runner.Runner) error {
	activeTrack := store.GetActiveTrack()
	lessons := cat.List(activeTrack)

	// Find lessons needing repetition (learning or completed but not mastered)
	var candidates []*catalog.Lesson
	for _, l := range lessons {
		p := store.GetLessonProgress(l.ID)
		if p.Status == state.StatusCompleted || p.Status == state.StatusLearning {
			candidates = append(candidates, l)
		}
	}

	if len(candidates) == 0 {
		// Check all tracks
		for _, track := range cat.Tracks() {
			for _, l := range cat.List(track) {
				p := store.GetLessonProgress(l.ID)
				if p.Status == state.StatusCompleted || p.Status == state.StatusLearning {
					candidates = append(candidates, l)
				}
			}
		}
	}

	if len(candidates) == 0 {
		fmt.Printf("%s No lessons currently need review in track '%s'.\nAll finished lessons are either mastered (3+ repetitions) or no lessons have been started yet!\nPractice next lesson with: %s\n",
			ui.Colorize(ui.Bold+ui.BrightYellow, "!"),
			activeTrack,
			ui.Colorize(ui.Bold+ui.BrightWhite, "vimbyte"),
		)
		return nil
	}

	// Pick the candidate with the lowest repetitions
	chosen := candidates[0]
	minReps := store.GetLessonProgress(chosen.ID).RepetitionCount
	for _, l := range candidates {
		reps := store.GetLessonProgress(l.ID).RepetitionCount
		if reps < minReps {
			chosen = l
			minReps = reps
		}
	}

	fmt.Printf("%s Review candidate: %s (#%02d - %s)\n",
		ui.Colorize(ui.Bold+ui.BrightCyan, "▶"),
		chosen.ID,
		chosen.Order,
		chosen.Title,
	)

	return handleInteractiveCard(cat, store, r, chosen.ID)
}

func handleTracks(cat *catalog.Catalog, store *state.Store) error {
	tracks := cat.Tracks()
	activeTrack := store.GetActiveTrack()
	stats := make(map[string]state.TrackStats)

	for _, track := range tracks {
		lessons := cat.List(track)
		st := state.TrackStats{
			Track: track,
			Total: len(lessons),
		}
		for _, l := range lessons {
			p := store.GetLessonProgress(l.ID)
			switch p.Status {
			case state.StatusUnseen:
				st.Unseen++
			case state.StatusLearning:
				st.Learning++
			case state.StatusCompleted:
				st.Completed++
			case state.StatusMastered:
				st.Mastered++
			}
		}
		if st.Total > 0 {
			st.PercentCompleted = float64(st.Completed+st.Mastered) / float64(st.Total) * 100
		}
		stats[track] = st
	}

	fmt.Print(ui.RenderTracksTable(tracks, activeTrack, stats))
	return nil
}

func handleSwitch(cat *catalog.Catalog, store *state.Store, track string) error {
	found := false
	for _, t := range cat.Tracks() {
		if strings.EqualFold(t, track) {
			track = t
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("track '%s' not found. Available tracks: %s",
			track, strings.Join(cat.Tracks(), ", "))
	}

	if err := store.SetActiveTrack(track); err != nil {
		return fmt.Errorf("persisting active track: %w", err)
	}

	fmt.Printf("%s Switched active track to '%s'\n",
		ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
		ui.Colorize(ui.Bold+ui.BrightWhite, track),
	)
	return nil
}

func handleList(cat *catalog.Catalog, store *state.Store, track string) error {
	lessons := cat.List(track)
	if len(lessons) == 0 {
		return fmt.Errorf("no lessons found for track '%s'. Available tracks: %s",
			track, strings.Join(cat.Tracks(), ", "))
	}

	fmt.Print(ui.RenderLessonsList(track, lessons, func(id string) state.LessonProgress {
		return store.GetLessonProgress(id)
	}))
	return nil
}

func handleReset(cat *catalog.Catalog, store *state.Store, track string, all bool) error {
	if all {
		if err := store.ResetAll(); err != nil {
			return fmt.Errorf("resetting all progress: %w", err)
		}
		fmt.Printf("%s Cleared all user progress across all tracks.\n",
			ui.Colorize(ui.Bold+ui.BrightGreen, "✓"))
		return nil
	}

	if track == "" {
		track = store.GetActiveTrack()
	}

	lessons := cat.List(track)
	if len(lessons) == 0 {
		return fmt.Errorf("track '%s' has no lessons to reset", track)
	}

	var ids []string
	for _, l := range lessons {
		ids = append(ids, l.ID)
	}

	if err := store.ResetTrack(ids); err != nil {
		return fmt.Errorf("resetting track progress: %w", err)
	}

	fmt.Printf("%s Reset progress for track '%s' (%d lessons).\n",
		ui.Colorize(ui.Bold+ui.BrightGreen, "✓"),
		track,
		len(ids),
	)
	return nil
}
