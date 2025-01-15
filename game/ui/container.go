package ui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// containerEmbed представляет контейнер для управления дочерними элементами UI.
type containerEmbed struct {
	children         []*child
	isDirty          bool
	frame            image.Rectangle
	touchIDs         []ebiten.TouchID
	calculatedWidth  int
	calculatedHeight int
}

// processEvent обрабатывает все события, такие как касания и движения мыши.
func (ct *containerEmbed) processEvent() {
	ct.handleTouchEvents()
	ct.handleMouseEvents()
}

// Draw отрисовывает все дочерние элементы контейнера.
func (ct *containerEmbed) Draw(screen *ebiten.Image) {
	for _, child := range ct.children {
		ct.drawChild(screen, child)
	}
}

// drawChild отрисовывает отдельный дочерний элемент.
func (ct *containerEmbed) drawChild(screen *ebiten.Image, child *child) {
	bounds := ct.computeBounds(child)
	if ct.shouldDrawChild(child) {
		ct.handleDraw(screen, bounds, child)
	}
	child.item.Draw(screen)
	ct.debugDraw(screen, bounds, child)
}

// computeBounds вычисляет границы дочернего элемента с учетом его позиции.
func (ct *containerEmbed) computeBounds(child *child) image.Rectangle {
	if child.absolute {
		return child.bounds
	}
	return child.bounds.Add(ct.frame.Min)
}

// handleDraw обрабатывает отрисовку дочернего элемента, если у него есть обработчик.
func (ct *containerEmbed) handleDraw(screen *ebiten.Image, bounds image.Rectangle, child *child) {
	if handler, ok := child.item.Handler.(DrawHandler); ok {
		handler.HandleDraw(screen, bounds)
		return
	}
	if handler, ok := child.item.Handler.(Drawer); ok {
		handler.Draw(screen, bounds, child.item)
	}
}

// shouldDrawChild проверяет, нужно ли отрисовывать дочерний элемент.
func (ct *containerEmbed) shouldDrawChild(child *child) bool {
	return !child.item.Hidden && child.item.Display != DisplayNone && child.item.Handler != nil
}

// debugDraw отрисовывает отладочную информацию, если включен режим отладки.
func (ct *containerEmbed) debugDraw(screen *ebiten.Image, bounds image.Rectangle, child *child) {
	if Debug {
		pos := fmt.Sprintf("(%d, %d)-(%d, %d):%s:%s", bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, child.item.TagName, child.item.ID)
		FillRect(screen, FillRectOpts{
			Color: color.RGBA{0, 0, 0, 200},
			Rect:  image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+len(pos)*6, bounds.Min.Y+12),
		})
		ebitenutil.DebugPrintAt(screen, pos, bounds.Min.X, bounds.Min.Y)
	}
}

// HandleJustPressedTouchID обрабатывает событие начала касания.
func (ct *containerEmbed) HandleJustPressedTouchID(touchID ebiten.TouchID, x, y int) bool {
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		childFrame := ct.childFrame(child)
		if child.item.Display == DisplayNone {
			continue
		}
		if child.HandleJustPressedTouchID(childFrame, touchID, x, y) || child.item.HandleJustPressedTouchID(touchID, x, y) {
			return true
		}
	}
	return false
}

// HandleJustReleasedTouchID обрабатывает событие окончания касания.
func (ct *containerEmbed) HandleJustReleasedTouchID(touchID ebiten.TouchID, x, y int) {
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		childFrame := ct.childFrame(child)
		child.HandleJustReleasedTouchID(childFrame, touchID, x, y)
		child.item.HandleJustReleasedTouchID(touchID, x, y)
	}
}

// handleMouse обрабатывает движение мыши.
func (ct *containerEmbed) handleMouse(x, y int) bool {
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		childFrame := ct.childFrame(child)
		if child.item.Display == DisplayNone {
			continue
		}
		if mouseHandler, ok := child.item.Handler.(MouseHandler); ok && isInside(childFrame, x, y) && mouseHandler.HandleMouse(x, y) {
			return true
		}
		if child.item.handleMouse(x, y) {
			return true
		}
	}
	return false
}

// handleMouseEnterLeave обрабатывает события входа и выхода мыши.
func (ct *containerEmbed) handleMouseEnterLeave(x, y int) bool {
	result := false
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		childFrame := ct.childFrame(child)
		if child.item.Display == DisplayNone {
			continue
		}
		if mouseHandler, ok := child.item.Handler.(MouseEnterLeaveHandler); ok {
			if !result && !child.isMouseEntered && isInside(childFrame, x, y) {
				result = mouseHandler.HandleMouseEnter(x, y)
				child.isMouseEntered = true
			}
			if child.isMouseEntered && !isInside(childFrame, x, y) {
				child.isMouseEntered = false
				mouseHandler.HandleMouseLeave()
			}
		}
		if child.item.handleMouseEnterLeave(x, y) {
			result = true
		}
	}
	return result
}

