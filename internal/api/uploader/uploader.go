/*
** FLICK PROJECT, 2026
** flick/internal/api/uploader/uploader
** File description:
** Shared uploader identification, used by the quota route and the tus package
 */

package uploader

import (
	"context"
	"fmt"

	"github.com/Flick-Corp/flick/internal/api/database"
	"github.com/jackc/pgx/v5/pgtype"
)

// Uploader: The resolved identity behind an upload request.
type Uploader struct {
	ID          pgtype.UUID
	IsAnonymous bool
	Blocked     bool
}

// ResolveUploaderByID: Resolve an uploader from its raw id against the database,
// the single source of truth shared by the quota route and the tus package
// (whose hooks expose headers, not an *http.Request).
//
// Params:
// - ctx (context.Context): The request context.
// - queries (*database.Queries): The database queries.
// - uploaderID (string): The raw X-Flick-User-ID value to validate.
//
// Returns:
// - result1 (Uploader): The resolved uploader; the zero value on error.
// - result2 (error): An error if the id is missing, invalid or unknown.
func ResolveUploaderByID(ctx context.Context, queries *database.Queries, uploaderID string) (Uploader, error) {
	if uploaderID == "" {
		return Uploader{}, fmt.Errorf("missing uploader id")
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(uploaderID); err != nil {
		return Uploader{}, fmt.Errorf("invalid user id %q: %w", uploaderID, err)
	}

	if _, err := queries.GetAnonymousUserByID(ctx, userUUID); err == nil {
		return Uploader{
			ID:          userUUID,
			IsAnonymous: true,
			Blocked:     false,
		}, nil
	}

	if user, err := queries.GetUserByID(ctx, userUUID); err == nil {
		return Uploader{
			ID:          userUUID,
			IsAnonymous: false,
			Blocked:     user.Blocked,
		}, nil
	}

	return Uploader{}, fmt.Errorf("unknown user id %q", uploaderID)
}
