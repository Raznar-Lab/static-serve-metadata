package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
)

type WebConfig struct {
	MaxAge         uint64   `env:"WEB_MAX_AGE" envDefault:"3600"`
	IndexFile      string   `env:"WEB_INDEX_FILE" envDefault:"index.html"`
	DataPath       string   `env:"WEB_DATA_PATH" envDefault:"web"`
	Port           string   `env:"WEB_PORT" envDefault:"8080"`
	Bind           string   `env:"WEB_BIND" envDefault:"0.0.0.0"`
	TrustedProxies []string `env:"WEB_TRUSTED_PROXIES" envSeparator:"," envDefault:""`
	ProxyHeader    string   `env:"WEB_PROXY_HEADER" envDefault:"X-Forwarded-For"`
}

type LanguageConfig struct {
	Country []string `json:"country"`
	Default bool     `json:"default"`
	Prefix  string   `json:"prefix"`
}

type SEOConfig struct {
	GeoHeader string                    `env:"SEO_GEO_HEADER"`
	DataPath  string                    `env:"SEO_DATA_PATH"`
	Languages map[string]LanguageConfig `env:"-"`
}

type AppConfig struct {
	WebConfig WebConfig
	SeoConfig SEOConfig
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{}

	if err := env.Parse(&cfg.WebConfig); err != nil {
		return nil, err
	}
	if err := env.Parse(&cfg.SeoConfig); err != nil {
		return nil, err
	}

	cfg.SeoConfig.Languages = map[string]LanguageConfig{}
	langList := strings.Split(os.Getenv("SEO_LANGUAGES"), ",")
	for _, lang := range langList {
		langKey := strings.ToUpper(strings.TrimSpace(lang))
		envKey := func(suffix string) string {
			return os.Getenv(fmt.Sprintf("SEO_LANG_%s_%s", langKey, suffix))
		}

		cfg.SeoConfig.Languages[strings.ToLower(langKey)] = LanguageConfig{
			Country: splitAndTrim(envKey("COUNTRY")),
			Default: strings.ToLower(envKey("DEFAULT")) == "true",
			Prefix:  envKey("PREFIX"),
		}

		fmt.Println(cfg.SeoConfig)
	}

	fmt.Println(cfg)
	return cfg, nil
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
