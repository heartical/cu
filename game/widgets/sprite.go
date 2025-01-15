package widgets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"cu/game/animation"
	"cu/game/sprites"
	"cu/game/ui"
)

// Sprite представляет собой виджет для отрисовки спрайта.
type Sprite struct{}

// Убедимся, что Sprite реализует интерфейс ui.Drawer.
var _ ui.Drawer = (*Sprite)(nil)

// Draw отрисовывает спрайт на экране в центре заданного прямоугольника.
func (s *Sprite) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	// Получаем спрайт из атрибутов View.
	sprite := sprites.Get(view.Attrs["sprite"])

	// Вычисляем центр прямоугольника для позиционирования спрайта.
	centerX := float64(frame.Min.X) + float64(frame.Dx())/2
	centerY := float64(frame.Min.Y) + float64(frame.Dy())/2

	// Отрисовываем спрайт с центрированием.
	animation.DrawSprite(screen, sprite, 0, centerX, centerY, 0, 1, 1, 0.5, 0.5)
}
