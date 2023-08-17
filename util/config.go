package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	DatabaseUrl          string        `mapstructure:"DATABASE_URL"`
	OriginAllowed        string        `mapstructure:"ORIGIN_ALLOWED"`
	B2ApplicationKeyId   string        `mapstructure:"B2_APPLICATION_KEY_ID"`
	BucketId             string        `mapstructure:"BUCKET_ID"`
	BucketName           string        `mapstructure:"BUCKET_NAME"`
	B2ApplicationKey     string        `mapstructure:"B2_APPLICATION_KEY"`
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	Domain               string        `mapstructure:"DOMAIN"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
