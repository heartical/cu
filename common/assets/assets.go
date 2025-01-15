package assets

import (
	"embed"
	"fmt"
	"io"
	"net/http"

	"github.com/tinne26/etxt"
)

//go:embed templates/index.html
var TemplateFS embed.FS

//go:embed *.html images/*.xml images/*.png fonts/*.ttf
var embeddedResources embed.FS

// IndexHTML содержит содержимое HTML-файла для главной страницы.
var IndexHTML []byte

// Renderer используется для рендеринга текста с использованием загруженных шрифтов.
var Renderer *etxt.Renderer

// FontName определяет имя шрифта по умолчанию.
const FontName = "x14y20pxScoreDozer"

// Asset содержит все загруженные ресурсы, такие как изображения, XML-файлы и шрифты.
type Asset struct {
	UIPackRPGSheetPNG   []byte
	UIPackRPGSheetXML   []byte
	UIPackSpaceSheetPNG []byte
	UIPackSpaceSheetXML []byte
	FontScoreDozer      []byte
}

// Init инициализирует и загружает все ресурсы из встроенной файловой системы.
// Возвращает экземпляр Asset и ошибку, если загрузка не удалась.
func Init() (*Asset, error) {
	assets := &Asset{}
	if err := loadEmbeddedResources(assets); err != nil {
		return nil, fmt.Errorf("failed to load embedded resources: %w", err)
	}
	return assets, nil
}

// loadEmbeddedResources загружает ресурсы из встроенной файловой системы в структуру Asset.
func loadEmbeddedResources(assets *Asset) error {
	resourceMappings := []struct {
		path string
		dest *[]byte
	}{
		{"_index.html", &IndexHTML},
		{"images/uipack_rpg_sheet.xml", &assets.UIPackRPGSheetXML},
		{"images/uipack_rpg_sheet.png", &assets.UIPackRPGSheetPNG},
		{"images/uipackSpace_sheet.xml", &assets.UIPackSpaceSheetXML},
		{"images/uipackSpace_sheet.png", &assets.UIPackSpaceSheetPNG},
		{"fonts/x14y20pxScoreDozer.ttf", &assets.FontScoreDozer},
	}

	for _, mapping := range resourceMappings {
		data, err := embeddedResources.ReadFile(mapping.path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", mapping.path, err)
		}
		*mapping.dest = data
	}

	return nil
}

// LoadRemoteResource загружает ресурс по указанному URL и возвращает его содержимое.
// Возвращает ошибку, если запрос не удался или статус ответа не OK.
func LoadRemoteResource(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	return body, nil
}

// LoadFonts инициализирует рендерер и загружает шрифт для использования в текстовом рендеринге.
// В случае ошибки загрузки шрифта программа завершается с паникой.
func LoadFonts(fontData []byte) {
	fontLib := etxt.NewFontLibrary()
	_, err := fontLib.ParseFontBytes(fontData)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %w", err))
	}

	Renderer = etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	Renderer.SetCacheHandler(glyphsCache.NewHandler())
	Renderer.SetFont(fontLib.GetFont(FontName))
	Renderer.SetAlign(etxt.YCenter, etxt.XCenter)
	Renderer.SetSizePx(14)
}
