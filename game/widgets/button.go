package widgets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"

	"cu/common/assets"
	"cu/game/animation"
	"cu/game/sprites"
	"cu/game/ui"
)

// Button представляет собой виджет кнопки, который может обрабатывать нажатия и отрисовывать текст.
type Button struct {
	Color     color.Color // Цвет текста кнопки.
	OnClick   func()      // Функция, вызываемая при нажатии на кнопку.
	mouseover bool        // Флаг, указывающий, находится ли курсор мыши над кнопкой.
	pressed   bool        // Флаг, указывающий, нажата ли кнопка.
}

// Убедимся, что Button реализует необходимые интерфейсы.
var (
	_ ui.ButtonHandler          = (*Button)(nil)
	_ ui.Drawer                 = (*Button)(nil)
	_ ui.MouseEnterLeaveHandler = (*Button)(nil)
)

// HandlePress обрабатывает событие нажатия на кнопку.
func (b *Button) HandlePress(x, y int, t ebiten.TouchID) {
	b.pressed = true
}

// HandleRelease обрабатывает событие отпускания кнопки.
func (b *Button) HandleRelease(x, y int, isCancel bool) {
	b.pressed = false
	if !isCancel && b.OnClick != nil {
		b.OnClick() // Вызываем обработчик нажатия, если событие не было отменено.
	}
}

// Draw отрисовывает кнопку и текст на ней.
func (b *Button) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	// Вычисляем центр кнопки для позиционирования спрайта и текста.
	centerX := float64(frame.Min.X + frame.Dx()/2)
	centerY := float64(frame.Min.Y + frame.Dy()/2)

	// Получаем спрайты для обычного и нажатого состояния кнопки.
	sprite := view.Attrs["sprite"]
	spritePressed := view.Attrs["sprite_pressed"]

	// Настраиваем параметры отрисовки спрайта.
	opts := animation.DrawOpts(centerX, centerY, 0, 1, 1, 0.5, 0.5)
	if b.mouseover {
		opts.ColorM.Scale(1.1, 1.1, 1.1, 1) // Увеличиваем яркость при наведении.
	}

	// Отрисовываем спрайт в зависимости от состояния кнопки.
	if b.pressed && spritePressed != "" {
		animation.DrawSpriteWithOpts(screen, sprites.Get(spritePressed), 0, opts, nil)
	} else if sprite != "" {
		animation.DrawSpriteWithOpts(screen, sprites.Get(sprite), 0, opts, nil)
	}

	// Настраиваем и отрисовываем текст на кнопке.
	assets.Renderer.SetAlign(etxt.YCenter, etxt.XCenter)
	assets.Renderer.SetTarget(screen)
	if b.Color != nil {
		assets.Renderer.SetColor(b.Color)
	} else {
		assets.Renderer.SetColor(color.White) // Используем белый цвет по умолчанию.
	}
	assets.Renderer.Draw(view.Text, int(centerX), int(centerY))
}

// HandleMouseEnter обрабатывает событие входа курсора мыши в область кнопки.
func (b *Button) HandleMouseEnter(x, y int) bool {
	b.mouseover = true
	return true
}

// HandleMouseLeave обрабатывает событие выхода курсора мыши из области кнопки.
func (b *Button) HandleMouseLeave() {
	b.mouseover = false
}
