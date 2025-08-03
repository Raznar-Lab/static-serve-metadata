package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"raznar.id/static-serve-metadata/config"
)

type Metadata struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type SEOData struct {
	URL      string     `json:"url"`
	Default  bool       `json:"default"`
	Template bool       `json:"template"`
	Metadata []Metadata `json:"metadata"`
	Title    string     `json:"title"`
	Prefix   string     `json:"prefix"`
}

type GroupSEO struct {
	SeoContents         []SEOData
	SeoDefaultContents  SEOData
	SeoTemplateContents SEOData
}

func (g GroupSEO) GetDataByURL(url string) SEOData {
	url = strings.TrimSuffix(url, "/")
	for _, data := range g.SeoContents {
		seoURL := strings.TrimSuffix(path.Join(data.Prefix, data.URL), "/")
		if seoURL == url {
			return data
		}
	}
	return g.SeoDefaultContents
}

func (g GroupSEO) GetTitle() string {
	return g.SeoTemplateContents.Title
}

func (s SEOData) GetTitle() string {
	return s.Title
}

func (s SEOData) CollectMetadataString() string {
	var tags []string
	for _, m := range s.Metadata {
		tags = append(tags, m.ToHTML())
	}
	return strings.Join(tags, "\n    ")
}

func (m Metadata) ToHTML() string {
	return fmt.Sprintf(`<meta name="%s" content="%s">`, m.Tag, m.Content)
}

func handleWeb(ac *config.AppConfig, mapSEO map[string]GroupSEO, fileContent []byte) fiber.Handler {
	defaultLang := getDefaultLang(ac)
	templateHTML := string(fileContent)

	return func(c *fiber.Ctx) error {
		lang := getLangCode(ac, c.Get(ac.SeoConfig.GeoHeader), c.Path())
		if lang == "" {
			lang = defaultLang
		}

		seoGroup := mapSEO[lang]
		seo := seoGroup.GetDataByURL(c.Path())

		// Compose metadata
		metadata := seoGroup.SeoTemplateContents.CollectMetadataString() + "\n" + seo.CollectMetadataString()

		// Replace placeholders
		content := strings.Replace(templateHTML, "<!-- seo header -->", metadata, 1)
		title := seo.GetTitle()
		if title == "" {
			title = seoGroup.GetTitle()
		}
		content = strings.Replace(content, "<!-- title -->", title, 1)

		c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", ac.WebConfig.MaxAge))
		c.Set("Content-Type", "text/html")
		return c.SendString(content)
	}
}

func getLangCode(ac *config.AppConfig, country, path string) string {
	for lang, config := range ac.SeoConfig.Languages {
		if contains(config.Country, country) || strings.HasPrefix(path, config.Prefix) {
			return lang
		}
	}
	return ""
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func getDefaultLang(ac *config.AppConfig) string {
	for lang, cfg := range ac.SeoConfig.Languages {
		if cfg.Default {
			return lang
		}
	}
	return "default"
}

func loadSEO(ac *config.AppConfig) (map[string]GroupSEO, error) {
	result := make(map[string]GroupSEO)

	for lang := range ac.SeoConfig.Languages {
		dir := path.Join(ac.SeoConfig.DataPath, lang)
		data, err := loadSeoContents(dir)
		if err != nil {
			return nil, fmt.Errorf("loading SEO data for %s: %w", lang, err)
		}

		group := GroupSEO{SeoContents: data}
		for _, seo := range data {
			if seo.Default && group.SeoDefaultContents.URL == "" {
				group.SeoDefaultContents = seo
			}
			if seo.Template && group.SeoTemplateContents.URL == "" {
				group.SeoTemplateContents = seo
			}
		}
		result[lang] = group
	}

	return result, nil
}

func loadSeoContents(dir string) ([]SEOData, error) {
	var contents []SEOData
	err := filepath.WalkDir(dir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		raw, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		var parsed []SEOData
		if err := json.Unmarshal(raw, &parsed); err == nil {
			contents = append(contents, parsed...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking dir %s: %w", dir, err)
	}
	return contents, nil
}

func RunWeb(ac *config.AppConfig) error {
	webConf := ac.WebConfig

	fileContent, err := os.ReadFile(path.Join(webConf.DataPath, webConf.IndexFile))
	if err != nil {
		return err
	}

	mapSEO, err := loadSEO(ac)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		TrustedProxies:            webConf.TrustedProxies,
		EnableTrustedProxyCheck:   len(webConf.TrustedProxies) > 0,
		ProxyHeader:               webConf.ProxyHeader,
	})

	app.Use(logger.New())

	handler := handleWeb(ac, mapSEO, fileContent)

	app.Get("/", handler)
	app.Static("/", webConf.DataPath)
	app.Get("*", handler)

	return app.Listen(fmt.Sprintf("%s:%s", webConf.Bind, webConf.Port))
}
