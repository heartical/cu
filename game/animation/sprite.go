package animation

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/rand"
)

// SpriteSize представляет размер спрайта в целых числах.
type SpriteSize struct {
	Width  int
	Height int
}

// SpriteSizeF представляет размер спрайта в числах с плавающей точкой.
type SpriteSizeF struct {
	Width  float64
	Height float64
}

// Sprite представляет спрайт, состоящий из нескольких кадров.
type Sprite struct {
	frames    []*image.Rectangle
	image     *ebiten.Image
	subImages []*ebiten.Image
	size      SpriteSize
	sizeF     SpriteSizeF
	length    int
	flippedH  bool
	flippedV  bool
	op        *ebiten.DrawImageOptions
	shaderOp  *ebiten.DrawRectShaderOptions
}

// NewSprite создает новый спрайт на основе изображения и кадров.
func NewSprite(img *ebiten.Image, frames []*image.Rectangle) *Sprite {
	subImages := make([]*ebiten.Image, len(frames))
	for i, frame := range frames {
		subImages[i] = img.SubImage(*frame).(*ebiten.Image)
	}

	size := SpriteSize{}
	sizeF := SpriteSizeF{}
	if len(frames) > 0 {
		size = SpriteSize{Width: frames[0].Dx(), Height: frames[0].Dy()}
		sizeF = SpriteSizeF{Width: float64(frames[0].Dx()), Height: float64(frames[0].Dy())}
	}

	return &Sprite{
		frames:    frames,
		image:     img,
		subImages: subImages,
		length:    len(frames),
		size:      size,
		sizeF:     sizeF,
		op:        &ebiten.DrawImageOptions{},
		shaderOp:  &ebiten.DrawRectShaderOptions{},
	}
}

// Size возвращает размер спрайта.
func (s *Sprite) Size() (int, int) {
	return s.size.Width, s.size.Height
}

// Width возвращает ширину спрайта.
func (s *Sprite) Width() int {
	return s.size.Width
}

// Height возвращает высоту спрайта.
func (s *Sprite) Height() int {
	return s.size.Height
}

// Length возвращает количество кадров спрайта.
func (s *Sprite) Length() int {
	return s.length
}

// RandomIndex возвращает случайный индекс кадра.
func (s *Sprite) RandomIndex() int {
	return rand.Intn(s.length)
}

// LoopIndex возвращает индекс кадра с учетом зацикливания.
func (s *Sprite) LoopIndex(index int) int {
	return index % s.length
}

// IsEnd проверяет, является ли индекс последним кадром.
func (s *Sprite) IsEnd(index int) bool {
	return index >= s.length-1
}

// FlipH переключает горизонтальное отражение спрайта.
func (s *Sprite) FlipH() {
	s.flippedH = !s.flippedH
}

// FlipV переключает вертикальное отражение спрайта.
func (s *Sprite) FlipV() {
	s.flippedV = !s.flippedV
}

// SetFlipH устанавливает горизонтальное отражение спрайта.
func (s *Sprite) SetFlipH(flipH bool) {
	s.flippedH = flipH
}

// SetFlipV устанавливает вертикальное отражение спрайта.
func (s *Sprite) SetFlipV(flipV bool) {
	s.flippedV = flipV
}

// Draw отрисовывает спрайт на экране с учетом переданных параметров.
func (s *Sprite) Draw(screen *ebiten.Image, index int, opts *DrawOptions) {
	op := s.op
	op.GeoM.Reset()
	op.ColorM = opts.ColorM
	op.CompositeMode = opts.CompositeMode

	w, h := s.sizeF.Width, s.sizeF.Height
	sx, sy := opts.ScaleX, opts.ScaleY
	ox, oy := opts.OriginX, opts.OriginY

	if s.flippedH {
		sx *= -1
		ox = 1 - ox
	}
	if s.flippedV {
		sy *= -1
		oy = 1 - oy
	}

	if sx != 1 || sy != 1 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(w*ox, h*oy)
	}

	if opts.Rotate != 0 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Rotate(opts.Rotate)
		op.GeoM.Translate(w*ox, h*oy)
	}

	op.GeoM.Translate(opts.X-w*ox, opts.Y-h*oy)
	screen.DrawImage(s.subImages[index], op)
}

// DrawWithShader отрисовывает спрайт с использованием шейдера.
func (s *Sprite) DrawWithShader(screen *ebiten.Image, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	op := s.shaderOp
	op.GeoM.Reset()
	op.CompositeMode = opts.CompositeMode
	op.Uniforms = shaderOpts.Uniforms

	w, h := s.sizeF.Width, s.sizeF.Height
	sx, sy := opts.ScaleX, opts.ScaleY
	ox, oy := opts.OriginX, opts.OriginY

	if opts.Rotate != 0 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Rotate(opts.Rotate)
		op.GeoM.Translate(w*ox, h*oy)
	}

	if s.flippedH {
		sx *= -1
	}
	if s.flippedV {
		sy *= -1
	}

	if sx != 1 || sy != 1 {
		op.GeoM.Translate(-w*ox, -h*oy)
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(w*ox, h*oy)
	}

	op.GeoM.Translate(opts.X-w*ox, opts.Y-h*oy)
	op.Images[0] = s.subImages[index]
	op.Images[1] = shaderOpts.Images[0]
	op.Images[2] = shaderOpts.Images[1]
	op.Images[3] = shaderOpts.Images[2]
	screen.DrawRectShader(int(w), int(h), shaderOpts.Shader, op)
}

// Clone создает копию спрайта.
func (s *Sprite) Clone() *Sprite {
	clone := *s
	clone.op = &ebiten.DrawImageOptions{}
	clone.shaderOp = &ebiten.DrawRectShaderOptions{}
	return &clone
}
