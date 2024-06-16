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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type tetrisLine = [gPlayAreaWidthInBlocks]int
type tetrisGrid = [gPlayAreaHeightInBlocks + gInvisibleLines]tetrisLine

// Interface for any tetris block
type tetrisBlock interface {
	setStyle(kindOfBlock int)
	setWaitPosition()
	setInitialPosition()
	// everything is done on a square basis
	rotateLeft(tetrisGrid)
	rotateRight(tetrisGrid)
	moveLeft(tetrisGrid)
	moveRight(tetrisGrid)
	moveDown(tetrisGrid) (stuck bool)
	// fill the grid with the block and tell
	// which lines may be completed (given as an interval)
	writeInGrid(*tetrisGrid) (linesToCheck [2]int)
	// draw on screen given coordinates in pixels for the
	// origin of the area
	draw(screen *ebiten.Image, x, y int)
}

func updateTBPosition(b tetrisBlock, rlMove int, dMove bool, grid tetrisGrid) (stuck bool) {

	if rlMove < 0 {
		b.moveLeft(grid)
	}

	if rlMove > 0 {
		b.moveRight(grid)
	}

	// try to move down or detect that the block is stuck
	if dMove {
		stuck = b.moveDown(grid)
	}

	return
}

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
}

func (t *tetris) init() {
	t.area = tetrisGrid{}
	t.currentBlock = getNewBlock()
	t.currentBlock.setInitialPosition()
	t.nextBlock = getNewBlock()
	t.autoDownFrame = 0
	t.autoDownFrameLimit = 30
	t.manualDownFrame = 0
	t.manualDownFrameLimit = 6
	t.lrMoveFrame = 0
	t.lrMoveFrameLimit = 6
	t.lrFirstMoveFrame = 0
	t.lrFirstMoveFrameLimit = 15
}

func (t *tetris) update(moveDownRequest, moveLeftRequest, moveRightRequest, rotateLeft, rotateRight bool) {

	if rotateLeft {
		t.currentBlock.rotateLeft(t.area)
	}

	if rotateRight {
		t.currentBlock.rotateRight(t.area)
	}

	// left/right movements of blocks handling
	xMove := 0
	if moveLeftRequest {
		xMove--
	}
	if moveRightRequest {
		xMove++
	}

	if !moveLeftRequest && !moveRightRequest {
		t.lrMoveFrame = 0
		t.lrFirstMoveFrame = 0
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
	}

	if moveDownRequest {
		manualDown = t.manualDownFrame == 0
		t.manualDownFrame++
		if t.manualDownFrame >= t.manualDownFrameLimit {
			t.manualDownFrame = 0
		}
	}

	// update position according to movements requests
	if updateTBPosition(t.currentBlock, xMove, autoDown || manualDown, t.area) {
		toCheck := t.currentBlock.writeInGrid(&t.area)

		t.checkLinesAndUpdate(toCheck)
		if t.lost() {
			t.init()
		}

		t.currentBlock = t.nextBlock
		t.currentBlock.setInitialPosition()
		t.nextBlock = getNewBlock()
	}
}

// check if the lines in toCheck are complete
// if so, remove them and update the grid
func (t *tetris) checkLinesAndUpdate(toCheck [2]int) {

	checkSize := toCheck[1] - toCheck[0] + 1
	remove := make([]bool, checkSize)
	toRemove := 0
	count := -1
	firstAvailable := toCheck[0] - 1

	// get the lines that will disapear
CheckLoop:
	for l := toCheck[1]; l >= toCheck[0]; l-- {
		count++
		for x := 0; x < len(t.area[l]); x++ {
			if t.area[l][x] == 0 {
				firstAvailable = l
				continue CheckLoop
			}
		}
		remove[count] = true
		toRemove++
	}

	if toRemove > 0 {
		// remove them from the grid from bottom to top

		// in the removal zone
		for y := toCheck[1]; y >= toCheck[0]; y-- {
			if firstAvailable >= 0 {
				t.area[y] = t.area[firstAvailable]
				firstAvailable--
				for firstAvailable >= toCheck[0] && remove[toCheck[1]-firstAvailable] {
					firstAvailable--
				}
			} else {
				t.area[y] = tetrisLine{}
			}
		}

		// above the removal zone
		for y := toCheck[0] - 1; y >= 0; y-- {
			if firstAvailable >= 0 {
				t.area[y] = t.area[firstAvailable]
				firstAvailable--
			} else {
				t.area[y] = tetrisLine{}
			}
		}

	}

}

