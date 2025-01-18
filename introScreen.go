package main

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) loadIntro() {
	goLangImageData, _, err := image.Decode(bytes.NewReader(goLangImage))
	if err != nil {
		log.Fatal(err)
	}

	gopherImageData, _, err := image.Decode(bytes.NewReader(gopherImage))
	if err != nil {
		log.Fatal(err)
	}

	thesimpledevImageData, _, err := image.Decode(bytes.NewReader(thesimpledevImage))
	if err != nil {
		log.Fatal(err)
	}

	aptImageData, _, err := image.Decode(bytes.NewReader(aptImage))
	if err != nil {
		log.Fatal(err)
	}

	g.intro = &Intro{
		images: []*ebiten.Image{
			ebiten.NewImageFromImage(goLangImageData),
			ebiten.NewImageFromImage(gopherImageData),
			ebiten.NewImageFromImage(thesimpledevImageData),
			ebiten.NewImageFromImage(aptImageData),
		},
	}

}

func (g *Game) drawIntro(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()

	switch g.introImageIndex {
	case 0:
		// Image 1: Fill as large as possible while preserving aspect ratio
		iw, ih := g.intro.images[0].Size()
		op := &ebiten.DrawImageOptions{}

		// Calculate a scale that keeps the image fully visible and preserves ratio
		scaleW := float64(screenWidth) / float64(iw)
		scaleH := float64(screenHeight) / float64(ih)
		scale := math.Min(scaleW, scaleH)

		op.GeoM.Scale(scale, scale)
		// Center on screen
		x := (float64(screenWidth) - float64(iw)*scale) / 2
		y := (float64(screenHeight) - float64(ih)*scale) / 2
		op.GeoM.Translate(x, y)

		screen.DrawImage(g.intro.images[0], op)

	case 1:
		// Image 2 above Image 3
		// 1) Draw Image 2 (scaled smaller so it’s not too big)
		iw2, ih2 := g.intro.images[1].Size()
		op2 := &ebiten.DrawImageOptions{}

		// Scale so it’s at most ~40% of the screen width
		maxWidth2 := 0.4 * float64(screenWidth)
		scale2 := maxWidth2 / float64(iw2)
		// Make sure it also fits in screen height
		if float64(ih2)*scale2 > float64(screenHeight)*0.4 {
			scale2 = (float64(screenHeight) * 0.4) / float64(ih2)
		}
		op2.GeoM.Scale(scale2, scale2)

		// Position slightly above the vertical center
		scaledW2 := float64(iw2) * scale2
		scaledH2 := float64(ih2) * scale2
		x2 := (float64(screenWidth) - scaledW2) / 2
		y2 := (float64(screenHeight)-scaledH2)/2 - 50
		op2.GeoM.Translate(x2, y2)

		screen.DrawImage(g.intro.images[1], op2)

		// 2) Draw Image 3 along bottom (as you mentioned, that part is “perfect”).
		iw3, ih3 := g.intro.images[2].Size()
		op3 := &ebiten.DrawImageOptions{}

		// Scale to fit the full screen width, preserving aspect ratio
		scale3 := float64(screenWidth) / float64(iw3)
		op3.GeoM.Scale(scale3, scale3)

		scaledW3 := float64(iw3) * scale3
		scaledH3 := float64(ih3) * scale3
		x3 := (float64(screenWidth) - scaledW3) / 2
		y3 := float64(screenHeight) - scaledH3
		op3.GeoM.Translate(x3, y3)

		screen.DrawImage(g.intro.images[2], op3)

	case 2:
		// Image 4: center it and ensure it fits fully on screen
		iw4, ih4 := g.intro.images[3].Size()
		op4 := &ebiten.DrawImageOptions{}

		// Scale so that the image is fully visible
		scaleW := float64(screenWidth) / float64(iw4)
		scaleH := float64(screenHeight) / float64(ih4)
		scale4 := math.Min(scaleW, scaleH)

		// Optionally shrink it a bit more so it isn't too large:
		// scale4 *= 0.8 // if you want breathing room
		op4.GeoM.Scale(scale4, scale4)

		scaledW4 := float64(iw4) * scale4
		scaledH4 := float64(ih4) * scale4
		x4 := (float64(screenWidth) - scaledW4) / 2
		y4 := (float64(screenHeight) - scaledH4) / 2
		op4.GeoM.Translate(x4, y4)

		screen.DrawImage(g.intro.images[3], op4)
	}
}

func (g *Game) updateIntro() {
	g.introTimer++
	// Every 120 frames (~2 seconds at 60FPS), move to the next image.
	if g.introTimer%120 == 0 {
		g.introImageIndex++
		// If we've shown all images, advance to LEVEL1
		if g.introImageIndex >= len(g.intro.images) {
			g.level = LEVEL1
		}
	}
}
