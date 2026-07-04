/*
** FLICK PROJECT, 2026
** flick/internal/api/routes/configure
** File description:
** Configuration route
 */

package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Flick-Corp/flick/internal/api/auth"
	"github.com/Flick-Corp/flick/internal/api/database"
	"github.com/Flick-Corp/flick/internal/api/logging"
	"github.com/Flick-Corp/flick/internal/api/path"
	"github.com/Flick-Corp/flick/internal/api/serverconfig"
)

// ServerConfigHandler: Sends the server config to the web if GET, if POST modifies the server config.
//
// Params:
// - queries (*database.Queries): The queries.
//
// Returns:
// - result1 (http.HandlerFunc): The HTTP handler function.
func ServerConfigHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, status, err := auth.RequireAdmin(r.Context(), queries, auth.GetTokenFromHTTPRequest(r)); err != nil {
			WriteError(w, status, err.Error())
			return
		}

		// GET method
		if r.Method == http.MethodGet {
			data, err := os.ReadFile(filepath.Join(path.GetFlickDir(), "server-config.json"))
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "Failed to read config")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		// POST method
		if r.Method == http.MethodPost {
			var newConf serverconfig.Configuration

			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&newConf); err != nil {
				WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
				return
			}

			if err := newConf.Validate(); err != nil {
				WriteError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
				return
			}

			serverconfig.Conf = newConf
			data, _ := json.MarshalIndent(serverconfig.Conf, "", "  ")
			if err := os.WriteFile(filepath.Join(path.GetFlickDir(), "server-config.json"), data, 0644); err != nil {
				logging.LogInfoError("Cannot save server configuration: %v", err)
				WriteError(w, http.StatusInternalServerError, "Failed to save config")
				return
			}
			logging.LogInfoSuccess("Server configuration updated")
			fmt.Fprint(w, "OK")
			return
		}

		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
