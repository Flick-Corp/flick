/*
** FLICK PROJECT, 2026
** flick/internal/api/logging/logging
** File description:
** Logging go file
 */

package logging

import (
	"fmt"
	"github.com/matteoepitech/flick/internal/api/utils"
	"time"
)

// printLogLabel: Print a label in this format <date> [<title>] <subtitle> > .
//
// Params:
// - title (string): The title.
// - titleColor (string): The title color.
// - subtitle (string): The subtitle title.
// - subtitleColor (string): The subtitle title color.
func printLogLabel(title string, titleColor string, subtitle string, subtitleColor string) {
	now := time.Now().Format("15:04:05")

	fmt.Printf(utils.Dim+"%s"+utils.Reset+" "+
		utils.Gray+"["+utils.Reset+titleColor+utils.Bold+"%s"+utils.Reset+utils.Gray+"]"+utils.Reset+" "+
		subtitleColor+utils.Bold+"%s"+utils.Reset+utils.Gray+" > "+utils.Reset,
		now, title, subtitle)
}

// LogInfoSuccess: Print a success log message.
//
// Params:
// - format (string): The format string (printf style).
// - args (...any): The arguments for formatting.
func LogInfoSuccess(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	printLogLabel("SUCCESS", utils.BrightGreen, "INFO", utils.BrightWhite)
	fmt.Printf("%s\n", msg)
}

// LogInfoError: Print an error log message.
//
// Params:
// - format (string): The format string (printf style).
// - args (...any): The arguments for formatting.
//
// Returns:
// - error: The formatted error.
func LogInfoError(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)

	printLogLabel("ERROR", utils.BrightRed, "INFO", utils.BrightWhite)
	fmt.Printf("%s\n", msg)
	return fmt.Errorf("%s", msg)
}

// LogInfo: Print a log message.
//
// Params:
// - format (string): The format string (printf style).
// - args (...any): The arguments for formatting.
func LogInfo(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	printLogLabel("INFO", utils.BrightBlue, "INFO", utils.BrightWhite)
	fmt.Printf("%s\n", msg)
}
