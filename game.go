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

import "github.com/loig/ebitenginegamejam2024/assets"

const (
	stateTitle int = iota
	statePlay
	stateBalance
	stateLost
	stateImprove
	stateWon
	stateControls
	stateCredits
)

type game struct {
	state       int
	firstPlay   bool
	currentPlay tetris
	level       int
	goalLevel   int
	balance     balancing
	numChoices  int
	audio       assets.SoundManager
	money       moneyHandler
	improv      improvements
	fog         fog
	titleSelect int
	titleFrame  int
	winFrame    int
}

func (g *game) init() {
	g.audio = assets.InitAudio()
	g.state = stateControls
	g.firstPlay = true
	g.numChoices = 3
	g.improv = setupImprovements()
	g.goalLevel = 11
}
