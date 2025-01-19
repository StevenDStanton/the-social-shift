package main

import (
	_ "embed"
	"image/color"
	"log"
	"math"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 1280
	screenHeight = 720

	fontSize = 16
	DPI      = 72
)

var (
	face       font.Face
	rows       = screenHeight / cellHeight
	cols       = screenWidth / cellWidth
	cellWidth  = fontSize
	cellHeight = fontSize
	divider    = int(math.Round(float64(cols) * 0.6))
)

type Player struct {
	x, y int
	walk []*audio.Player
}

type Level int

const (
	LEVEL0 Level = iota
	LEVEL1
	LEVEL2
	LEVEL3
	LEVEL4
)

type Intro struct {
	images []*ebiten.Image
}

type Game struct {
	level            Level
	intro            *Intro
	grid             [][]rune
	audioContext     *audio.Context
	player           *Player
	movementCooldown int
	introImageIndex  int
	introTimer       int
	isLevelLoaded    bool
}

// Update handles the game logic.
func (g *Game) Update() error {

	switch {
	case g.level == LEVEL0:
		g.updateIntro()
		return nil
	case g.level == LEVEL1:
		g.updateLevel()
		return nil
	default:
		return nil
	}

}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	switch {
	case g.level == LEVEL0:
		g.drawIntro(screen)
	default:
		g.drawLevel(screen)
	}

}

// Layout defines the screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("The Social Shift")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	// Now convert that parsed font into a usable Face.

	game := &Game{
		audioContext: audio.NewContext(44100),
	}

	game.loadIntro()
	game.loadLevel()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
