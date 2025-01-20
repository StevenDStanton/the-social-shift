package main

import (
	"log"

	_ "image/png"

	"github.com/StevenDStanton/the-social-shift/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Title        = "The Social Contract"
	ScreenWidth  = 1280
	ScreenHeight = 720
)

func main() {

	ebiten.SetWindowTitle(Title)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	g := game.New(ScreenWidth, ScreenHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
