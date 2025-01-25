package intro

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Intro struct {
	images     []*ebiten.Image
	timer      int
	imageIndex int
	Level      Level
	Game       Game
}

type Level interface {
	LoadLevel()
}

type Game interface {
	RemoveComponent(c interface{})
}

func New() *Intro {

	i := &Intro{}
	i.loadIntro()
	return i
}

//go:embed img/1.png
var goLangImage []byte

//go:embed img/2.png
var gopherImage []byte

//go:embed img/3.png
var thesimpledevImage []byte

//go:embed img/4.png
var aptImage []byte

func (i *Intro) loadIntro() {
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

	i.images = []*ebiten.Image{
		ebiten.NewImageFromImage(goLangImageData),
		ebiten.NewImageFromImage(gopherImageData),
		ebiten.NewImageFromImage(thesimpledevImageData),
		ebiten.NewImageFromImage(aptImageData),
	}

}

func (i *Intro) Update() {
	i.timer++
	if i.timer%120 == 0 {
		i.imageIndex++
		if i.imageIndex >= len(i.images)+1 {
			i.Level.LoadLevel()
			i.Game.RemoveComponent(i)
		}
	}
}

func (i *Intro) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()

	switch i.imageIndex {
	case 0:
		poweredText := "Powered By Ebitengine"
		face := basicfont.Face7x13

		tw := text.BoundString(face, poweredText).Dx()
		th := text.BoundString(face, poweredText).Dy()

		textScale := 4.0
		scaledTextW := float64(tw) * textScale
		scaledTextH := float64(th) * textScale

		const spacing = 20.0

		iw, ih := i.images[0].Size()

		maxImageScaleW := float64(screenWidth) / float64(iw)

		maxAvailableHForImage := float64(screenHeight) - (scaledTextH + spacing)
		maxImageScaleH := maxAvailableHForImage / float64(ih)

		imageScale := math.Min(maxImageScaleW, maxImageScaleH)
		if imageScale < 0 {
			imageScale = 0.5
		}

		scaledImageW := float64(iw) * imageScale
		scaledImageH := float64(ih) * imageScale

		imageX := (float64(screenWidth) - scaledImageW) / 2
		imageY := (float64(screenHeight) - (scaledImageH + spacing + scaledTextH)) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(imageScale, imageScale)
		op.GeoM.Translate(imageX, imageY)
		screen.DrawImage(i.images[0], op)

		textImg := ebiten.NewImage(tw, th)
		text.Draw(textImg, poweredText, face, 0, th, color.White)

		txtOp := &ebiten.DrawImageOptions{}
		txtOp.GeoM.Scale(textScale, textScale)

		textX := (float64(screenWidth) - scaledTextW) / 2
		textY := imageY + scaledImageH + spacing - 150

		txtOp.GeoM.Translate(textX, textY)
		screen.DrawImage(textImg, txtOp)

	case 1:
		iw2, ih2 := i.images[1].Size()
		op2 := &ebiten.DrawImageOptions{}

		maxWidth2 := 0.4 * float64(screenWidth)
		scale2 := maxWidth2 / float64(iw2)

		if float64(ih2)*scale2 > float64(screenHeight)*0.4 {
			scale2 = (float64(screenHeight) * 0.4) / float64(ih2)
		}
		op2.GeoM.Scale(scale2, scale2)

		scaledW2 := float64(iw2) * scale2
		scaledH2 := float64(ih2) * scale2
		x2 := (float64(screenWidth) - scaledW2) / 2
		y2 := (float64(screenHeight)-scaledH2)/2 - 50
		op2.GeoM.Translate(x2, y2)

		screen.DrawImage(i.images[1], op2)

		iw3, ih3 := i.images[2].Size()
		op3 := &ebiten.DrawImageOptions{}

		scale3 := float64(screenWidth) / float64(iw3)
		op3.GeoM.Scale(scale3, scale3)

		scaledW3 := float64(iw3) * scale3
		scaledH3 := float64(ih3) * scale3
		x3 := (float64(screenWidth) - scaledW3) / 2
		y3 := float64(screenHeight) - scaledH3
		op3.GeoM.Translate(x3, y3)

		screen.DrawImage(i.images[2], op3)

	case 2:
		iw4, ih4 := i.images[3].Size()
		op4 := &ebiten.DrawImageOptions{}

		scaleW := float64(screenWidth) / float64(iw4)
		scaleH := float64(screenHeight) / float64(ih4)
		scale4 := math.Min(scaleW, scaleH)

		op4.GeoM.Scale(scale4, scale4)

		scaledW4 := float64(iw4) * scale4
		scaledH4 := float64(ih4) * scale4
		x4 := (float64(screenWidth) - scaledW4) / 2
		y4 := (float64(screenHeight) - scaledH4) / 2
		op4.GeoM.Translate(x4, y4)

		screen.DrawImage(i.images[3], op4)

	case 3:
		face := basicfont.Face7x13
		creditLines := []string{
			"Credits",
			"==========",
			"Coded by TheSimpleDev",
			"Footstep sounds by Kenney (kenney.nl)",
			"Awesome Gopher Art by MariaLetta",
		}

		lineHeight := text.BoundString(face, "A").Dy() + 4
		totalLines := len(creditLines)
		totalHeight := lineHeight * totalLines
		startY := (screenHeight - totalHeight) / 2

		for i, line := range creditLines {
			lineWidth := text.BoundString(face, line).Dx()
			startX := (screenWidth - lineWidth) / 2
			x := startX
			y := startY + i*lineHeight

			text.Draw(screen, line, face, x, y, color.White)
		}
	}
}
