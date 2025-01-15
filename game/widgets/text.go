package widgets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tinne26/etxt"

	"cu/common/assets"
	"cu/game/ui"
)

// Text представляет собой виджет для отрисовки текста.
type Text struct {
	Color     color.Color    // Цвет текста.
	Shadow    bool           // Флаг, указывающий, нужно ли отрисовывать тень.
	HorzAlign etxt.HorzAlign // Горизонтальное выравнивание текста.
	VertAlign etxt.VertAlign // Вертикальное выравнивание текста.
	Text      string         // Текст для отрисовки (если пуст, используется текст из View).
}

// Убедимся, что Text реализует интерфейс ui.Drawer.
var _ ui.Drawer = (*Text)(nil)

// Draw отрисовывает текст на экране с учетом настроек выравнивания, цвета и тени.
func (t *Text) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	// Отрисовка тени, если она включена.
	if t.Shadow {
		shadowWidth := float64(len(view.Text)*6 + 4)
		ebitenutil.DrawRect(screen, float64(frame.Min.X), float64(frame.Min.Y), shadowWidth, float64(frame.Dy()), color.RGBA{0, 0, 0, 50})
	}

	// Вычисляем позицию текста в зависимости от выравнивания.
	x, y := frame.Min.X+frame.Dx()/2, frame.Min.Y+frame.Dy()/2
	if t.HorzAlign == etxt.Left {
		x = frame.Min.X
	}
	if t.VertAlign == etxt.Top {
		y = frame.Min.Y
	}

	// Устанавливаем цвет текста.
	if t.Color != nil {
		assets.Renderer.SetColor(t.Color)
	} else {
		assets.Renderer.SetColor(color.White) // Используем белый цвет по умолчанию.
	}

	// Настраиваем выравнивание и отрисовываем текст.
	assets.Renderer.SetAlign(t.VertAlign, t.HorzAlign)
	assets.Renderer.SetTarget(screen)
	if t.Text == "" {
		assets.Renderer.Draw(view.Text, x, y) // Используем текст из View, если текст виджета пуст.
	} else {
		assets.Renderer.Draw(t.Text, x, y) // Используем текст виджета.
	}
}
