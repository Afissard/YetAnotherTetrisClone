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

type tetrisBlock struct {
	x, y   int           // position of upper left corner in squares
	r      int           // rotation state id
	states [4][4][4]bool // possible rotation states of the block
	style  int           // style of the block (for drawing)
}

func (t *tetrisBlock) setInitialPosition() {
	t.x = 3
	t.y = 1
}

// x and y are given in squares
func (t tetrisBlock) isInValidPosition(grid tetrisGrid) bool {

	for yRel, line := range t.states[t.r] {
		yAbs := t.y + yRel
		for xRel, square := range line {
			xAbs := t.x + xRel
			if square {
				if yAbs >= len(grid) ||
					xAbs < 0 ||
					xAbs >= len(grid[yAbs]) ||
					grid[yAbs][xAbs] != 0 {
					return false
				}
			}
		}
	}

	return true
}

func (t *tetrisBlock) moveDown(grid tetrisGrid) (stuck bool) {
	t.y++
	if !t.isInValidPosition(grid) {
		t.y--
		stuck = true
	}
	return
}

func (t *tetrisBlock) moveLeft(grid tetrisGrid) {
	t.x--
	if !t.isInValidPosition(grid) {
		t.x++
	}
}

func (t *tetrisBlock) moveRight(grid tetrisGrid) {
	t.x++
	if !t.isInValidPosition(grid) {
		t.x--
	}
}

func (b *tetrisBlock) updatePosition(rlMove int, dMove bool, grid tetrisGrid) (stuck bool) {

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

func (t *tetrisBlock) rotateLeft(grid tetrisGrid) {
	t.r = (t.r + 3) % 4
	if !t.isInValidPosition(grid) {
		t.r = (t.r + 1) % 4
	}
}

func (t *tetrisBlock) rotateRight(grid tetrisGrid) {
	t.r = (t.r + 1) % 4
	if !t.isInValidPosition(grid) {
		t.r = (t.r + 3) % 4
	}
}

func (t tetrisBlock) writeInGrid(grid *tetrisGrid) (toCheck [2]int) {

	yMin := len(grid)
	yMax := 0

	for yRel, line := range t.states[t.r] {
		yAbs := t.y + yRel
		for xRel, square := range line {
			if square {
				xAbs := t.x + xRel
				grid[yAbs][xAbs] = t.style
				if yAbs < yMin {
					yMin = yAbs
				}
				if yAbs > yMax {
					yMax = yAbs
				}
			}
		}
	}

	return [2]int{yMin, yMax}
}

// xFrom, yFrom in pixels
func (t tetrisBlock) draw(screen *ebiten.Image, xFrom, yFrom int) {

	for yRel, line := range t.states[t.r] {
		yAbs := t.y + yRel
		for xRel, square := range line {
			if square {
				xAbs := t.x + xRel
				vector.DrawFilledRect(screen, float32(xFrom+xAbs*gSquareSideSize), float32(yFrom+yAbs*gSquareSideSize), float32(gSquareSideSize), float32(gSquareSideSize), color.Gray{Y: 40}, false)
			}
		}
	}
}
