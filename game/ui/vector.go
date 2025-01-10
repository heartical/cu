package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type FillRectangleOptions struct {
	Bounds image.Rectangle
	Color  color.Color
}

type DrawRectangleOptions struct {
	Bounds      image.Rectangle
	Color       color.Color
	StrokeWidth float32
}

func FillRectangle(target *ebiten.Image, opts FillRectangleOptions) {
	bounds := opts.Bounds
	col := opts.Color

	vector.DrawFilledRect(target,
		float32(bounds.Min.X),
		float32(bounds.Min.Y),
		float32(bounds.Dx()),
		float32(bounds.Dy()),
		col,
		false,
	)
}

func DrawRectangleOutline(target *ebiten.Image, opts DrawRectangleOptions) {
	bounds := opts.Bounds
	strokeWidth := opts.StrokeWidth

	vector.StrokeRect(target,
		float32(bounds.Min.X),
		float32(bounds.Min.Y),
		float32(bounds.Dx()),
		float32(bounds.Dy()),
		strokeWidth,
		opts.Color,
		false,
	)
}
