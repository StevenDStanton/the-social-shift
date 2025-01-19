package player

import (
	"bytes"
	"log"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"golang.org/x/exp/rand"
)

//go:embed assets/audio/footstep_carpet_000.wav
var step1 []byte

//go:embed assets/audio/footstep_carpet_001.wav
var step2 []byte

//go:embed assets/audio/footstep_carpet_002.wav
var step3 []byte

type Player struct {
	x, y   int
	symbol rune
	walk   []*audio.Player
}

func New(ac *audio.Context) *Player {
	d1, err := wav.Decode(ac, bytes.NewReader(step1))
	if err != nil {
		log.Fatal(err)
	}
	p1, err := ac.NewPlayer(d1)
	if err != nil {
		log.Fatal(err)
	}

	d2, err := wav.Decode(ac, bytes.NewReader(step2))
	if err != nil {
		log.Fatal(err)
	}
	p2, err := ac.NewPlayer(d2)
	if err != nil {
		log.Fatal(err)
	}

	d3, err := wav.Decode(ac, bytes.NewReader(step3))
	if err != nil {
		log.Fatal(err)
	}
	p3, err := ac.NewPlayer(d3)
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		x:      10,
		y:      10,
		symbol: '@',
		walk:   []*audio.Player{p1, p2, p3},
	}

}

func (p *Player) GetPosition() (int, int) {
	return p.x, p.y
}

func (p *Player) GetSymbol() rune {
	return p.symbol
}

func (p *Player) Move(x, y int) {
	p.x += x
	p.y += y
	p.PlayFootstep()
}

func (p *Player) PlayFootstep() {
	index := rand.Intn(len(p.walk))
	sound := p.walk[index]
	sound.SetVolume(0.5)
	sound.Rewind()
	sound.Play()
}
