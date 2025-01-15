package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// FillRectOpts содержит параметры для отрисовки заполненного прямоугольника.
type FillRectOpts struct {
	Rect  image.Rectangle // Прямоугольник для отрисовки.
	Color color.Color     // Цвет заливки.
}

// DrawRectOpts содержит параметры для отрисовки контура прямоугольника.
type DrawRectOpts struct {
	Rect        image.Rectangle // Прямоугольник для отрисовки.
	Color       color.Color     // Цвет контура.
	StrokeWidth float32         // Толщина контура.
}

// FillRect отрисовывает заполненный прямоугольник на целевом изображении.
func FillRect(target *ebiten.Image, opts FillRectOpts) {
	vector.DrawFilledRect(
		target,
		float32(opts.Rect.Min.X),
		float32(opts.Rect.Min.Y),
		float32(opts.Rect.Dx()),
		float32(opts.Rect.Dy()),
		opts.Color,
		false, // Не использовать антиалиасинг для повышения производительности.
	)
}

// DrawRectOutline отрисовывает контур прямоугольника на целевом изображении.
func DrawRectOutline(target *ebiten.Image, opts DrawRectOpts) {
	vector.StrokeRect(
		target,
		float32(opts.Rect.Min.X),
		float32(opts.Rect.Min.Y),
		float32(opts.Rect.Dx()),
		float32(opts.Rect.Dy()),
		opts.StrokeWidth,
		opts.Color,
		false, // Не использовать антиалиасинг для повышения производительности.
	)
}
