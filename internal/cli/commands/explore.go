/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/explore
** File description:
** Registers the `flick explore` subcommand, delegating the interactive group
** explorer to the explore subpackage.
 */

package commands

import (
	"fmt"

	"github.com/Flick-Corp/flick/internal/cli/commands/explore"
	"github.com/Flick-Corp/flick/internal/cli/config"
	"github.com/spf13/cobra"
)

// exploreCmd: the `flick explore` subcommand.
var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Browse your groups and their files interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := config.EnsureCredentials()
		if err != nil {
			return fmt.Errorf("failed to load credentials: %w", err)
		}
		if creds.Token == "" {
			return fmt.Errorf("you are not logged in")
		}
		return explore.Run(creds.Token)
	},
}

// init: Register the explore subcommand.
func init() {
	rootCmd.AddCommand(exploreCmd)
}
