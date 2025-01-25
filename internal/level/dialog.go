package level

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/txt/level1_dialog.json
var level_1_dialog string

type DialogState struct {
	Intro    []string `json:"intro"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	ID      string   `json:"id"`
	Symbol  rune     `json:"symbol"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
	Text    string   `json:"text,omitempty"`
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	Text    string   `json:"text"`
	Effects []string `json:"effects"`
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type alias Entity
	var tmp struct {
		Symbol string `json:"symbol"`
		*alias
	}
	tmp.alias = (*alias)(e)

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	runes := []rune(tmp.Symbol)
	if len(runes) == 0 {
		e.Symbol = 0
	} else {
		e.Symbol = runes[0]
	}

	return nil
}

func (l *Level) loadDialog() {
	var dialogState DialogState
	err := json.Unmarshal([]byte(levelDialogs[l.level]), &dialogState)
	if err != nil {
		log.Fatal(err)
	}

	l.dialogState = &dialogState
	l.levelIntroDialog = dialogState.Intro
	l.dialogActive = true
	l.selectedDialog = 0
	l.showingIntro = true
	l.entities = make(EntityMap)

	for _, entity := range dialogState.Entities {
		key := strconv.Itoa(entity.Y) + strconv.Itoa(entity.X)
		l.entities[key] = entity
		Obstacles[entity.Symbol] = entity.Symbol
		l.MapGrid[entity.Y][entity.X] = entity.Symbol
		if entity.ID == "Player" {
			l.Player.SetPosition(entity.X, entity.Y)
		}
	}

	l.setDialog(l.levelIntroDialog)

	// for y := ROW_DIVIDER; y < ROWS; y++ {
	// 	for x := 0; x < COLS; x++ {
	// 		l.TheGrid[y][x] = '.'
	// 	}
	// }

}

func (l *Level) updateDialog() {
	l.selectCooldown--
	if l.selectCooldown > 0 {
		return
	}
	l.selectCooldown = 10

	if l.showingIntro {
		l.showIntro()
	}
}

func (l *Level) showIntro() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
		l.showingIntro = false
		l.clearDialogArea()
	}
}

func (l *Level) clearDialogArea() {
	for y := 0; y < ROW_DIVIDER; y++ {
		for x := COL_DIVIDER; x < COLS; x++ {
			l.TheGrid[y][x] = LEVEL_EMPTY
		}
	}
}

func (l *Level) setDialog(text []string) {
	dialRow := 0
	for _, textLine := range text {
		lineStart := 0
		for lineStart < len(textLine) {
			lineEnd := lineStart + (COLS - COL_DIVIDER - DIALOG_INDENT)
			if lineEnd > len(textLine) {
				lineEnd = len(textLine)
			}

			// Compute the row in the grid based on dialog row plus indent
			gridRow := dialRow + DIALOG_INDENT

			// Copy each character of this chunk into the grid
			for dialCol := 0; dialCol < lineEnd-lineStart; dialCol++ {
				gridCol := COL_DIVIDER + DIALOG_INDENT + dialCol
				l.TheGrid[gridRow][gridCol] = rune(textLine[lineStart+dialCol])
			}

			// Move down one dialog row
			dialRow++
			lineStart = lineEnd

			// If we've reached the bottom of the dialog region, stop
			if gridRow+1 >= ROW_DIVIDER {
				return
			}
		}
	}
}
