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
	"math/rand"
)

const (
	leftChoice int = iota
	rigthChoice
	numChoices
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
	levels     [numBalances]int
	maxLevels  [numBalances]int
	choices    [numChoices]int
	numChoices int
}

func newBalance() balancing {

	b := balancing{}
	for i := 0; i < numChoices; i++ {
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

func (b *balancing) setChoice(choicePosition int) {
	b.levels[b.choices[choicePosition]]++
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
