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

func (g *game) Draw(screen *ebiten.Image) {

	// draw backgrounds from left to right

	var leftShift float32

	vector.DrawFilledRect(screen, leftShift, 0, float32(gMultFactor*gPlayAreaSide), float32(gMultFactor*gPlayAreaHeight), color.Gray{Y: 128}, false)

	leftShift += float32(gMultFactor * gPlayAreaSide)

	vector.DrawFilledRect(screen, leftShift, 0, float32(gMultFactor*gPlayAreaWidth), float32(gMultFactor*gPlayAreaHeight), color.White, false)

	leftShift += float32(gMultFactor * gPlayAreaWidth)

	vector.DrawFilledRect(screen, leftShift, 0, float32(gMultFactor*gPlayAreaSide), float32(gMultFactor*gPlayAreaHeight), color.Gray{Y: 128}, false)

	leftShift += float32(gMultFactor * gPlayAreaSide)

	// draw info boxes from top to bottom

	leftShift += float32(gMultFactor * gInfoLeftSide)

	var topShift float32 = float32(gMultFactor * gInfoTop)

	vector.DrawFilledRect(screen, leftShift, topShift, float32(gMultFactor*gInfoWidth), float32(gMultFactor*gInfoSmallBoxHeight), color.Gray{Y: 128}, false)

	topShift += float32(gMultFactor * (gInfoSmallBoxHeight + gScoreToLevel))

	vector.DrawFilledRect(screen, leftShift, topShift, float32(gMultFactor*gInfoWidth), float32(gMultFactor*gInfoBoxHeight), color.Gray{Y: 128}, false)

	topShift += float32(gMultFactor * (gInfoBoxHeight + gLevelToLines))

	vector.DrawFilledRect(screen, leftShift, topShift, float32(gMultFactor*gInfoWidth), float32(gMultFactor*gInfoBoxHeight), color.Gray{Y: 128}, false)

	topShift += float32(gMultFactor * (gInfoBoxHeight + gLinesToNext))
	leftShift += float32(gMultFactor * (gInfoWidth - gNextBoxSide))

	vector.DrawFilledRect(screen, leftShift, topShift, float32(gMultFactor*gNextBoxSide), float32(gMultFactor*gNextBoxSide), color.Gray{Y: 128}, false)

	// draw current play

	g.currentPlay.draw(screen)

}
