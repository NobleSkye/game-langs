package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var screenWidth = 800
var screenHeight = 600

type Game interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

type MainMenu struct {
	selection int
	games     []string
}

func (m *MainMenu) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		m.selection++
		if m.selection >= len(m.games) {
			m.selection = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		m.selection--
		if m.selection < 0 {
			m.selection = len(m.games) - 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch m.selection {
		case 0:
			// Add cases for additional games here
		}
	}
	return nil
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i, game := range m.games {
		if i == m.selection {
			screen.Fill(color.RGBA{255, 255, 0, 255}) // Highlight selected item
		}
		ebitenutil.DebugPrintAt(screen, game, 100, 50+i*30)
	}
}

func (m *MainMenu) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var currentGame Game = &MainMenu{
	selection: 0,
	games:     []string{"Pong", "Game 2", "Game 3"}, // Add more games here
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game Collection")
	if err := ebiten.RunGame(currentGame); err != nil {
		fmt.Println(err)
	}
}
