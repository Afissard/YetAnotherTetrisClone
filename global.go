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

const (
	gWidth  int = gMultFactor * (gPlayAreaWidth + 2*gPlayAreaSide + gInfoLeftSide + gInfoWidth + gInfoRightSide)
	gHeight int = gMultFactor * gPlayAreaHeight

	gSquareSideSize int = 8 // basic block size in pixels

	gPlayAreaWidthInBlocks  int = 10
	gPlayAreaHeightInBlocks int = 18

	gPlayAreaWidth  int = gPlayAreaWidthInBlocks * gSquareSideSize  // width of play area in pixels
	gPlayAreaHeight int = gPlayAreaHeightInBlocks * gSquareSideSize // height of play area in pixels
	gPlayAreaSide   int = 10                                        // play area side shift in pixels

	gInfoLeftSide       int = 6                   // info left side shift in pixels
	gInfoWidth          int = 6 * gSquareSideSize // info width in pixels
	gInfoRightSide      int = 6                   // info right side shift in pixels
	gInfoBoxHeight      int = 3 * gSquareSideSize // height of standard info box in pixels
	gNextMargin         int = 4
	gNextBoxSide        int = 4*gSquareSideSize + gNextMargin*2 // side of box displaying next piece in pixels
	gInfoSmallBoxHeight int = 2 * gSquareSideSize               // height of small info box
	gInfoTop            int = 4                                 // info top shift in pixel
	gScoreToLevel       int = 24                                // distance from score bottom to level top
	gLevelToLines       int = 2                                 // distance from level bottom to lines top
	gLinesToNext        int = 6                                 // distance from lines bottom to next top

	gMultFactor int = 1 // multiplicative factor in case increasing size is needed

	gInvisibleLines int = 3 // number of hidden lines above the grid
)
