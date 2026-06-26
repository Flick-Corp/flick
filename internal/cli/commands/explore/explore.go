/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/explore/explore
** File description:
** Entry point for the interactive group explorer: wires the Bubble Tea program
** that browses groups, navigates folders, downloads and manages files.
 */

package explore

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

// getTermSize: query the terminal dimensions at startup so that the view can
// correctly limit visible items even before the first WindowSizeMsg arrives.
//
// Returns:
// - h (int): Terminal height in rows (24 if unknown).
// - w (int): Terminal width  in cols  (80 if unknown).
func getTermSize() (h, w int) {
	h, w = 24, 80 // safe fallback
	if fd := int(os.Stdin.Fd()); term.IsTerminal(fd) {
		tw, th, err := term.GetSize(fd)
		if err == nil {
			w, h = tw, th
		}
	}
	return
}

// Run: Launch the interactive group explorer for the given session token.
//
// Params:
// - token (string): The session token used to authenticate API calls.
//
// Returns:
// - result1 (error): An error if the program failed.
func Run(token string) error {
	h, w := getTermSize()
	model := exploreModel{
		token:  token,
		mode:   modeGroups,
		status: "Loading...",
		height: h,
		width:  w,
	}
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("failed to run explorer: %w", err)
	}
	return nil
}
