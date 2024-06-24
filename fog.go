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

const (
	fogFramesPerLine  int = 60
	fogHoldFrames     int = 20
	fogDecreaseFactor int = 4
)

type fog struct {
	hiddenLines        int
	currentHiddenLines int
	protectionLevel    int
	frame              int
	decreasing         bool
}

func (f *fog) reset(hiddenLines int, protectionLevel int) {
	f.hiddenLines = hiddenLines
	f.currentHiddenLines = hiddenLines
	f.protectionLevel = protectionLevel
	f.frame = 0
	f.decreasing = false
}

func (f *fog) update() {
	if f.protectionLevel > 0 {
		f.frame++
		if !f.decreasing && f.currentHiddenLines >= f.hiddenLines {
			if f.frame >= fogHoldFrames {
				f.frame = 0
				f.decreasing = true
			}
			return
		}
		if f.frame >= fogFramesPerLine {
			f.frame = 0
			if f.decreasing {
				f.currentHiddenLines--
				if f.hiddenLines-fogDecreaseFactor*f.protectionLevel >= f.currentHiddenLines {
					f.decreasing = false
				}
			} else {
				f.currentHiddenLines++
			}
		}
	}
}

func (f fog) draw(screen *ebiten.Image) {

	y := float64(gPlayAreaHeight - f.currentHiddenLines*gSquareSideSize)

	if (f.decreasing && f.currentHiddenLines > 0) || (!f.decreasing && f.currentHiddenLines < f.hiddenLines) {
		yDec := (float64(f.frame) / float64(fogFramesPerLine)) * float64(gSquareSideSize)
		if f.decreasing {
			y += yDec
		} else {
			y -= yDec
		}
	}

	if y > 0 {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(gPlayAreaSide), y)
		screen.DrawImage(assets.ImageFog, &options)
	}
}
