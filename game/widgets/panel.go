package widgets

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"

	"cu/common/assets"
	"cu/game/animation"
	"cu/game/sprites"
	"cu/game/ui"
)

// Panel представляет собой виджет панели, который может отрисовываться и обрабатывать события мыши и касания.
type Panel struct {
	Color     color.Color // Цвет текста на панели.
	OnClick   func()      // Функция, вызываемая при нажатии на панель.
	mouseover bool        // Флаг, указывающий, находится ли курсор мыши над панелью.
	pressed   bool        // Флаг, указывающий, нажата ли панель.
}

// Убедимся, что Panel реализует необходимые интерфейсы.
var (
	_ ui.ButtonHandler          = (*Panel)(nil)
	_ ui.NotButton              = (*Panel)(nil)
	_ ui.Drawer                 = (*Panel)(nil)
	_ ui.MouseEnterLeaveHandler = (*Panel)(nil)
)

// Draw отрисовывает панель и текст на ней.
func (p *Panel) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	panelName := view.Attrs["sprite"] // Получаем имя панели из атрибутов View.
	border := sprites.Get(fmt.Sprintf("%s_top_left", panelName)).Width()
	top := sprites.Get(fmt.Sprintf("%s_top", panelName)).Height()
	fborder := float64(border)

	// Отрисовка центральной части панели.
	centerSprite := sprites.Get(fmt.Sprintf("%s_center", panelName))
	x := float64(frame.Min.X) + fborder
	for x < float64(frame.Max.X)-fborder {
		y := float64(frame.Min.Y) + float64(top)
		for y < float64(frame.Max.Y)-fborder {
			opts := animation.DrawOpts(x, y, 0, 1, 1, 0, 0)
			p.drawSprite(screen, centerSprite, opts)
			y += float64(centerSprite.Height())
		}
		x += float64(centerSprite.Width())
	}

	// Отрисовка углов и краев панели.
	p.drawCorner(screen, fmt.Sprintf("%s_top_left", panelName), float64(frame.Min.X), float64(frame.Min.Y))
	p.drawEdge(screen, fmt.Sprintf("%s_top", panelName), float64(frame.Min.X+border), float64(frame.Min.Y), true)
	p.drawCorner(screen, fmt.Sprintf("%s_top_right", panelName), float64(frame.Max.X-border), float64(frame.Min.Y))
	p.drawEdge(screen, fmt.Sprintf("%s_left", panelName), float64(frame.Min.X), float64(frame.Min.Y+border), false)
	p.drawEdge(screen, fmt.Sprintf("%s_right", panelName), float64(frame.Max.X-border), float64(frame.Min.Y+border), false)
	p.drawCorner(screen, fmt.Sprintf("%s_bottom_left", panelName), float64(frame.Min.X), float64(frame.Max.Y-border))
	p.drawEdge(screen, fmt.Sprintf("%s_bottom", panelName), float64(frame.Min.X+border), float64(frame.Max.Y-border), true)
	p.drawCorner(screen, fmt.Sprintf("%s_bottom_right", panelName), float64(frame.Max.X-border), float64(frame.Max.Y-border))

	// Отрисовка текста на панели, если он задан.
	if view.Text != "" {
		x, y := float64(frame.Min.X+frame.Dx()/2), float64(frame.Min.Y+frame.Dy()/2)
		assets.Renderer.SetAlign(etxt.YCenter, etxt.XCenter)
		assets.Renderer.SetTarget(screen)
		if p.Color != nil {
			assets.Renderer.SetColor(p.Color)
		} else {
			assets.Renderer.SetColor(color.White) // Используем белый цвет по умолчанию.
		}
		assets.Renderer.Draw(view.Text, int(x), int(y))
	}
}

// drawCorner отрисовывает угол панели.
func (p *Panel) drawCorner(screen *ebiten.Image, spriteName string, x, y float64) {
	sprite := sprites.Get(spriteName)
	opts := animation.DrawOpts(x, y, 0, 1, 1, 0, 0)
	p.drawSprite(screen, sprite, opts)
}

// drawEdge отрисовывает край панели.
func (p *Panel) drawEdge(screen *ebiten.Image, spriteName string, startX, startY float64, isHorizontal bool) {
	sprite := sprites.Get(spriteName)
	if isHorizontal {
		for x := startX; x < float64(screen.Bounds().Max.X)-float64(sprite.Width()); x += float64(sprite.Width()) {
			opts := animation.DrawOpts(x, startY, 0, 1, 1, 0, 0)
			p.drawSprite(screen, sprite, opts)
		}
	} else {
		for y := startY; y < float64(screen.Bounds().Max.Y)-float64(sprite.Height()); y += float64(sprite.Height()) {
			opts := animation.DrawOpts(startX, y, 0, 1, 1, 0, 0)
			p.drawSprite(screen, sprite, opts)
		}
	}
}

// drawSprite отрисовывает спрайт с учетом состояния панели (нажата или наведена мышь).
func (p *Panel) drawSprite(screen *ebiten.Image, sprite *animation.Sprite, opts *animation.DrawOptions) {
	if p.IsButton() {
		if p.pressed {
			opts.ColorM.Scale(0.9, 0.9, 0.9, 1) // Уменьшаем яркость при нажатии.
		} else if p.mouseover {
			opts.ColorM.Scale(1.1, 1.1, 1.1, 1) // Увеличиваем яркость при наведении.
		}
	}
	animation.DrawSpriteWithOpts(screen, sprite, 0, opts, nil)
}

// IsButton возвращает true, если панель является кнопкой (имеет обработчик OnClick).
func (p *Panel) IsButton() bool {
	return p.OnClick != nil
}

// HandleMouseEnter обрабатывает событие входа курсора мыши в область панели.
func (p *Panel) HandleMouseEnter(x, y int) bool {
	p.mouseover = true
	return true
}

// HandleMouseLeave обрабатывает событие выхода курсора мыши из области панели.
func (p *Panel) HandleMouseLeave() {
	p.mouseover = false
}

// HandlePress обрабатывает событие нажатия на панель.
func (p *Panel) HandlePress(x, y int, t ebiten.TouchID) {
	p.pressed = true
}

// HandleRelease обрабатывает событие отпускания панели.
func (p *Panel) HandleRelease(x, y int, isCancel bool) {
	p.pressed = false
	if !isCancel && p.OnClick != nil {
		p.OnClick() // Вызываем обработчик нажатия, если событие не было отменено.
	}
}
