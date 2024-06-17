/*
A game for Ebitengine game jam 2024

# Copyright (C) 2024 Loïg Jezequel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

type tetrisLine = [gPlayAreaWidthInBlocks]int
type tetrisGrid = [gPlayAreaHeightInBlocks + gInvisibleLines]tetrisLine

// Structure for one tetris game
type tetris struct {
	area                  tetrisGrid
	currentBlock          tetrisBlock
	nextBlock             tetrisBlock
	autoDownFrame         int
	autoDownFrameLimit    int
	manualDownFrame       int
	manualDownFrameLimit  int
	lrMoveFrame           int
	lrMoveFrameLimit      int
	lrFirstMoveFrame      int
	lrFirstMoveFrameLimit int
	manualMoveAllowed     bool
	numLines              int
	dropLenght            int
	deathLines            int
	// animation and lines removal handling
	toCheck                          [2]int
	toRemove                         [4]bool
	toRemoveNum                      int
	firstAvailable                   int
	removeLineAnimationFrame         int
	removeLineAnimationStep          int
	removeLineAnimationStepNumFrames int
}

func (t *tetris) init(level int, balance balancing) {
	if level == 0 {
		t.area = tetrisGrid{}
		t.currentBlock = getNewBlock()
		t.currentBlock.setInitialPosition()
		t.nextBlock = getNewBlock()
	}
	t.autoDownFrame = 0
	t.autoDownFrameLimit = balance.getSpeed()
	t.manualDownFrame = 0
	t.manualDownFrameLimit = 4
	t.lrMoveFrame = 0
	t.lrMoveFrameLimit = 6
	t.lrFirstMoveFrame = 0
	t.lrFirstMoveFrameLimit = 15
	t.manualMoveAllowed = true
	t.numLines = 0
	t.dropLenght = 0
	t.deathLines = balance.getDeathLines()
	t.toCheck = [2]int{}
	t.toRemove = [4]bool{}
	t.toRemoveNum = 0
	t.removeLineAnimationFrame = 0
	t.removeLineAnimationStep = 0
	t.removeLineAnimationStepNumFrames = 8
}

func (t *tetris) setUpNext() (dead bool) {
	dead = t.lost()

	t.currentBlock = t.nextBlock
	t.currentBlock.setInitialPosition()
	t.nextBlock = getNewBlock()

	t.manualMoveAllowed = false

	return
}

func (t *tetris) update(moveDownRequest, moveLeftRequest, moveRightRequest, rotateLeft, rotateRight bool, level int) (dead bool, scoreIncrease int, playSounds [assets.NumSounds]bool) {

	if t.removeLineAnimationStep > 0 {

		t.removeLineAnimationFrame++
		if t.removeLineAnimationFrame >= t.removeLineAnimationStepNumFrames {
			t.removeLineAnimationStep++
			t.removeLineAnimationFrame = 0
		}

		if t.removeLineAnimationStep < 8 {
			return
		}

		t.removeLineAnimationStep = 0

		// lines removal animation and effects
		playSounds[assets.SoundLinesFallingID] = true
		t.removeLines()

		switch t.toRemoveNum {
		case 1:
			scoreIncrease += 40 * (level + 1)
		case 2:
			scoreIncrease += 100 * (level + 1)
		case 3:
			scoreIncrease += 300 * (level + 1)
		case 4:
			scoreIncrease += 1200 * (level + 1)
		}
		t.numLines += t.toRemoveNum

		t.toRemove = [4]bool{}
		t.toRemoveNum = 0
		t.toCheck = [2]int{}

		dead = t.setUpNext()

		return
	}

	if rotateLeft && !rotateRight {
		playSounds[assets.SoundRotationID] = t.currentBlock.rotateLeft(t.area)
	}

	if rotateRight && !rotateLeft {
		playSounds[assets.SoundRotationID] = t.currentBlock.rotateRight(t.area)
	}

	mayAllowManualMoves := false

	// left/right movements of blocks handling
	xMove := 0
	if moveLeftRequest {
		xMove--
	}
	if moveRightRequest {
		xMove++
	}

	if !moveLeftRequest && !moveRightRequest {
		mayAllowManualMoves = true
		t.lrMoveFrame = 0
		t.lrFirstMoveFrame = 0
	}

	if !t.manualMoveAllowed {
		xMove = 0
	}

	if xMove != 0 {
		if t.lrMoveFrame > 0 || (t.lrFirstMoveFrame > 0 && t.lrFirstMoveFrame < t.lrFirstMoveFrameLimit) {
			xMove = 0
		}
		t.lrMoveFrame++
		if t.lrMoveFrame >= t.lrMoveFrameLimit {
			t.lrMoveFrame = 0
		}
		if t.lrFirstMoveFrame < t.lrFirstMoveFrameLimit {
			t.lrFirstMoveFrame++
		}
	}

	// automatic down movement of blocks handling
	autoDown := false
	t.autoDownFrame++

	if t.autoDownFrame >= t.autoDownFrameLimit {
		autoDown = true
		t.autoDownFrame = 0
	}

	// manual down movement of blocks handling
	manualDown := false

	if !moveDownRequest {
		t.manualDownFrame = 0
		t.manualMoveAllowed = t.manualMoveAllowed || mayAllowManualMoves
		t.dropLenght = 0
	}

	if moveDownRequest && t.manualMoveAllowed {
		manualDown = t.manualDownFrame == 0
		t.manualDownFrame++
		if t.manualDownFrame >= t.manualDownFrameLimit {
			t.manualDownFrame = 0
		}
	}

	if manualDown {
		t.dropLenght++
	}

	// update position according to movements requests
	var stuck bool
	stuck, playSounds[assets.SoundLeftRightID] = t.currentBlock.updatePosition(xMove, autoDown || manualDown, t.area)
	if stuck {
		playSounds[assets.SoundTouchGroundID] = true

		t.toCheck = t.currentBlock.writeInGrid(&t.area)

		scoreIncrease = t.dropLenght

		t.toRemoveNum, t.firstAvailable, t.toRemove = t.checkLines()

		if t.toRemoveNum > 0 {
			t.removeLineAnimationStep = 1
			playSounds[assets.SoundLinesVanishingID] = true
			return
		}

		dead = t.setUpNext()
	}

	return
}

// check if the lines in toCheck are complete
// if so, remove them and update the grid
func (t tetris) checkLines() (toRemoveNum int, firstAvailable int, toRemove [4]bool) {

	count := -1
	firstAvailable = t.toCheck[0] - 1

	// get the lines that will disapear
CheckLoop:
	for l := t.toCheck[0]; l <= t.toCheck[1]; l++ {
		count++
		for x := 0; x < len(t.area[l]); x++ {
			if t.area[l][x] == 0 {
				firstAvailable = l
				continue CheckLoop
			}
		}
		toRemove[count] = true
		toRemoveNum++
	}

	return
}

func (t *tetris) removeLines() {

	// remove them from the grid from bottom to top

	// in the removal zone
	for y := t.toCheck[1]; y >= t.toCheck[0]; y-- {
		if t.firstAvailable >= 0 {
			t.area[y] = t.area[t.firstAvailable]
			t.firstAvailable--
			for t.firstAvailable >= t.toCheck[0] && t.toRemove[t.firstAvailable-t.toCheck[0]] {
				t.firstAvailable--
			}
		} else {
			t.area[y] = tetrisLine{}
		}
	}

	// above the removal zone
	for y := t.toCheck[0] - 1; y >= 0; y-- {
		if t.firstAvailable >= 0 {
			t.area[y] = t.area[t.firstAvailable]
			t.firstAvailable--
		} else {
			t.area[y] = tetrisLine{}
		}
	}

}

// check if their is anything in the above area
// which would mean that the game is lost
func (t tetris) lost() bool {
	for _, line := range t.area[:gInvisibleLines+t.deathLines] {
		for _, v := range line {
			if v != 0 {
				return true
			}
		}
	}
	return false
}

func getNewBlock() (block tetrisBlock) {

	switch rand.Intn(7) {
	case 0:
		block = getIBlock()
	case 1:
		block = getOBlock()
	case 2:
		block = getJBlock()
	case 3:
		block = getLBlock()
	case 4:
		block = getSBlock()
	case 5:
		block = getTBlock()
	case 6:
		block = getZBlock()
	}

	return

}

func (t tetris) draw(screen *ebiten.Image) {

	xNextOrigin := gPlayAreaSide + gPlayAreaWidth + gPlayAreaSide + gInfoLeftSide + gInfoShiftNext + gNextMargin
	yNextOrigin := gInfoTop + gInfoSmallBoxHeight + gScoreToLevel + gInfoBoxHeight + gLevelToLines + gInfoBoxHeight + gLinesToNext + gNextMargin

	t.nextBlock.draw(screen, xNextOrigin, yNextOrigin)

	xOrigin := gPlayAreaSide
	yOrigin := gSquareSideSize * -gInvisibleLines

	if t.removeLineAnimationStep == 0 {
		t.currentBlock.draw(screen, xOrigin, yOrigin)
	}

	for y, line := range t.area {
		for x, style := range line {
			if style != noStyle {

				// removal animation
				if t.removeLineAnimationStep%2 == 1 {
					if y >= t.toCheck[0] && y <= t.toCheck[1] &&
						t.toRemove[y-t.toCheck[0]] {
						if t.removeLineAnimationStep == 7 {
							continue
						}
						style = breakStyle
					}
				}

				options := ebiten.DrawImageOptions{}
				options.GeoM.Translate(float64(xOrigin+x*gSquareSideSize), float64(yOrigin+y*gSquareSideSize))
				screen.DrawImage(assets.ImageSquares.SubImage(image.Rect((style-1)*gSquareSideSize, 0, style*gSquareSideSize, gSquareSideSize)).(*ebiten.Image), &options)
			}
		}
	}

}
