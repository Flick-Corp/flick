/*
** FLICK PROJECT, 2026
** flick/internal/utils/archive/archive
** File description:
** Shared zip-archiving helpers used to bundle local files/folders into a single
** archive before uploading them to the server (used by the upload command and
** the interactive explorer).
 */

package archive

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// archiveRoot: Choose the base against which src's entries are named inside the
// archive. A local relative path keeps its full structure, so "dir1/a.txt" and
// "dir2/a.txt" stay distinct.
//
// Params:
// - src (string): The file or directory path passed on the command line.
//
// Returns:
// - result1 (string): The base directory to compute archive-relative paths from.
func archiveRoot(src string) string {
	if filepath.IsLocal(src) {
		return "."
	}
	return filepath.Dir(src)
}

// ToTemp: Build a single zip archive of every src into a temporary file
// and return its path.
//
// Params:
// - srcs ([]string): The files and/or directories to archive together.
//
// Returns:
// - result1 (string): The path to the temporary zip file.
// - result2 (error): An error if occured.
func ToTemp(srcs []string) (string, error) {
	tmp, err := os.CreateTemp("", "flick-upload-*.zip")
	if err != nil {
		return "", fmt.Errorf("Failure: Cannot create temp archive: %w", err)
	}
	defer tmp.Close()

	zw := zip.NewWriter(tmp)
	for _, src := range srcs {
		root := archiveRoot(src)
		if err := addToZip(zw, root, src); err != nil {
			zw.Close()
			os.Remove(tmp.Name())
			return "", fmt.Errorf("Failure: Cannot build archive: %w", err)
		}
	}

	if err := zw.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("Failure: Cannot finalize archive: %w", err)
	}
	return tmp.Name(), nil
}

// addToZip: Add file(s) to the zip archive, storing each one with relative path.
//
// Params:
// - zw (*zip.Writer): The zip writer to add entries to.
// - root (string): The base directory used to compute relative paths.
// - path (string): The current file or directory to add.
//
// Returns:
// - result1 (error): An error if occured.
func addToZip(zw *zip.Writer, root string, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := addToZip(zw, root, filepath.Join(path, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}

	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}

	w, err := zw.Create(filepath.ToSlash(relPath))
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

// RandomName: A random uuid-style name for the uploaded archive.
//
// Returns:
// - result1 (string): The "<uuid>.zip" archive name.
func RandomName() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("flick-%d.zip", time.Now().UnixNano())
	}

	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x.zip", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
