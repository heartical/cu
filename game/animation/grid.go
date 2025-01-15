package animation

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"regexp"
	"strconv"
)

// assertPositiveInteger проверяет, что значение является положительным целым числом.
// В противном случае программа завершается с фатальной ошибкой.
func assertPositiveInteger(value int, name string) {
	if value < 1 {
		log.Fatalf("%s должен быть положительным числом, получено: %d", name, value)
	}
}

// assertSize проверяет, что значение не превышает заданный лимит.
// В противном случае программа завершается с фатальной ошибкой.
func assertSize(value, limit int, name string) {
	if value > limit {
		log.Fatalf("%s должен быть <= %d, получено: %d", name, limit, value)
	}
}

// frameCache представляет кэш кадров для сетки.
type frameCache map[string]map[int]map[int]*image.Rectangle

var (
	framesCache     frameCache
	intervalMatcher *regexp.Regexp
)

func init() {
	framesCache = make(frameCache)
	intervalMatcher = regexp.MustCompile(`^(\d+)-(\d+)$`)
}

// Grid представляет сетку кадров для анимации.
type Grid struct {
	frameWidth, frameHeight int    // Ширина и высота одного кадра.
	imageWidth, imageHeight int    // Ширина и высота изображения.
	left, top               int    // Отступы от края изображения.
	width, height           int    // Количество кадров по ширине и высоте.
	border                  int    // Расстояние между кадрами.
	key                     string // Уникальный ключ для кэширования кадров.
}

// NewGrid создает новую сетку кадров.
func NewGrid(frameWidth, frameHeight, imageWidth, imageHeight int, args ...int) *Grid {
	assertPositiveInteger(frameWidth, "frameWidth")
	assertPositiveInteger(frameHeight, "frameHeight")
	assertPositiveInteger(imageWidth, "imageWidth")
	assertPositiveInteger(imageHeight, "imageHeight")
	assertSize(frameWidth, imageWidth, "frameWidth")
	assertSize(frameHeight, imageHeight, "frameHeight")

	left, top, border := 0, 0, 0
	switch len(args) {
	case 3:
		border = args[2]
		fallthrough
	case 2:
		top = args[1]
		fallthrough
	case 1:
		left = args[0]
	}

	grid := &Grid{
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		left:        left,
		top:         top,
		width:       imageWidth / frameWidth,
		height:      imageHeight / frameHeight,
		border:      border,
	}

	grid.key = generateGridKey(grid.frameWidth, grid.frameHeight, grid.imageWidth, grid.imageHeight, grid.left, grid.top)

	return grid
}

// generateGridKey генерирует уникальный ключ для сетки на основе её параметров.
func generateGridKey(args ...int) string {
	var buffer bytes.Buffer
	separator := ""
	for _, arg := range args {
		buffer.WriteString(separator)
		buffer.WriteString(strconv.Itoa(arg))
		separator = "-"
	}
	return buffer.String()
}

// createFrame создает новый кадр на основе координат (x, y).
func (g *Grid) createFrame(x, y int) *image.Rectangle {
	frameWidth, frameHeight := g.frameWidth, g.frameHeight
	x0 := g.left + (x-1)*frameWidth + x*g.border
	y0 := g.top + (y-1)*frameHeight + y*g.border
	return &image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x0 + frameWidth, Y: y0 + frameHeight},
	}
}

// getOrCreateFrame возвращает кадр из кэша или создает новый, если он отсутствует.
func (g *Grid) getOrCreateFrame(x, y int) *image.Rectangle {
	if x < 1 || x > g.width || y < 1 || y > g.height {
		log.Fatalf("Кадр с координатами x=%d, y=%d не существует", x, y)
	}

	key := g.key
	if _, ok := framesCache[key]; !ok {
		framesCache[key] = make(map[int]map[int]*image.Rectangle)
	}
	if _, ok := framesCache[key][x]; !ok {
		framesCache[key][x] = make(map[int]*image.Rectangle)
	}
	if _, ok := framesCache[key][x][y]; !ok {
		framesCache[key][x][y] = g.createFrame(x, y)
	}

	return framesCache[key][x][y]
}

// GetFrames возвращает список кадров в зависимости от переданных аргументов.
func (g *Grid) GetFrames(args ...interface{}) []*image.Rectangle {
	var frames []*image.Rectangle

	if len(args) == 0 {
		// Возвращаем все кадры, если аргументы не переданы.
		for y := 1; y <= g.height; y++ {
			for x := 1; x <= g.width; x++ {
				frames = append(frames, g.getOrCreateFrame(x, y))
			}
		}
		return frames
	}

	// Обрабатываем переданные интервалы.
	for i := 0; i < len(args); i += 2 {
		minX, maxX, stepX := parseInterval(args[i])
		minY, maxY, stepY := parseInterval(args[i+1])

		for y := minY; (stepY > 0 && y <= maxY) || (stepY < 0 && y >= maxY); y += stepY {
			for x := minX; (stepX > 0 && x <= maxX) || (stepX < 0 && x >= maxX); x += stepX {
				frames = append(frames, g.getOrCreateFrame(x, y))
			}
		}
	}

	return frames
}

// Width возвращает количество кадров по ширине сетки.
func (g *Grid) Width() int {
	return g.width
}

// Height возвращает количество кадров по высоте сетки.
func (g *Grid) Height() int {
	return g.height
}

// Frames является псевдонимом для GetFrames.
func (g *Grid) Frames(args ...interface{}) []*image.Rectangle {
	return g.GetFrames(args...)
}

// G является псевдонимом для GetFrames.
func (g *Grid) G(args ...interface{}) []*image.Rectangle {
	return g.GetFrames(args...)
}
