package widgets

import (
	"engine/game/animation"
	"engine/game/sprites"
	"engine/game/ui"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bar struct {
	Value float64
}

var (
	_ ui.Drawer = (*Bar)(nil)
)

func (b *Bar) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	b.drawBlackBar(screen, frame)
	b.drawBar(screen, frame, view)
}

func (b *Bar) drawBlackBar(screen *ebiten.Image, frame image.Rectangle) {
	x, y := float64(frame.Min.X), float64(frame.Min.Y)
	spr := sprites.Get("barBack_horizontalLeft.png")
	animation.DrawSprite(screen, spr, 0, x, y, 0, 1, 1, 0, 0)
	x += float64(spr.W())

	spr = sprites.Get("barBack_horizontalMid.png")
	for x < float64(frame.Max.X)-float64(spr.W()) {
		animation.DrawSprite(screen, spr, 0, x, y, 0, 1, 1, 0, 0)
		x += float64(spr.W())
	}

	spr = sprites.Get("barBack_horizontalRight.png")
	animation.DrawSprite(screen, spr, 0, float64(frame.Max.X), y, 0, 1, 1, 1, 0)
}

func (b *Bar) drawBar(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	maxX := frame.Min.X + int(b.Value*float64(frame.Dx()))
	color := view.Attrs["color"]

	x, y := float64(frame.Min.X), float64(frame.Min.Y)
	spr := sprites.Get(fmt.Sprintf("bar%s_horizontalLeft.png", color))
	animation.DrawSprite(screen, spr, 0, x, y, 0, 1, 1, 0, 0)
	x += float64(spr.W())

	spr = sprites.Get(fmt.Sprintf("bar%s_horizontalMid.png", color))
	for x < float64(maxX)-float64(spr.W()) {
		animation.DrawSprite(screen, spr, 0, x, y, 0, 1, 1, 0, 0)
		x += float64(spr.W())
	}

	spr = sprites.Get(fmt.Sprintf("bar%s_horizontalRight.png", color))
	animation.DrawSprite(screen, spr, 0, float64(maxX), y, 0, 1, 1, 1, 0)
}