// check if their is anything in the above area
// which would mean that the game is lost
func (t tetris) lost() bool {
	for _, line := range t.area[:gInvisibleLines+1] {
		for _, v := range line {
			if v != 0 {
				return true
			}
		}
	}
	return false
}

func getNewBlock() tetrisBlock {

	block := squareBlock{}

	block.setStyle(normalKind)
	block.setWaitPosition()

	return &block
}

func (t tetris) draw(screen *ebiten.Image) {

	xNextOrigin := gMultFactor * (gPlayAreaSide + gPlayAreaWidth + gPlayAreaSide + gInfoLeftSide + (gInfoWidth - gNextBoxSide) + gNextMargin)
	yNextOrigin := gMultFactor * (gInfoTop + gInfoSmallBoxHeight + gScoreToLevel + gInfoBoxHeight + gLevelToLines + gInfoBoxHeight + gLinesToNext + gNextMargin)

	t.nextBlock.draw(screen, xNextOrigin, yNextOrigin)

	xOrigin := gMultFactor * gPlayAreaSide
	yOrigin := gMultFactor * gSquareSideSize * -gInvisibleLines

	t.currentBlock.draw(screen, xOrigin, yOrigin)

	for y, line := range t.area {
		for x, v := range line {
			if v != 0 {
				vector.DrawFilledRect(screen, float32(xOrigin+x*gSquareSideSize), float32(yOrigin+y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 40}, false)
			}
		}
	}

}

// styles for blocks
const (
	normalKind int = iota
)

const (
	noStyle int = iota
	squareBlockStyle
)

// Square block
//
//	o#
//	##
type squareBlock struct {
	x, y  int // position of "o" in squares
	style int // style for constituting squares
}

func (b *squareBlock) setStyle(kindOfBlock int) {
	b.style = squareBlockStyle
}

func (b *squareBlock) setWaitPosition() {
	b.x = 1
	b.y = 1
}

func (b *squareBlock) setInitialPosition() {
	b.x = 4
	b.y = 2
}

func (b *squareBlock) rotateLeft(grid tetrisGrid)  {}
func (b *squareBlock) rotateRight(grid tetrisGrid) {}

func (b *squareBlock) moveLeft(grid tetrisGrid) {

	if b.x-1 < 0 || (b.y >= 0 && grid[b.y][b.x-1] != 0) || (b.y+1 >= 0 && grid[b.y+1][b.x-1] != 0) {
		return
	}

	b.x--
}

func (b *squareBlock) moveRight(grid tetrisGrid) {

	if b.x+2 >= len(grid[0]) || (b.y >= 0 && grid[b.y][b.x+2] != 0) || (b.y+1 >= 0 && grid[b.y+1][b.x+2] != 0) {
		return
	}

	b.x++
}

func (b *squareBlock) moveDown(grid tetrisGrid) (stuck bool) {

	if b.y+2 >= len(grid) || grid[b.y+2][b.x] != 0 || grid[b.y+2][b.x+1] != 0 {
		return true
	}

	b.y++
	return
}

func (b *squareBlock) writeInGrid(grid *tetrisGrid) (toCheck [2]int) {

	grid[b.y][b.x] = b.style
	grid[b.y+1][b.x] = b.style
	grid[b.y][b.x+1] = b.style
	grid[b.y+1][b.x+1] = b.style

	return [2]int{b.y, b.y + 1}
}

func (b *squareBlock) draw(screen *ebiten.Image, x, y int) {

	vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(2*gSquareSideSize), float32(2*gSquareSideSize), color.Gray{Y: 220}, false)

}