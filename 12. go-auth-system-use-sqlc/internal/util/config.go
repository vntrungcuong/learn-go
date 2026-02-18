// Config utils
package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APIServerPort        int           `mapstructure:"API_SERVER_PORT"`
	PostgreSQLConn       string        `mapstructure:"POSTGRES_CONN"`
	RedisAddr            string        `mapstructure:"REDIS_ADDR"`
	JWTSecret            string        `mapstructure:"JWT_SECRET"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
