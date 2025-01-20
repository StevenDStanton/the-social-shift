package level

var (
	cameraX int
	cameraY int
)

func (l *Level) viewWidth() int {
	return T_DIVIDER
}

func (l *Level) viewHeight() int {
	return B_DIVIDER
}

func (l *Level) UpdateCamera(playerX, playerY int) {
	halfWidth := l.viewWidth() / 2
	halfHeight := l.viewHeight() / 2
	desiredX := playerX - halfWidth
	desiredY := playerY - halfHeight

	mapHeight := len(l.MapGrid)
	if mapHeight == 0 {
		return
	}
	mapWidth := len(l.MapGrid[0])

	viewW := l.viewWidth()
	viewH := l.viewHeight()

	if mapWidth <= viewW {
		cameraX = 0
	} else {
		if desiredX < 0 {
			desiredX = 0
		}
		maxCameraX := mapWidth - viewW
		if desiredX > maxCameraX {
			desiredX = maxCameraX
		}
		cameraX = desiredX
	}

	if mapHeight <= viewH {
		cameraY = 0
	} else {
		if desiredY < 0 {
			desiredY = 0
		}
		maxCameraY := mapHeight - viewH
		if desiredY > maxCameraY {
			desiredY = maxCameraY
		}
		cameraY = desiredY
	}
}
