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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

const (
	improveLife int = iota
	improveHold
	improveResetAutoDown
	improveHideMove
	numImprove
)

type improvements struct {
	prices  [numImprove][]int
	levels  [numImprove]int
	current int
}

func setupImprovements() (imp improvements) {
	imp.prices[improveLife] = []int{1, 2, 3}
	imp.prices[improveHold] = []int{1}
	imp.prices[improveResetAutoDown] = []int{1}
	imp.prices[improveHideMove] = []int{1, 2, 3}
	return
}

func (g game) drawStateImprove(screen *ebiten.Image) {

	yStart := 128
	xSeparator := 100

	drawMoney(screen, gWidth/2, yStart, g.money.money, true)

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(gWidth-gImproveTextWidth)/2, float64(yStart))

	for i := 0; i <= numImprove; i++ {
		options.GeoM.Translate(0, float64(gImproveTextHeight))
		screen.DrawImage(assets.ImageImprovements.SubImage(image.Rect(0, i*gImproveTextHeight, gImproveTextWidth, (i+1)*gImproveTextWidth)).(*ebiten.Image), &options)

		if i != numImprove {
			if len(g.improv.prices[i]) > g.improv.levels[i] {
				drawMoney(screen, (gWidth+gImproveTextWidth)/2+xSeparator, yStart+(i+1)*gImproveTextHeight+gImproveTextHeight/2, g.improv.prices[i][g.improv.levels[i]], false)
			} else {
				shift := gImproveTextWidth + xSeparator
				options.GeoM.Translate(float64(shift), 0)
				screen.DrawImage(assets.ImageMax, &options)
				options.GeoM.Translate(-float64(shift), 0)
			}
		}

		if g.improv.current == i {
			options.GeoM.Translate(float64(-gImproveTextHeight), 0)
			screen.DrawImage(assets.ImageImprovementsArrow, &options)
			options.GeoM.Translate(float64(gImproveTextHeight), 0)
		}
	}

}

func (g *game) updateStateImprove() bool {

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
			g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.improv.current = (g.improv.current + 1) % (numImprove + 1)
		for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
			g.improv.current = (g.improv.current + 1) % (numImprove + 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.improv.current == numImprove {
			return true
		}

		if g.improv.prices[g.improv.current][g.improv.levels[g.improv.current]] <= g.money.money {
			g.money.money -= g.improv.prices[g.improv.current][g.improv.levels[g.improv.current]]
			g.improv.levels[g.improv.current]++
		}
	}

	return false
}
