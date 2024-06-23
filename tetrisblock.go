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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

type tetrisBlock struct {
	x, y   int           // position of upper left corner in squares
	r      int           // rotation state id
	states [4][4][4]bool // possible rotation states of the block
	style  int           // style of the block (for drawing)
	id     int8          // identifier of the block for randomisation
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

func (t *tetrisBlock) moveLeft(grid tetrisGrid) bool {
	t.x--
	if !t.isInValidPosition(grid) {
		t.x++
		return false
	}
	return true
}

func (t *tetrisBlock) moveRight(grid tetrisGrid) bool {
	t.x++
	if !t.isInValidPosition(grid) {
		t.x--
		return false
	}
	return true
}

func (b *tetrisBlock) updatePosition(rlMove int, dMove bool, grid tetrisGrid) (stuck bool, lrMoved bool) {

	if rlMove < 0 {
		lrMoved = b.moveLeft(grid)
	}

	if rlMove > 0 {
		lrMoved = b.moveRight(grid)
	}

	// try to move down or detect that the block is stuck
	if dMove {
		stuck = b.moveDown(grid)
	}

	return
}

func (t *tetrisBlock) rotateLeft(grid tetrisGrid) bool {
	t.r = (t.r + 3) % 4
	if !t.isInValidPosition(grid) {
		t.r = (t.r + 1) % 4
		return false
	}
	return true
}

func (t *tetrisBlock) rotateRight(grid tetrisGrid) bool {
	t.r = (t.r + 1) % 4
	if !t.isInValidPosition(grid) {
		t.r = (t.r + 3) % 4
		return false
	}
	return true
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
func (t tetrisBlock) draw(screen *ebiten.Image, gray uint8, xFrom, yFrom int) {

	for yRel, line := range t.states[t.r] {
		yAbs := t.y + yRel
		for xRel, square := range line {
			if square {
				xAbs := t.x + xRel

				options := ebiten.DrawImageOptions{}
				options.ColorScale.ScaleWithColor(color.Gray{gray})
				options.GeoM.Translate(float64(xFrom+xAbs*gSquareSideSize), float64(yFrom+yAbs*gSquareSideSize))
				screen.DrawImage(assets.ImageSquares.SubImage(image.Rect((t.style-1)*gSquareSideSize, 0, t.style*gSquareSideSize, gSquareSideSize)).(*ebiten.Image), &options)
			}
		}
	}
}

// check if
func canReplace(atX, atY int, preferedBlock, otherBlock tetrisBlock, grid tetrisGrid) bool {

	if preferedBlock.id >= 0 {
		preferedBlock.x = atX
		preferedBlock.y = atY
		return preferedBlock.isInValidPosition(grid)
	}

	otherBlock.x = atX
	otherBlock.y = atY

	return otherBlock.isInValidPosition(grid)
}
