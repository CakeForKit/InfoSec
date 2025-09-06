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
	rotors    []Rotor
	reflector Reflector
}

func NewEnigma(rotors []Rotor, reflector Reflector) Enigma {
	return &enigma{
		rotors:    rotors,
		reflector: reflector,
	}
}

func (e *enigma) EncryptText(text []byte) []byte {
	resText := make([]byte, len(text))
	// var strBuilder strings.Builder
	// strBuilder.Grow(len(text))
	for i, v := range text {
		resText[i] = e.EncryptAlpha(v)
	}
	return resText
}

func (e *enigma) EncryptAlpha(alpha byte) byte {
	Nrotors := len(e.rotors)
	e.rotors[0].Rotate()
	nextA := e.rotors[0].Transform(alpha, 0)
	lastRing := e.rotors[0].GetRing()

	// fmt.Printf("%c -> %c\n", alpha+'A', nextA+'A')
	// var tmp byte
	for i := 1; i < Nrotors; i++ {
		if e.rotors[i-1].GetRing() == e.rotors[i-1].GetSteppingPos() {
			e.rotors[i].Rotate()
		}
		// tmp = nextA
		nextA = e.rotors[i].Transform(nextA, lastRing)
		lastRing = e.rotors[i].GetRing()
		// fmt.Printf("%c -> %c\n", tmp+'A', nextA+'A')
	}

	// tmp = nextA
	nextA = e.reflector.Transform(nextA, lastRing, -1)
	// fmt.Printf("ref: %c -> %c\n", tmp+'A', nextA+'A')

	lastRing = 0
	for i := len(e.rotors) - 1; i >= 0; i-- {
		// tmp = nextA
		nextA = e.rotors[i].TransformBack(nextA, lastRing)
		lastRing = e.rotors[i].GetRing()
		// fmt.Printf("%c -> %c\n", tmp+'A', nextA+'A')
	}
	// tmp = nextA
	nextA = byte((int(nextA) - int(lastRing) + alphabetSize) % alphabetSize)
	// fmt.Printf("%c -> %c\n", tmp+'A', nextA+'A')
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
