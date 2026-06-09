/*
** FLICK PROJECT, 2026
** flick/internal/api/routes/stats
** File description:
** Stats route handler
 */

package routes

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/matteoepitech/flick/internal/api/code"
)

// Variables that will contain the total uploads/downloads.
var (
	totalUploads   atomic.Uint64
	totalDownloads atomic.Uint64
)

// IncUploads: Increment the total uploads counter.
func IncUploads() {
	totalUploads.Add(1)
}

// IncDownloads: Increment the total downloads counter.
func IncDownloads() {
	totalDownloads.Add(1)
}

// Uploads: Read the total uploads counter.
//
// Returns:
// - result1 (uint64): Total uploads.
func Uploads() uint64 {
	return totalUploads.Load()
}

// Downloads: Read the total downloads counter.
//
// Retruns:
// - result1(uint64): Total downloads.
func Downloads() uint64 {
	return totalDownloads.Load()
}

// SendStats: Build the stats handler.
//
// Returns:
// - http.HandlerFunc: The handler function.
func SendStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		payload := map[string]any{
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
			"activeCodes":    code.Cache.ItemCount(),
			"totalUploads":   Uploads(),
			"totalDownloads": Downloads(),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}
}
