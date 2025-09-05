package main

type Rotor interface {
	Rotate()
	Transform(alpha byte, nextRing byte) byte
	TransformBack(alpha byte, nextRing byte) byte
	GetRing() byte
	SetRing(ring byte)
	GetSteppingPos() byte
}

type rotor struct {
	permutation   []byte
	rePermutation []byte
	ring          byte
	steppingPos   byte
}

func NewRotor(permutation []byte, steppingPos byte, ring byte) Rotor {
	var rePermutation []byte
	for i, v := range permutation {
		rePermutation[v] = byte(i)
	}
	return &rotor{
		permutation:   permutation,
		steppingPos:   steppingPos,
		rePermutation: rePermutation,
		ring:          ring,
	}
}

func (r *rotor) Rotate() {
	r.ring = byte((int(r.ring) + 1) % alphabetSize)
}

func (r *rotor) Transform(alpha byte, prevRing byte) byte {
	// fmt.Printf("Transform: alpha=%c, prevRing=%c\n", alpha+'A', prevRing+'A')
	a := int(alpha)
	pr := int(prevRing)
	intputAlpha := (a + (int(r.ring) - pr + alphabetSize)) % alphabetSize
	return r.permutation[intputAlpha]
}

func (r *rotor) TransformBack(alpha byte, nextRing byte) byte {
	// fmt.Printf("TransformBack: alpha=%c, nextRing=%c\n", alpha+'A', nextRing+'A')
	a := int(alpha)
	pr := int(nextRing)
	intputAlpha := (a - (pr - int(r.ring)) + alphabetSize) % alphabetSize
	return r.rePermutation[intputAlpha]
}

func (r *rotor) GetRing() byte {
	return r.ring
}

func (r *rotor) SetRing(ring byte) {
	r.ring = ring
}

func (r *rotor) GetSteppingPos() byte {
	return r.steppingPos
}
