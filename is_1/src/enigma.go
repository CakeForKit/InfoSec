package main

import (
	"errors"
)

var (
	ErrLenPoses = errors.New("len(poses) != len(rotors)")
)

type Enigma interface {
	EncryptAlpha(alpha byte) byte
	EncryptText(text []byte) []byte
	SetRotorPositions(poses []byte) error
}

type enigma struct {
	switchingPanel Rotor
	rotors         []Rotor
	reflector      Reflector
}

func NewEnigma(switchingPanel Rotor, rotors []Rotor, reflector Reflector) Enigma {
	return &enigma{
		switchingPanel: switchingPanel,
		rotors:         rotors,
		reflector:      reflector,
	}
}

func (e *enigma) EncryptText(text []byte) []byte {
	resText := make([]byte, len(text))
	for i, v := range text {
		resText[i] = e.EncryptAlpha(v)
	}
	return resText
}

func (e *enigma) EncryptAlpha(alpha byte) byte {
	alpha = e.switchingPanel.SwitchTo(alpha)
	Nrotors := len(e.rotors)
	e.rotors[0].Rotate()
	nextA := e.rotors[0].Transform(alpha, 0)
	lastRing := e.rotors[0].GetRing()

	for i := 1; i < Nrotors; i++ {
		if e.rotors[i-1].GetRing() == e.rotors[i-1].GetSteppingPos() {
			e.rotors[i].Rotate()
		}
		nextA = e.rotors[i].Transform(nextA, lastRing)
		lastRing = e.rotors[i].GetRing()
	}

	nextA = e.reflector.Transform(nextA, lastRing, -1)

	lastRing = 0
	for i := len(e.rotors) - 1; i >= 0; i-- {
		nextA = e.rotors[i].TransformBack(nextA, lastRing)
		lastRing = e.rotors[i].GetRing()
	}
	nextA = byte((int(nextA) - int(lastRing) + alphabetSize) % alphabetSize)
	nextA = e.switchingPanel.SwitchFrom(nextA)

	return nextA
}

func (e *enigma) SetRotorPositions(poses []byte) error {
	if len(poses) != len(e.rotors) {
		return ErrLenPoses
	}
	for i := 0; i < len(e.rotors); i++ {
		e.rotors[i].SetRing(poses[i])
	}
	return nil
}
