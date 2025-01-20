package level

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/img/1.png
var goLangImage []byte

//go:embed assets/img/2.png
var gopherImage []byte

//go:embed assets/img/3.png
var thesimpledevImage []byte

//go:embed assets/img/4.png
var aptImage []byte

//go:embed assets/txt/level1.txt
var level1 string

type LevelNumber int

type Intro struct {
	images []*ebiten.Image
}

const (
	LEVEL0 LevelNumber = iota
	LEVEL1
)

const (
	fontSize      = 16
	DPI           = 72
	COLS          = 80
	ROWS          = 45
	T_DIVIDER     = 48
	B_DIVIDER     = 40
	LEVEL_EMPTY   = ' '
	LEVEL_WALL    = '#'
	LEVEL_DIVIDER = '|'
	PLAYER_SYMBOL = '@'
)

var (
	Obstacles = map[rune]rune{
		LEVEL_WALL:    LEVEL_WALL,
		LEVEL_DIVIDER: LEVEL_DIVIDER,
	}

	allLevelsData = []string{
		level1,
	}
)

type Level struct {
	LevelNumber     LevelNumber
	UIGrid          [][]rune
	MapGrid         [][]rune
	isLevelLoaded   bool
	face            font.Face
	introImageIndex int
	introTimer      int
	intro           *Intro
	PlayerStartX    int
	PlayerStartY    int
}

func New(x, y int) *Level {

	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     DPI,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	l := &Level{face: face, PlayerStartX: x, PlayerStartY: y}

	l.loadIntro()

	return l
}

func (l *Level) Update() {

	if l.LevelNumber == LEVEL0 {
		l.updateIntro()
		return
	}

	if l.isLevelLoaded {
		return
	}
	l.LoadActiveLevel()
	l.initUIGrid()
	l.isLevelLoaded = true
}

func (l *Level) Draw(screen *ebiten.Image) {
	switch {
	case l.LevelNumber == LEVEL0:
		l.drawIntro(screen)
	default:
		l.drawLevel(screen)
	}
}

func (l *Level) drawMap(screen *ebiten.Image) {
	viewW := l.viewWidth()
	viewH := l.viewHeight()

	for screenY := 0; screenY < viewH; screenY++ {
		for screenX := 0; screenX < viewW; screenX++ {
			mapY := cameraY + screenY
			mapX := cameraX + screenX

			if mapY < 0 || mapY >= len(l.MapGrid) {
				continue
			}
			if mapX < 0 || mapX >= len(l.MapGrid[mapY]) {
				continue
			}

			ch := l.MapGrid[mapY][mapX]

			pixelX := screenX * fontSize
			pixelY := (screenY + 1) * fontSize

			text.Draw(screen, string(ch), l.face, pixelX, pixelY, color.White)
		}
	}
}

func (l *Level) initUIGrid() {
	l.UIGrid = make([][]rune, ROWS)
	for y := 0; y < ROWS; y++ {
		l.UIGrid[y] = make([]rune, COLS)
		for x := 0; x < COLS; x++ {
			l.UIGrid[y][x] = LEVEL_EMPTY
		}
	}

	for y := 0; y < ROWS; y++ {
		l.UIGrid[y][T_DIVIDER] = LEVEL_DIVIDER
	}

}

func (l *Level) drawUI(screen *ebiten.Image) {
	for y := 0; y < len(l.UIGrid); y++ {
		for x := 0; x < len(l.UIGrid[y]); x++ {
			ch := l.UIGrid[y][x]
			px := x * fontSize
			py := (y + 1) * fontSize
			text.Draw(screen, string(ch), l.face, px, py, color.White)
		}
	}
}

func (l *Level) drawLevel(screen *ebiten.Image) {
	l.drawMap(screen)
	l.drawUI(screen)
}

func (l *Level) IsWalkable(x, y int) bool {
	tile := l.MapGrid[y][x]

	_, exists := Obstacles[tile]
	if exists {
		log.Println("Obstacle", string(tile), " found at", x, y)
	}
	return !exists
}

func (l *Level) LoadActiveLevel() error {
	index := int(l.LevelNumber) - 1

	if index < 0 || index >= len(allLevelsData) {
		return fmt.Errorf("level %d not found", l.LevelNumber)
	}

	mapData := allLevelsData[index]

	lines := strings.Split(strings.TrimSpace(mapData), "\n")

	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	grid := make([][]rune, height)

	for y := 0; y < height; y++ {
		grid[y] = make([]rune, width)
		for x, ch := range lines[y] {
			grid[y][x] = ch
		}
	}

	l.MapGrid = grid

	l.MapGrid[l.PlayerStartY][l.PlayerStartX] = PLAYER_SYMBOL

	return nil
}

func (l *Level) UpdateBoard(x, y, dx, dy int) {
	// Note that since languages like Go use
	// row-major indexing, we need to swap x and y
	// This may seem confusing if you are used to
	// column-major indexing used in many game engines
	l.MapGrid[y][x] = LEVEL_EMPTY
	l.MapGrid[dy][dx] = PLAYER_SYMBOL

}

func (l *Level) loadIntro() {
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

	l.intro = &Intro{
		images: []*ebiten.Image{
			ebiten.NewImageFromImage(goLangImageData),
			ebiten.NewImageFromImage(gopherImageData),
			ebiten.NewImageFromImage(thesimpledevImageData),
			ebiten.NewImageFromImage(aptImageData),
		},
	}

}

func (l *Level) updateIntro() {
	l.introTimer++
	if l.introTimer%120 == 0 {
		l.introImageIndex++
		if l.introImageIndex >= len(l.intro.images)+1 {
			l.LevelNumber = LEVEL1
		}
	}
}

func (l *Level) drawIntro(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()

	switch l.introImageIndex {
	case 0:
		poweredText := "Powered By Ebitengine"
		face := basicfont.Face7x13

		tw := text.BoundString(face, poweredText).Dx()
		th := text.BoundString(face, poweredText).Dy()

		textScale := 4.0
		scaledTextW := float64(tw) * textScale
		scaledTextH := float64(th) * textScale

		const spacing = 20.0

		iw, ih := l.intro.images[0].Size()

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
		screen.DrawImage(l.intro.images[0], op)

		textImg := ebiten.NewImage(tw, th)
		text.Draw(textImg, poweredText, face, 0, th, color.White)

		txtOp := &ebiten.DrawImageOptions{}
		txtOp.GeoM.Scale(textScale, textScale)

		textX := (float64(screenWidth) - scaledTextW) / 2
		textY := imageY + scaledImageH + spacing - 150

		txtOp.GeoM.Translate(textX, textY)
		screen.DrawImage(textImg, txtOp)

	case 1:
		iw2, ih2 := l.intro.images[1].Size()
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

		screen.DrawImage(l.intro.images[1], op2)

		iw3, ih3 := l.intro.images[2].Size()
		op3 := &ebiten.DrawImageOptions{}

		scale3 := float64(screenWidth) / float64(iw3)
		op3.GeoM.Scale(scale3, scale3)

		scaledW3 := float64(iw3) * scale3
		scaledH3 := float64(ih3) * scale3
		x3 := (float64(screenWidth) - scaledW3) / 2
		y3 := float64(screenHeight) - scaledH3
		op3.GeoM.Translate(x3, y3)

		screen.DrawImage(l.intro.images[2], op3)

	case 2:
		iw4, ih4 := l.intro.images[3].Size()
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

		screen.DrawImage(l.intro.images[3], op4)

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
