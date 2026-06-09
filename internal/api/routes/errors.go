/*
** FLICK PROJECT, 2026
** flick/internal/api/routes/errors
** File description:
** Shared JSON error responses so clients (web UI especially)
 */

package routes

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse: The JSON body returned for any error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError: Writes a JSON error response with the given status code.
//
// Params:
// - w (http.ResponseWriter): The response writer.
// - status (int): The HTTP status code.
// - message (string): The string error message.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
