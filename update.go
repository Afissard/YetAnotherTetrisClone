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

	switch g.state {
	case stateTitle:
		if g.updateStateTitle() {
			g.state++
			g.balance = newBalance()
			g.currentPlay.init(g.level, g.balance)
		}
	case statePlay:
		if g.updateStatePlay() {
			g.state = stateTitle
			g.firstPlay = false
			g.level = 0
		}
		if g.currentPlay.numLines >= g.balance.getGoalLines() {
			g.state++
			g.level++
			g.balance.getChoice()
		}
	case stateBalance:
		if g.updateStateBalance() {
			g.state = statePlay
			g.currentPlay.init(g.level, g.balance)
		}
	}

	return nil
}

func (g *game) updateStateTitle() (end bool) {
	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	return
}

func (g *game) updateStatePlay() bool {
	dead, score, sounds := g.currentPlay.update(
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.IsKeyPressed(ebiten.KeyRight),
		inpututil.IsKeyJustPressed(ebiten.KeySpace),
		inpututil.IsKeyJustPressed(ebiten.KeyEnter),
		g.level,
	)

	g.score += score
	g.audio.NextSounds = sounds

	return dead
}

func (g *game) updateStateBalance() (end bool) {

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.choice = (g.choice + g.balance.numChoices - 1) % g.balance.numChoices
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.choice = (g.choice + 1) % g.balance.numChoices
	}

	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	if end {
		g.balance.setChoice(g.choice)
	}

	return
}
