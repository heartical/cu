package ui

import (
	"fmt"
	"image"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// View представляет собой компонент UI, который может содержать другие компоненты и управлять их отрисовкой и обновлением.
type View struct {
	Left, Top, Width, Height, MarginLeft, MarginTop, MarginRight, MarginBottom int
	Right, Bottom                                                              *int
	WidthInPct, HeightInPct, Grow, Shrink                                      float64
	Position                                                                   Position
	Direction                                                                  Direction
	Wrap                                                                       FlexWrap
	Justify                                                                    Justify
	AlignItems                                                                 AlignItem
	AlignContent                                                               AlignContent
	Display                                                                    Display
	ID, Raw, TagName, Text                                                     string
	Attrs                                                                      map[string]string
	Hidden                                                                     bool
	Handler                                                                    Handler
	containerEmbed
	flexEmbed
	lock      sync.Mutex
	hasParent bool
	parent    *View
}

// Update обновляет состояние View и его дочерних элементов.
func (v *View) Update() {
	if v.isDirty {
		v.startLayout()
	}
	if !v.hasParent {
		v.processHandler()
	}
	for _, child := range v.children {
		child.item.Update()
		child.item.processHandler()
	}
	if !v.hasParent {
		v.processEvent()
	}
}

// processHandler обрабатывает обновления, если View имеет обработчик.
func (v *View) processHandler() {
	if u, ok := v.Handler.(UpdateHandler); ok {
		u.HandleUpdate()
	} else if u, ok := v.Handler.(Updater); ok {
		u.Update(v)
	}
}

// startLayout запускает процесс компоновки View и его дочерних элементов.
func (v *View) startLayout() {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.hasParent {
		v.frame = image.Rect(v.Left, v.Top, v.Left+v.Width, v.Top+v.Height)
	}
	v.flexEmbed.View = v

	for _, child := range v.children {
		if child.item.Position == PositionStatic {
			child.item.startLayout()
		}
	}

	v.layout(v.frame.Dx(), v.frame.Dy(), &v.containerEmbed)
	v.isDirty = false
}

// UpdateWithSize обновляет размеры View и запускает его обновление.
func (v *View) UpdateWithSize(width, height int) {
	if !v.hasParent && (v.Width != width || v.Height != height) {
		v.Height, v.Width = height, width
		v.isDirty = true
	}
	v.Update()
}

// Layout помечает View как "грязный", что приводит к перекомпоновке при следующем обновлении.
func (v *View) Layout() {
	v.isDirty = true
	if v.hasParent {
		v.parent.isDirty = true
	}
}

// Draw отрисовывает View и его дочерние элементы на экране.
func (v *View) Draw(screen *ebiten.Image) {
	if v.isDirty {
		v.startLayout()
	}
	if !v.hasParent {
		v.handleDrawRoot(screen, v.frame)
	}
	if !v.Hidden && v.Display != DisplayNone {
		v.containerEmbed.Draw(screen)
	}
	if Debug && !v.hasParent && v.Display != DisplayNone {
		debugBorders(screen, v.containerEmbed)
	}
}

// AddTo добавляет View к родительскому элементу.
func (v *View) AddTo(parent *View) *View {
	if v.hasParent {
		panic("view already has a parent")
	}
	parent.AddChild(v)
	return v
}

// AddChild добавляет дочерние элементы к View.
func (v *View) AddChild(views ...*View) *View {
	for _, view := range views {
		v.addChild(view)
	}
	return v
}

// RemoveChild удаляет дочерний элемент из View.
func (v *View) RemoveChild(cv *View) bool {
	for i, child := range v.children {
		if child.item == cv {
			v.children = append(v.children[:i], v.children[i+1:]...)
			v.isDirty = true
			cv.hasParent, cv.parent = false, nil
			return true
		}
	}
	return false
}

// RemoveAll удаляет все дочерние элементы из View.
func (v *View) RemoveAll() {
	v.isDirty = true
	for _, child := range v.children {
		child.item.hasParent, child.item.parent = false, nil
	}
	v.children = []*child{}
}

// PopChild удаляет и возвращает последний дочерний элемент.
func (v *View) PopChild() *View {
	if len(v.children) == 0 {
		return nil
	}
	c := v.children[len(v.children)-1]
	v.children = v.children[:len(v.children)-1]
	v.isDirty = true
	c.item.hasParent, c.item.parent = false, nil
	return c.item
}

// addChild добавляет дочерний элемент к View.
func (v *View) addChild(cv *View) *View {
	child := &child{item: cv, handledTouchID: -1}
	v.children = append(v.children, child)
	v.isDirty = true
	cv.hasParent, cv.parent = true, v
	return v
}

// isWidthFixed проверяет, фиксирована ли ширина View.
func (v *View) isWidthFixed() bool {
	return v.Width != 0 || v.WidthInPct != 0
}

// width возвращает ширину View.
func (v *View) width() int {
	if v.Width == 0 {
		return v.calculatedWidth
	}
	return v.Width
}

// isHeightFixed проверяет, фиксирована ли высота View.
func (v *View) isHeightFixed() bool {
	return v.Height != 0 || v.HeightInPct != 0
}

// height возвращает высоту View.
func (v *View) height() int {
	if v.Height == 0 {
		return v.calculatedHeight
	}
	return v.Height
}

// getChildren возвращает список дочерних элементов View.
func (v *View) getChildren() []*View {
	if v == nil || v.children == nil {
		return nil
	}
	ret := make([]*View, len(v.children))
	for i, child := range v.children {
		ret[i] = child.item
	}
	return ret
}

// GetByID ищет View по его ID.
func (v *View) GetByID(id string) (*View, bool) {
	if v.ID == id {
		return v, true
	}
	for _, child := range v.children {
		if view, ok := child.item.GetByID(id); ok {
			return view, true
		}
	}
	return nil, false
}

// MustGetByID ищет View по его ID и вызывает панику, если View не найден.
func (v *View) MustGetByID(id string) *View {
	view, ok := v.GetByID(id)
	if !ok {
		panic("view not found")
	}
	return view
}

// SetLeft устанавливает левую позицию View.
func (v *View) SetLeft(left int) { v.Left = left; v.Layout() }

// SetRight устанавливает правую позицию View.
func (v *View) SetRight(right int) { v.Right = Int(right); v.Layout() }

// SetTop устанавливает верхнюю позицию View.
func (v *View) SetTop(top int) { v.Top = top; v.Layout() }

// SetBottom устанавливает нижнюю позицию View.
func (v *View) SetBottom(bottom int) { v.Bottom = Int(bottom); v.Layout() }

// SetWidth устанавливает ширину View.
func (v *View) SetWidth(width int) { v.Width = width; v.Layout() }

// SetHeight устанавливает высоту View.
func (v *View) SetHeight(height int) { v.Height = height; v.Layout() }

// SetMarginLeft устанавливает левый отступ View.
func (v *View) SetMarginLeft(marginLeft int) { v.MarginLeft = marginLeft; v.Layout() }

// SetMarginTop устанавливает верхний отступ View.
func (v *View) SetMarginTop(marginTop int) { v.MarginTop = marginTop; v.Layout() }

// SetMarginRight устанавливает правый отступ View.
func (v *View) SetMarginRight(marginRight int) { v.MarginRight = marginRight; v.Layout() }

// SetMarginBottom устанавливает нижний отступ View.
func (v *View) SetMarginBottom(marginBottom int) { v.MarginBottom = marginBottom; v.Layout() }

// SetPosition устанавливает позицию View.
func (v *View) SetPosition(position Position) { v.Position = position; v.Layout() }

// SetDirection устанавливает направление компоновки View.
func (v *View) SetDirection(direction Direction) { v.Direction = direction; v.Layout() }

// SetWrap устанавливает режим переноса для View.
func (v *View) SetWrap(wrap FlexWrap) { v.Wrap = wrap; v.Layout() }

// SetJustify устанавливает выравнивание по главной оси для View.
func (v *View) SetJustify(justify Justify) { v.Justify = justify; v.Layout() }

// SetAlignItems устанавливает выравнивание по поперечной оси для View.
func (v *View) SetAlignItems(alignItems AlignItem) { v.AlignItems = alignItems; v.Layout() }

// SetAlignContent устанавливает выравнивание содержимого для View.
func (v *View) SetAlignContent(alignContent AlignContent) { v.AlignContent = alignContent; v.Layout() }

// SetGrow устанавливает коэффициент растяжения для View.
func (v *View) SetGrow(grow float64) { v.Grow = grow; v.Layout() }

// SetShrink устанавливает коэффициент сжатия для View.
func (v *View) SetShrink(shrink float64) { v.Shrink = shrink; v.Layout() }

// SetDisplay устанавливает режим отображения для View.
func (v *View) SetDisplay(display Display) { v.Display = display; v.Layout() }

// SetHidden скрывает или показывает View.
func (v *View) SetHidden(hidden bool) { v.Hidden = hidden; v.Layout() }

// Config возвращает конфигурацию View.
func (v *View) Config() ViewConfig {
	cfg := ViewConfig{
		TagName: v.TagName, ID: v.ID, Left: v.Left, Right: v.Right, Top: v.Top, Bottom: v.Bottom,
		Width: v.Width, Height: v.Height, MarginLeft: v.MarginLeft, MarginTop: v.MarginTop,
		MarginRight: v.MarginRight, MarginBottom: v.MarginBottom, Position: v.Position,
		Direction: v.Direction, Wrap: v.Wrap, Justify: v.Justify, AlignItems: v.AlignItems,
		AlignContent: v.AlignContent, Grow: v.Grow, Shrink: v.Shrink, children: []ViewConfig{},
	}
	for _, child := range v.getChildren() {
		cfg.children = append(cfg.children, child.Config())
	}
	return cfg
}

// handleDrawRoot обрабатывает отрисовку корневого View.
func (v *View) handleDrawRoot(screen *ebiten.Image, b image.Rectangle) {
	if h, ok := v.Handler.(DrawHandler); ok {
		h.HandleDraw(screen, b)
	} else if h, ok := v.Handler.(Drawer); ok {
		h.Draw(screen, b, v)
	}
}

// ViewConfig представляет конфигурацию View.
type ViewConfig struct {
	TagName, ID                                                                string
	Left, Top, Width, Height, MarginLeft, MarginTop, MarginRight, MarginBottom int
	Right, Bottom                                                              *int
	Position                                                                   Position
	Direction                                                                  Direction
	Wrap                                                                       FlexWrap
	Justify                                                                    Justify
	AlignItems                                                                 AlignItem
	AlignContent                                                               AlignContent
	Grow, Shrink                                                               float64
	children                                                                   []ViewConfig
}

// Tree возвращает строковое представление дерева View.
func (cfg ViewConfig) Tree() string {
	return cfg.tree("")
}

// tree возвращает строковое представление дерева View с отступами.
func (cfg ViewConfig) tree(indent string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s<%s ", indent, cfg.TagName))
	if cfg.ID != "" {
		sb.WriteString(fmt.Sprintf("id=\"%s\" ", cfg.ID))
	}
	sb.WriteString("style=\"")
	styles := []string{
		fmt.Sprintf("left: %d", cfg.Left), fmt.Sprintf("right: %d", *cfg.Right),
		fmt.Sprintf("top: %d", cfg.Top), fmt.Sprintf("bottom: %d", *cfg.Bottom),
		fmt.Sprintf("width: %d", cfg.Width), fmt.Sprintf("height: %d", cfg.Height),
		fmt.Sprintf("margin-left: %d", cfg.MarginLeft), fmt.Sprintf("margin-top: %d", cfg.MarginTop),
		fmt.Sprintf("margin-right: %d", cfg.MarginRight), fmt.Sprintf("margin-bottom: %d", cfg.MarginBottom),
		fmt.Sprintf("position: %s", cfg.Position), fmt.Sprintf("direction: %s", cfg.Direction),
		fmt.Sprintf("wrap: %s", cfg.Wrap), fmt.Sprintf("justify: %s", cfg.Justify),
		fmt.Sprintf("align-items: %s", cfg.AlignItems), fmt.Sprintf("align-content: %s", cfg.AlignContent),
		fmt.Sprintf("grow: %f", cfg.Grow), fmt.Sprintf("shrink: %f", cfg.Shrink),
	}
	sb.WriteString(strings.Join(styles, "; "))
	sb.WriteString("\">\n")
	for _, child := range cfg.children {
		sb.WriteString(child.tree(indent + "  "))
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("%s</%s>", indent, cfg.TagName))
	sb.WriteString("\n")
	return sb.String()
}
