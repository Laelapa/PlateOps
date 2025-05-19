package rt

import (
	"context"
	"time"

	"github.com/Laelapa/GoHome/logging"
	"github.com/Laelapa/PlateOps/internal/repository"
	"github.com/Laelapa/PlateOps/util/typeconvert"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type Params struct {
	userID    uuid.UUID
	token     string
	expiresAt time.Time
	userAgent string
	ipAddress string
}

func RegisterNewToken(queries *repository.Queries, logger *logging.Logger, params Params) error {

	rParams := repository.CreateRefreshTokenParams{
		Token:     params.token,
		UserID:    typeconvert.GoogleUUIDToPgtypeUUID(params.userID),
		ExpiresAt: typeconvert.TimeToPgtypeTimestamp(params.expiresAt),
		UserAgent: params.userAgent,
		IpAddress: params.ipAddress,
	}

	rt, err := queries.CreateRefreshToken(context.Background(), rParams)
	if err != nil {
		logger.LogAppError(
			"Repository error registering refresh token",
			err,
			zap.String("token", params.token),
		)
		return err
	}
	logger.LogAppInfo(
		"Registered new refresh token",
		zap.String("token", rt.Token),
		zap.String("user_id", typeconvert.PgtypeUUIDToString(rt.UserID)),
		zap.Time("expires_at", rt.ExpiresAt.Time),
		zap.String("user_agent", rt.UserAgent),
		zap.String("ip_address", rt.IpAddress),
	)
	return nil
}
