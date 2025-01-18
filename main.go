package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

// Embed the image and audio assets.

//go:embed assets/audio/Stamp_old_3_16b_.wav
var stampWav []byte

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

type Game struct {
	grid                  [][]rune
	mousePressedLastFrame bool
	audioContext          *audio.Context
	audioPlayer           *audio.Player
	player                *Player
	movementCooldown      int
}

// Update handles the game logic.
func (g *Game) Update() error {
	g.Movement()
	mousePress(g)
	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for y := 0; y < len(g.grid); y++ {
		for x := 0; x < len(g.grid[y]); x++ {
			char := g.grid[y][x]
			text.Draw(screen, string(char), face, x*cellWidth, (y+1)*cellHeight, color.White)

		}
	}

}

// Layout defines the screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("The Social Shift")

	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	// Now convert that parsed font into a usable Face.
	face, err = opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    24, // Font size in points
		DPI:     72, // Typical screen DPI
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple ASCII map: a box with walls (#) and a player (@).

	player := &Player{x: 10, y: 10}

	game := &Game{player: player, movementCooldown: 15}

	game.initGrid()

	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) initializeAudio() error {
	if g.audioContext != nil && g.audioPlayer != nil {
		return nil // Already initialized
	}

	// Initialize audio context
	g.audioContext = audio.NewContext(44100)
	log.Println("Audio context initialized")

	// Load the embedded audio
	audioBuffer := bytes.NewReader(stampWav)
	wavStream, err := wav.DecodeWithSampleRate(44100, audioBuffer)
	if err != nil {
		return err
	}

	g.audioPlayer, err = g.audioContext.NewPlayer(wavStream)
	if err != nil {
		return err
	}
	log.Println("Audio player created successfully")

	return nil
}

func mousePress(g *Game) error {
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mousePressed && !g.mousePressedLastFrame {
		x, y := ebiten.CursorPosition()

		log.Printf("Mouse clicked at (%d, %d)", x, y)
		// Initialize audio on first click
		if err := g.initializeAudio(); err != nil {
			log.Printf("Audio initialization failed: %v", err)
			return nil // Continue running even if audio fails
		}

		// Play the audio
		if g.audioPlayer != nil {
			g.audioPlayer.Rewind()
			g.audioPlayer.Play()
			log.Println("Audio played")
		}
	}
	g.mousePressedLastFrame = mousePressed

	return nil
}

func (g *Game) initGrid() {
	g.grid = make([][]rune, rows)

	for y := 0; y < rows; y++ {
		g.grid[y] = make([]rune, cols)
		for x := 0; x < cols; x++ {
			switch {
			// Top or bottom boundary
			case y == 0 || y == rows-1:
				g.grid[y][x] = '#'
			// Left or right boundary
			case x == 0 || x == cols-1:
				g.grid[y][x] = '#'
			// Divider between play area and UI panel
			case x == divider:
				g.grid[y][x] = '|'
			// Everything else defaults to space
			default:
				g.grid[y][x] = ' '
			}
		}
	}

	g.grid[g.player.x][g.player.y] = '@'

}

func (g *Game) Movement() {

	if g.movementCooldown > 0 {
		g.movementCooldown--
		return
	}

	moved := false
	// Move the player '@' based on key presses
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.movePlayer(0, -1)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.movePlayer(-1, 0)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.movePlayer(0, 1)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.movePlayer(1, 0)
		moved = true
	}

	if moved {
		g.movementCooldown = 15
	}
}

func (g *Game) movePlayer(dx, dy int) {
	newX := g.player.x + dx
	newY := g.player.y + dy

	if g.grid[newY][newX] == '#' || g.grid[newY][newX] == '|' {
		return
	}

	// Move the player
	g.grid[g.player.y][g.player.x] = ' '
	g.player.x = newX
	g.player.y = newY
	g.grid[g.player.y][g.player.x] = '@'
}