// handleMouseButtonLeftPressed обрабатывает нажатие левой кнопки мыши.
func (ct *containerEmbed) handleMouseButtonLeftPressed(x, y int) bool {
	result := false
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		childFrame := ct.childFrame(child)
		if child.item.Display == DisplayNone {
			continue
		}
		if mouseLeftClickHandler, ok := child.item.Handler.(MouseLeftButtonHandler); ok && !result && isInside(childFrame, x, y) {
			result = mouseLeftClickHandler.HandleJustPressedMouseButtonLeft(x, y)
			child.isMouseLeftButtonHandler = true
		}
		if button, ok := child.item.Handler.(ButtonHandler); ok && !result && isInside(childFrame, x, y) {
			if !child.isButtonPressed {
				child.isButtonPressed = true
				child.isMouseLeftButtonHandler = true
				result = true
				button.HandlePress(x, y, -1)
			}
		}
		if !result && child.item.handleMouseButtonLeftPressed(x, y) {
			result = true
		}
	}
	return result
}

// handleMouseButtonLeftReleased обрабатывает отпускание левой кнопки мыши.
func (ct *containerEmbed) handleMouseButtonLeftReleased(x, y int) {
	for i := len(ct.children) - 1; i >= 0; i-- {
		child := ct.children[i]
		if mouseLeftClickHandler, ok := child.item.Handler.(MouseLeftButtonHandler); ok && child.isMouseLeftButtonHandler {
			child.isMouseLeftButtonHandler = false
			mouseLeftClickHandler.HandleJustReleasedMouseButtonLeft(x, y)
		}
		if button, ok := child.item.Handler.(ButtonHandler); ok && child.isButtonPressed && child.isMouseLeftButtonHandler {
			child.isButtonPressed = false
			child.isMouseLeftButtonHandler = false
			button.HandleRelease(x, y, !isInside(ct.childFrame(child), x, y))
		}
		child.item.handleMouseButtonLeftReleased(x, y)
	}
}

// isInside проверяет, находится ли точка (x, y) внутри прямоугольника.
func isInside(r *image.Rectangle, x, y int) bool {
	return r.Min.X <= x && x <= r.Max.X && r.Min.Y <= y && y <= r.Max.Y
}

// handleTouchEvents обрабатывает события касания.
func (ct *containerEmbed) handleTouchEvents() {
	justPressedTouchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	for _, touchID := range justPressedTouchIDs {
		x, y := ebiten.TouchPosition(touchID)
		recordTouchPosition(touchID, x, y)
		ct.HandleJustPressedTouchID(touchID, x, y)
		ct.touchIDs = append(ct.touchIDs, touchID)
	}

	for _, touchID := range ct.touchIDs {
		if inpututil.IsTouchJustReleased(touchID) {
			pos := lastTouchPosition(touchID)
			ct.HandleJustReleasedTouchID(touchID, pos.X, pos.Y)
		} else {
			x, y := ebiten.TouchPosition(touchID)
			recordTouchPosition(touchID, x, y)
		}
	}
}

// handleMouseEvents обрабатывает события мыши.
func (ct *containerEmbed) handleMouseEvents() {
	x, y := ebiten.CursorPosition()
	ct.handleMouse(x, y)
	ct.handleMouseEnterLeave(x, y)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ct.handleMouseButtonLeftPressed(x, y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ct.handleMouseButtonLeftReleased(x, y)
	}
}

// setFrame устанавливает границы контейнера.
func (ct *containerEmbed) setFrame(frame image.Rectangle) {
	ct.frame = frame
	ct.isDirty = true
}

// childFrame возвращает границы дочернего элемента с учетом его позиции.
func (ct *containerEmbed) childFrame(c *child) *image.Rectangle {
	if !c.absolute {
		r := c.bounds.Add(ct.frame.Min)
		return &r
	}
	return &c.bounds
}

// touchPosition хранит координаты касания.
type touchPosition struct {
	X, Y int
}

// touchPositions хранит последние позиции касаний.
var touchPositions = make(map[ebiten.TouchID]touchPosition)

// recordTouchPosition записывает позицию касания.
func recordTouchPosition(t ebiten.TouchID, x, y int) {
	touchPositions[t] = touchPosition{x, y}
}

// lastTouchPosition возвращает последнюю записанную позицию касания.
func lastTouchPosition(t ebiten.TouchID) *touchPosition {
	if pos, ok := touchPositions[t]; ok {
		return &pos
	}
	return &touchPosition{0, 0}
}
