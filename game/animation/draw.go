package animation

import "github.com/hajimehoshi/ebiten/v2"

// drawOpts содержит глобальные параметры отрисовки по умолчанию.
var drawOpts = DrawOpts(0, 0, 0, 1, 1, 1, 0.5, 0.5)

// DrawSprite отрисовывает спрайт на экране с указанными параметрами.
func DrawSprite(screen *ebiten.Image, sprite *Sprite, index int, x, y, rotate, scaleX, scaleY, originX, originY float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rotate)
	drawOpts.SetScale(scaleX, scaleY)
	drawOpts.SetOrigin(originX, originY)
	sprite.Draw(screen, index, drawOpts)
}

// DrawSpriteWithOpts отрисовывает спрайт на экране с использованием переданных параметров отрисовки и шейдера.
func DrawSpriteWithOpts(screen *ebiten.Image, sprite *Sprite, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	if shaderOpts != nil {
		sprite.DrawWithShader(screen, index, opts, shaderOpts)
	} else {
		sprite.Draw(screen, index, opts)
	}
}

// DrawAnimation отрисовывает анимацию на экране с указанными параметрами.
func DrawAnimation(screen *ebiten.Image, animation *Animation, x, y, rotate, scaleX, scaleY, originX, originY float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rotate)
	drawOpts.SetScale(scaleX, scaleY)
	drawOpts.SetOrigin(originX, originY)
	animation.Draw(screen, drawOpts)
}

// DrawAnimationWithOpts отрисовывает анимацию на экране с использованием переданных параметров отрисовки и шейдера.
func DrawAnimationWithOpts(screen *ebiten.Image, animation *Animation, opts *DrawOptions, shaderOpts *ShaderOptions) {
	if shaderOpts != nil {
		animation.DrawWithShader(screen, opts, shaderOpts)
	} else {
		animation.Draw(screen, opts)
	}
}
