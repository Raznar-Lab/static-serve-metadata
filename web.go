package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	logger_middleware "github.com/gofiber/fiber/v2/middleware/logger"
	"raznar.id/static-serve-metadata/config"
	"raznar.id/static-serve-metadata/logger"
)

type Metadata struct {
	Tag     string `json:"tag"`
	Content string
}

type SEOData struct {
	URL      string     `json:"url"`
	Default  bool       `json:"default"`
	Metadata []Metadata `json:"metadata"`
}

type GroupSEO struct {

	// key: lang code
	SeoContents        []SEOData
	SeoDefaultContents SEOData
}

func (g GroupSEO) GetDataByURL(url string) SEOData {
	for _, ctn := range g.SeoContents {
		if strings.Contains(url, ctn.URL) {
			return ctn
		}
	}

	return g.SeoDefaultContents
}

func (s SEOData) IsEmpty() bool {
	return s.URL == ""
}

func (s SEOData) CollectMetadataString() string {
	metadataList := []string{}
	for _, mtd := range s.Metadata {
		metadataList = append(metadataList, mtd.ConvertToHTML())
	}

	return strings.Join(metadataList, "\n")
}

func (s Metadata) ConvertToHTML() string {
	return fmt.Sprintf("<meta name=\"%s\" content=\"%s\">", s.Tag, s.Content)
}

func handleWeb(ac *config.AppConfig, mapSEO map[string]GroupSEO, fileContent []byte) func(c *fiber.Ctx) (err error) {
	defaultLang := getDefaultLang(ac)
	return func(c *fiber.Ctx) (err error) {
		fileCtn := string(fileContent)
		geoHeader := c.Get(ac.SeoConfig.GeoHeader)
		wPath := c.Path()

		langCode := getLangCode(ac, geoHeader)
		if langCode == "" {
			langCode = defaultLang
		}

		groupSEO := mapSEO[langCode]
		logger.System.DebugInfo("lang ", langCode)
		seoData := groupSEO.GetDataByURL(wPath)
		fileCtn = strings.Replace(fileCtn, "<!-- seo header -->", seoData.CollectMetadataString(), 1)

		c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", ac.WebConfig.MaxAge))
		c.Set("Content-Type", "text/html")
		return c.SendString(fileCtn)
	}
}

func getLangCode(ac *config.AppConfig, country string) (lang string) {
	for k, v := range ac.SeoConfig.Languages {
		if slices.Contains(v.Country, country) {
			lang = k
			return
		}
	}

	return
}

func getDefaultLang(ac *config.AppConfig) (defaultLang string) {
	for k, v := range ac.SeoConfig.Languages {
		if v.Default {
			defaultLang = k
			break
		}
	}

	return
}

func loadSEO(ac *config.AppConfig) (map[string]GroupSEO, error) {
	groupSeo := make(map[string]GroupSEO)

	for lang := range maps.Keys(ac.SeoConfig.Languages) {
		langDirectory := path.Join(ac.SeoConfig.DataPath, lang)

		seoContents, err := loadSeoContents(langDirectory)
		if err != nil {
			return groupSeo, err
		}

		groupSeo[lang] = GroupSEO{SeoContents: seoContents}
	}

	for lang, content := range groupSeo {
		for _, seo := range content.SeoContents {
			if seo.Default {
				content.SeoDefaultContents = seo
				break
			}
		}

		groupSeo[lang] = content
	}

	return groupSeo, nil
}

func loadSeoContents(directory string) ([]SEOData, error) {
	var seoContents []SEOData

	// Walk through all files and directories
	err := filepath.WalkDir(directory, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %w", filePath, err)
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			logger.System.LogWarn("Failed to read file:", filePath, "Error:", err)
			return nil // Continue with the next file
		}

		// Parse JSON
		var seoDataList []SEOData
		if err := json.Unmarshal(fileContent, &seoDataList); err != nil {
			logger.System.LogWarn("Failed to parse JSON in file:", filePath, "Error:", err)
			return nil // Continue with the next file
		}

		logger.System.DebugInfo("filename:", d.Name())
		logger.System.DebugInfo("data:", seoDataList)

		seoContents = append(seoContents, seoDataList...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", directory, err)
	}

	return seoContents, nil
}


func RunWeb(ac *config.AppConfig) (err error) {
	webConf := ac.WebConfig
	fConf := fiber.Config{}
	fConf.TrustedProxies = webConf.TrustedProxies

	if len(fConf.TrustedProxies) > 0 {
		fConf.EnableTrustedProxyCheck = true
		fConf.ProxyHeader = webConf.ProxyHeader
	}

	fileContent, err := os.ReadFile(path.Join(webConf.DataPath, webConf.IndexFile))
	if err != nil {
		return
	}

	mapSeo, err := loadSEO(ac)
	if err != nil {
		return
	}

	webApp := fiber.New(fConf)
	webApp.Use(logger_middleware.New())

	webHandler := handleWeb(ac, mapSeo, fileContent)
	webApp.Get("/", webHandler)
	webApp.Static("/", webConf.DataPath)
	webApp.Get("*", webHandler)

	err = webApp.Listen(fmt.Sprintf("%s:%s", webConf.Bind, webConf.Port))
	return
}
