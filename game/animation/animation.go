package animation

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// imageCache представляет кэш для хранения подкадров изображений.
type imageCache map[*ebiten.Image]map[*image.Rectangle]*ebiten.Image

var (
	imageCacheMap imageCache
	defaultDelta  = time.Millisecond * 16 // Дефолтное значение для обновления анимации.
)

func init() {
	imageCacheMap = make(imageCache)
}

// parseDurations парсит длительности кадров из различных типов данных.
func parseDurations(durations interface{}, frameCount int) []time.Duration {
	result := make([]time.Duration, frameCount)
	switch val := durations.(type) {
	case time.Duration:
		for i := range result {
			result[i] = val
		}
	case []time.Duration:
		copy(result, val)
	case []interface{}:
		for i := range val {
			result[i] = parseDuration(val[i])
		}
	case map[string]time.Duration:
		for key, duration := range val {
			min, max, step := parseInterval(key)
			for i := min; i <= max; i += step {
				result[i-1] = duration
			}
		}
	case map[string]interface{}:
		for key, duration := range val {
			min, max, step := parseInterval(key)
			for i := min; i <= max; i += step {
				result[i-1] = parseDuration(duration)
			}
		}
	default:
		log.Fatalf("failed to parse durations: type=%T val=%+v", durations, durations)
	}
	return result
}

// parseDuration парсит длительность из различных типов данных.
func parseDuration(value interface{}) time.Duration {
	switch val := value.(type) {
	case time.Duration:
		return val
	case int, float64:
		return time.Millisecond * time.Duration(val.(int))
	default:
		log.Fatalf("failed to parse duration value: %+v", value)
	}
	return 0
}

// parseIntervals преобразует длительности кадров в интервалы времени.
func parseIntervals(durations []time.Duration) ([]time.Duration, time.Duration) {
	result := []time.Duration{0}
	var total time.Duration
	for _, v := range durations {
		total += v
		result = append(result, total)
	}
	return result, total
}

// Status представляет состояние анимации.
type Status int

const (
	Playing Status = iota // Анимация воспроизводится.
	Paused                // Анимация на паузе.
)

// Animation представляет анимацию, состоящую из кадров спрайта.
type Animation struct {
	sprite        *Sprite
	position      int
	timer         time.Duration
	durations     []time.Duration
	intervals     []time.Duration
	totalDuration time.Duration
	onLoop        OnLoop
	status        Status
}

// OnLoop представляет функцию обратного вызова, которая вызывается при завершении цикла анимации.
type OnLoop func(*Animation, int)

// Nop — пустая функция обратного вызова.
func Nop(*Animation, int) {}

// Pause приостанавливает анимацию при завершении цикла.
func Pause(anim *Animation, _ int) {
	anim.Pause()
}

// PauseAtEnd приостанавливает анимацию в конце цикла.
func PauseAtEnd(anim *Animation, _ int) {
	anim.PauseAtEnd()
}

// PauseAtStart приостанавливает анимацию в начале цикла.
func PauseAtStart(anim *Animation, _ int) {
	anim.PauseAtStart()
}

// NewAnimation создает новую анимацию на основе спрайта и длительностей кадров.
func NewAnimation(sprite *Sprite, durations interface{}, onLoop ...OnLoop) *Animation {
	durs := parseDurations(durations, sprite.Length())
	intervals, total := parseIntervals(durs)
	ol := Nop
	if len(onLoop) > 0 {
		ol = onLoop[0]
	}
	return &Animation{
		sprite:        sprite,
		durations:     durs,
		intervals:     intervals,
		totalDuration: total,
		onLoop:        ol,
		status:        Playing,
	}
}

// New создает новую анимацию на основе изображения, кадров и длительностей.
func New(img *ebiten.Image, frames []*image.Rectangle, durations interface{}, onLoop ...OnLoop) *Animation {
	return NewAnimation(NewSprite(img, frames), durations, onLoop...)
}

// Clone создает копию анимации.
func (anim *Animation) Clone() *Animation {
	return &Animation{
		sprite:        anim.sprite,
		position:      anim.position,
		timer:         anim.timer,
		durations:     anim.durations,
		intervals:     anim.intervals,
		totalDuration: anim.totalDuration,
		onLoop:        anim.onLoop,
		status:        anim.status,
	}
}

