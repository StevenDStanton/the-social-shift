package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/rand"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

func (g *Game) loadLevel() {

	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err = opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple ASCII map: a box with walls (#) and a player (@).

	d1, err := wav.Decode(g.audioContext, bytes.NewReader(step1))
	if err != nil {
		log.Fatal(err)
	}
	p1, err := g.audioContext.NewPlayer(d1)
	if err != nil {
		log.Fatal(err)
	}

	d2, err := wav.Decode(g.audioContext, bytes.NewReader(step2))
	if err != nil {
		log.Fatal(err)
	}
	p2, err := g.audioContext.NewPlayer(d2)
	if err != nil {
		log.Fatal(err)
	}

	d3, err := wav.Decode(g.audioContext, bytes.NewReader(step3))
	if err != nil {
		log.Fatal(err)
	}
	p3, err := g.audioContext.NewPlayer(d3)
	if err != nil {
		log.Fatal(err)
	}

	g.player = &Player{
		x:    10,
		y:    10,
		walk: []*audio.Player{p1, p2, p3},
	}

}

func (g *Game) drawLevel(screen *ebiten.Image) {
	for y := 0; y < len(g.grid); y++ {
		for x := 0; x < len(g.grid[y]); x++ {
			char := g.grid[y][x]
			text.Draw(screen, string(char), face, x*cellWidth, (y+1)*cellHeight, color.White)

		}
	}
}

func (g *Game) updateLevel() {
	g.Movement()
	switch {
	case g.level == LEVEL1:
		if !g.levelLoaded {
			g.initGrid()
			g.levelLoaded = true
		}
	}
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
		index := rand.Intn(len(g.player.walk))
		sound := g.player.walk[index]
		sound.SetVolume(0.5)

		sound.Rewind() // If you want it to start from the beginning each time
		sound.Play()
	}
}
