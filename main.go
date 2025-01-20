package main

import (
	"log"

	_ "image/png"

	"github.com/StevenDStanton/the-social-shift/internal/game"
	"github.com/StevenDStanton/the-social-shift/internal/level"
	"github.com/StevenDStanton/the-social-shift/internal/player"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	title        = "The Social Contract"
	screenWidth  = 1280
	screenHeight = 720
	playerStartX = 10
	playerStartY = 10
)

func main() {

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := game.New(screenWidth, screenHeight)

	l := level.New(playerStartX, playerStartY)

	p := player.New(g.AudioContext, playerStartX, playerStartY)

	p.Level = l

	g.Components = append(g.Components, p)
	g.Components = append(g.Components, l)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
