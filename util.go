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

// draw a number right alligned in a rectangle which top right is given by (x, y) in pixels
func drawNumberAt(screen *ebiten.Image, gray uint8, x, y int, num int) {

	if num < 0 {
		num = 0
	}

	options := ebiten.DrawImageOptions{}

	options.ColorScale.ScaleWithColor(color.Gray{gray})

	options.GeoM.Translate(float64(x), float64(y))

	atLeastOnce := true
	for num > 0 || atLeastOnce {
		atLeastOnce = false
		digit := num % 10
		num = num / 10

		options.GeoM.Translate(float64(-gSquareSideSize), float64(0))
		screen.DrawImage(assets.ImageDigits.SubImage(image.Rect(digit*gSquareSideSize, 0, (digit+1)*gSquareSideSize, gSquareSideSize)).(*ebiten.Image), &options)
	}

}

// remove one element from a slice of int
func removeElement(t []int, pos int) []int {
	t[pos], t[len(t)-1] = t[len(t)-1], t[pos]
	return t[:len(t)-1]
}
