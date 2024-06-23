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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

func (g *game) Update() (err error) {

	g.audio.NextSounds = [assets.NumSounds]bool{}

	betterRotation := g.improv.levels[improveResetAutoDown] > 0
	canHold := g.improv.levels[improveHold] > 0
	life := g.improv.levels[improveLife]

	switch g.state {
	case stateTitle:
		if g.updateStateTitle() {
			g.state = statePlay
			g.balance = newBalance(g.numChoices)
			g.currentPlay.init(g.level, g.balance, g.level, 0, betterRotation, canHold, life, life)
		}
	case statePlay:
		if g.updateStatePlay() {
			g.state = stateLost
			g.money.addScore(g.currentPlay.score)
		}
		if g.currentPlay.numLines >= g.balance.getGoalLines() {
			g.state = stateBalance
			g.level++
			g.balance.getChoice()
		}
	case stateBalance:
		if g.balance.update() {
			g.state = statePlay
			g.currentPlay.init(g.level, g.balance, g.level, g.currentPlay.score, betterRotation, canHold, life, g.currentPlay.currentLife)
		}
	case stateLost:
		if g.money.update() {
			g.state = stateImprove
			g.firstPlay = false
			g.level = 0
		}
		g.currentPlay.score = g.money.score
	case stateImprove:
		if g.updateStateImprove() {
			g.state = stateTitle
		}
	}

	return nil
}

func (g *game) updateStateTitle() (end bool) {
	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	return
}

func (g *game) updateStatePlay() bool {
	dead, sounds := g.currentPlay.update(
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.IsKeyPressed(ebiten.KeyRight),
		inpututil.IsKeyJustPressed(ebiten.KeyUp),
		inpututil.IsKeyJustPressed(ebiten.KeySpace),
		inpututil.IsKeyJustPressed(ebiten.KeyEnter),
		g.level,
	)

	g.audio.NextSounds = sounds

	return dead
}
