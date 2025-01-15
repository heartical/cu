package ui

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// child представляет дочерний элемент UI, который может обрабатывать события касания, свайпа и нажатия кнопок.
type child struct {
	absolute                 bool
	item                     *View
	bounds                   image.Rectangle
	isButtonPressed          bool
	isMouseLeftButtonHandler bool
	isMouseEntered           bool
	handledTouchID           ebiten.TouchID
	swipe
}

// swipe содержит информацию о свайпе.
type swipe struct {
	downX, downY int
	upX, upY     int
	downTime     time.Time
	upTime       time.Time
	swipeDir     SwipeDirection
	swipeTouchID ebiten.TouchID
}

// HandleJustPressedTouchID обрабатывает событие начала касания.
func (c *child) HandleJustPressedTouchID(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	result := false
	if c.checkButtonHandlerStart(frame, touchID, x, y) {
		result = true
	}
	if !result && c.checkTouchHandlerStart(frame, touchID, x, y) {
		result = true
	}
	c.checkSwipeHandlerStart(frame, touchID, x, y)
	return result
}

// HandleJustReleasedTouchID обрабатывает событие окончания касания.
func (c *child) HandleJustReleasedTouchID(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	c.checkTouchHandlerEnd(frame, touchID, x, y)
	c.checkButtonHandlerEnd(frame, touchID, x, y)
	c.checkSwipeHandlerEnd(frame, touchID, x, y)
}

// checkTouchHandlerStart проверяет и обрабатывает начало касания для TouchHandler.
func (c *child) checkTouchHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	if handler, ok := c.item.Handler.(TouchHandler); ok && isInside(frame, x, y) {
		if handler.HandleJustPressedTouchID(touchID, x, y) {
			c.handledTouchID = touchID
			return true
		}
	}
	return false
}

// checkTouchHandlerEnd проверяет и обрабатывает окончание касания для TouchHandler.
func (c *child) checkTouchHandlerEnd(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	if handler, ok := c.item.Handler.(TouchHandler); ok && c.handledTouchID == touchID {
		handler.HandleJustReleasedTouchID(touchID, x, y)
		c.handledTouchID = -1
	}
}

// checkSwipeHandlerStart проверяет и обрабатывает начало свайпа.
func (c *child) checkSwipeHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	if _, ok := c.item.Handler.(SwipeHandler); ok && isInside(frame, x, y) {
		c.swipeTouchID = touchID
		c.downTime = time.Now()
		c.downX, c.downY = x, y
		return true
	}
	return false
}

// checkSwipeHandlerEnd проверяет и обрабатывает окончание свайпа.
func (c *child) checkSwipeHandlerEnd(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	if handler, ok := c.item.Handler.(SwipeHandler); ok && c.swipeTouchID == touchID {
		c.swipeTouchID = -1
		c.upTime = time.Now()
		c.upX, c.upY = x, y
		if c.checkSwipe() {
			handler.HandleSwipe(c.swipeDir)
			return true
		}
	}
	return false
}

const (
	swipeThresholdDist = 50.0                   // Минимальное расстояние для распознавания свайпа.
	swipeThresholdTime = time.Millisecond * 300 // Максимальное время для распознавания свайпа.
)

// checkSwipe проверяет, был ли выполнен свайп.
func (c *child) checkSwipe() bool {
	duration := c.upTime.Sub(c.downTime)
	if duration > swipeThresholdTime {
		return false
	}

	deltaX := float64(c.downX - c.upX)
	if math.Abs(deltaX) >= swipeThresholdDist {
		if deltaX > 0 {
			c.swipeDir = SwipeLeft
		} else {
			c.swipeDir = SwipeRight
		}
		return true
	}

	deltaY := float64(c.downY - c.upY)
	if math.Abs(deltaY) >= swipeThresholdDist {
		if deltaY > 0 {
			c.swipeDir = SwipeUp
		} else {
			c.swipeDir = SwipeDown
		}
		return true
	}

	return false
}

// checkButtonHandlerStart проверяет и обрабатывает начало нажатия кнопки.
func (c *child) checkButtonHandlerStart(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) bool {
	if button, ok := c.item.Handler.(ButtonHandler); ok {
		if notButton, ok := c.item.Handler.(NotButton); ok && !notButton.IsButton() {
			return false
		}
		if isInside(frame, x, y) && !c.isButtonPressed {
			c.isButtonPressed = true
			c.handledTouchID = touchID
			button.HandlePress(x, y, touchID)
			return true
		} else if c.handledTouchID == touchID {
			c.handledTouchID = -1
		}
	}
	return false
}

// checkButtonHandlerEnd проверяет и обрабатывает окончание нажатия кнопки.
func (c *child) checkButtonHandlerEnd(frame *image.Rectangle, touchID ebiten.TouchID, x, y int) {
	if button, ok := c.item.Handler.(ButtonHandler); ok && c.handledTouchID == touchID && c.isButtonPressed {
		c.isButtonPressed = false
		c.handledTouchID = -1
		isCancel := !(x == 0 && y == 0) && !isInside(frame, x, y)
		button.HandleRelease(x, y, isCancel)
	}
}

// isInside проверяет, находится ли точка (x, y) внутри прямоугольника frame.
// func isInside(frame *image.Rectangle, x, y int) bool {
// 	return frame != nil && x >= frame.Min.X && x <= frame.Max.X && y >= frame.Min.Y && y <= frame.Max.Y
// }
