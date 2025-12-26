// Package config provides configuration management for the application.
package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServicePort     string
	JwtSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() (*Config, error) {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "5050"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JwtSecet is required")
	}

	accessTTLStr := os.Getenv("ACCESS_TOKEN_TTL_MINUTES")
	if accessTTLStr == "" {
		accessTTLStr = "15"
	}

	accessTTLMinutes, err := strconv.Atoi(accessTTLStr)
	if err != nil || accessTTLMinutes <= 0 {
		return nil, errors.New("invalid access token TTL minutes")
	}

	refreshTTLStr := os.Getenv("REFRESH_TOKEN_TTL_MINUTES")
	if refreshTTLStr == "" {
		refreshTTLStr = "10080" // 7 days
	}

	refreshTTLMinutes, err := strconv.Atoi(refreshTTLStr)
	if err != nil || refreshTTLMinutes <= 0 {
		return nil, errors.New("invalid REFRESH_TOKEN_TTL_MINUTES")
	}

	return &Config{
		ServicePort:     port,
		JwtSecret:       jwtSecret,
		AccessTokenTTL:  time.Duration(accessTTLMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshTTLMinutes) * time.Minute,
	}, nil
}
