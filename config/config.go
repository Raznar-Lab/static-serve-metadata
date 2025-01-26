package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type WebConfig struct {
	MaxAge         uint64   `yaml:"max_age"`
	IndexFile      string   `yaml:"index_file"`
	DataPath       string   `yaml:"data_path"`
	Port           string   `yaml:"port"`
	Bind           string   `yaml:"bind"`
	TrustedProxies []string `yaml:"trusted_proxies"`
	ProxyHeader    string   `yaml:"proxy_header"`
}

type LanguageConfig struct {
	Country []string `yaml:"country"`
	Default bool     `yaml:"default"`
}

type SEOConfig struct {
	GeoHeader string `yaml:"geo_header"`
	Languages map[string]LanguageConfig `yaml:"languages"`
	DataPath  string                    `yaml:"data_path"`
}

type AppConfig struct {
	filepath  string
	WebConfig WebConfig `yaml:"web"`
	SeoConfig SEOConfig `yaml:"seo"`
}

var defaultAppConfig = &AppConfig{}

func (appConfig *AppConfig) Load() (err error) {
	if !appConfig.IsExists() {
		defaultAppConfig.filepath = appConfig.filepath
		appConfig = defaultAppConfig
		err = appConfig.Save()
		if err != nil {
			return
		}
	}

	configData, err := os.ReadFile(appConfig.filepath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(configData, appConfig)
	if err != nil {
		return
	}

	err = appConfig.Parse()
	return
}

// parse environment variables
func (appConfig *AppConfig) Parse() (err error) {
	// add parser logic

	return
}

func (appConfig *AppConfig) IsExists() bool {
	// generate default
	_, err := os.Stat(appConfig.filepath)
	return err == nil || os.IsExist(err)
}

func (appConfig *AppConfig) Save() (err error) {
	// write default
	configData, err := yaml.Marshal(defaultAppConfig)
	os.WriteFile(appConfig.filepath, configData, 0644)

	return
}

func New(filepath string, load bool) (appConfig *AppConfig, err error) {
	appConfig = &AppConfig{
		filepath: filepath,
	}

	if load {
		err = appConfig.Load()
	}

	return
}
