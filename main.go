package main

import (
	_ "embed"
	"log"

	_ "image/png"

	"github.com/StevenDStanton/the-social-shift/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetWindowTitle(game.Title)
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	g := game.New()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
