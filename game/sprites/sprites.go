package sprites

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"cu/game/animation"
)

// sprites хранит загруженные спрайты.
var sprites = make(map[string]*animation.Sprite)

// textureAtlas представляет структуру XML-атласа текстур.
type textureAtlas struct {
	SubTextures []subTexture `xml:"SubTexture"`
}

// subTexture представляет отдельную текстуру в атласе.
type subTexture struct {
	Name   string `xml:"name,attr"`
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// LoadOpts содержит параметры для загрузки спрайтов.
type LoadOpts struct {
	PanelOpts map[string]PanelOpts
}

// LoadSprites загружает спрайты из XML-атласа и изображения.
func LoadSprites(xmlData, imgData []byte, opts LoadOpts) {
	var atlas textureAtlas
	if err := xml.Unmarshal(xmlData, &atlas); err != nil {
		panic(fmt.Sprintf("не удалось распарсить XML: %v", err))
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgData))
	if err != nil {
		panic(fmt.Sprintf("не удалось загрузить изображение: %v", err))
	}

	rects := make(map[string]image.Rectangle)
	for _, subTex := range atlas.SubTextures {
		rect := image.Rect(subTex.X, subTex.Y, subTex.X+subTex.Width, subTex.Y+subTex.Height)
		grid := animation.NewGrid(subTex.Width, subTex.Height, img.Bounds().Dx(), img.Bounds().Dy(), subTex.X, subTex.Y)
		sprites[subTex.Name] = animation.NewSprite(img, grid.Frames())
		rects[subTex.Name] = rect
	}

	for key, panelOpts := range opts.PanelOpts {
		rect, ok := rects[key]
		if !ok {
			panic(fmt.Sprintf("панель не найдена: %s", key))
		}
		panels := createPanels(img, rect, panelOpts)
		for panelKey, panel := range panels {
			sprites[fmt.Sprintf("%s_%s", key, panelKey)] = panel
		}
	}
}

// Get возвращает спрайт по имени.
func Get(name string) *animation.Sprite {
	sprite, ok := sprites[name]
	if !ok {
		panic(fmt.Sprintf("спрайт не найден: %s", name))
	}
	return sprite
}
