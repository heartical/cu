package main

import (
	"engine/assets"
	"engine/game/sprites"
	"engine/game/text"
	"engine/game/ui"
	"engine/game/widgets"
	"sync"

	"errors"
	"image/color"

	// "syscall/js"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	screen screen
	gameUI *ui.View

	initOnce sync.Once
}

type screen struct {
	Width  int
	Height int
}

func (g *Game) Update() error {
	g.initOnce.Do(func() { g.setupUI() })
	g.gameUI.UpdateWithSize(g.screen.Width, g.screen.Height)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{63, 124, 182, 255})
	g.gameUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screen.Width = outsideWidth
	g.screen.Height = outsideHeight
	return g.screen.Width, g.screen.Height
}

func NewGame() (*Game, error) {

	if err := assets.Init(); err != nil {
		return nil, err
	}

	text.LoadFonts(assets.Font)
	sprites.LoadSprites(
		assets.UIPackRPGSheetXML,
		assets.UIPackRPGSheetPNG,
		sprites.LoadOpts{
			PanelOpts: map[string]sprites.PanelOpts{
				"panelInset_beige.png": {
					Border: 32,
					Center: 36,
				},
				"panel_brown.png": {
					Border: 32,
					Center: 36,
				},
			},
		})
	sprites.LoadSprites(
		assets.UIPackSpaceSheetXML,
		assets.UIPackSpaceSheetPNG,
		sprites.LoadOpts{
			PanelOpts: map[string]sprites.PanelOpts{
				"glassPanel_corners.png": {
					Border: 40,
					Center: 20,
				},
				"glassPanel_projection.png": {
					Border: 20,
					Center: 10,
				},
			},
		})
	game := &Game{}
	return game, nil
}

func init() {
	ui.RegisterComponents(ui.ComponentsMap{
		"panel":  &widgets.Panel{},
		"sprite": &widgets.Sprite{},
	})
}

const playGameText = "Do you play game?"

func (g *Game) setupUI() {
	ui.Debug = true

	d := time.Duration(0)
	c := 0

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
					Width:  45,
					Height: 49,
					Handler: &widgets.Button{
						Color:   color.RGBA{210, 178, 144, 255},
						OnClick: func() { println("button clicked") },
					}}
			},
			"panel-button": func() *ui.View {
				return &ui.View{
					Width:   100,
					Height:  50,
					Handler: &widgets.Panel{OnClick: func() { println("button clicked") }},
				}
			},
		},
		Handler: ui.NewHandler(ui.HandlerOpts{
			Update: func(v *ui.View) {
				d += time.Second / 60
				switch {
				case c < len(playGameText) && d > time.Millisecond*100:
					c = c + 1
					d = 0
				case d > time.Millisecond*1000:
					c = 0
					d = 0
				}
			},
		}),
	})
}

const defaultWindowWidth = 720
const defaultWindowHeight = 1280

func main() {
	game, _ := NewGame()
	// jsRoomID := js.Global().Get("roomID")

	ebiten.SetScreenClearedEveryFrame(true)

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)

	deviceScaleFactor := ebiten.Monitor().DeviceScaleFactor()
	if deviceScaleFactor == 0.0 {
		deviceScaleFactor = 1.0
	}

	screenWidth, screenHeight := ebiten.Monitor().Size()
	if deviceScaleFactor == 1.0 && ((screenWidth > 3000 && screenHeight > 2000) || (screenWidth > 2000 && screenHeight > 3000)) {
		deviceScaleFactor = 1.15
	}

	windowWidth := int(float64(defaultWindowWidth) * deviceScaleFactor)
	windowHeight := int(float64(defaultWindowHeight) * deviceScaleFactor)
	ebiten.SetWindowSize(windowWidth, windowHeight)

	// if jsRoomID.Type() == js.TypeString {
	// 	fmt.Println(jsRoomID.String())
	// 	// game.roomID = jsRoomID.String()
	// } else {
	// 	fmt.Println("Unknown Room ID")
	// 	// game.roomID = "Unknown Room ID"
	// }

	if err := ebiten.RunGame(game); err != nil && !errors.Is(err, ebiten.Termination) {
		panic(err)
	}
}
