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

var (
	TypeRotor256_1 []byte            = []byte{250, 15, 200, 112, 123, 32, 183, 50, 193, 146, 218, 114, 90, 234, 139, 71, 43, 6, 173, 79, 0, 142, 141, 5, 184, 109, 154, 91, 140, 121, 241, 195, 203, 214, 81, 38, 144, 49, 34, 204, 151, 7, 164, 240, 155, 163, 110, 207, 67, 20, 166, 245, 63, 101, 9, 199, 103, 74, 80, 157, 153, 128, 181, 45, 33, 3, 73, 239, 119, 242, 130, 220, 134, 254, 191, 31, 102, 187, 231, 127, 2, 150, 39, 13, 201, 178, 251, 177, 76, 131, 113, 215, 180, 16, 243, 66, 95, 30, 159, 223, 137, 253, 160, 108, 209, 172, 225, 53, 216, 78, 210, 229, 58, 198, 222, 162, 104, 255, 196, 89, 94, 87, 72, 26, 70, 136, 85, 235, 208, 40, 93, 125, 27, 44, 84, 4, 236, 106, 8, 42, 244, 156, 116, 10, 23, 111, 29, 12, 249, 206, 213, 194, 97, 54, 75, 211, 227, 169, 129, 147, 230, 202, 36, 41, 83, 161, 35, 132, 185, 99, 143, 11, 86, 174, 122, 68, 28, 14, 148, 168, 21, 96, 175, 1, 224, 126, 170, 17, 135, 189, 57, 46, 61, 190, 92, 19, 226, 124, 212, 217, 237, 167, 77, 138, 105, 55, 149, 115, 232, 228, 117, 186, 179, 176, 64, 65, 192, 145, 252, 182, 56, 37, 233, 205, 52, 120, 47, 22, 82, 24, 118, 219, 248, 158, 69, 221, 197, 246, 188, 152, 48, 25, 59, 133, 107, 238, 62, 60, 98, 88, 51, 100, 18, 171, 165, 247}
	TypeRotor256_2 []byte            = []byte{130, 154, 218, 69, 97, 33, 64, 225, 84, 51, 0, 136, 180, 157, 96, 103, 216, 49, 118, 36, 83, 166, 70, 59, 249, 108, 27, 91, 172, 119, 135, 156, 200, 117, 31, 147, 232, 11, 132, 214, 221, 24, 129, 34, 215, 22, 40, 46, 77, 78, 139, 245, 223, 179, 63, 23, 228, 184, 251, 5, 111, 144, 158, 244, 124, 163, 19, 162, 68, 105, 35, 178, 89, 174, 234, 233, 54, 227, 61, 188, 204, 80, 241, 86, 176, 193, 107, 39, 92, 109, 246, 113, 146, 7, 79, 4, 100, 62, 26, 151, 110, 206, 104, 247, 32, 242, 250, 201, 58, 171, 65, 99, 127, 67, 169, 222, 142, 253, 52, 131, 12, 29, 168, 211, 175, 195, 194, 55, 114, 43, 56, 240, 85, 93, 224, 181, 187, 177, 121, 230, 90, 150, 42, 145, 10, 101, 198, 20, 138, 192, 254, 38, 45, 14, 53, 37, 66, 76, 1, 235, 239, 141, 149, 3, 21, 237, 102, 72, 75, 190, 208, 167, 137, 15, 226, 120, 252, 16, 196, 209, 44, 95, 71, 191, 220, 219, 8, 207, 112, 47, 170, 13, 159, 126, 229, 134, 57, 87, 74, 238, 30, 213, 164, 197, 243, 248, 122, 133, 161, 9, 152, 6, 88, 199, 50, 82, 125, 217, 128, 48, 18, 28, 210, 173, 255, 153, 203, 98, 165, 236, 160, 140, 205, 202, 155, 116, 73, 186, 123, 81, 231, 25, 106, 60, 183, 2, 17, 212, 94, 185, 143, 182, 189, 148, 41, 115}
	TypeRotor1     [alphabetSize]int = [alphabetSize]int{4, 10, 12, 5, 11, 6, 3, 16, 21, 25, 13, 19, 14, 22, 24, 7, 23, 20, 18, 15, 0, 8, 1, 17, 2, 9}
	TypeRotor2     [alphabetSize]int = [alphabetSize]int{0, 9, 3, 10, 18, 8, 17, 20, 23, 1, 11, 7, 22, 19, 12, 2, 16, 6, 25, 13, 15, 24, 5, 21, 14, 4}
	TypeRotor3     [alphabetSize]int = [alphabetSize]int{1, 3, 5, 7, 9, 11, 2, 15, 17, 19, 23, 21, 25, 13, 24, 4, 8, 22, 6, 0, 10, 12, 20, 18, 16, 14}
	TypeRotor4     [alphabetSize]int = [alphabetSize]int{4, 18, 14, 21, 15, 25, 9, 0, 24, 16, 20, 8, 17, 7, 23, 11, 13, 5, 19, 6, 10, 3, 2, 12, 22, 1}
	TypeRotor5     [alphabetSize]int = [alphabetSize]int{21, 25, 1, 17, 6, 8, 19, 24, 20, 15, 18, 3, 13, 7, 11, 23, 0, 22, 12, 9, 16, 14, 5, 4, 2, 10}
	TypeRotor6     [alphabetSize]int = [alphabetSize]int{9, 15, 6, 21, 14, 20, 12, 5, 24, 16, 1, 4, 13, 7, 25, 17, 3, 10, 0, 18, 23, 11, 8, 2, 19, 22}
	TypeRotor7     [alphabetSize]int = [alphabetSize]int{13, 25, 9, 7, 6, 17, 2, 23, 12, 24, 18, 22, 1, 14, 20, 5, 0, 8, 21, 11, 15, 4, 10, 16, 3, 19}
	TypeRotor8     [alphabetSize]int = [alphabetSize]int{5, 10, 16, 7, 19, 11, 23, 14, 2, 1, 9, 15, 3, 25, 17, 0, 12, 4, 18, 22, 13, 8, 20, 24, 6, 21}
	TypeBetaRotor  [alphabetSize]int = [alphabetSize]int{11, 4, 24, 9, 21, 2, 13, 8, 23, 22, 15, 1, 16, 12, 3, 17, 19, 0, 10, 25, 6, 5, 20, 7, 14, 18}
	TypeGammaRotor [alphabetSize]int = [alphabetSize]int{5, 18, 14, 10, 0, 13, 20, 4, 17, 7, 12, 1, 19, 8, 24, 2, 22, 11, 16, 15, 25, 23, 21, 6, 9, 3}
)

