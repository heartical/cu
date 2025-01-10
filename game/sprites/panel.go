package sprites

import (
	"engine/game/animation"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PanelOpts struct {
	Border int
	Center int
}

func createPanels(img *ebiten.Image, r image.Rectangle, opts PanelOpts) map[string]*animation.Sprite {
	ret := map[string]*animation.Sprite{}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	border := opts.Border
	center := opts.Center
	cx, cy := r.Min.X+r.Dx()/2, r.Min.Y+r.Dy()/2

	// top left
	g := animation.NewGrid(border, border, w, h, r.Min.X, r.Min.Y)
	ret["top_left"] = animation.NewSprite(img, g.Frames())
	// top
	g = animation.NewGrid(center, border, w, h, cx-center/2, r.Min.Y)
	ret["top"] = animation.NewSprite(img, g.Frames())
	// top right
	g = animation.NewGrid(border, border, w, h, r.Min.X+r.Dx()-border, r.Min.Y)
	ret["top_right"] = animation.NewSprite(img, g.Frames())
	// left
	g = animation.NewGrid(border, center, w, h, r.Min.X, cy-center/2)
	ret["left"] = animation.NewSprite(img, g.Frames())
	// center
	g = animation.NewGrid(center, center, w, h, cx-center/2, cy-center/2)
	ret["center"] = animation.NewSprite(img, g.Frames())
	// right
	g = animation.NewGrid(border, center, w, h, r.Min.X+r.Dx()-border, cy-center/2)
	ret["right"] = animation.NewSprite(img, g.Frames())
	// bottom left
	g = animation.NewGrid(border, border, w, h, r.Min.X, r.Min.Y+r.Dy()-border)
	ret["bottom_left"] = animation.NewSprite(img, g.Frames())
	// bottom
	g = animation.NewGrid(center, border, w, h, cx-center/2, r.Max.Y-border)
	ret["bottom"] = animation.NewSprite(img, g.Frames())
	// bottom right
	g = animation.NewGrid(border, border, w, h, r.Min.X+r.Dx()-border, r.Min.Y+r.Dy()-border)
	ret["bottom_right"] = animation.NewSprite(img, g.Frames())

	return ret
}
