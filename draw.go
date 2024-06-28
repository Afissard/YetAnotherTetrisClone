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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

func (g *game) Draw(screen *ebiten.Image) {

	switch g.state {
	case stateControls:
		screen.DrawImage(assets.ImageControls, &ebiten.DrawImageOptions{})
	case stateCredits:
		screen.DrawImage(assets.ImageCredits, &ebiten.DrawImageOptions{})
	case stateTitle:
		if g.firstPlay {
			screen.DrawImage(assets.ImageTitle1, &ebiten.DrawImageOptions{})
		} else {
			screen.DrawImage(assets.ImageTitle2, &ebiten.DrawImageOptions{})
		}
		if g.titleSelect == 0 {
			drawArrow(screen, gWidth/2-150, 3*gHeight/4+20, math.Pi/2, g.titleFrame)
		} else {
			drawArrow(screen, gWidth/2-250, 3*gHeight/4+128, math.Pi/2, g.titleFrame)
		}
	case statePlay:
		g.drawPlay(screen, 255)
	case stateBalance:
		g.drawPlay(screen, 100)
		g.balance.draw(screen)
	case stateLost:
		g.drawPlay(screen, 100)
		g.money.draw(screen)
	case stateImprove:
		g.drawShop(screen)
		g.drawStateImprove(screen)
	case stateWon:
		screen.DrawImage(assets.ImageWin, &ebiten.DrawImageOptions{})
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(gWidth/2)-175, float64(gHeight/2)-140)
		options.GeoM.Translate(0, -float64(gAnimRocket[g.winFrame%len(gAnimRocket)]))
		screen.DrawImage(assets.ImageRocket, &options)
	}

}

func (g game) drawShop(screen *ebiten.Image) {

	screen.DrawImage(assets.ImageShopBack, &ebiten.DrawImageOptions{})

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(gWidth-gShopTitleWidth)/2, float64(gTitleMargin))
	screen.DrawImage(assets.ImageShopTitle, &options)

}

func (g game) drawPlay(screen *ebiten.Image, gray uint8) {
	// draw background
	options := ebiten.DrawImageOptions{}
	options.ColorScale.ScaleWithColor(color.Gray{gray})
	screen.DrawImage(assets.ImageBack, &options)
	// draw death lines
	g.drawDeathLines(screen, gray)

	// draw current play
	g.currentPlay.draw(screen, gray)
	// draw number of lines destroyed
	drawNumberAt(screen, gray, gWidth-gXLinesFromRightSide+gMultFactor, gYLinesFromTop, g.currentPlay.numLines, g.balance.getGoalLines())
	// draw score
	drawNumberAt(screen, gray, gWidth-gXScoreFromRightSide+gMultFactor, gYScoreFromTop, g.currentPlay.score, -1)
	// draw level
	drawNumberAt(screen, gray, gWidth-gXLevelFromRightSide+gMultFactor, gYLevelFromTop, g.level+1, g.goalLevel)
	// hide lines
	g.fog.draw(screen, gray)
}

func (g game) drawDeathLines(screen *ebiten.Image, gray uint8) {
	// death lines
	options := ebiten.DrawImageOptions{}

	options.ColorScale.ScaleWithColor(color.Gray{gray})
	scaling := float64(gSquareSideSize) / float64(gDangerSide)
	options.GeoM.Scale(scaling, scaling)
	options.GeoM.Translate(float64(gPlayAreaSide), 0)
	mult := 1
	for line := 0; line < g.currentPlay.deathLines; line++ {
		for pos := 0; pos < gPlayAreaWidthInBlocks; pos++ {
			screen.DrawImage(assets.ImageDanger, &options)
			options.GeoM.Translate(float64(mult*gSquareSideSize), 0)
		}
		mult = -mult
		options.GeoM.Translate(float64(mult*gSquareSideSize), float64(gSquareSideSize))
	}
}
