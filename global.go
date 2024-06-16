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
	gWidth  int = gPlayAreaWidth + 2*gPlayAreaSide + gInfoLeftSide + gInfoWidth + gInfoRightSide
	gHeight int = gPlayAreaHeight

	gSquareSideSize int = 8 // basic block size in pixels

	gPlayAreaWidthInBlocks  int = 10
	gPlayAreaHeightInBlocks int = 18

	gPlayAreaWidth  int = gPlayAreaWidthInBlocks * gSquareSideSize  // width of play area in pixels
	gPlayAreaHeight int = gPlayAreaHeightInBlocks * gSquareSideSize // height of play area in pixels
	gPlayAreaSide   int = 9                                         // play area side shift in pixels

	gInfoLeftSide       int = 8  // info left side shift in pixels
	gInfoWidth          int = 46 // info width in pixels
	gInfoRightSide      int = 8  // info right side shift in pixels
	gInfoBoxHeight      int = 22 // height of standard info box in pixels
	gNextMargin         int = 5
	gNextBoxSide        int = 42 // side of box displaying next piece in pixels
	gInfoSmallBoxHeight int = 14 // height of small info box
	gInfoTop            int = 6  // info top shift in pixel
	gScoreToLevel       int = 25 // distance from score bottom to level top
	gLevelToLines       int = 2  // distance from level bottom to lines top
	gLinesToNext        int = 8  // distance from lines bottom to next top
	gInfoShiftNext      int = 7  // distance from info side to next

	gInvisibleLines int = 3 // number of hidden lines above the grid
)
