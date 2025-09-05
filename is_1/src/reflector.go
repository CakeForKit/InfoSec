package main

type Reflector interface {
	Transform(alpha byte, nextRing byte, dir int) byte
}

type reflector struct {
	permutation []byte
}

func NewReflector(permutation []byte) Reflector {
	return &reflector{
		permutation: permutation,
	}
}

func (r *reflector) Transform(alpha byte, nextRing byte, dir int) byte {
	a := int(alpha)
	ring := int(nextRing)
	input := (a + dir*ring + alphabetSize) % alphabetSize
	return r.permutation[input]
}
