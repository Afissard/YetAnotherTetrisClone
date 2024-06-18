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
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

const (
	balanceGoalLines int = iota
	balanceSpeed
	balanceHiddenLines
	balanceDeathLines
	numBalances
)

const (
	maxLevelGoalLines   = 2
	maxLevelSpeed       = 4
	maxLevelHiddenLines = 9
	maxLevelDeathLines  = 8
)

type balancing struct {
	levels          [numBalances]int
	maxLevels       [numBalances]int
	choice          int
	choiceDirection int
	choices         []int
	numChoices      int
	inTransition    bool
	transitionFrame int
}

func (b *balancing) update() (end bool) {

	if b.inTransition {
		b.transitionFrame++
		if b.transitionFrame >= gChoiceSelectionNumFrame {
			b.inTransition = false
			b.transitionFrame = 0
			if b.choiceDirection < 0 {
				b.choice = (b.choice + 1) % b.numChoices
			} else {
				b.choice = (b.choice + b.numChoices - 1) % b.numChoices
			}
		}
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		b.choiceDirection = 1
		b.inTransition = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		b.choiceDirection = -1
		b.inTransition = true
	}

	end = inpututil.IsKeyJustPressed(ebiten.KeyEnter)

	if end {
		b.setChoice(b.choices[b.choice])
	}

	return
}

func (b balancing) draw(screen *ebiten.Image) {

	r := float64(gHeight / 5)
	cX, cY := gWidth/2, 2*gHeight/3

	angleShift := float64(b.choiceDirection) * float64(b.transitionFrame) / float64(gChoiceSelectionNumFrame) * (math.Pi * 2) / float64(b.numChoices)
	if !b.inTransition {
		angleShift = 0
	}

	currentX, currentY := math.Cos(math.Pi/2+angleShift)*r, -math.Sin(math.Pi/2+angleShift)*r
	currentX += float64(cX - gChoiceSize/2)
	currentY += float64(cY - gChoiceSize/2)

	// current choice
	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(currentX, currentY)
	if !b.inTransition {
		screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(numBalances*gChoiceSize, 0, (numBalances+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)
	}
	screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(b.choices[b.choice]*gChoiceSize, 0, (b.choices[b.choice]+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)

	// other choices
	for i := 0; i < b.numChoices-1; i++ {
		// find the choice to display
		displayNum := (b.choice + i + 1) % b.numChoices
		theChoice := b.choices[displayNum]

		//find the position to display it
		angle := float64(i+1)*(math.Pi*2)/float64(b.numChoices) + math.Pi/2 + angleShift
		x, y := math.Cos(angle)*r, -math.Sin(angle)*r
		x += float64(cX - gChoiceSize/2)
		y += float64(cY - gChoiceSize/2)

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(x, y)
		screen.DrawImage(assets.ImageMalus.SubImage(image.Rect(theChoice*gChoiceSize, 0, (theChoice+1)*gChoiceSize, gChoiceSize)).(*ebiten.Image), &options)
	}

}

func newBalance(numChoices int) balancing {

	b := balancing{}

	b.choices = make([]int, numChoices)
	for i := range b.choices {
		b.choices[i] = -1
	}

	b.maxLevels[balanceGoalLines] = maxLevelGoalLines
	b.maxLevels[balanceSpeed] = maxLevelSpeed
	b.maxLevels[balanceHiddenLines] = maxLevelHiddenLines
	b.maxLevels[balanceDeathLines] = maxLevelDeathLines
	return b
}

func (b *balancing) getChoice() {

	possibleChoices := make([]int, 0, 2*numBalances)

BalanceLoop:
	for c := 0; c < numBalances; c++ {
		if b.levels[c] < b.maxLevels[c] {
			possibleChoices = append(possibleChoices, c)
			for _, oldChoice := range b.choices {
				if c == oldChoice {
					continue BalanceLoop
				}
			}
			possibleChoices = append(possibleChoices, c)
		}
	}

	choice := 0
	for ; len(possibleChoices) > 0 && choice < len(b.choices); choice++ {
		take := rand.Intn(len(possibleChoices))
		b.choices[choice] = possibleChoices[take]

		possibleChoices[take], possibleChoices[len(possibleChoices)-1] = possibleChoices[len(possibleChoices)-1], possibleChoices[take]
		possibleChoices = possibleChoices[:len(possibleChoices)-1]

		if take-1 >= 0 && take-1 < len(possibleChoices) {
			if b.choices[choice] == possibleChoices[take-1] {
				possibleChoices[take-1], possibleChoices[len(possibleChoices)-1] = possibleChoices[len(possibleChoices)-1], possibleChoices[take-1]
				possibleChoices = possibleChoices[:len(possibleChoices)-1]
			}
		}

		if take+1 >= 0 && take+1 < len(possibleChoices) {
			if b.choices[choice] == possibleChoices[take+1] {
				possibleChoices[take+1], possibleChoices[len(possibleChoices)-1] = possibleChoices[len(possibleChoices)-1], possibleChoices[take+1]
				possibleChoices = possibleChoices[:len(possibleChoices)-1]
			}
		}
	}

	b.numChoices = choice

	for ; choice < len(b.choices); choice++ {
		b.choices[choice] = -1
	}

}

func (b *balancing) setChoice(choice int) {
	b.levels[choice]++
}

func (b balancing) getDeathLines() (numLines int) {
	const maxDeathLines int = gPlayAreaHeightInBlocks / 2

	numLines = b.levels[balanceDeathLines] + 1
	if numLines > maxDeathLines {
		numLines = maxDeathLines
	}
	return
}

func (b balancing) getHiddenLines() (numLines int) {
	const hiddenFactor int = 2

	numLines = hiddenFactor * b.levels[balanceHiddenLines]
	if numLines > gPlayAreaHeightInBlocks {
		numLines = gPlayAreaHeightInBlocks
	}

	return
}

func (b balancing) getGoalLines() int {
	var goalLines [maxLevelGoalLines + 1]int = [maxLevelGoalLines + 1]int{
		1, 15, 20,
	}

	if b.levels[balanceGoalLines] < len(goalLines) {
		return goalLines[b.levels[balanceGoalLines]]
	}
	return goalLines[len(goalLines)-1]
}

func (b balancing) getSpeed() int {
	var speedLevels [maxLevelSpeed + 1]int = [maxLevelSpeed + 1]int{
		45, 30, 20, 10, 5,
	}

	if b.levels[balanceSpeed] < len(speedLevels) {
		return speedLevels[b.levels[balanceSpeed]]
	}
	return speedLevels[len(speedLevels)-1]
}
