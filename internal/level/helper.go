package level

import (
	"strconv"
)

func (l *Level) IsWalkable(x, y int) bool {
	tile := l.MapGrid[y][x]

	if l.showingItem {
		l.showingItem = false
		l.clearDialogArea()
	}

	_, exists := Obstacles[tile]
	if exists {
		key := strconv.Itoa(y) + strconv.Itoa(x)
		item, interactionExists := l.entities[key]
		if interactionExists {
			l.showingItem = true
			l.showingIntro = false
			l.clearDialogArea()
			l.setDialog([]string{item.Text})
		} else {
			l.showingItem = true
			l.showingIntro = false
			l.clearDialogArea()
			l.setDialog([]string{"Its a wall, you probably", "shouldn't lick it", "It doesn't look very clean"})
		}

	}

	return !exists
}

func (l *Level) UpdateBoard(x, y, dx, dy int) {
	// If no entity interaction, perform normal movement
	l.MapGrid[y][x] = LEVEL_EMPTY
	l.MapGrid[dy][dx] = PLAYER_SYMBOL
}
