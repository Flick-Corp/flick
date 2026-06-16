/*
** FLICK PROJECT, 2026
** flick/internal/api/routes/account/oauth/device_token
** File description:
** Device authorization flow (CLI polling)
 */

package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/matteoepitech/flick/internal/api/database"
	"github.com/matteoepitech/flick/internal/api/routes"
)

// DeviceTokenRequest: The JSON body request.
type DeviceTokenRequest struct {
	DeviceCode string `json:"device_code" validate:"required"`
}

// DeviceTokenResponse: The JSON body reponse when the device is approved.
type DeviceTokenResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

// DeviceTokenPending: The JSON body returned while the device is still waiting.
type DeviceTokenPendingResponse struct {
	Status string `json:"status"`
}

// DeviceTokenHandler: Check if the user_code has been approved or not.
//
// Params:
// - queries (*database.Queries): The database queries.
//
// Returns:
// - result1 (http.HandlerFunc): The handler function.
func DeviceTokenHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			routes.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var request DeviceTokenRequest
		var validate = validator.New()

		if err := decoder.Decode(&request); err != nil {
			routes.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
			return
		}
		if err := validate.Struct(request); err != nil {
			routes.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		auth, err := queries.GetDeviceAuthorizationByDeviceCode(r.Context(), request.DeviceCode)
		if err != nil {
			routes.WriteError(w, http.StatusNotFound, "Invalid device code: "+err.Error())
			return
		}

		switch auth.Status {
		case database.OauthStatusPending:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(DeviceTokenPendingResponse{Status: "pending"})
			return
		case database.OauthStatusDenied:
			routes.WriteError(w, http.StatusForbidden, "Authorization denied")
			return
		case database.OauthStatusApproved:
			if auth.SessionToken == nil {
				routes.WriteError(w, http.StatusInternalServerError, "Approved device has no session token")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(DeviceTokenResponse{
				Token:  *auth.SessionToken,
				UserID: auth.UserID.String(),
			})
			return
		default:
			routes.WriteError(w, http.StatusInternalServerError, "Unknown authorization status")
			return
		}
	}
}
