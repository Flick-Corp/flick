/*
** FLICK PROJECT, 2026
** flick/internal/api/utils/data
** File description:
** Data utils source file
 */

package data

import (
	"os"
	"path/filepath"

	"github.com/matteoepitech/flick/internal/api/path"
)

// DeleteDataDirWithCode: Delete the data directory of a code from disk.
//
// Params:
// - code (string): The code to delete.
//
// Returns:
// - result1 (error): Error if removing the data directory fails.
func DeleteDataDirWithCode(code string) error {
	dataDir := filepath.Join(path.GetDataDir(), code)
	return os.RemoveAll(dataDir)
}
