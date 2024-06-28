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

	gSquareSideSize int = 8 * gMultFactor // basic block size in pixels

	gPlayAreaWidthInBlocks  int = 10
	gPlayAreaHeightInBlocks int = 18

	gPlayAreaWidth  int = gPlayAreaWidthInBlocks * gSquareSideSize  // width of play area in pixels
	gPlayAreaHeight int = gPlayAreaHeightInBlocks * gSquareSideSize // height of play area in pixels
	gPlayAreaSide   int = 9 * gMultFactor                           // play area side shift in pixels

	gInfoLeftSide       int = 8 * gMultFactor  // info left side shift in pixels
	gInfoWidth          int = 46 * gMultFactor // info width in pixels
	gInfoRightSide      int = 8 * gMultFactor  // info right side shift in pixels
	gInfoBoxHeight      int = 22 * gMultFactor // height of standard info box in pixels
	gNextMargin         int = 5 * gMultFactor
	gNextBoxSide        int = 42 * gMultFactor // side of box displaying next piece in pixels
	gInfoSmallBoxHeight int = 14 * gMultFactor // height of small info box
	gInfoTop            int = 6 * gMultFactor  // info top shift in pixel
	gScoreToLevel       int = 25 * gMultFactor // distance from score bottom to level top
	gLevelToLines       int = 2 * gMultFactor  // distance from level bottom to lines top
	gLinesToNext        int = 8 * gMultFactor  // distance from lines bottom to next top
	gInfoShiftNext      int = 7 * gMultFactor  // distance from info side to next

	gXLinesFromRightSide int = 20 * gMultFactor     // distance from right of screen to right of lines
	gYLinesFromTop       int = 80 * gMultFactor     // distance from top of screen to top of lines
	gXScoreFromRightSide int = 12 * gMultFactor     // distance from right of screen to right of score
	gYScoreFromTop       int = 25 * gMultFactor     // distance from top of screen to top of score
	gXLevelFromRightSide int = gXLinesFromRightSide // distance from right of screen to right of level
	gYLevelFromTop       int = 56 * gMultFactor     // distance from top of screen to top of score

	gInvisibleLines int = 3 // number of hidden lines above the grid

	gMultFactor int = 8 // multiply the size of old graphics

	gChoiceSize      int = 300 // size in pixels of the side of a balancing choice
	gChoiceLevelSize int = 70  // size in pixels of the icone giving the level of a balancing choice

	gChoiceSelectionNumFrame int = 30 // number of frames for changing balancing choice

	gSpeedLevels int = 21

	gInvisibleNumFrames int = 60 // num frames for one step of invisibility

	gCoinSideSize int = 128 // size of the side of the coin image in pixels

	gImproveTextWidth  int = 218 // width of text for improvements in pixels
	gImproveTextHeight int = 218 // height of text for improvements in pixels

	gDangerSide int = 218 // width/height of the danger zone symbol in pixels

	gLevelCompleteWidth int = 762 // width of the title of the end of level screen in pixels
	gYouLoseWidth       int = 385 // width of the title of the lose screen in pixels
	gShopTitleWidth     int = 662 // width of the title of the shop screen in pixels

	gTitleMargin int = 20 // margin on top of the end of level/lose/shop screens in pixels

	gMaxWidth int = 199 // width of "maxed" in pixels

	gArrowWidth  int = 60 // width of the improvement selection arrow in pixels
	gArrowHeight int = 51 // height of the improvement selection arrow in pixels

	gContinueHeight int = 32  // height of the continue button in the shop
	gContinueWidth  int = 323 // width of the continue button in the shop

	gMoneyBackHeight int = 195 // height in pixels of the background for displaying money

	gHoldSide int = 181 // size in pixel of the side of the hold block

)

var gSpeeds [gSpeedLevels]int = [gSpeedLevels]int{
	53, 49, 45, 41, 37, 33, 28, 22, 17, 11, 10,
	9, 8, 7, 6, 6, 5, 5, 4, 4, 3,
}
