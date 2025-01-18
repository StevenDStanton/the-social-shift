package main

import (
	_ "embed"
	"image/color"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

// Embed the image and audio assets.

//go:embed assets/audio/Stamp_old_3_16b_.wav
var stampWav []byte

//go:embed assets/img/1.png
var goLangImage []byte

//go:embed assets/img/2.png
var gopherImage []byte

//go:embed assets/img/3.png
var thesimpledevImage []byte

//go:embed assets/img/4.png
var aptImage []byte

const (
	screenWidth  = 1280
	screenHeight = 720

	cellWidth  = 24
	cellHeight = 24

	rows    = 45 // 720 / 16
	cols    = 80 // 1280 / 16
	divider = 50 // Column where the UI panel starts

)

var face font.Face

type Player struct {
	x, y int
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
	level                 Level
	intro                 *Intro
	grid                  [][]rune
	mousePressedLastFrame bool
	audioContext          *audio.Context
	audioPlayer           *audio.Player
	player                *Player
	movementCooldown      int
	introImageIndex       int
	introTimer            int
}

// Update handles the game logic.
func (g *Game) Update() error {

	mousePress(g)

	switch {
	case g.level == LEVEL0:
		g.updateIntro()
		return nil
	case g.level == LEVEL1:
		g.Movement()
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
	case g.level == LEVEL1:
		g.drawLevel(screen)
	default:
		log.Fatal("Invalid level")
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

	game := &Game{}

	game.loadIntro()
	game.loadLevel()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
