package game

import (
	_ "embed"
	"math"

	"image/color"

	"github.com/StevenDStanton/the-social-shift/internal/player"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed assets/img/1.png
var goLangImage []byte

//go:embed assets/img/2.png
var gopherImage []byte

//go:embed assets/img/3.png
var thesimpledevImage []byte

//go:embed assets/img/4.png
var aptImage []byte

type Level int
type LevelStruct rune

const (
	LEVEL0 Level = iota
	LEVEL1
	LEVEL2
	LEVEL3
	LEVEL4
)

const (
	fontSize = 16
	DPI      = 72
	COOLDOWN = 15

	LEVEL_EMPTY   = ' '
	LEVEL_WALL    = '#'
	LEVEL_DIVIDER = '|'
)

var (
	face font.Face

	Obstacles = map[rune]rune{
		LEVEL_WALL:    LEVEL_WALL,
		LEVEL_DIVIDER: LEVEL_DIVIDER,
	}
)

type BoardMove struct {
	x int
	y int
}

type GameDisplay struct {
	rows         int
	cols         int
	divider      int
	screenWidth  int
	screenHeight int
}

type Game struct {
	display          GameDisplay
	level            Level
	intro            *Intro
	grid             [][]rune
	audioContext     *audio.Context
	player           *player.Player
	movementCooldown int
	introImageIndex  int
	introTimer       int
	isLevelLoaded    bool
}

func New(screenWidth, screenHeight int) *Game {

	ac := audio.NewContext(44100)

	g := &Game{
		audioContext: ac,
		player:       player.New(ac),
		display: GameDisplay{
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
			rows:         screenHeight / fontSize,
			cols:         screenWidth / fontSize,
			divider:      int(math.Round(float64(screenWidth/fontSize) * 0.6)),
		},
	}

	g.loadIntro()
	g.loadLevel()

	return g

}

// Update handles the game logic.
func (g *Game) Update() error {

	switch {
	case g.level == LEVEL0:
		g.updateIntro()
		return nil
	case g.level == LEVEL1:
		g.updateLevel()
		return nil
	default:
		return nil
	}

}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	switch {
	case g.level == LEVEL0:
		g.drawIntro(screen)
	default:
		g.drawLevel(screen)
	}

}

// Layout defines the screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.display.screenWidth, g.display.screenHeight
}
