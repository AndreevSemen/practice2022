package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	DBUri                  string
	JWTSecret              string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	RefreshTokenLength     int
	ServerPort             int
	HTTPSMode              bool
	HTTPSCrtFile           string
	HTTPSKeyFile           string
}

func GetConfig() (Config, error) {
	var (
		cfg Config
		ok  bool
		err error
	)

	cfg.DBUri, ok = os.LookupEnv("DB_URI")
	if !ok {
		err := errors.New("DB_URI is not set")
		return Config{}, err
	}

	serverPort, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		err := errors.New("DB_URI is not set")
		return Config{}, err
	}

	cfg.ServerPort, err = strconv.Atoi(serverPort)
	if err != nil {
		err = errors.Wrap(err, "parse SERVER_PORT")
		return Config{}, err
	}

	cfg.JWTSecret, ok = os.LookupEnv("JWT_SECRET")
	if !ok {
		err := errors.New("JWT_SECRET is not set")
		return Config{}, err
	}

	accessExpiration, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRATION")
	if !ok {
		err = errors.Wrap(err, "ACCESS_TOKEN_EXPIRATION is not set")
		return Config{}, err
	}

	cfg.AccessTokenExpiration, err = time.ParseDuration(accessExpiration)
	if err != nil {
		err = errors.Wrap(err, "parse ACCESS_TOKEN_EXPIRATION")
		return Config{}, err
	}

	refreshExpiration, ok := os.LookupEnv("REFRESH_TOKEN_EXPIRATION")
	if !ok {
		err = errors.Wrap(err, "REFRESH_TOKEN_EXPIRATION is not set")
		return Config{}, err
	}

	cfg.RefreshTokenExpiration, err = time.ParseDuration(refreshExpiration)
	if err != nil {
		err = errors.Wrap(err, "parse REFRESH_TOKEN_EXPIRATION")
		return Config{}, err
	}

	_, ok = os.LookupEnv("HTTPS_MODE")
	if ok {
		cfg.HTTPSMode = true

		cfg.HTTPSCrtFile, ok = os.LookupEnv("HTTPS_CRT")
		if !ok {
			err := errors.New("HTTPS_CRT is not set")
			return Config{}, err
		}

		cfg.HTTPSKeyFile, ok = os.LookupEnv("HTTPS_KEY")
		if !ok {
			err := errors.New("HTTPS_KEY is not set")
			return Config{}, err
		}
	}

	return cfg, nil
}
