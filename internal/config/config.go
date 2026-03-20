package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	tokenAuthURL  = "https://auth.tesla.com/oauth2/v3/authorize"
	tokenURL      = "https://auth.tesla.com/oauth2/v3/token"
	defaultPort   = "5000"
	defaultScopes = "offline_access vehicle_device_data"
	oauthCallback = "/oauth/callback"
)

type Config struct {
	TeslaClientID       string
	TeslaClientSecret   string
	TeslaRedirectURI    string
	AppBaseURL          string
	TeslaVIN            string
	ShortcutBearerToken string
	TeslaBaseURL        string
	Port                string
	Scopes              []string
	TeslaAuthURL        string
	TeslaTokenURL       string
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		TeslaClientID:       strings.TrimSpace(os.Getenv("TESLA_CLIENT_ID")),
		TeslaClientSecret:   strings.TrimSpace(os.Getenv("TESLA_CLIENT_SECRET")),
		AppBaseURL:          strings.TrimSpace(os.Getenv("APP_BASE_URL")),
		TeslaVIN:            strings.TrimSpace(os.Getenv("TESLA_VIN")),
		ShortcutBearerToken: strings.TrimSpace(os.Getenv("SHORTCUT_BEARER_TOKEN")),
		TeslaBaseURL:        strings.TrimSpace(os.Getenv("TESLA_BASE_URL")),
		Port:                strings.TrimSpace(os.Getenv("PORT")),
		TeslaAuthURL:        tokenAuthURL,
		TeslaTokenURL:       tokenURL,
	}

	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	if cfg.AppBaseURL != "" && !strings.HasPrefix(cfg.AppBaseURL, "http://") && !strings.HasPrefix(cfg.AppBaseURL, "https://") {
		cfg.AppBaseURL = "https://" + cfg.AppBaseURL
	}

	scopeValue := strings.TrimSpace(os.Getenv("TESLA_SCOPES"))
	if scopeValue == "" {
		scopeValue = defaultScopes
	}
	cfg.Scopes = strings.Fields(scopeValue)
	cfg.TeslaRedirectURI = strings.TrimRight(cfg.AppBaseURL, "/") + oauthCallback

	missing := make([]string, 0, 6)
	if cfg.TeslaClientID == "" {
		missing = append(missing, "TESLA_CLIENT_ID")
	}
	if cfg.TeslaClientSecret == "" {
		missing = append(missing, "TESLA_CLIENT_SECRET")
	}
	if cfg.AppBaseURL == "" {
		missing = append(missing, "APP_BASE_URL")
	}
	if cfg.TeslaVIN == "" {
		missing = append(missing, "TESLA_VIN")
	}
	if cfg.ShortcutBearerToken == "" {
		missing = append(missing, "SHORTCUT_BEARER_TOKEN")
	}
	if cfg.TeslaBaseURL == "" {
		missing = append(missing, "TESLA_BASE_URL")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}
