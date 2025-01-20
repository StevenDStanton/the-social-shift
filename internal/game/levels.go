package game

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
		Size:    fontSize,
		DPI:     DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (g *Game) drawLevel(screen *ebiten.Image) {
	for y := 0; y < len(g.grid); y++ {
		for x := 0; x < len(g.grid[y]); x++ {
			char := g.grid[y][x]
			text.Draw(screen, string(char), face, x*fontSize, (y+1)*fontSize, color.White)

		}
	}
}

func (g *Game) updateLevel() {
	g.Movement()

	if g.isLevelLoaded {
		return
	}

	switch {
	case g.level == LEVEL1:
		g.initGrid()

	}

	g.isLevelLoaded = true

}

func (g *Game) initGrid() {
	g.grid = make([][]rune, g.display.rows)

	for y := 0; y < g.display.rows; y++ {
		g.grid[y] = make([]rune, g.display.cols)
		for x := 0; x < g.display.cols; x++ {
			switch {
			// Top or bottom boundary
			case y == 0 || y == g.display.rows-1:
				g.grid[y][x] = '#'
			// Left or right boundary
			case x == 0 || x == g.display.cols-1:
				g.grid[y][x] = '#'
			// Divider between play area and UI panel
			case x == g.display.divider:
				g.grid[y][x] = '|'
			// Everything else defaults to space
			default:
				g.grid[y][x] = ' '
			}
		}
	}
	x, y := g.player.GetPosition()
	g.grid[x][y] = g.player.GetSymbol()

}

func (g *Game) canPlayerMoveToPosition(move BoardMove) bool {
	currentX, currentY := g.player.GetPosition()
	futureX := currentX + move.x
	futureY := currentY + move.y
	futureLoc := g.grid[futureY][futureX]

	_, exists := Obstacles[futureLoc]

	return !exists
}

func (g *Game) Movement() {

	if g.movementCooldown > 0 {
		g.movementCooldown--
		return
	}

	canMove := false
	currentMove := BoardMove{0, 0}
	// Move the player '@' based on key presses
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		currentMove = BoardMove{0, -1}
		canMove = g.canPlayerMoveToPosition(currentMove)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		currentMove = BoardMove{-1, 0}
		canMove = g.canPlayerMoveToPosition(currentMove)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		currentMove = BoardMove{0, 1}
		canMove = g.canPlayerMoveToPosition(currentMove)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		currentMove = BoardMove{1, 0}
		canMove = g.canPlayerMoveToPosition(currentMove)
	}

	if canMove {
		g.movementCooldown = COOLDOWN
		g.updateBoard(currentMove)
	}
}

func (g *Game) updateBoard(currentMove BoardMove) {
	playerPastX, playerPastY := g.player.GetPosition()
	g.player.Move(currentMove.x, currentMove.y)
	playerNewX, playerNewY := g.player.GetPosition()
	g.grid[playerPastY][playerPastX] = LEVEL_EMPTY
	g.grid[playerNewY][playerNewX] = g.player.GetSymbol()

}
