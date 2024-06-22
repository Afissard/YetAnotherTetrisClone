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

//go:embed malus.png
var imageMalusBytes []byte
var ImageMalus *ebiten.Image

//go:embed level.png
var imageLevelBytes []byte
var ImageLevel *ebiten.Image

//go:embed coin.png
var imageCoinBytes []byte
var ImageCoin *ebiten.Image

//go:embed bigdigits.png
var imageBigdigitsBytes []byte
var ImageBigdigits *ebiten.Image

func Load(mult int) {
	var err error

	imageDecoded, _, err := image.Decode(bytes.NewReader(imageSquaresBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageSquares = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageSquares = resize(ImageSquares, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageBackBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageBack = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageBack = resize(ImageBack, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageDigitsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageDigits = ebiten.NewImageFromImage(imageDecoded)
	// resize
	ImageDigits = resize(ImageDigits, mult)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageMalusBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageMalus = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageLevelBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageLevel = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageCoinBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageCoin = ebiten.NewImageFromImage(imageDecoded)

	imageDecoded, _, err = image.Decode(bytes.NewReader(imageBigdigitsBytes))
	if err != nil {
		log.Fatal(err)
	}
	ImageBigdigits = ebiten.NewImageFromImage(imageDecoded)
}

func resize(img *ebiten.Image, mult int) (res *ebiten.Image) {

	res = ebiten.NewImage(img.Bounds().Dx()*mult, img.Bounds().Dy()*mult)

	options := ebiten.DrawImageOptions{}
	options.GeoM.Scale(float64(mult), float64(mult))

	res.DrawImage(img, &options)

	return

}
