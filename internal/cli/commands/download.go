/*
** FLICK PROJECT, 2026
** flick/internal/cli/commands/download
** File description:
** Download flick source
 */

package commands

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matteoepitech/flick/internal/cli/config"
	"github.com/matteoepitech/flick/internal/cli/network"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// doDownloadRequest: Do the download request on the server.
//
// Params:
// - req (*http.Request): The request HTTP.
//
// Returns:
// - result1 (error): An error occured.
func doDownloadRequest(req *http.Request) error {
	resp, err := network.SharedClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failure: Cannot access the server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failure: %s", serverErrorMessage(body, resp.Status))
	}

	totalSizeStr := resp.Header.Get("X-Total-Size")
	totalSize, _ := strconv.ParseInt(totalSizeStr, 10, 64)
	if totalSize <= 0 {
		totalSize = -1
	}
	bar := progressbar.DefaultBytes(totalSize, "Downloading")

	contentType := resp.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("Failure: invalid Content-Type header: %w", err)
	}

	boundary, ok := params["boundary"]
	if !ok {
		return fmt.Errorf("Failure: missing multipart boundary in response")
	}

	reader := multipart.NewReader(resp.Body, boundary)

	// Uploads are always stored as a zip archive, so every "file" part is one
	// archive that we extract into the current directory.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if part.FormName() != "file" {
			continue
		}

		if err := downloadArchive(part, bar); err != nil {
			return err
		}
	}

	return nil
}

// downloadArchive: Buffer the archive to a temp file then extract into the current directory.
//
// Params:
// - part (io.Reader): The multipart file stream carrying the zip.
// - bar (*progressbar.ProgressBar): The progress bar to feed while downloading.
//
// Returns:
// - result1 (error): An error if occured.
func downloadArchive(part io.Reader, bar *progressbar.ProgressBar) error {
	tmp, err := os.CreateTemp("", "flick-download-*.zip")
	if err != nil {
		return fmt.Errorf("Failure: Cannot create temp archive: %w", err)
	}
	defer os.Remove(tmp.Name())

	proxyReader := io.TeeReader(part, bar)
	if _, err := io.Copy(tmp, proxyReader); err != nil {
		tmp.Close()
		return fmt.Errorf("Failure: Cannot download the archive: %w", err)
	}
	tmp.Close()

	return extractZip(tmp.Name(), ".")
}

// extractZip: Extract a zip archive into dest.
//
// Params:
// - zipPath (string): The path to the zip file on disk.
// - dest (string): The destination directory.
//
// Returns:
// - result1 (error): An error if occured.
func extractZip(zipPath string, dest string) error {
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("Failure: Cannot open the archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(absDest, f.Name)

		if target != absDest && !strings.HasPrefix(target, absDest+string(os.PathSeparator)) {
			return fmt.Errorf("Failure: unsafe path in archive: %q", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if err := writeZipEntry(f, target); err != nil {
			return err
		}
	}
	return nil
}

// writeZipEntry: Copy a single zip entry to target on disk.
//
// Params:
// - f (*zip.File): The zip entry to extract.
// - target (string): The destination path on disk.
//
// Returns:
// - result1 (error): An error if occured.
func writeZipEntry(f *zip.File, target string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// RunDownload: Run the download command.
//
// Params:
// - cmd (*cobra.Command): The command.
//
// Returns:
// - result1 (error): An error if occured.
func RunDownload(cmd *cobra.Command, args []string) error {
	var code string

	fmt.Printf("Specify the code: ")
	fmt.Scan(&code)
	fmt.Printf("Searching the code %s...\n", code)

	body := &bytes.Buffer{}

	req, err := http.NewRequest("GET", config.Conf.APIBaseURL()+"/download?code="+code, body)
	if err != nil {
		return fmt.Errorf("Failure: Cannot create the request for the server.")
	}

	return doDownloadRequest(req)
}
