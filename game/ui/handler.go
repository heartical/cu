package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Handler представляет интерфейс для обработки событий.
type Handler interface{}

// Drawer определяет интерфейс для отрисовки элементов UI.
type Drawer interface {
	Draw(screen *ebiten.Image, frame image.Rectangle, v *View)
}

// Updater определяет интерфейс для обновления элементов UI.
type Updater interface {
	Update(v *View)
}

// DrawHandler определяет интерфейс для обработки событий отрисовки.
type DrawHandler interface {
	HandleDraw(screen *ebiten.Image, frame image.Rectangle)
}

// UpdateHandler определяет интерфейс для обработки событий обновления.
type UpdateHandler interface {
	HandleUpdate()
}

// ButtonHandler определяет интерфейс для обработки событий нажатия и отпускания кнопки.
type ButtonHandler interface {
	HandlePress(x, y int, t ebiten.TouchID)
	HandleRelease(x, y int, isCancel bool)
}

// NotButton определяет интерфейс для элементов, которые не являются кнопками.
type NotButton interface {
	IsButton() bool
}

// TouchHandler определяет интерфейс для обработки событий касания.
type TouchHandler interface {
	HandleJustPressedTouchID(touch ebiten.TouchID, x, y int) bool
	HandleJustReleasedTouchID(touch ebiten.TouchID, x, y int)
}

// MouseHandler определяет интерфейс для обработки событий мыши.
type MouseHandler interface {
	HandleMouse(x, y int) bool
}

// MouseLeftButtonHandler определяет интерфейс для обработки событий левой кнопки мыши.
type MouseLeftButtonHandler interface {
	HandleJustPressedMouseButtonLeft(x, y int) bool
	HandleJustReleasedMouseButtonLeft(x, y int)
}

// MouseEnterLeaveHandler определяет интерфейс для обработки событий входа и выхода мыши.
type MouseEnterLeaveHandler interface {
	HandleMouseEnter(x, y int) bool
	HandleMouseLeave()
}

// SwipeDirection определяет направление свайпа.
type SwipeDirection int

const (
	SwipeLeft SwipeDirection = iota
	SwipeRight
	SwipeUp
	SwipeDown
)

// SwipeHandler определяет интерфейс для обработки событий свайпа.
type SwipeHandler interface {
	HandleSwipe(dir SwipeDirection)
}

// handler реализует интерфейсы для обработки событий.
type handler struct {
	opts HandlerOpts
}

// HandlerOpts содержит функции для обработки событий.
type HandlerOpts struct {
	Update        func(v *View)                                              // Функция для обновления.
	Draw          func(screen *ebiten.Image, frame image.Rectangle, v *View) // Функция для отрисовки.
	HandlePress   func(x, y int, t ebiten.TouchID)                           // Функция для обработки нажатия.
	HandleRelease func(x, y int, isCancel bool)                              // Функция для обработки отпускания.
}

// NewHandler создает новый обработчик событий.
func NewHandler(opts HandlerOpts) Handler {
	return &handler{opts: opts}
}

// Update вызывает функцию обновления, если она задана.
func (h *handler) Update(v *View) {
	if h.opts.Update != nil {
		h.opts.Update(v)
	}
}

// Draw вызывает функцию отрисовки, если она задана.
func (h *handler) Draw(screen *ebiten.Image, frame image.Rectangle, v *View) {
	if h.opts.Draw != nil {
		h.opts.Draw(screen, frame, v)
	}
}

// HandlePress вызывает функцию обработки нажатия, если она задана.
func (h *handler) HandlePress(x, y int, t ebiten.TouchID) {
	if h.opts.HandlePress != nil {
		h.opts.HandlePress(x, y, t)
	}
}

// HandleRelease вызывает функцию обработки отпускания, если она задана.
func (h *handler) HandleRelease(x, y int, isCancel bool) {
	if h.opts.HandleRelease != nil {
		h.opts.HandleRelease(x, y, isCancel)
	}
}
