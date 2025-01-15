package main

import (
	"context"
	"errors"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"cu/common/assets"
	"cu/common/e2e"
	"cu/game/sprites"
	"cu/game/ui"
	"cu/game/widgets"
)

// Game представляет собой основную структуру игры.
type Game struct {
	screen   screen
	gameUI   *ui.View
	initOnce sync.Once
}

// screen содержит размеры экрана.
type screen struct {
	Width  int
	Height int
}

// Update обновляет состояние игры.
func (g *Game) Update() error {
	g.initOnce.Do(g.setupUI) // Инициализация UI при первом вызове.
	g.gameUI.UpdateWithSize(g.screen.Width, g.screen.Height)
	return nil
}

// Draw отрисовывает игровой экран.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{63, 124, 182, 255}) // Заливаем экран цветом.
	g.gameUI.Draw(screen)                      // Отрисовываем UI.
}

// Layout задает размеры окна игры.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screen.Width = outsideWidth
	g.screen.Height = outsideHeight
	return g.screen.Width, g.screen.Height
}

// NewGame создает новый экземпляр игры и инициализирует ресурсы.
func NewGame() (*Game, error) {
	resources, err := assets.Init()
	if err != nil {
		return nil, err
	}

	// Загружаем шрифты.
	assets.LoadFonts(resources.FontScoreDozer)

	// Загружаем спрайты из RPG-пака.
	sprites.LoadSprites(
		resources.UIPackRPGSheetXML,
		resources.UIPackRPGSheetPNG,
		sprites.LoadOpts{
			PanelOpts: map[string]sprites.PanelOpts{
				"panelInset_beige.png": {Border: 32, Center: 36},
				"panel_brown.png":      {Border: 32, Center: 36},
			},
		},
	)

	// Загружаем спрайты из Space-пака.
	sprites.LoadSprites(
		resources.UIPackSpaceSheetXML,
		resources.UIPackSpaceSheetPNG,
		sprites.LoadOpts{
			PanelOpts: map[string]sprites.PanelOpts{
				"glassPanel_corners.png":    {Border: 40, Center: 20},
				"glassPanel_projection.png": {Border: 20, Center: 10},
			},
		},
	)

	return &Game{}, nil
}

// init регистрирует компоненты UI.
func init() {
	ui.RegisterComponents(ui.ComponentsMap{
		"panel":  &widgets.Panel{},
		"sprite": &widgets.Sprite{},
	})
}

// setupUI инициализирует пользовательский интерфейс игры.
func (g *Game) setupUI() {
	ui.Debug = true // Включаем отладочный режим.

	var duration time.Duration
	var charIndex int

	g.gameUI = ui.Parse(string(assets.IndexHTML), &ui.ParseOptions{
		Width:  g.screen.Width,
		Height: g.screen.Height,
		Components: ui.ComponentsMap{
			"panel": &widgets.Panel{},
			"gauge-text": func() *ui.View {
				return &ui.View{
					Width:   180,
					Height:  20,
					Handler: &widgets.Text{Color: color.RGBA{50, 48, 41, 255}},
				}
			},
			"gauge": func() ui.Handler { return &widgets.Bar{Value: .8} },
			"button": func() ui.Handler {
				return &widgets.Button{OnClick: func() { println("button clicked") }}
			},
			"bottom-button": func() *ui.View {
				return &ui.View{
					Handler: &widgets.Button{
						Color:   color.RGBA{210, 178, 144, 255},
						OnClick: func() { println("button clicked") },
					},
				}
			},
		},
		Handler: ui.NewHandler(ui.HandlerOpts{
			Update: func(v *ui.View) {
				duration += time.Second / 60
				switch {
				case charIndex < len(playGameText) && duration > time.Millisecond*100:
					charIndex++
					duration = 0
				case duration > time.Millisecond*1000:
					charIndex = 0
					duration = 0
				}
			},
		}),
	})
}

const (
	defaultWindowWidth  = 720
	defaultWindowHeight = 1280
	playGameText        = "Do you play game?"
)

func main() {
	game, err := NewGame()
	if err != nil {
		panic(err)
	}

	// Настройка параметров окна и игрового цикла.
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)

	// Определяем масштаб экрана.
	deviceScaleFactor := ebiten.Monitor().DeviceScaleFactor()
	if deviceScaleFactor == 0.0 {
		deviceScaleFactor = 1.0
	}

	screenWidth, screenHeight := ebiten.Monitor().Size()
	if deviceScaleFactor == 1.0 && ((screenWidth > 3000 && screenHeight > 2000) || (screenWidth > 2000 && screenHeight > 3000)) {
		deviceScaleFactor = 1.15
	}

	// Устанавливаем размеры окна с учетом масштаба.
	windowWidth := int(float64(defaultWindowWidth) * deviceScaleFactor)
	windowHeight := int(float64(defaultWindowHeight) * deviceScaleFactor)
	ebiten.SetWindowSize(windowWidth, windowHeight)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Освобождаем ресурсы контекста

	// Инициализируем FSM и запускаем его
	e2e.E2EE(ctx)

	// Запускаем игровой цикл.
	if err := ebiten.RunGame(game); err != nil && !errors.Is(err, ebiten.Termination) {
		panic(err)
	}
}
