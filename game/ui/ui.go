package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Debug      = false
	debugColor = color.RGBA{0xff, 0, 0, 0xff}
	// debugColorShift = colorm.ColorM{}
)

func debugBorders(screen *ebiten.Image, root containerEmbed) {
	queue := []containerEmbed{root}
	// renderColor := resetDebugColor()
	for len(queue) > 0 {
		levelSize := len(queue)
		for levelSize != 0 {
			curr := queue[0]
			queue = queue[1:]
			DrawRectangleOutline(screen, DrawRectangleOptions{
				Bounds: curr.frame,
				// Color:       renderColor,
				Color:       debugColor,
				StrokeWidth: 2,
			})
			for _, c := range curr.children {
				if c.item.Display == DisplayNone {
					continue
				}
				queue = append(queue, c.item.containerEmbed)
			}
			levelSize--
		}
		// renderColor = rotateDebugColor()
	}
}

// func rotateDebugColor() color.Color {
// 	debugColorShift = debugColorShift.Concat(colorm.RotateHue(1.66))
// 	return debugColorShift.Apply(debugColor)
// }

// func resetDebugColor() color.Color {
// 	debugColorShift = colorm.ColorM{}
// 	return debugColor
// }
