package configuration

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
	"log/slog"
	"os"
)

type Config struct {
	DBUrl  string `mapstructure:"AWS_ENDPOINT_URL_DYNAMODB"`
	Secret string `mapstructure:"SECRET"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigFile("app.env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("error", err)
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		slog.Error("error", err)
		return cfg, err
	}

	os.Setenv("AWS_ENDPOINT_URL_DYNAMODB", cfg.DBUrl)

	if len(cfg.Secret) == 0 {
		secret, err := generateKey()
		if err != nil {
			slog.Error("error", err)
			return cfg, err
		}

		cfg.Secret = secret
		viper.Set("AWS_ENDPOINT_URL_DYNAMODB", "http://localhost:8000")
		viper.Set("SECRET", secret)

		err = viper.WriteConfigAs(viper.ConfigFileUsed())
		if err != nil {
			slog.Error("error", err)
			return cfg, err
		}
	}

	return cfg, nil
}

func generateKey() (string, error) {
	fmt.Println("Please Enter Your Password: ")

	var pw string

	fmt.Scanln(&pw)

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(pw), salt, 1, 32, sha256.New)

	return hex.EncodeToString(key), nil
}