var (
	TypeSteppingPos10 byte = 'R'
	TypeSteppingPos1  int  = 'R' - 'A'
	TypeSteppingPos2  int  = 'F' - 'A'
	TypeSteppingPos3  int  = 'W' - 'A'
	TypeSteppingPos4  int  = 'K' - 'A'
	TypeSteppingPos5  int  = 0
)

var (
	Reflector256   = []byte{138, 202, 153, 170, 248, 147, 144, 215, 230, 219, 172, 132, 246, 187, 239, 169, 255, 218, 237, 244, 139, 200, 221, 155, 174, 168, 222, 242, 150, 188, 231, 190, 145, 223, 241, 158, 208, 184, 128, 236, 243, 180, 207, 225, 161, 159, 211, 238, 196, 210, 251, 141, 185, 183, 143, 216, 195, 175, 214, 140, 182, 201, 192, 148, 212, 157, 134, 129, 130, 227, 232, 189, 181, 204, 151, 220, 229, 233, 137, 162, 249, 213, 206, 234, 197, 179, 173, 224, 209, 177, 135, 199, 194, 186, 253, 136, 165, 226, 240, 235, 156, 163, 154, 167, 146, 228, 198, 254, 164, 152, 245, 217, 160, 193, 205, 203, 250, 131, 171, 252, 149, 166, 142, 133, 191, 178, 176, 247, 38, 67, 68, 117, 11, 123, 66, 90, 95, 78, 0, 20, 59, 51, 122, 54, 6, 32, 104, 5, 63, 120, 28, 74, 109, 2, 102, 23, 100, 65, 35, 45, 112, 44, 79, 101, 108, 96, 121, 103, 25, 15, 3, 118, 10, 86, 24, 57, 126, 89, 125, 85, 41, 72, 60, 53, 37, 52, 93, 13, 29, 71, 31, 124, 62, 113, 92, 56, 48, 84, 106, 91, 21, 61, 1, 115, 73, 114, 82, 42, 36, 88, 49, 46, 64, 81, 58, 7, 55, 111, 17, 9, 75, 22, 26, 33, 87, 43, 97, 69, 105, 76, 8, 30, 70, 77, 83, 99, 39, 18, 47, 14, 98, 34, 27, 40, 19, 110, 12, 127, 4, 80, 116, 50, 119, 94, 107, 16}
	ReflectorB     = [alphabetSize]int{24, 17, 20, 7, 16, 18, 11, 3, 15, 23, 13, 6, 14, 10, 12, 8, 4, 1, 5, 25, 2, 22, 21, 9, 0, 19}
	ReflectorC     = [alphabetSize]int{5, 21, 15, 9, 8, 0, 14, 3, 4, 7, 17, 25, 23, 22, 6, 2, 19, 10, 20, 16, 18, 1, 13, 12, 24, 11}
	ReflectorBDunn = [alphabetSize]int{4, 13, 10, 16, 0, 20, 6, 24, 19, 17, 22, 25, 21, 1, 18, 23, 3, 9, 14, 8, 5, 12, 11, 15, 7, 2}
	ReflectorCDunn = [alphabetSize]int{17, 3, 14, 1, 9, 13, 19, 10, 21, 4, 22, 25, 23, 20, 2, 18, 7, 0, 16, 6, 12, 8, 11, 24, 15, 5}
)

func GenerateReflector_test() {
	ref := GenerateReflector(alphabetSize)
	tmp := make([]byte, alphabetSize)
	for i := range alphabetSize {
		tmp[i] = byte(i)
	}
	fmt.Println(tmp)
	fmt.Println(ref)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: enigma input.txt output.txt")
		os.Exit(1)
	}
	inputFileName := os.Args[1]
	outputFIlename := os.Args[2]

	// inputFileName := "../data/input.txt"
	// outputFIlename := "../data/output.txt"

	inputData, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rotors := []Rotor{
		NewRotor(TypeRotor256_1, 'Q', 0),
		NewRotor(TypeRotor256_2, 'U', 0),
		NewRotor(TypeRotor256_2, '8', 0),
	}
	reflector := NewReflector(Reflector256)
	enigm := NewEnigma(rotors, reflector)

	posRing := []byte{'Q', '8', '8'}
	enigm.SetRotorPositions(posRing)
	// text1 := enigm.EncryptText([]byte("ABC - abc 12345678 |||| !!"))
	// text1 := enigm.EncryptText([]byte{97})
	encryptedText := enigm.EncryptText(inputData)

	// enigm.SetRotorPositions(posRing)
	// text2 := enigm.EncryptText(text1)
	// fmt.Printf("%s\n%s\n", text1, text2)

	err = os.WriteFile(outputFIlename, encryptedText, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