// SetOnLoop устанавливает функцию обратного вызова при завершении цикла анимации.
func (anim *Animation) SetOnLoop(onLoop OnLoop) {
	anim.onLoop = onLoop
}

// IsEnd проверяет, завершена ли анимация.
func (anim *Animation) IsEnd() bool {
	return anim.status == Paused && anim.position == anim.sprite.Length()-1
}

// seekFrameIndex ищет индекс кадра на основе текущего времени.
func seekFrameIndex(intervals []time.Duration, timer time.Duration) int {
	low, high := 0, len(intervals)-2
	for low <= high {
		mid := (low + high) / 2
		if timer >= intervals[mid+1] {
			low = mid + 1
		} else if timer < intervals[mid] {
			high = mid - 1
		} else {
			return mid
		}
	}
	return low
}

// Update обновляет анимацию с использованием дефолтного времени.
func (anim *Animation) Update() {
	anim.UpdateWithDelta(defaultDelta)
}

// UpdateWithDelta обновляет анимацию с учетом прошедшего времени.
func (anim *Animation) UpdateWithDelta(elapsed time.Duration) {
	if anim.status != Playing || anim.sprite.Length() <= 1 {
		return
	}
	anim.timer += elapsed
	loops := anim.timer / anim.totalDuration
	if loops != 0 {
		anim.timer -= anim.totalDuration * loops
		anim.onLoop(anim, int(loops))
	}
	anim.position = seekFrameIndex(anim.intervals, anim.timer)
}

// SetDurations устанавливает новые длительности кадров.
func (anim *Animation) SetDurations(durations interface{}) {
	durs := parseDurations(durations, anim.sprite.Length())
	anim.durations = durs
	anim.intervals, anim.totalDuration = parseIntervals(durs)
	anim.timer = 0
}

// Status возвращает текущее состояние анимации.
func (anim *Animation) Status() Status {
	return anim.status
}

// Pause приостанавливает анимацию.
func (anim *Animation) Pause() {
	anim.status = Paused
}

// Position возвращает текущий кадр анимации.
func (anim *Animation) Position() int {
	return anim.position + 1
}

// Durations возвращает длительности кадров анимации.
func (anim *Animation) Durations() []time.Duration {
	return anim.durations
}

// TotalDuration возвращает общую длительность анимации.
func (anim *Animation) TotalDuration() time.Duration {
	return anim.totalDuration
}

// Size возвращает размер спрайта анимации.
func (anim *Animation) Size() (int, int) {
	return anim.sprite.Size()
}

// Width возвращает ширину спрайта анимации.
func (anim *Animation) Width() int {
	return anim.sprite.Width()
}

// Height возвращает высоту спрайта анимации.
func (anim *Animation) Height() int {
	return anim.sprite.Height()
}

// Timer возвращает текущее время анимации.
func (anim *Animation) Timer() time.Duration {
	return anim.timer
}

// Sprite возвращает спрайт анимации.
func (anim *Animation) Sprite() *Sprite {
	return anim.sprite
}

// GoToFrame переходит к указанному кадру анимации.
func (anim *Animation) GoToFrame(position int) {
	anim.position = position - 1
	anim.timer = anim.intervals[anim.position]
}

// PauseAtEnd приостанавливает анимацию в конце.
func (anim *Animation) PauseAtEnd() {
	anim.position = anim.sprite.Length() - 1
	anim.timer = anim.totalDuration
	anim.Pause()
}

// PauseAtStart приостанавливает анимацию в начале.
func (anim *Animation) PauseAtStart() {
	anim.position = 0
	anim.timer = 0
	anim.status = Paused
}

// Resume возобновляет воспроизведение анимации.
func (anim *Animation) Resume() {
	anim.status = Playing
}

// Draw отрисовывает текущий кадр анимации на экране.
func (anim *Animation) Draw(screen *ebiten.Image, opts *DrawOptions) {
	anim.sprite.Draw(screen, anim.position, opts)
}

// DrawWithShader отрисовывает текущий кадр анимации с использованием шейдера.
func (anim *Animation) DrawWithShader(screen *ebiten.Image, opts *DrawOptions, shaderOpts *ShaderOptions) {
	anim.sprite.DrawWithShader(screen, anim.position, opts, shaderOpts)
}
