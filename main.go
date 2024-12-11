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
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/loig/ebitenginegamejam2024/assets"
)

func init() {
	flag.IntVar(&selectedKeyBind, "k", 0, "Select the keybind you want to use:\n- 1 for wasd\n- 0 or nothing for default")
	flag.Parse()
}

func main() {

	g := game{}
	g.init()

	ebiten.SetWindowTitle("Yet Another Tetris Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//ebiten.SetWindowSize(640, 576)

	assets.Load(gMultFactor)

	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}

}
