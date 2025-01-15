package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Debug включает или отключает отладочный режим.
	Debug = false

	// debugColor определяет цвет отладочных границ.
	debugColor = color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}
)

// debugBorders отрисовывает границы всех элементов UI для отладки.
func debugBorders(screen *ebiten.Image, root containerEmbed) {
	if !Debug {
		return // Если отладочный режим выключен, ничего не делаем.
	}

	queue := []containerEmbed{root} // Очередь для обхода элементов в ширину.
	for len(queue) > 0 {
		levelSize := len(queue) // Количество элементов на текущем уровне.
		for levelSize > 0 {
			current := queue[0] // Берем первый элемент из очереди.
			queue = queue[1:]   // Удаляем его из очереди.

			// Отрисовываем границу текущего элемента.
			DrawRectOutline(screen, DrawRectOpts{
				Rect:        current.frame,
				Color:       debugColor,
				StrokeWidth: 2,
			})

			// Добавляем дочерние элементы в очередь, если они не скрыты.
			for _, child := range current.children {
				if child.item.Display != DisplayNone {
					queue = append(queue, child.item.containerEmbed)
				}
			}

			levelSize-- // Уменьшаем счетчик элементов на текущем уровне.
		}
	}
}
