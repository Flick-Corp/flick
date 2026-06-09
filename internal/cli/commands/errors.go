/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/errors
** File description:
** Shared helpers to surface server error messages in the CLI.
 */

package commands

import (
	"encoding/json"
	"strings"
)

// serverErrorMessage: Extracts the human-readable error from a server response.
//
// Params:
// - body ([]byte): The raw response body.
// - status (string): The HTTP status line, used as fallback.
//
// Returns:
// - result1 (string): The error message to display.
func serverErrorMessage(body []byte, status string) string {
	var parsed struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error != "" {
		return parsed.Error
	}
	if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
		return trimmed
	}
	return status
}
