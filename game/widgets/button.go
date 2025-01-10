package widgets

import (
	"engine/game/animation"
	"engine/game/sprites"
	"engine/game/text"
	"engine/game/ui"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type Button struct {
	Color   color.Color
	OnClick func()

	mouseover bool
	pressed   bool
}

var (
	_ ui.ButtonHandler          = (*Button)(nil)
	_ ui.Drawer                 = (*Button)(nil)
	_ ui.MouseEnterLeaveHandler = (*Button)(nil)
)

func (b *Button) HandlePress(x, y int, t ebiten.TouchID) {
	b.pressed = true
}

func (b *Button) HandleRelease(x, y int, isCancel bool) {
	b.pressed = false
	if !isCancel {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	x, y := float64(frame.Min.X+frame.Dx()/2), float64(frame.Min.Y+frame.Dy()/2)

	sprite := view.Attrs["sprite"]
	spritePressed := view.Attrs["sprite_pressed"]

	opts := animation.DrawOpts(x, y, 0, 1, 1, .5, .5)
	if b.mouseover {
		opts.ColorM.Scale(1.1, 1.1, 1.1, 1)
	}
	if b.pressed && spritePressed != "" {
		animation.DrawSpriteWithOpts(screen, sprites.Get(spritePressed), 0, opts, nil)
	} else if sprite != "" {
		animation.DrawSpriteWithOpts(screen, sprites.Get(sprite), 0, opts, nil)
	}

	text.R.SetAlign(etxt.YCenter, etxt.XCenter)
	text.R.SetTarget(screen)
	if b.Color != nil {
		text.R.SetColor(b.Color)
	} else {
		text.R.SetColor(color.White)
	}
	text.R.Draw(view.Text, int(x), int(y))
}

func (b *Button) HandleMouseEnter(x, y int) bool {
	b.mouseover = true
	return true
}

func (b *Button) HandleMouseLeave() {
	b.mouseover = false
}
