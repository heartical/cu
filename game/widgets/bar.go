package widgets

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"cu/game/animation"
	"cu/game/sprites"
	"cu/game/ui"
)

// Bar представляет собой виджет горизонтальной полосы, которая может отображать значение.
type Bar struct {
	Value float64 // Текущее значение полосы (от 0 до 1).
}

// Убедимся, что Bar реализует интерфейс ui.Drawer.
var _ ui.Drawer = (*Bar)(nil)

// Draw отрисовывает полосу на экране.
func (b *Bar) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	b.drawBackground(screen, frame)
	b.drawForeground(screen, frame, view)
}

// drawBackground отрисовывает фон полосы.
func (b *Bar) drawBackground(screen *ebiten.Image, frame image.Rectangle) {
	x, y := float64(frame.Min.X), float64(frame.Min.Y)

	// Отрисовка левой части фона.
	leftSprite := sprites.Get("barBack_horizontalLeft.png")
	animation.DrawSprite(screen, leftSprite, 0, x, y, 0, 1, 1, 0, 0)
	x += float64(leftSprite.Width())

	// Отрисовка средней части фона.
	midSprite := sprites.Get("barBack_horizontalMid.png")
	for x < float64(frame.Max.X)-float64(midSprite.Width()) {
		animation.DrawSprite(screen, midSprite, 0, x, y, 0, 1, 1, 0, 0)
		x += float64(midSprite.Width())
	}

	// Отрисовка правой части фона.
	rightSprite := sprites.Get("barBack_horizontalRight.png")
	animation.DrawSprite(screen, rightSprite, 0, float64(frame.Max.X), y, 0, 1, 1, 1, 0)
}

// drawForeground отрисовывает передний план полосы в зависимости от значения.
func (b *Bar) drawForeground(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	maxX := frame.Min.X + int(b.Value*float64(frame.Dx())) // Вычисляем конечную позицию полосы.
	color := view.Attrs["color"]                           // Получаем цвет из атрибутов View.

	x, y := float64(frame.Min.X), float64(frame.Min.Y)

	// Отрисовка левой части переднего плана.
	leftSprite := sprites.Get(fmt.Sprintf("bar%s_horizontalLeft.png", color))
	animation.DrawSprite(screen, leftSprite, 0, x, y, 0, 1, 1, 0, 0)
	x += float64(leftSprite.Width())

	// Отрисовка средней части переднего плана.
	midSprite := sprites.Get(fmt.Sprintf("bar%s_horizontalMid.png", color))
	for x < float64(maxX)-float64(midSprite.Width()) {
		animation.DrawSprite(screen, midSprite, 0, x, y, 0, 1, 1, 0, 0)
		x += float64(midSprite.Width())
	}

	// Отрисовка правой части переднего плана.
	rightSprite := sprites.Get(fmt.Sprintf("bar%s_horizontalRight.png", color))
	animation.DrawSprite(screen, rightSprite, 0, float64(maxX), y, 0, 1, 1, 1, 0)
}
