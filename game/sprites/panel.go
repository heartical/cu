package sprites

import (
	"cu/game/animation"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// PanelOpts содержит параметры для создания панелей.
type PanelOpts struct {
	Border int // Толщина границы панели.
	Center int // Размер центральной части панели.
}

// createPanels создает панели из изображения на основе заданных параметров.
func createPanels(img *ebiten.Image, rect image.Rectangle, opts PanelOpts) map[string]*animation.Sprite {
	panels := make(map[string]*animation.Sprite)
	// width, height := img.Bounds().Dx(), img.Bounds().Dy()
	border, center := opts.Border, opts.Center
	centerX, centerY := rect.Min.X+rect.Dx()/2, rect.Min.Y+rect.Dy()/2

	// createPanel создает одну панель и добавляет её в карту.
	createPanel := func(x, y, width, height int, key string) {
		grid := animation.NewGrid(width, height, width, height, x, y)
		panels[key] = animation.NewSprite(img, grid.Frames())
	}

	// Создание панелей для каждой части UI.
	createPanel(rect.Min.X, rect.Min.Y, border, border, "top_left")
	createPanel(centerX-center/2, rect.Min.Y, center, border, "top")
	createPanel(rect.Min.X+rect.Dx()-border, rect.Min.Y, border, border, "top_right")
	createPanel(rect.Min.X, centerY-center/2, border, center, "left")
	createPanel(centerX-center/2, centerY-center/2, center, center, "center")
	createPanel(rect.Min.X+rect.Dx()-border, centerY-center/2, border, center, "right")
	createPanel(rect.Min.X, rect.Min.Y+rect.Dy()-border, border, border, "bottom_left")
	createPanel(centerX-center/2, rect.Max.Y-border, center, border, "bottom")
	createPanel(rect.Min.X+rect.Dx()-border, rect.Min.Y+rect.Dy()-border, border, border, "bottom_right")

	return panels
}
