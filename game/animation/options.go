package animation

import "github.com/hajimehoshi/ebiten/v2"

// DrawOptions содержит параметры для отрисовки спрайта.
type DrawOptions struct {
	X, Y             float64              // Позиция спрайта на экране.
	Rotate           float64              // Угол поворота спрайта.
	ScaleX, ScaleY   float64              // Масштабирование спрайта по осям X и Y.
	OriginX, OriginY float64              // Точка вращения и масштабирования спрайта.
	ColorM           ebiten.ColorM        // Матрица цветовых преобразований.
	CompositeMode    ebiten.CompositeMode // Режим композиции (наложение спрайта).
}

// SetPos устанавливает позицию спрайта.
func (opts *DrawOptions) SetPos(x, y float64) {
	opts.X, opts.Y = x, y
}

// SetRot устанавливает угол поворота спрайта.
func (opts *DrawOptions) SetRot(rotate float64) {
	opts.Rotate = rotate
}

// SetOrigin устанавливает точку вращения и масштабирования спрайта.
func (opts *DrawOptions) SetOrigin(originX, originY float64) {
	opts.OriginX, opts.OriginY = originX, originY
}

// SetScale устанавливает масштабирование спрайта.
func (opts *DrawOptions) SetScale(scaleX, scaleY float64) {
	opts.ScaleX, opts.ScaleY = scaleX, scaleY
}

// Reset сбрасывает все параметры отрисовки к значениям по умолчанию.
func (opts *DrawOptions) Reset() {
	opts.X, opts.Y = 0, 0
	opts.Rotate = 0
	opts.ScaleX, opts.ScaleY = 1, 1
	opts.OriginX, opts.OriginY = 0, 0
	opts.ColorM.Reset()
	opts.CompositeMode = ebiten.CompositeModeSourceOver
}

// ResetValues сбрасывает параметры отрисовки и устанавливает новые значения.
func (opts *DrawOptions) ResetValues(x, y, rotate, scaleX, scaleY, originX, originY float64) {
	opts.Reset()
	opts.SetPos(x, y)
	opts.SetRot(rotate)
	opts.SetScale(scaleX, scaleY)
	opts.SetOrigin(originX, originY)
}

// ShaderOptions содержит параметры для отрисовки с использованием шейдера.
type ShaderOptions struct {
	Uniforms map[string]interface{} // Uniform-переменные для шейдера.
	Shader   *ebiten.Shader         // Шейдер для отрисовки.
	Images   [3]*ebiten.Image       // Изображения, передаваемые в шейдер.
}

// DrawOpts создает новый экземпляр DrawOptions с заданными параметрами.
// Аргументы args могут содержать: [rotate, scaleX, scaleY, originX, originY].
func DrawOpts(x, y float64, args ...float64) *DrawOptions {
	rotate, scaleX, scaleY, originX, originY := 0.0, 1.0, 1.0, 0.0, 0.0

	// Обработка аргументов в зависимости от их количества.
	switch len(args) {
	case 5:
		originY = args[4]
		fallthrough
	case 4:
		originX = args[3]
		fallthrough
	case 3:
		scaleY = args[2]
		fallthrough
	case 2:
		scaleX = args[1]
		fallthrough
	case 1:
		rotate = args[0]
	}

	return &DrawOptions{
		X:             x,
		Y:             y,
		Rotate:        rotate,
		ScaleX:        scaleX,
		ScaleY:        scaleY,
		OriginX:       originX,
		OriginY:       originY,
		ColorM:        ebiten.ColorM{},
		CompositeMode: ebiten.CompositeModeSourceOver,
	}
}
