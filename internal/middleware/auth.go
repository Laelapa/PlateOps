package middleware

import (
	"net/http"

	"github.com/Laelapa/PlateOps/auth/tokenauthority"
	"github.com/Laelapa/PlateOps/util/ctxutils"
	"github.com/Laelapa/PlateOps/util/net"
	"github.com/Laelapa/PlateOps/util/typeconvert"

	"github.com/Laelapa/GoHome/logging"
	"go.uber.org/zap"
)

func AuthenticateWithJWT(
	tokenAuthority *tokenauthority.TokenAuthority,
	logger *logging.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check for the existence of the Authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.LogRequestWarn("Unauthorized request: missing Authorization header", r)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check for the Bearer token format of the Authorization header and extract the JWT.
			tokenString, err := net.StripBearer(authHeader)
			if err != nil {
				logger.LogRequestWarn("Unauthorized request: invalid Authorization header", r)
				logger.LogAppWarn("Unauthorized request: invalid Authorization header", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate the JWT and extract the userID.
			subject, err := tokenAuthority.ValidateJWT(tokenString)
			if err != nil {
				logger.LogRequestWarn("Unauthorized request: invalid JWT", r)
				logger.LogAppWarn("Unauthorized request: invalid JWT", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Convert the subject to a pgtype.UUID.
			pgUUID := typeconvert.StringToPgtypeUUID(subject)
			if !pgUUID.Valid {
				logger.LogRequestWarn("Unauthorized request: invalid user ID in JWT", r)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Store the userID in the request context.
			ctx := ctxutils.SetUserIDInContext(r.Context(), pgUUID)

			logger.LogRequestInfo("Request authenticated with JWT", r)
			logger.LogAppInfo(
				"Request authenticated with JWT",
				zap.String("userID", typeconvert.PgtypeUUIDToString(pgUUID)),
			)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
