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
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

const (
	numArrowBlinkFrame int = 30
)

const (
	improveLife int = iota
	improveHold
	improveResetAutoDown
	improveHideMove
	numImprove
)

type improvements struct {
	prices          [numImprove][]int
	levels          [numImprove]int
	current         int
	arrowBlinkFrame int
}

func setupImprovements() (imp improvements) {
	imp.prices[improveLife] = []int{10, 50, 150}
	imp.prices[improveHold] = []int{150}
	imp.prices[improveResetAutoDown] = []int{300}
	imp.prices[improveHideMove] = []int{20, 75, 250}
	return
}

func (i *improvements) reset() {
	i.arrowBlinkFrame = 0
	i.current = 0
	for i.current != numImprove && i.levels[i.current] >= len(i.prices[i.current]) {
		i.current = (i.current + 1) % (numImprove + 1)
	}
}

func drawMaxed(screen *ebiten.Image, x, y int) {
	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(assets.ImageMax, &options)
}

func drawArrow(screen *ebiten.Image, x, y int, rotate float64, blink int) {
	if blink < 2*numArrowBlinkFrame/3 {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Rotate(rotate)
		options.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(assets.ImageImprovementsArrow, &options)
	}
}

func drawContinue(screen *ebiten.Image, x, y int, selected bool, blink int) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(assets.ImageContinue, &options)

	if selected {
		drawArrow(screen, x-5, y+(gContinueHeight-gArrowWidth)/2, math.Pi/2, blink)
		drawArrow(screen, x+gContinueWidth+5, y+(gContinueHeight-gArrowWidth)/2+gArrowWidth, -math.Pi/2, blink)
	}
}

func drawShopText(screen *ebiten.Image, x, y int, selection int) {
	if selection < numImprove {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(assets.ImageTextShop.SubImage(image.Rect(0, selection*gTextMalusHeight, gTextMalusWidth, (selection+1)*gTextMalusHeight)).(*ebiten.Image), &options)
	}
}

func (g game) drawStateImprove(screen *ebiten.Image) {

	yStart := 256 + gCoinSideSize
	xSeparator := 70

	drawMoney(screen, gWidth/2, yStart-gCoinSideSize, g.money.money, true, 1)

	drawContinue(screen, (gWidth-gContinueWidth)/2, gHeight-gContinueHeight-gTitleMargin, g.improv.current == numImprove, g.improv.arrowBlinkFrame)

	drawShopText(screen, (gWidth-gTextMalusWidth)/2, gHeight-gContinueHeight-gTitleMargin-gTextMalusHeight-gTitleMargin, g.improv.current)

	x := (gWidth - (4*gImproveTextWidth + 3*xSeparator)) / 2
	y := yStart

	for i := 0; i < numImprove; i++ {

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(assets.ImageImprovements.SubImage(image.Rect(0, i*gImproveTextHeight, gImproveTextWidth, (i+1)*gImproveTextWidth)).(*ebiten.Image), &options)

		if i != numImprove {
			if len(g.improv.prices[i]) > g.improv.levels[i] {
				drawMoney(screen, x+3*gImproveTextWidth/5, y+gImproveTextHeight, g.improv.prices[i][g.improv.levels[i]], false, 0.4)
			} else {
				drawMaxed(screen, x+(gImproveTextWidth-gMaxWidth)/2, y+gImproveTextHeight-14)
			}
		}

		if g.improv.current == i {
			drawArrow(screen, x+(gImproveTextWidth-gArrowWidth)/2, y+gImproveTextHeight+40, 0, g.improv.arrowBlinkFrame)
		}

		x += gImproveTextWidth + xSeparator
	}

}

func (g *game) updateStateImprove() bool {

	g.improv.arrowBlinkFrame++
	if g.improv.arrowBlinkFrame >= numArrowBlinkFrame {
		g.improv.arrowBlinkFrame = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.audio.NextSounds[assets.SoundMenuMoveID] = true
		g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
			g.improv.current = (g.improv.current + numImprove) % (numImprove + 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.audio.NextSounds[assets.SoundMenuMoveID] = true
		g.improv.current = (g.improv.current + 1) % (numImprove + 1)
		for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
			g.improv.current = (g.improv.current + 1) % (numImprove + 1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.audio.NextSounds[assets.SoundMenuMoveID] = true
		if g.improv.current != numImprove {
			g.improv.current = numImprove
		} else {
			g.improv.current = 0
			for g.improv.current != numImprove && g.improv.levels[g.improv.current] >= len(g.improv.prices[g.improv.current]) {
				g.improv.current = (g.improv.current + 1) % (numImprove + 1)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.improv.current == numImprove {
			g.audio.NextSounds[assets.SoundMenuConfirmID] = true
			return true
		}

		if g.improv.levels[g.improv.current] < len(g.improv.prices[g.improv.current]) {
			if g.improv.prices[g.improv.current][g.improv.levels[g.improv.current]] <= g.money.money {
				g.money.money -= g.improv.prices[g.improv.current][g.improv.levels[g.improv.current]]
				g.improv.levels[g.improv.current]++
				g.audio.NextSounds[assets.SoundBuyID] = true
			} else {
				g.audio.NextSounds[assets.SoundMenuNoID] = true
			}
		} else {
			g.audio.NextSounds[assets.SoundMenuNoID] = true
		}
	}

	return false
}
