package level

import (
	"image/color"
	"log"
	"strings"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/txt/level1_map.txt
var level_1_map string

type EntityMap map[string]Entity

type Player interface {
	SetPosition(x, y int)
}

const (
	fontSize    = 16
	DPI         = 72
	COLS        = 80
	ROWS        = 45
	COL_DIVIDER = 48
	ROW_DIVIDER = 40
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
		level_1_map,
	}
	levelDialogs = []string{
		level_1_dialog,
	}
)

type Level struct {
	face             font.Face
	TheGrid          [][]rune
	MapGrid          [][]rune
	level            int
	dialogState      *DialogState
	levelIntroDialog []string
	dialogActive     bool
	selectedDialog   int
	showingIntro     bool
	showingItem      bool
	entities         EntityMap
	Player           Player
	selectCooldown   int
}

func New() *Level {
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

	l := &Level{face: face}
	l.configureGrid()
	return l

}

func (l *Level) Draw(screen *ebiten.Image) {
	for y := 0; y < len(l.TheGrid); y++ {
		for x := 0; x < len(l.TheGrid[y]); x++ {
			ch := l.TheGrid[y][x]
			px := x * fontSize
			py := (y + 1) * fontSize

			textColor := color.RGBA{255, 255, 255, 1}

			// if l.dialogActive && y == l.selectedDialog && x >= CENTER_DIVIDER {
			// 	textColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
			// }

			text.Draw(screen, string(ch), l.face, px, py, textColor)
		}
	}
}

func (l *Level) Update() {
	l.updateDialog()
}

func (l *Level) LoadLevel() {
	log.Println("Loading level", l.level)
	l.resetLevelState()
	l.loadMap()
	l.loadDialog()
	l.updateGridFromCamera()

}

func (l *Level) configureGrid() {
	l.TheGrid = make([][]rune, ROWS)
	for row := 0; row < ROWS; row++ {
		l.TheGrid[row] = make([]rune, COLS)
	}
}

func (l *Level) resetLevelState() {
	for y := 0; y < ROWS; y++ {
		for x := 0; x < COLS; x++ {
			l.TheGrid[y][x] = LEVEL_EMPTY
		}
	}

	l.showingIntro = true
}

func (l *Level) loadMap() {
	mapData := levelMaps[l.level]
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

	for _, entity := range l.entities {
		grid[entity.Y][entity.X] = entity.Symbol
	}

	l.MapGrid = grid

}
