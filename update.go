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
)

func (g *game) Update() (err error) {

	score, sounds := g.currentPlay.update(
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.IsKeyPressed(ebiten.KeyRight),
		inpututil.IsKeyJustPressed(ebiten.KeySpace),
		inpututil.IsKeyJustPressed(ebiten.KeyEnter),
		g.level,
	)

	g.score += score
	g.audio.NextSounds = sounds

	return nil
}
