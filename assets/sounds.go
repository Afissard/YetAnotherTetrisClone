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
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed rotation.wav
var soundRotationBytes []byte
var soundRotation []byte

//go:embed leftright.wav
var soundLeftRightBytes []byte
var soundLeftRight []byte

//go:embed touchground.wav
var soundTouchGroundBytes []byte
var soundTouchGround []byte

//go:embed linesvanishing.wav
var soundLinesVanishingBytes []byte
var soundLinesVanishing []byte

//go:embed linesfalling.wav
var soundLinesFallingBytes []byte
var soundLinesFalling []byte

const (
	SoundRotationID int = iota
	SoundLeftRightID
	SoundTouchGroundID
	SoundLinesVanishingID
	SoundLinesFallingID
	NumSounds
)

type SoundManager struct {
	audioContext *audio.Context
	NextSounds   [NumSounds]bool
}

// play requested sounds
func (s SoundManager) PlaySounds() {
	for sound, play := range s.NextSounds {
		if play {
			s.playSound(sound)
		}
	}
}

// play a sound
func (s SoundManager) playSound(sound int) {
	var soundBytes []byte
	switch sound {
	case SoundRotationID:
		soundBytes = soundRotation
	case SoundLeftRightID:
		soundBytes = soundLeftRight
	case SoundTouchGroundID:
		soundBytes = soundTouchGround
	case SoundLinesVanishingID:
		soundBytes = soundLinesVanishing
	case SoundLinesFallingID:
		soundBytes = soundLinesFalling
	}

	if len(soundBytes) > 0 {
		soundPlayer := s.audioContext.NewPlayerFromBytes(soundBytes)
		soundPlayer.Play()
	}
}

// decode music and sounds
func InitAudio() (manager SoundManager) {

	var error error
	manager.audioContext = audio.NewContext(44100)

	// sounds
	sound, error := wav.DecodeWithSampleRate(manager.audioContext.SampleRate(), bytes.NewReader(soundRotationBytes))
	if error != nil {
		log.Panic("Audio problem:", error)
	}
	soundRotation, error = io.ReadAll(sound)
	if error != nil {
		log.Panic("Audio problem:", error)
	}

	sound, error = wav.DecodeWithSampleRate(manager.audioContext.SampleRate(), bytes.NewReader(soundLeftRightBytes))
	if error != nil {
		log.Panic("Audio problem:", error)
	}
	soundLeftRight, error = io.ReadAll(sound)
	if error != nil {
		log.Panic("Audio problem:", error)
	}

	sound, error = wav.DecodeWithSampleRate(manager.audioContext.SampleRate(), bytes.NewReader(soundTouchGroundBytes))
	if error != nil {
		log.Panic("Audio problem:", error)
	}
	soundTouchGround, error = io.ReadAll(sound)
	if error != nil {
		log.Panic("Audio problem:", error)
	}

	sound, error = wav.DecodeWithSampleRate(manager.audioContext.SampleRate(), bytes.NewReader(soundLinesVanishingBytes))
	if error != nil {
		log.Panic("Audio problem:", error)
	}
	soundLinesVanishing, error = io.ReadAll(sound)
	if error != nil {
		log.Panic("Audio problem:", error)
	}

	sound, error = wav.DecodeWithSampleRate(manager.audioContext.SampleRate(), bytes.NewReader(soundLinesFallingBytes))
	if error != nil {
		log.Panic("Audio problem:", error)
	}
	soundLinesFalling, error = io.ReadAll(sound)
	if error != nil {
		log.Panic("Audio problem:", error)
	}

	return
}
