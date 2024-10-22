package logging
import (
	"go.uber.org/zap"
)

func NewLogger(env string) (*zap.SugaredLogger, error) {

	logger, err := zap.NewDevelopment()

	if env == "prod" || env == "production" {
		logger, err = zap.NewProduction() // TODO: setup production output
	}

	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	return logger.Sugar(), nil
}
