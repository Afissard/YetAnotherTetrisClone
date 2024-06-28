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

	// play sounds
	g.audio.PlaySounds()
	g.audio.NextSounds = [assets.NumSounds]bool{}

	if g.state != stateControls && g.state != stateWon {
		g.audio.UpdateMusic(0.7)
	}

	betterRotation := g.improv.levels[improveResetAutoDown] > 0
	canHold := g.improv.levels[improveHold] > 0
	life := g.improv.levels[improveLife]*2 - 1

	switch g.state {
	case stateControls:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.audio.NextSounds[assets.SoundMenuConfirmID] = true
			g.state = stateTitle
			g.titleFrame = 0
		}
	case stateCredits:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.audio.NextSounds[assets.SoundMenuConfirmID] = true
			g.state = stateTitle
			g.titleFrame = 0
		}
	case stateTitle:
		g.titleFrame++
		if g.titleFrame >= numArrowBlinkFrame {
			g.titleFrame = 0
		}
		if g.updateStateTitle() {
			if g.titleSelect == 0 {
				g.firstPlay = false
				g.state = statePlay
				g.balance = newBalance(g.numChoices)
				g.currentPlay.init(g.level, g.balance, g.level, 0, betterRotation, canHold, life, life)
				g.fog.reset(g.balance.getHiddenLines(), g.improv.levels[improveHideMove])
			} else {
				g.state = stateCredits
			}
		}
	case statePlay:
		if g.updateStatePlay() {
			g.state = stateLost
			g.money.addScore(g.currentPlay.score)
		}
		if !g.currentPlay.inAnimation && g.currentPlay.numLines >= g.balance.getGoalLines() {
			if g.level+1 >= g.goalLevel {
				g.state = stateWon
				g.audio.NextSounds[assets.SoundBuyID] = true
				g.audio.StopMusic()
				return nil
			}
			g.state = stateBalance
			g.balance.getChoice()
		}
	case stateBalance:
		finished, playSounds := g.balance.update()
		g.audio.NextSounds = playSounds
		if finished {
			g.state = statePlay
			g.level++
			g.currentPlay.init(g.level, g.balance, g.level, g.currentPlay.score, betterRotation, canHold, life, g.currentPlay.currentLife)
			g.fog.reset(g.balance.getHiddenLines(), g.improv.levels[improveHideMove])
		}
	case stateLost:
		finished, playSounds := g.money.update()
		g.audio.NextSounds = playSounds
		if finished {
			g.state = stateImprove
			g.level = 0
			g.improv.reset()
		}
		g.currentPlay.score = g.money.score
	case stateImprove:
		if g.updateStateImprove() {
			g.state = stateTitle
			g.titleFrame = 0
		}
	case stateWon:
		g.winFrame++
		if g.winFrame >= len(gAnimRocket) {
			g.winFrame = 0
			g.audio.NextSounds[assets.SoundTouchGroundID] = true
		}
		if g.winFrame == 16 {
			g.audio.NextSounds[assets.SoundRocketID] = true
		}
	}

	return nil
}

func (g *game) updateStateTitle() (end bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.audio.NextSounds[assets.SoundMenuMoveID] = true
		g.titleSelect = (g.titleSelect + 1) % 2
	}

	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	g.audio.NextSounds[assets.SoundMenuConfirmID] = end
	return
}

func (g *game) updateStatePlay() bool {
	sounds := g.currentPlay.update(
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.IsKeyPressed(ebiten.KeyRight),
		inpututil.IsKeyJustPressed(ebiten.KeyUp),
		inpututil.IsKeyJustPressed(ebiten.KeyAlt),
		inpututil.IsKeyJustPressed(ebiten.KeySpace),
		g.level,
	)

	g.audio.NextSounds = sounds

	g.fog.update()

	return g.currentPlay.dead && !g.currentPlay.inAnimation
}
