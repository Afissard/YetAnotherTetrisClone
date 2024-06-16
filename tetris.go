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
	"math/rand"

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
	for l := toCheck[0]; l <= toCheck[1]; l++ {
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
				for firstAvailable >= toCheck[0] && remove[firstAvailable-toCheck[0]] {
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

	var block tetrisBlock

	switch rand.Intn(4) {
	case 0:
		block = &tBlock{}
	case 1:
		block = &squareBlock{}
	case 2:
		block = &lBlock{}
	case 3:
		block = &jBlock{}
	}

	block.setStyle(normalKind)
	block.setWaitPosition()

	return block
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
	tBlockStyle
	lBlockStyle
	jBlockStyle
)

// J block
/*
//  r =
//	 0    1     2   3
//	      #   #     ##
//  #o#   o   #o#   o
//    #  ##         #
*/
type jBlock struct {
	x, y  int // position of "o" in squares
	style int // style for constituting squares
	r     int // rotation state
}

func (b *jBlock) setStyle(kindOfBlock int) {
	b.style = jBlockStyle
}

func (b *jBlock) setWaitPosition() {
	b.x = 1
	b.y = 1
}

func (b *jBlock) setInitialPosition() {
	b.x = 4
	b.y = 2
}

func (b jBlock) canRotateInGridLimits(grid tetrisGrid) bool {
	return b.r == 0 ||
		(b.r == 1 && b.x+1 < len(grid[b.y])) ||
		(b.r == 2 && b.y+1 < len(grid)) ||
		(b.r == 3 && b.x-1 >= 0)
}

func (b jBlock) canRotateTo(grid tetrisGrid, r int) bool {
	return (r == 1 && grid[b.y-1][b.x] == 0 && grid[b.y+1][b.x] == 0 && grid[b.y+1][b.x-1] == 0) ||
		(r == 2 && grid[b.y-1][b.x-1] == 0 && grid[b.y][b.x+1] == 0 && grid[b.y][b.x-1] == 0) ||
		(r == 3 && grid[b.y-1][b.x] == 0 && grid[b.y+1][b.x] == 0 && grid[b.y-1][b.x+1] == 0) ||
		(r == 0 && grid[b.y+1][b.x+1] == 0 && grid[b.y][b.x-1] == 0 && grid[b.y][b.x+1] == 0)
}

func (b *jBlock) rotateLeft(grid tetrisGrid) {

	if !b.canRotateInGridLimits(grid) {
		return
	}

	goal := (b.r + 3) % 4

	if b.canRotateTo(grid, goal) {
		b.r = goal
	}
}

func (b *jBlock) rotateRight(grid tetrisGrid) {

	if !b.canRotateInGridLimits(grid) {
		return
	}

	goal := (b.r + 1) % 4

	if b.canRotateTo(grid, goal) {
		b.r = goal
	}
}

func (b *jBlock) moveLeft(grid tetrisGrid) {

	// grid limit
	if (b.r != 3 && b.x-2 < 0) ||
		(b.r == 3 && b.x-1 < 0) {
		return
	}

	// top block
	if (b.r == 1 && grid[b.y-1][b.x-1] != 0) ||
		(b.r == 2 && grid[b.y-1][b.x-2] != 0) ||
		(b.r == 3 && grid[b.y-1][b.x-1] != 0) {
		return
	}

	// center block
	if ((b.r == 0 || b.r == 2) && grid[b.y][b.x-2] != 0) ||
		((b.r == 1 || b.r == 3) && grid[b.y][b.x-1] != 0) {
		return
	}

	// bottom block
	if (b.r == 0 && grid[b.y+1][b.x] != 0) ||
		(b.r == 1 && grid[b.y+1][b.x-2] != 0) ||
		(b.r == 3 && grid[b.y+1][b.x-1] != 0) {
		return
	}

	b.x--
}

func (b *jBlock) moveRight(grid tetrisGrid) {

	// grid limit
	if (b.r != 1 && b.x+2 >= len(grid[b.y])) ||
		(b.r == 1 && b.x+1 >= len(grid[b.y])) {
		return
	}

	// top block
	if (b.r == 1 && grid[b.y-1][b.x+1] != 0) ||
		(b.r == 2 && grid[b.y-1][b.x] != 0) ||
		(b.r == 3 && grid[b.y-1][b.x+2] != 0) {
		return
	}

	// center block
	if ((b.r == 0 || b.r == 2) && grid[b.y][b.x+2] != 0) ||
		((b.r == 1 || b.r == 3) && grid[b.y][b.x+1] != 0) {
		return
	}

	// bottom block
	if (b.r == 0 && grid[b.y+1][b.x+2] != 0) ||
		(b.r == 1 && grid[b.y+1][b.x+1] != 0) ||
		(b.r == 3 && grid[b.y+1][b.x+1] != 0) {
		return
	}

	b.x++
}

func (b *jBlock) moveDown(grid tetrisGrid) (stuck bool) {

	// grid limit
	if b.r != 2 && b.y+2 >= len(grid) ||
		b.r == 2 && b.y+1 >= len(grid) {
		return true
	}

	// left block
	if (b.r == 0 && grid[b.y+1][b.x-1] != 0) ||
		(b.r == 1 && grid[b.y+2][b.x-1] != 0) ||
		(b.r == 2 && grid[b.y+1][b.x-1] != 0) {
		return true
	}

	// center block
	if (b.r == 0 || b.r == 2) && grid[b.y+1][b.x] != 0 {
		return true
	}
	if (b.r == 1 || b.r == 3) && grid[b.y+2][b.x] != 0 {
		return true
	}

	// right block
	if (b.r == 0 && grid[b.y+2][b.x+1] != 0) ||
		(b.r == 2 && grid[b.y+1][b.x+1] != 0) ||
		(b.r == 3 && grid[b.y][b.x+1] != 0) {
		return true
	}

	b.y++
	return
}

func (b *jBlock) writeInGrid(grid *tetrisGrid) (toCheck [2]int) {

	grid[b.y][b.x] = b.style

	if b.r == 1 || b.r == 3 {
		// top and bottom
		grid[b.y-1][b.x] = b.style
		grid[b.y+1][b.x] = b.style
	}

	if b.r == 0 || b.r == 2 {
		// left and right
		grid[b.y][b.x-1] = b.style
		grid[b.y][b.x+1] = b.style
	}

	// diag bottom left
	if b.r == 1 {
		grid[b.y+1][b.x-1] = b.style
	}

	// diag top left
	if b.r == 2 {
		grid[b.y-1][b.x-1] = b.style
	}

	// diag top right
	if b.r == 3 {
		grid[b.y-1][b.x+1] = b.style
	}

	// diag bottom right
	if b.r == 0 {
		grid[b.y+1][b.x+1] = b.style
	}

	toCheck = [2]int{b.y, b.y}
	if b.r != 0 {
		toCheck[0] -= 1
	}
	if b.r != 2 {
		toCheck[1] += 1
	}

	return
}

func (b *jBlock) draw(screen *ebiten.Image, x, y int) {
	// center
	vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)

	if b.r == 1 || b.r == 3 {
		// top
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
		// bottom
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	if b.r == 0 || b.r == 2 {
		// right
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
		// left
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag bottom left
	if b.r == 1 {
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag top left
	if b.r == 2 {
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag top right
	if b.r == 3 {
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag bottom right
	if b.r == 0 {
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}
}

// L block
/*
//  r =
//	 0    1     2   3
//	     ##     #   #
//  #o#   o   #o#   o
//  #     #         ##
*/
type lBlock struct {
	x, y  int // position of "o" in squares
	style int // style for constituting squares
	r     int // rotation state
}

func (b *lBlock) setStyle(kindOfBlock int) {
	b.style = lBlockStyle
}

func (b *lBlock) setWaitPosition() {
	b.x = 1
	b.y = 1
}

func (b *lBlock) setInitialPosition() {
	b.x = 4
	b.y = 2
}

func (b lBlock) canRotateInGridLimits(grid tetrisGrid) bool {
	return b.r == 0 ||
		(b.r == 1 && b.x+1 < len(grid[b.y])) ||
		(b.r == 2 && b.y+1 < len(grid)) ||
		(b.r == 3 && b.x-1 >= 0)
}

func (b lBlock) canRotateTo(grid tetrisGrid, r int) bool {
	return (r == 1 && grid[b.y-1][b.x-1] == 0 && grid[b.y-1][b.x] == 0 && grid[b.y+1][b.x] == 0) ||
		(r == 2 && grid[b.y-1][b.x+1] == 0 && grid[b.y][b.x+1] == 0 && grid[b.y][b.x-1] == 0) ||
		(r == 3 && grid[b.y-1][b.x] == 0 && grid[b.y+1][b.x] == 0 && grid[b.y+1][b.x+1] == 0) ||
		(r == 0 && grid[b.y+1][b.x-1] == 0 && grid[b.y][b.x-1] == 0 && grid[b.y][b.x+1] == 0)
}

func (b *lBlock) rotateLeft(grid tetrisGrid) {

	if !b.canRotateInGridLimits(grid) {
		return
	}

	goal := (b.r + 3) % 4

	if b.canRotateTo(grid, goal) {
		b.r = goal
	}
}

func (b *lBlock) rotateRight(grid tetrisGrid) {

	if !b.canRotateInGridLimits(grid) {
		return
	}

	goal := (b.r + 1) % 4

	if b.canRotateTo(grid, goal) {
		b.r = goal
	}
}

func (b *lBlock) moveLeft(grid tetrisGrid) {

	// grid limit
	if (b.r != 3 && b.x-2 < 0) ||
		(b.r == 3 && b.x-1 < 0) {
		return
	}

	// top block
	if (b.r == 1 && grid[b.y-1][b.x-2] != 0) ||
		(b.r == 2 && grid[b.y-1][b.x] != 0) ||
		(b.r == 3 && grid[b.y-1][b.x-1] != 0) {
		return
	}

	// center block
	if ((b.r == 0 || b.r == 2) && grid[b.y][b.x-2] != 0) ||
		((b.r == 1 || b.r == 3) && grid[b.y][b.x-1] != 0) {
		return
	}

	// bottom block
	if (b.r == 0 && grid[b.y+1][b.x-2] != 0) ||
		(b.r == 1 && grid[b.y+1][b.x-1] != 0) ||
		(b.r == 3 && grid[b.y+1][b.x-1] != 0) {
		return
	}

	b.x--
}

func (b *lBlock) moveRight(grid tetrisGrid) {

	// grid limit
	if (b.r != 1 && b.x+2 >= len(grid[b.y])) ||
		(b.r == 1 && b.x+1 >= len(grid[b.y])) {
		return
	}

	// top block
	if (b.r == 1 && grid[b.y-1][b.x+1] != 0) ||
		(b.r == 2 && grid[b.y-1][b.x+2] != 0) ||
		(b.r == 3 && grid[b.y-1][b.x+1] != 0) {
		return
	}

	// center block
	if ((b.r == 0 || b.r == 2) && grid[b.y][b.x+2] != 0) ||
		((b.r == 1 || b.r == 3) && grid[b.y][b.x+1] != 0) {
		return
	}

	// bottom block
	if (b.r == 0 && grid[b.y+1][b.x] != 0) ||
		(b.r == 1 && grid[b.y+1][b.x+1] != 0) ||
		(b.r == 3 && grid[b.y+1][b.x+2] != 0) {
		return
	}

	b.x++
}

func (b *lBlock) moveDown(grid tetrisGrid) (stuck bool) {

	// grid limit
	if b.r != 2 && b.y+2 >= len(grid) ||
		b.r == 2 && b.y+1 >= len(grid) {
		return true
	}

	// left block
	if (b.r == 0 && grid[b.y+2][b.x-1] != 0) ||
		(b.r == 1 && grid[b.y][b.x-1] != 0) ||
		(b.r == 2 && grid[b.y+1][b.x-1] != 0) {
		return true
	}

	// center block
	if (b.r == 0 || b.r == 2) && grid[b.y+1][b.x] != 0 {
		return true
	}
	if (b.r == 1 || b.r == 3) && grid[b.y+2][b.x] != 0 {
		return true
	}

	// right block
	if (b.r == 0 && grid[b.y+1][b.x+1] != 0) ||
		(b.r == 2 && grid[b.y+1][b.x+1] != 0) ||
		(b.r == 3 && grid[b.y+2][b.x+1] != 0) {
		return true
	}

	b.y++
	return
}

func (b *lBlock) writeInGrid(grid *tetrisGrid) (toCheck [2]int) {

	grid[b.y][b.x] = b.style

	if b.r == 1 || b.r == 3 {
		// top and bottom
		grid[b.y-1][b.x] = b.style
		grid[b.y+1][b.x] = b.style
	}

	if b.r == 0 || b.r == 2 {
		// left and right
		grid[b.y][b.x-1] = b.style
		grid[b.y][b.x+1] = b.style
	}

	// diag bottom left
	if b.r == 0 {
		grid[b.y+1][b.x-1] = b.style
	}

	// diag top left
	if b.r == 1 {
		grid[b.y-1][b.x-1] = b.style
	}

	// diag top right
	if b.r == 2 {
		grid[b.y-1][b.x+1] = b.style
	}

	// diag bottom right
	if b.r == 3 {
		grid[b.y+1][b.x+1] = b.style
	}

	toCheck = [2]int{b.y, b.y}
	if b.r != 0 {
		toCheck[0] -= 1
	}
	if b.r != 2 {
		toCheck[1] += 1
	}
	return
}

func (b *lBlock) draw(screen *ebiten.Image, x, y int) {
	// center
	vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)

	if b.r == 1 || b.r == 3 {
		// top
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
		// bottom
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	if b.r == 0 || b.r == 2 {
		// right
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
		// left
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag bottom left
	if b.r == 0 {
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag top left
	if b.r == 1 {
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag top right
	if b.r == 2 {
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// diag bottom right
	if b.r == 3 {
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}
}

// T block
/*
//	r =
//	 0    1   2   3
//	      #   #   #
//	#o#  #o  #o#  o#
//	 #    #       #
*/
type tBlock struct {
	x, y  int // position of "o" in squares
	style int // style for constituting squares
	r     int // rotation state
}

func (b *tBlock) setStyle(kindOfBlock int) {
	b.style = tBlockStyle
}

func (b *tBlock) setWaitPosition() {
	b.x = 1
	b.y = 1
}

func (b *tBlock) setInitialPosition() {
	b.x = 4
	b.y = 2
}

func (b tBlock) canRotate(grid tetrisGrid) bool {

	return (b.r == 0 && grid[b.y-1][b.x] == 0) ||
		(b.r == 1 && b.x+1 < len(grid[b.y]) && grid[b.y][b.x+1] == 0) ||
		(b.r == 2 && b.y+1 < len(grid) && grid[b.y+1][b.x] == 0) ||
		(b.r == 3 && b.x-1 >= 0 && grid[b.y][b.x-1] == 0)
}

func (b *tBlock) rotateLeft(grid tetrisGrid) {

	if !b.canRotate(grid) {
		return
	}

	b.r = (b.r + 3) % 4
}

func (b *tBlock) rotateRight(grid tetrisGrid) {

	if !b.canRotate(grid) {
		return
	}

	b.r = (b.r + 1) % 4
}

func (b *tBlock) moveLeft(grid tetrisGrid) {
	// area border
	if (b.r != 3 && b.x-2 < 0) ||
		(b.r == 3 && b.x-1 < 0) {
		return
	}

	// blocked on top
	if b.r != 0 &&
		grid[b.y-1][b.x-1] != 0 {
		return
	}

	// blocked on bottom
	if b.r != 2 &&
		grid[b.y+1][b.x-1] != 0 {
		return
	}

	// blocked on center
	if (b.r == 3 && grid[b.y][b.x-1] != 0) ||
		(b.r != 3 && grid[b.y][b.x-2] != 0) {
		return
	}

	b.x--
}

func (b *tBlock) moveRight(grid tetrisGrid) {

	// area border
	if (b.r != 1 && b.x+2 >= len(grid[b.y])) ||
		(b.r == 1 && b.x+1 >= len(grid[b.y])) {
		return
	}

	// blocked on top
	if b.r != 0 &&
		grid[b.y-1][b.x+1] != 0 {
		return
	}

	// blocked on bottom
	if b.r != 2 &&
		grid[b.y+1][b.x+1] != 0 {
		return
	}

	// blocked on center
	if (b.r == 1 && grid[b.y][b.x+1] != 0) ||
		(b.r != 1 && grid[b.y][b.x+2] != 0) {
		return
	}

	b.x++
}

func (b *tBlock) moveDown(grid tetrisGrid) (stuck bool) {
	if (b.r != 2 && b.y+2 >= len(grid)) ||
		(b.r == 2 && b.y+1 >= len(grid)) ||
		(b.r != 3 && grid[b.y+1][b.x-1] != 0) ||
		(b.r != 2 && grid[b.y+2][b.x] != 0) ||
		(b.r == 2 && grid[b.y+1][b.x] != 0) ||
		(b.r != 1 && grid[b.y+1][b.x+1] != 0) {
		return true
	}

	b.y++
	return
}

func (b *tBlock) writeInGrid(grid *tetrisGrid) (toCheck [2]int) {

	grid[b.y][b.x] = b.style

	// top square
	if b.r != 0 {
		grid[b.y-1][b.x] = b.style
	}

	// right square
	if b.r != 1 {
		grid[b.y][b.x+1] = b.style
	}

	// bottom square
	if b.r != 2 {
		grid[b.y+1][b.x] = b.style
	}

	// left square
	if b.r != 3 {
		grid[b.y][b.x-1] = b.style
	}

	toCheck = [2]int{b.y, b.y}
	if b.r != 0 {
		toCheck[0] -= 1
	}

	if b.r != 2 {
		toCheck[1] += 1
	}
	return
}

func (b *tBlock) draw(screen *ebiten.Image, x, y int) {
	// center
	vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)

	// top
	if b.r != 0 {
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y-1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	//right
	if b.r != 1 {
		vector.DrawFilledRect(screen, float32(x+(b.x+1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// bottom
	if b.r != 2 {
		vector.DrawFilledRect(screen, float32(x+b.x*gSquareSideSize), float32(y+(b.y+1)*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}

	// left
	if b.r != 3 {
		vector.DrawFilledRect(screen, float32(x+(b.x-1)*gSquareSideSize), float32(y+b.y*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 220}, false)
	}
}

// Square block
/*
//	o#
//	##
*/
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

	if b.x+2 >= len(grid[b.y]) || (b.y >= 0 && grid[b.y][b.x+2] != 0) || (b.y+1 >= 0 && grid[b.y+1][b.x+2] != 0) {
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
