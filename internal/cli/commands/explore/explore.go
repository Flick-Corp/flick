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

	tea "github.com/charmbracelet/bubbletea"
)

// Run: Launch the interactive group explorer for the given session token.
//
// Params:
// - token (string): The session token used to authenticate API calls.
//
// Returns:
// - result1 (error): An error if the program failed.
func Run(token string) error {
	model := exploreModel{
		token:  token,
		mode:   modeGroups,
		status: "Loading...",
	}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		return fmt.Errorf("failed to run explorer: %w", err)
	}
	return nil
}
