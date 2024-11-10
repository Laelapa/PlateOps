package auth

type JWTConfig struct {
	JWTSecret        string
	Issuer           string
	ExpiresInSeconds int
}

type RTConfig struct {
	TokenSizeInBytes int
	ExpiresInHours int
}

type Authenticator struct {
	jwtConfig JWTConfig
	rtConfig RTConfig
}

func New(
	jwtConfig JWTConfig,
	rtConfig RTConfig,
) (*Authenticator, error) {

	// TODO: verify configs

	return &Authenticator{
		jwtConfig: jwtConfig,
		rtConfig: rtConfig,
	}, nil
}
