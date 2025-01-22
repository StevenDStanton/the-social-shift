package level

func (l *Level) Update() {
	if l.LevelNumber == LEVEL0 {
		l.updateIntro()
		return
	}

	l.updateDialog()
}

func (l *Level) UpdateBoard(x, y, dx, dy int) {
	// Note that since languages like Go use
	// row-major indexing, we need to swap x and y
	// This may seem confusing if you are used to
	// column-major indexing used in many game engines
	l.MapGrid[y][x] = LEVEL_EMPTY
	l.MapGrid[dy][dx] = PLAYER_SYMBOL

}
