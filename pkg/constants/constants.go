package constants

import "github.com/wafi04/vazzuniversebackend/pkg/config"

var (
	JWT_SECRET = []byte(config.LoadEnv("JWT_SECRET"))
)
