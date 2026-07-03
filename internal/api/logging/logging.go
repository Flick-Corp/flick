/*
** FLICK PROJECT, 2026
** flick/internal/api/logging/logging
** File description:
** Logging go file
 */

package logging

import (
	"fmt"
	"time"

	"github.com/Flick-Corp/flick/internal/utils/colors"
)

// printLogLabel: Print a label in this format <date> [<title>] > .
//
// Params:
// - title (string): The title.
// - titleColor (string): The title color.
func printLogLabel(title string, titleColor string) {
	now := time.Now().Format("15:04:05")

	fmt.Printf(colors.Dim+"%s"+colors.Reset+" ", now)
	fmt.Printf(colors.Gray + "[" + colors.Reset)
	fmt.Printf(titleColor+colors.Bold+"%s"+colors.Reset, title)
	fmt.Printf(colors.Gray + "]" + colors.Reset)
	fmt.Printf(colors.Gray + " > " + colors.Reset)
}

// LogInfoSuccess: Print a success log message.
//
// Params:
// - format (string): The format string (printf style).
// - args (...any): The arguments for formatting.
func LogInfoSuccess(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	printLogLabel("SUCCESS", colors.BrightGreen)
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

	printLogLabel("ERROR", colors.BrightRed)
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

	printLogLabel("INFO", colors.BrightBlue)
	fmt.Printf("%s\n", msg)
}
