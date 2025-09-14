package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const alphabetSize = 256

func GenerateRotor(alphabetSize int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	permutation := r.Perm(alphabetSize)

	result := make([]byte, alphabetSize)
	for i, v := range permutation {
		result[i] = byte(v)
	}

	return result
}

func GenerateReflector(alphabetSize int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var half int = alphabetSize / 2
	firstPerm := r.Perm(half)
	secPerm := r.Perm(half)

	reflector := make([]byte, alphabetSize)
	for i := range alphabetSize / 2 {
		reflector[firstPerm[i]] = byte(half + secPerm[i])
		reflector[half+secPerm[i]] = byte(firstPerm[i])
	}
	return reflector
}

func GenerateReflector_test() {
	ref := GenerateReflector(alphabetSize)
	for _, v := range ref {
		fmt.Printf("%d, ", v)
	}
}

func GenerateRotor_test() {
	ref := GenerateRotor(alphabetSize)
	for _, v := range ref {
		fmt.Printf("%d, ", v)
	}
}

// func main() {
// 	GenerateRotor_test()
// }

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: enigma input.txt output.txt")
		os.Exit(1)
	}
	inputFileName := os.Args[1]  // "../data/input.txt"
	outputFIlename := os.Args[2] // "../data/output.txt"

	inputData, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switchingPanel := NewRotor(TypeRotor256_1, '0', 0)
	rotors := []Rotor{
		NewRotor(TypeRotor256_1, '1', 0),
		NewRotor(TypeRotor256_2, '-', 0),
		NewRotor(TypeRotor256_3, ' ', 0),
	}
	reflector := NewReflector(Reflector256_2)
	enigm := NewEnigma(switchingPanel, rotors, reflector)

	posRing := []byte{'Q', '8', '8'}
	enigm.SetRotorPositions(posRing)
	encryptedText := enigm.EncryptText(inputData)
	err = os.WriteFile(outputFIlename, encryptedText, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
