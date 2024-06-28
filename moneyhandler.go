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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/loig/ebitenginegamejam2024/assets"
)

const (
	scoreToMoney             int = 100
	scoreUnitPerFrame        int = 5
	scoreCountStep           int = 3
	scoreCoinAnimationFrames int = 30
)

type moneyHandler struct {
	displayMoney       int
	money              int
	score              int
	count              int
	nextCoin           int
	previousMoney      int
	scoreReduction     int
	coins              []coinAnimator
	firstAvailableCoin int
	numActive          int
}

type coinAnimator struct {
	frame          int
	startX, startY int
	goalX, goalY   int
	x, y           float64
	scale          float64
	active         bool
}

func newCoinAnimator(xFrom, yFrom, xTo, yTo int) coinAnimator {
	return coinAnimator{
		startX: xFrom, startY: yFrom,
		goalX: xTo, goalY: yTo,
		x: float64(xFrom), y: float64(yFrom),
		active: true,
	}
}

func (c *coinAnimator) update() bool {
	c.frame++
	if c.frame < scoreCoinAnimationFrames/3 {
		c.scale += 1 / float64(scoreCoinAnimationFrames/3)
		if c.scale > 1 {
			c.scale = 1
		}
	}

	/*
		else if c.frame >= 2*scoreCoinAnimationFrames/3 {
			c.scale -= 1 / float64(scoreCoinAnimationFrames/3)
			if c.scale < 0 {
				c.scale = 0
			}
		}
	*/

	c.x += float64(c.goalX-c.startX) / float64(scoreCoinAnimationFrames)
	c.y += float64(c.goalY-c.startY) / float64(scoreCoinAnimationFrames)

	if c.frame >= scoreCoinAnimationFrames {
		c.active = false
		return true
	}

	return false
}

func (c coinAnimator) draw(screen *ebiten.Image) {
	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(-float64(gCoinSideSize/2), -float64(gCoinSideSize/2))
	options.GeoM.Scale(c.scale, c.scale)
	options.GeoM.Translate(float64(c.x), float64(c.y))
	screen.DrawImage(assets.ImageCoin, &options)
}

func (m *moneyHandler) addScore(score int) {
	m.displayMoney = m.money
	m.previousMoney = m.money
	m.money += score / scoreToMoney
	m.score = score
	m.count = 0
	m.nextCoin = 0
	m.scoreReduction = scoreUnitPerFrame
	m.firstAvailableCoin = 0
	for i := range m.coins {
		m.coins[i].active = false
	}
	m.numActive = 0
}

func (m *moneyHandler) update() (finished bool, playSounds [assets.NumSounds]bool) {

	if m.score > 0 {
		if m.score < m.scoreReduction {
			m.scoreReduction = m.score
		}
		m.score -= m.scoreReduction
		m.count += m.scoreReduction
		m.nextCoin += m.scoreReduction

		if m.count/100 >= scoreCountStep {
			m.scoreReduction++
			m.count = 0
		}

		if m.nextCoin/100 > 0 {
			m.nextCoin -= 100
			m.numActive++
			theCoin := newCoinAnimator(gWidth-gXScoreFromRightSide+gMultFactor-gSquareSideSize/2, gYScoreFromTop+gSquareSideSize/2, gWidth/2, 3*gHeight/4)
			if m.firstAvailableCoin >= len(m.coins) {
				m.coins = append(m.coins, theCoin)
				m.firstAvailableCoin++
			} else {
				m.coins[m.firstAvailableCoin] = theCoin
				m.firstAvailableCoin++
				for m.firstAvailableCoin < len(m.coins) {
					if !m.coins[m.firstAvailableCoin].active {
						break
					}
					m.firstAvailableCoin++
				}
			}
		}
	}

	for i := range m.coins {
		if m.coins[i].active {
			if m.coins[i].update() {
				m.displayMoney++
				m.numActive--
				playSounds[assets.SoundCoinID] = true
				if i < m.firstAvailableCoin {
					m.firstAvailableCoin = i
				}
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if m.score <= 0 {
			playSounds[assets.SoundMenuConfirmID] = m.numActive <= 0
			return m.numActive <= 0, playSounds
		}
		m.nextCoin += m.score
		m.score = 0
		m.displayMoney += m.nextCoin / 100
		m.nextCoin = 0
	}

	return false, playSounds
}

func (m moneyHandler) draw(screen *ebiten.Image) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(gWidth-gYouLoseWidth)/2, float64(gTitleMargin))
	screen.DrawImage(assets.ImageYouLose, &options)

	options = ebiten.DrawImageOptions{}
	options.GeoM.Translate(0, float64(3*gHeight/4-gMoneyBackHeight/2))
	screen.DrawImage(assets.ImageMoneyBack, &options)

	drawMoney(screen, gWidth/2, 3*gHeight/4, m.displayMoney, true, 1)

	for _, c := range m.coins {
		if c.active {
			c.draw(screen)
		}
	}

}

func drawMoney(screen *ebiten.Image, x, y int, money int, symbolFirst bool, scaling float64) {
	options := ebiten.DrawImageOptions{}

	num := money
	displaySize := 0
	for num/10 > 0 {
		displaySize++
		num = num / 10
	}
	if displaySize == 0 {
		displaySize = 1
	}
	displaySize++

	options.GeoM.Scale(scaling, scaling)
	options.GeoM.Translate(float64(x)+float64(displaySize*gCoinSideSize)*scaling/2, float64(y)-float64(gCoinSideSize)*scaling/2)

	if !symbolFirst {
		options.GeoM.Translate(float64(-gCoinSideSize)*scaling, float64(0))
		screen.DrawImage(assets.ImageCoin, &options)
	}

	atLeastOnce := true
	num = money
	for num > 0 || atLeastOnce {
		atLeastOnce = false
		digit := num % 10
		num = num / 10

		options.GeoM.Translate(float64(-gCoinSideSize)*scaling, float64(0))
		screen.DrawImage(assets.ImageBigdigits.SubImage(image.Rect(digit*gCoinSideSize, 0, (digit+1)*gCoinSideSize, gCoinSideSize)).(*ebiten.Image), &options)
	}

	if symbolFirst {
		options.GeoM.Translate(float64(-gCoinSideSize)*scaling, float64(0))
		screen.DrawImage(assets.ImageCoin, &options)
	}
}
