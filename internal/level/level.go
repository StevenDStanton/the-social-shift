package level

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
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

//go:embed assets/txt/level1_map.txt
var level1Map string

//go:embed assets/txt/level1_dialog.txt
var level1Dialog string

type LevelNumber int

type Intro struct {
	images []*ebiten.Image
}

const (
	LEVEL0 LevelNumber = iota
	LEVEL1
)

const (
	fontSize       = 16
	DPI            = 72
	COLS           = 80
	ROWS           = 45
	CENTER_DIVIDER = 48
	B_DIVIDER      = 40
	DEBUG_MODE     = false
)

const (
	LEVEL_EMPTY     = ' '
	LEVEL_WALL      = '#'
	PLAYER_SYMBOL   = '@'
	DIALOG_SELECTOR = '>'
	DIALOG_INDENT   = 2
)

var (
	Obstacles = map[rune]rune{
		LEVEL_WALL: LEVEL_WALL,
	}

	levelMaps = []string{
		level1Map,
	}
	levelDialogs = []string{
		level1Dialog,
	}
)

type Level struct {
	LevelNumber     LevelNumber
	TheGrid         [][]rune
	MapGrid         [][]rune
	face            font.Face
	introImageIndex int
	introTimer      int
	intro           *Intro
	PlayerStartX    int
	PlayerStartY    int
	dialogLines     []string
	selectedDialog  int
	dialogActive    bool
	selectCooldown  int
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

	l.TheGrid = make([][]rune, ROWS)
	for row := 0; row < ROWS; row++ {
		l.TheGrid[row] = make([]rune, COLS)
	}

	return l
}

func (l *Level) IsWalkable(x, y int) bool {
	tile := l.MapGrid[y][x]

	_, exists := Obstacles[tile]
	if exists {
		log.Println("Obstacle", string(tile), " found at", x, y)
	}
	return !exists
}

func (l *Level) loadLevel(index LevelNumber) error {
	l.LevelNumber = index

	if int(index) < 0 || int(index)-1 >= len(levelMaps) {
		fmt.Println("Level not found")
		return fmt.Errorf("level %d not found", l.LevelNumber)
	}

	l.clearTheGrid()
	l.loadMap()
	l.loadDialog()
	l.updateGridFromCamera()

	return nil
}

func (l *Level) clearTheGrid() {
	for y := 0; y < B_DIVIDER; y++ {
		for x := 0; x < COLS; x++ {
			l.TheGrid[y][x] = LEVEL_EMPTY
		}
	}
}

func (l *Level) loadMap() {
	index := int(l.LevelNumber) - 1
	mapData := levelMaps[index]
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
}

func (l *Level) loadDialog() {
	index := int(l.LevelNumber) - 1
	if index >= 0 && index < len(levelDialogs) {
		l.dialogLines = strings.Split(strings.TrimSpace(levelDialogs[index]), "\n")
		l.dialogActive = true
		l.selectedDialog = 0
	}

	lineIndex := 0
	charIndex := 0
	for y := 0; y < B_DIVIDER; y++ {
		for x := CENTER_DIVIDER; x < COLS; x++ {
			if lineIndex < len(l.dialogLines) {
				l.TheGrid[y][x] = rune(l.dialogLines[lineIndex][charIndex])
				charIndex++
				if charIndex >= len(l.dialogLines[lineIndex]) {
					charIndex = 0
					lineIndex++
					break
				}
			}
		}
	}

}

func (l *Level) updateDialog() {
	l.selectCooldown--
	if !l.dialogActive || len(l.dialogLines) == 0 || l.selectCooldown > 0 {
		return
	}
	l.selectCooldown = 10

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		l.selectedDialog = (l.selectedDialog + 1) % len(l.dialogLines)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.selectedDialog = (l.selectedDialog - 1 + len(l.dialogLines)) % len(l.dialogLines)
	}

	for y := 0; y < B_DIVIDER; y++ {
		for x := CENTER_DIVIDER; x < COLS; x++ {
			l.TheGrid[y][x] = LEVEL_EMPTY
		}
	}

	// Update dialog area with current selection
	for i, line := range l.dialogLines {
		y := i
		if y >= B_DIVIDER {
			break
		}

		// Add selector for the currently selected option
		if i == l.selectedDialog {
			l.TheGrid[y][CENTER_DIVIDER] = DIALOG_SELECTOR
		}

		// Add indentation and the dialog text
		startX := CENTER_DIVIDER + DIALOG_INDENT
		for x, ch := range line {
			if startX+x >= COLS {
				break
			}
			l.TheGrid[y][startX+x] = ch
		}
	}
}
