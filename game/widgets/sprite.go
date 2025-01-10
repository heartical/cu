package widgets

import (
	"engine/game/animation"
	"engine/game/sprites"
	"engine/game/ui"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
}

var (
	_ ui.Drawer = (*Sprite)(nil)
)

func (t *Sprite) Draw(screen *ebiten.Image, frame image.Rectangle, view *ui.View) {
	sprite := view.Attrs["sprite"]
	spr := sprites.Get(sprite)
	x, y := float64(frame.Min.X)+float64(frame.Dx())/2, float64(frame.Min.Y)+float64(frame.Dy())/2
	animation.DrawSprite(screen, spr, 0, x, y, 0, 1, 1, .5, .5)
}
