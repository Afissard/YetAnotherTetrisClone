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

// styles for blocks
const (
	noStyle int = iota // it is important that noStyle is 0
	iBlockStyle
	oBlockStyle
	jBlockStyle
	lBlockStyle
	sBlockStyle
	tBlockStyle
	zBlockStyle
	breakStyle
)

func getIBlock() tetrisBlock {
	return tetrisBlock{
		id:    2,
		style: iBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{false, false, false, false},
				{true, true, true, true},
				{false, false, false, false}},
			{{false, true, false, false},
				{false, true, false, false},
				{false, true, false, false},
				{false, true, false, false}},
			{{false, false, false, false},
				{false, false, false, false},
				{true, true, true, true},
				{false, false, false, false}},
			{{false, true, false, false},
				{false, true, false, false},
				{false, true, false, false},
				{false, true, false, false}},
		},
	}
}

func getOBlock() tetrisBlock {
	return tetrisBlock{
		id:    3,
		style: oBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false}},
			{{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false}},
			{{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false}},
			{{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false}},
		},
	}
}

func getJBlock() tetrisBlock {
	return tetrisBlock{
		id:    1,
		style: jBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{true, true, true, false},
				{false, false, true, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{false, true, false, false},
				{true, true, false, false},
				{false, false, false, false}},
			{{true, false, false, false},
				{true, true, true, false},
				{false, false, false, false},
				{false, false, false, false}},
			{{false, true, true, false},
				{false, true, false, false},
				{false, true, false, false},
				{false, false, false, false}},
		},
	}
}

func getLBlock() tetrisBlock {
	return tetrisBlock{
		id:    0,
		style: lBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{true, true, true, false},
				{true, false, false, false},
				{false, false, false, false}},
			{{true, true, false, false},
				{false, true, false, false},
				{false, true, false, false},
				{false, false, false, false}},
			{{false, false, true, false},
				{true, true, true, false},
				{false, false, false, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{false, true, false, false},
				{false, true, true, false},
				{false, false, false, false}},
		},
	}
}

func getSBlock() tetrisBlock {
	return tetrisBlock{
		id:    5,
		style: sBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{false, true, true, false},
				{true, true, false, false},
				{false, false, false, false}},
			{{true, false, false, false},
				{true, true, false, false},
				{false, true, false, false},
				{false, false, false, false}},
			{{false, false, false, false},
				{false, true, true, false},
				{true, true, false, false},
				{false, false, false, false}},
			{{true, false, false, false},
				{true, true, false, false},
				{false, true, false, false},
				{false, false, false, false}},
		},
	}
}

func getTBlock() tetrisBlock {
	return tetrisBlock{
		id:    6,
		style: tBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{true, true, true, false},
				{false, true, false, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{true, true, false, false},
				{false, true, false, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{true, true, true, false},
				{false, false, false, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{false, true, true, false},
				{false, true, false, false},
				{false, false, false, false}},
		},
	}
}

func getZBlock() tetrisBlock {
	return tetrisBlock{
		id:    4,
		style: zBlockStyle,
		states: [4][4][4]bool{
			{{false, false, false, false},
				{true, true, false, false},
				{false, true, true, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{true, true, false, false},
				{true, false, false, false},
				{false, false, false, false}},
			{{false, false, false, false},
				{true, true, false, false},
				{false, true, true, false},
				{false, false, false, false}},
			{{false, true, false, false},
				{true, true, false, false},
				{true, false, false, false},
				{false, false, false, false}},
		},
	}
}
