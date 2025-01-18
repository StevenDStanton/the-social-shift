package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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
		Size:    24, // Font size in points
		DPI:     72, // Typical screen DPI
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple ASCII map: a box with walls (#) and a player (@).

	g.player = &Player{x: 10, y: 10}

	g.initGrid()

	if err != nil {
		log.Fatalf("Failed to initialize game: %v", err)
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
	}
}
