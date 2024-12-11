package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	defaultKeys = keyboardMap{
		enter: ebiten.KeyEnter,
		alt:   ebiten.KeyAlt,
		space: ebiten.KeySpace,
		up:    ebiten.KeyUp,
		down:  ebiten.KeyDown,
		left:  ebiten.KeyLeft,
		right: ebiten.KeyRight,
	}

	// equivalent of zqsd for azerty keyboard
	wasdKeys = keyboardMap{
		enter: ebiten.KeyEnter,
		alt:   ebiten.KeyAlt,
		space: ebiten.KeySpace,
		up:    ebiten.KeyZ,
		down:  ebiten.KeyS,
		left:  ebiten.KeyA,
		right: ebiten.KeyD,
	}
)

type keyboardMap struct {
	enter ebiten.Key
	alt   ebiten.Key
	space ebiten.Key
	up    ebiten.Key
	down  ebiten.Key
	left  ebiten.Key
	right ebiten.Key
}

type KeyboardInputs struct {
	kmap  keyboardMap
	enter bool
	alt   bool
	space bool
	up    bool
	down  bool
	left  bool
	right bool
}

func (k *KeyboardInputs) update() {
	k.enter = inpututil.IsKeyJustPressed(k.kmap.enter)
	k.alt = inpututil.IsKeyJustPressed(k.kmap.alt)
	k.space = inpututil.IsKeyJustPressed(k.kmap.space)
	k.up = inpututil.IsKeyJustPressed(k.kmap.up)
	k.down = ebiten.IsKeyPressed(k.kmap.down)
	k.left = ebiten.IsKeyPressed(k.kmap.left)
	k.right = ebiten.IsKeyPressed(k.kmap.right)
}
