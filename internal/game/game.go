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
	Title        = "The Social Shift"
	ScreenWidth  = 1280
	ScreenHeight = 720

	fontSize = 16
	DPI      = 72
	COOLDOWN = 15

	LEVEL_EMPTY   = ' '
	LEVEL_WALL    = '#'
	LEVEL_DIVIDER = '|'
)

var (
	face       font.Face
	rows       = ScreenHeight / cellHeight
	cols       = ScreenWidth / cellWidth
	cellWidth  = fontSize
	cellHeight = fontSize
	divider    = int(math.Round(float64(cols) * 0.6))
	Obstacles  = map[rune]rune{
		LEVEL_WALL:    LEVEL_WALL,
		LEVEL_DIVIDER: LEVEL_DIVIDER,
	}
)

type BoardMove struct {
	x int
	y int
}

type Game struct {
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

func New() *Game {
	ac := audio.NewContext(44100)
	g := &Game{
		audioContext: ac,
		player:       player.New(ac),
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
	return ScreenWidth, ScreenHeight
}
