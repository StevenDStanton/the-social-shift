package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type LevelStruct rune

type Updatable interface {
	Update()
}

type Drawable interface {
	Draw(screen *ebiten.Image)
}

type Game struct {
	AudioContext *audio.Context
	screenWidth  int
	screenHeight int
	Components   []interface{}
}

func New(screenWidth, screenHeight int) *Game {
	g := &Game{
		AudioContext: audio.NewContext(44100),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
	return g
}

func (g *Game) Update() error {
	for _, c := range g.Components {
		if updatable, ok := c.(Updatable); ok {
			updatable.Update()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, c := range g.Components {
		if updatable, ok := c.(Drawable); ok {
			updatable.Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
