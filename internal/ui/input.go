package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

// Action represents an interactive single-key command triggered by the user.
type Action int

const (
	ActionUnknown Action = iota
	ActionDrill
	ActionDone
	ActionNext
	ActionPrev
	ActionToggleSolution
	ActionQuit
)

// RenderPromptBar returns the interactive single-key prompt string.
func RenderPromptBar(showSolution bool) string {
	solLabel := "[s] Solution"
	if showSolution {
		solLabel = "[s] Hide Solution"
	}

	d := Colorize(Bold+BrightGreen, "[d]") + " Drill in Vim"
	enter := Colorize(Bold+BrightCyan, "[Enter]") + " Mark Done"
	s := Colorize(Bold+BrightYellow, solLabel)
	n := Colorize(Bold+BrightBlue, "[n]") + " Next"
	p := Colorize(Bold+BrightMagenta, "[p]") + " Prev"
	q := Colorize(Bold+BrightBlack, "[q]") + " Quit"

	return fmt.Sprintf(" %s   %s   %s   %s   %s   %s\n", d, enter, s, n, p, q)
}

// IsTerminal returns true if standard input is attached to an interactive terminal.
func IsTerminal() bool {
	term := os.Getenv("TERM")
	if term == "dumb" {
		return false
	}
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&termios)),
		0, 0, 0,
	)
	return err == 0
}

// ReadSingleKey reads one key from the terminal without requiring Enter.
func ReadSingleKey() (Action, error) {
	if !IsTerminal() {
		// Non-terminal input (pipe or redirected input)
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			return ActionQuit, err
		}
		trimmed := strings.ToLower(strings.TrimSpace(line))
		switch trimmed {
		case "d", "drill", "y", "yes":
			return ActionDrill, nil
		case "", "done", "enter":
			return ActionDone, nil
		case "n", "next":
			return ActionNext, nil
		case "p", "prev":
			return ActionPrev, nil
		case "s", "sol", "solution":
			return ActionToggleSolution, nil
		case "q", "quit", "exit":
			return ActionQuit, nil
		default:
			return ActionUnknown, nil
		}
	}

	// Interactive terminal: read raw single byte from /dev/tty
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		var b [1]byte
		n, err := os.Stdin.Read(b[:])
		if err != nil || n == 0 {
			return ActionQuit, err
		}
		return mapByteToAction(b[0]), nil
	}
	defer tty.Close()

	// Capture current stty configuration
	saveCmd := exec.Command("stty", "-F", "/dev/tty", "-g")
	origStateBytes, err := saveCmd.Output()
	origState := strings.TrimSpace(string(origStateBytes))

	if err == nil && origState != "" {
		// Put terminal into cbreak mode (unbuffered, no local echo)
		cbreakCmd := exec.Command("stty", "-F", "/dev/tty", "cbreak", "-echo")
		_ = cbreakCmd.Run()

		defer func() {
			restoreCmd := exec.Command("stty", "-F", "/dev/tty", origState)
			_ = restoreCmd.Run()
		}()
	}

	buf := make([]byte, 8)
	n, err := tty.Read(buf)
	if err != nil || n == 0 {
		if err == io.EOF {
			return ActionQuit, nil
		}
		return ActionQuit, err
	}

	return mapByteToAction(buf[0]), nil
}

func mapByteToAction(b byte) Action {
	switch b {
	case 'd', 'D':
		return ActionDrill
	case '\n', '\r':
		return ActionDone
	case 'n', 'N':
		return ActionNext
	case 'p', 'P':
		return ActionPrev
	case 's', 'S':
		return ActionToggleSolution
	case 'q', 'Q', 3, 27: // 'q', Ctrl-C (\x03), ESC (\x1b)
		return ActionQuit
	default:
		return ActionUnknown
	}
}
