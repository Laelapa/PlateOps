package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util"
	"github.com/Laelapa/PlateOps/util/typeconvert"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TODO: Possibly relocate DeepError to a more general package if it ends up being needed elsewhere.

var (
	ErrTokenNotExist       = errors.New("token does not exist")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenRevoked        = errors.New("token revoked")
	ErrTokenLoggedOut      = errors.New("token logged out")
	ErrTokenMismatch       = errors.New("token does not belong to this user")
	ErrFingerprintMismatch = errors.New("token fingerprint mismatch")
	ErrDatabaseFailure     = errors.New("database operation failed")
)

// VerifyRefreshToken checks if the refresh token is valid and matches the provided user & fingerprints.
// If all checks pass, nil is returned, indicating the token is valid.
// If any of these checks fail, an appropriate error is returned.
// This function does not provide any form of sanitization or validation of the input parameters.
// It is assumed that the caller has already validated the input before calling it.
// TODO: Consider adding one more parameter for clock drift lenience.
// TODO: Consider dropping userID from the request parameters to simplify request requirements.
func VerifyRefreshToken(
	ctx context.Context, // the request's context
	q *repository.Queries,
	refreshToken string,
	userID uuid.UUID,
	userAgent string,
	ipAddress string,
	requestTime time.Time,
) error {

	// Fetch the token from the database.
	dbToken, err := q.GetToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTokenNotExist
		}

		return &util.DeepError{
			BusinessErr:  ErrDatabaseFailure,
			TechnicalErr: err,
		}
	}

	// Check for expiration.
	if dbToken.ExpiresAt.Time.Before(requestTime) {
		return ErrTokenExpired
	}
	// Check for revocation.
	if dbToken.RevokedAt.Valid {
		return ErrTokenRevoked
	}
	// Check for logged out status.
	if dbToken.LoggedOutAt.Valid {
		return ErrTokenLoggedOut
	}

	// Cross-check request info with the database.
	if userID != typeconvert.PgtypeUUIDToGoogleUUID(dbToken.UserID) {
		return ErrTokenMismatch
	}

	// Cross-check fingerprints.
	if dbToken.UserAgent != userAgent {
		return ErrFingerprintMismatch
	}

	if dbToken.IpAddress != ipAddress {
		return ErrFingerprintMismatch
	}

	return nil
}
