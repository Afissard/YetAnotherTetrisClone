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
	imp.prices[improveLife] = []int{5, 2, 3}
	imp.prices[improveHold] = []int{12}
	imp.prices[improveResetAutoDown] = []int{11}
	imp.prices[improveHideMove] = []int{9, 2, 3}
	return
}

func (g game) drawStateImprove(screen *ebiten.Image) {

	yStart := 128
	xSeparator := 100
	ySeparator := xSeparator

	//drawMoney(screen, gWidth/2, yStart, g.money.money, true)

	x := (gWidth - (2*gImproveTextWidth + xSeparator)) / 2
	y := yStart

	xTranslation := 0
	firstLine := true
	for i := 0; i <= numImprove; i++ {

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(assets.ImageImprovements.SubImage(image.Rect(0, i*gImproveTextHeight, gImproveTextWidth, (i+1)*gImproveTextWidth)).(*ebiten.Image), &options)

		if i != numImprove {
			if len(g.improv.prices[i]) > g.improv.levels[i] {
				drawMoney(screen, x+3*gImproveTextWidth/5, y+gImproveTextHeight, g.improv.prices[i][g.improv.levels[i]], false, 0.5)
			} else {
				options.GeoM.Translate(float64(gImproveTextWidth-gMaxWidth)/2, float64(gImproveTextHeight)-10)
				screen.DrawImage(assets.ImageMax, &options)
			}
		}

		if g.improv.current == i {
			options.GeoM.Translate(float64(-gArrowWidth), float64(gImproveTextHeight-gArrowHeight)/2)
			screen.DrawImage(assets.ImageImprovementsArrow, &options)
		}

		if !firstLine || i < numImprove/2-1 {
			x += gImproveTextWidth + xSeparator
			xTranslation += gImproveTextWidth + xSeparator
		} else {
			x -= xTranslation
			y += gImproveTextHeight + ySeparator
			firstLine = false
		}

	}

}

func (g *game) updateStateImprove() bool {

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
			g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
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
