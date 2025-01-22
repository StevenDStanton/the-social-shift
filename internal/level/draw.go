package level

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (l *Level) Draw(screen *ebiten.Image) {
	switch {
	case l.LevelNumber == LEVEL0:
		l.drawIntro(screen)
	default:
		l.drawGrid(screen)
	}
}

func (l *Level) drawGrid(screen *ebiten.Image) {
	for y := 0; y < len(l.TheGrid); y++ {
		for x := 0; x < len(l.TheGrid[y]); x++ {
			ch := l.TheGrid[y][x]
			px := x * fontSize
			py := (y + 1) * fontSize

			textColor := color.RGBA{0x80, 0x80, 0x80, 0xff}
			if l.dialogActive && y == l.selectedDialog && x >= CENTER_DIVIDER {
				textColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
			}

			text.Draw(screen, string(ch), l.face, px, py, textColor)
		}
	}
}
