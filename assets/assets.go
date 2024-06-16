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
package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed squares.png
var imageSquaresBytes []byte
var ImageSquares *ebiten.Image

//go:embed back.png
var imageBackBytes []byte
var ImageBack *ebiten.Image

//go:embed digits.png
var imageDigitsBytes []byte
var ImageDigits *ebiten.Image

func Load() {
	var err error

	imageDecoded, _, err := image.Decode(bytes.NewReader(imageSquaresBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageSquares = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageBackBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageBack = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageDigitsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageDigits = ebiten.NewImageFromImage(imageDecoded)
}
