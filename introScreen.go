package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
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
		// ---------------------------------------------
		// 0) Choose how large you want the text
		// ---------------------------------------------
		poweredText := "Powered By Ebitengine"
		face := basicfont.Face7x13

		// Measure unscaled text
		tw := text.BoundString(face, poweredText).Dx()
		th := text.BoundString(face, poweredText).Dy()

		// We want to scale the text so it's bigger
		textScale := 4.0
		scaledTextW := float64(tw) * textScale
		scaledTextH := float64(th) * textScale

		// Decide how much vertical space you want between image and text
		const spacing = 20.0

		// ---------------------------------------------
		// 1) Figure out how to scale the image
		//    so that (image height + spacing + text height) fits on screen
		// ---------------------------------------------
		iw, ih := g.intro.images[0].Size()

		// How wide can the image be? Full screen width
		maxImageScaleW := float64(screenWidth) / float64(iw)

		// How tall can the image be? It's the screen minus room for text + spacing
		// (so the image plus spacing plus text won't exceed screen height).
		maxAvailableHForImage := float64(screenHeight) - (scaledTextH + spacing)
		maxImageScaleH := maxAvailableHForImage / float64(ih)

		// We pick the smaller of these two to preserve the aspect ratio
		imageScale := math.Min(maxImageScaleW, maxImageScaleH)
		if imageScale < 0 {
			// If our text is so large that maxImageScaleH < 0, just default
			// to something small to avoid a negative scale (edge case).
			imageScale = 0.5
		}

		// ---------------------------------------------
		// 2) Draw the image centered horizontally, top aligned so that
		//    we have room below for the text
		// ---------------------------------------------
		scaledImageW := float64(iw) * imageScale
		scaledImageH := float64(ih) * imageScale

		imageX := (float64(screenWidth) - scaledImageW) / 2
		imageY := (float64(screenHeight) - (scaledImageH + spacing + scaledTextH)) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(imageScale, imageScale)
		op.GeoM.Translate(imageX, imageY)
		screen.DrawImage(g.intro.images[0], op)

		// ---------------------------------------------
		// 3) Draw the text below the image
		// ---------------------------------------------
		// Make an offscreen image for the text, drawn at normal size
		textImg := ebiten.NewImage(tw, th)
		text.Draw(textImg, poweredText, face, 0, th, color.White)

		// We'll scale that offscreen text when drawing onto `screen`
		txtOp := &ebiten.DrawImageOptions{}
		txtOp.GeoM.Scale(textScale, textScale)

		// Position the top of the text (scaled) spacing pixels below the image
		textX := (float64(screenWidth) - scaledTextW) / 2
		textY := imageY + scaledImageH + spacing - 150

		txtOp.GeoM.Translate(textX, textY)
		screen.DrawImage(textImg, txtOp)

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
