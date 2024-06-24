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
	"github.com/loig/ebitenginegamejam2024/assets"
)

func (g *game) Draw(screen *ebiten.Image) {

	switch g.state {
	case statePlay:
		g.drawPlay(screen, 255)
	case stateBalance:
		g.drawPlay(screen, 100)
		g.balance.draw(screen)
	case stateLost:
		g.drawPlay(screen, 100)
		g.money.draw(screen)
	case stateImprove:
		g.drawStateImprove(screen)
	}

	// play sounds
	g.audio.PlaySounds()

}

func (g game) drawPlay(screen *ebiten.Image, gray uint8) {
	// draw background
	options := ebiten.DrawImageOptions{}
	options.ColorScale.ScaleWithColor(color.Gray{gray})
	screen.DrawImage(assets.ImageBack, &options)
	// draw current play
	g.currentPlay.draw(screen, gray)
	// draw number of lines destroyed
	drawNumberAt(screen, gray, gWidth-gXLinesFromRightSide+gMultFactor, gYLinesFromTop, g.balance.getGoalLines()-g.currentPlay.numLines)
	// draw score
	drawNumberAt(screen, gray, gWidth-gXScoreFromRightSide+gMultFactor, gYScoreFromTop, g.currentPlay.score)
	// draw level
	drawNumberAt(screen, gray, gWidth-gXLevelFromRightSide+gMultFactor, gYLevelFromTop, g.level)
	// hide lines
	g.fog.draw(screen)
	// death lines
	if !g.firstPlay || g.level > 0 {
		vector.StrokeRect(screen, float32(gPlayAreaSide), 0, float32(gPlayAreaWidth), float32(g.currentPlay.deathLines*gSquareSideSize), 2, color.Black, false)
	}
}
