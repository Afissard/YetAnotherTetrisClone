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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

func (g *game) Draw(screen *ebiten.Image) {

	// draw background

	options := ebiten.DrawImageOptions{}
	screen.DrawImage(assets.ImageBack, &options)

	// draw current play

	g.currentPlay.draw(screen)

	// draw number of lines destroyed
	drawNumberAt(screen, gWidth-gXLinesFromRightSide+1, gYLinesFromTop, g.currentPlay.numLines)

	// draw score
	drawNumberAt(screen, gWidth-gXScoreFromRightSide+1, gYScoreFromTop, g.score)

	// draw level
	drawNumberAt(screen, gWidth-gXLevelFromRightSide+1, gYLevelFromTop, g.level)

}
