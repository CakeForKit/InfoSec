package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// DES реализация алгоритма DES
type DES struct {
	subkeys [16]uint64
}

// НовыйDES создает новый экземпляр DES с заданным ключом
func NewDES(key uint64) *DES {
	des := &DES{}
	des.generateSubkeys(key)
	return des
}

// Начальная перестановка IP
func initialPermutation(block uint64) uint64 {
	ipTable := [64]int{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7,
	}
	return permute(block, ipTable[:], 64)
}

// Конечная перестановка IP^-1
func finalPermutation(block uint64) uint64 {
	fpTable := [64]int{
		40, 8, 48, 16, 56, 24, 64, 32,
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25,
	}
	return permute(block, fpTable[:], 64)
}

// Расширяющая перестановка E
func expansionPermutation(right uint32) uint64 {
	eTable := [48]int{
		32, 1, 2, 3, 4, 5,
		4, 5, 6, 7, 8, 9,
		8, 9, 10, 11, 12, 13,
		12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21,
		20, 21, 22, 23, 24, 25,
		24, 25, 26, 27, 28, 29,
		28, 29, 30, 31, 32, 1,
	}
	return uint64(permute32(uint64(right), eTable[:], 32))
}

// Перестановка P
func pPermutation(data uint32) uint32 {
	pTable := [32]int{
		16, 7, 20, 21, 29, 12, 28, 17,
		1, 15, 23, 26, 5, 18, 31, 10,
		2, 8, 24, 14, 32, 27, 3, 9,
		19, 13, 30, 6, 22, 11, 4, 25,
	}
	return uint32(permute32(uint64(data), pTable[:], 32))
}

// S-блоки (подстановка)
func sBoxSubstitution(data uint64) uint32 {
	sBoxes := [8][4][16]uint8{
		// S1
		{
			{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
			{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
			{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
			{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13},
		},
		// S2
		{
			{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
			{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
			{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
			{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
		},
		// S3
		{
			{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
			{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
			{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
			{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
		},
		// S4
		{
			{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
			{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
			{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
			{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14},
		},
		// S5
		{
			{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
			{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
			{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
			{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
		},
		// S6
		{
			{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
			{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
			{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
			{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
		},
		// S7
		{
			{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
			{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
			{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
			{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12},
		},
		// S8
		{
			{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
			{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
			{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
			{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11},
		},
	}

	var result uint32
	for i := 0; i < 8; i++ {
		// Берем 6 бит для текущего S-блока
		chunk := (data >> (42 - 6*i)) & 0x3F

		// Вычисляем строку и столбец
		row := ((chunk & 0x20) >> 4) | (chunk & 0x01)
		col := (chunk >> 1) & 0x0F

		// Получаем значение из S-блока
		sVal := sBoxes[i][row][col]

		// Добавляем к результату
		result = (result << 4) | uint32(sVal)
	}

	return result
}

// Функция Фейстеля
func feistelFunction(right uint32, subkey uint64) uint32 {
	// Расширяющая перестановка
	expanded := expansionPermutation(right)

	// XOR с подключом
	xored := expanded ^ subkey

	// S-блоки
	substituted := sBoxSubstitution(xored)

	// Перестановка P
	return pPermutation(substituted)
}

// Генерация подключей
func (des *DES) generateSubkeys(key uint64) {
	// PC-1 перестановка (удаление битов четности)
	pc1Table := [56]int{
		57, 49, 41, 33, 25, 17, 9, 1,
		58, 50, 42, 34, 26, 18, 10, 2,
		59, 51, 43, 35, 27, 19, 11, 3,
		60, 52, 44, 36, 63, 55, 47, 39,
		31, 23, 15, 7, 62, 54, 46, 38,
		30, 22, 14, 6, 61, 53, 45, 37,
		29, 21, 13, 5, 28, 20, 12, 4,
	}

	// PC-2 перестановка
	pc2Table := [48]int{
		14, 17, 11, 24, 1, 5, 3, 28,
		15, 6, 21, 10, 23, 19, 12, 4,
		26, 8, 16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55, 30, 40,
		51, 45, 33, 48, 44, 49, 39, 56,
		34, 53, 46, 42, 50, 36, 29, 32,
	}

	// Сдвиги для каждого раунда
	shifts := [16]int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}

	// Начальная перестановка PC-1
	permutedKey := permute(key, pc1Table[:], 64)

	// Разделение на C0 и D0
	c := (permutedKey >> 28) & 0x0FFFFFFF
	d := permutedKey & 0x0FFFFFFF

	// Генерация 16 подключей
	for i := 0; i < 16; i++ {
		// Циклический сдвиг
		c = ((c << shifts[i]) | (c >> (28 - shifts[i]))) & 0x0FFFFFFF
		d = ((d << shifts[i]) | (d >> (28 - shifts[i]))) & 0x0FFFFFFF

		// Объединение и перестановка PC-2
		combined := (c << 28) | d
		des.subkeys[i] = permute(combined, pc2Table[:], 56)
	}
}

// Шифрование одного блока
func (d *DES) encryptBlock(block uint64) uint64 {
	// Начальная перестановка
	block = initialPermutation(block)

	// Разделение на левую и правую части
	left := uint32(block >> 32)
	right := uint32(block & 0xFFFFFFFF)

	// 16 раундов Фейстеля
	for i := 0; i < 16; i++ {
		nextLeft := right
		fResult := feistelFunction(right, d.subkeys[i])
		nextRight := left ^ fResult

		left = nextLeft
		right = nextRight
	}

	// Объединение (поменять местами после последнего раунда)
	combined := (uint64(right) << 32) | uint64(left)

	// Конечная перестановка
	return finalPermutation(combined)
}

// Дешифрование одного блока
func (d *DES) decryptBlock(block uint64) uint64 {
	// Используем подключи в обратном порядке
	temp := d.subkeys
	for i := 0; i < 8; i++ {
		d.subkeys[i], d.subkeys[15-i] = d.subkeys[15-i], d.subkeys[i]
	}

	result := d.encryptBlock(block)

	// Восстанавливаем порядок подключей
	d.subkeys = temp

	return result
}

// Вспомогательная функция для перестановки
func permute(data uint64, table []int, inputSize int) uint64 {
	var result uint64
	for i, pos := range table {
		bit := (data >> (inputSize - pos)) & 1
		result |= bit << (uint(len(table)) - 1 - uint(i))
	}
	return result
}

// Вспомогательная функция для перестановки 32-битных данных
func permute32(data uint64, table []int, inputSize int) uint64 {
	var result uint64
	for i, pos := range table {
		bit := (data >> (inputSize - pos)) & 1
		result |= bit << (uint(len(table)) - 1 - uint(i))
	}
	return result
}

// Шифрование файла
func (d *DES) encryptFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Буфер для чтения (8 байт = 64 бита)
	buffer := make([]byte, 8)

	for {
		n, err := inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		// Дополнение последнего блока если нужно
		if n < 8 {
			for i := n; i < 8; i++ {
				buffer[i] = 0 // Дополнение нулями
			}
		}

		// Конвертация байт в uint64
		block := binary.BigEndian.Uint64(buffer)

		// Шифрование блока
		encryptedBlock := d.encryptBlock(block)

		// Запись зашифрованного блока
		encryptedBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(encryptedBytes, encryptedBlock)

		_, err = outputFile.Write(encryptedBytes)
		if err != nil {
			return err
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

// Дешифрование файла
func (d *DES) decryptFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Буфер для чтения (8 байт = 64 бита)
	buffer := make([]byte, 8)

	for {
		n, err := inputFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		// Конвертация байт в uint64
		block := binary.BigEndian.Uint64(buffer)

		// Дешифрование блока
		decryptedBlock := d.decryptBlock(block)

		// Запись расшифрованного блока
		decryptedBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(decryptedBytes, decryptedBlock)

		// Для последнего блока убираем дополнение
		if err == io.EOF || n < 8 {
			decryptedBytes = decryptedBytes[:n]
		}

		_, err = outputFile.Write(decryptedBytes)
		if err != nil {
			return err
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func readKeyFromFile(keyFilePath string) (uint64, error) {
	// Чтение файла с ключом
	keyData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения файла ключа: %v", err)
	}

	// Очистка ключа от пробелов и переводов строк
	keyStr := strings.TrimSpace(string(keyData))
	keyStr = strings.ReplaceAll(keyStr, " ", "")
	keyStr = strings.ReplaceAll(keyStr, "\n", "")
	keyStr = strings.ReplaceAll(keyStr, "\r", "")
	keyStr = strings.ReplaceAll(keyStr, "0x", "")
	keyStr = strings.ReplaceAll(keyStr, "0X", "")

	// Проверка длины ключа
	if len(keyStr) != 16 {
		return 0, fmt.Errorf("неправильная длина ключа: ожидается 16 hex-символов (64 бита), получено %d", len(keyStr))
	}

	// Парсинг hex-строки
	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга hex-ключа: %v", err)
	}

	// Конвертация в uint64 (BigEndian)
	if len(keyBytes) != 8 {
		return 0, fmt.Errorf("неправильный размер ключа: ожидается 8 байт, получено %d", len(keyBytes))
	}

	return binary.BigEndian.Uint64(keyBytes), nil
}

func generateAndSaveKey(keyFilePath string) (uint64, error) {
	// Генерация случайного ключа
	keyBytes := make([]byte, 8)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("ошибка генерации ключа: %v", err)
	}

	key := binary.BigEndian.Uint64(keyBytes)

	// Сохранение ключа в файл в hex-формате
	keyHex := hex.EncodeToString(keyBytes)
	err = os.WriteFile(keyFilePath, []byte(keyHex), 0600) // права только для владельца
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения ключа: %v", err)
	}

	return key, nil
}

func main() {
	// fmt.Printf("len(os.Args) = %d \n", len(os.Args))
	if len(os.Args) < 3 {
		fmt.Println("Использование:")
		fmt.Println("  Шифрование: go run des.go encrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Дешифрование: go run des.go decrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Генерация ключа: go run des.go genkey <файл_ключа>")
		fmt.Println("")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "genkey":
		if len(os.Args) < 3 {
			fmt.Println("Использование: go run des.go genkey <файл_ключа>")
			return
		}
		keyFilePath := os.Args[2]

		key, err := generateAndSaveKey(keyFilePath)
		if err != nil {
			fmt.Printf("Ошибка генерации ключа: %v\n", err)
			return
		}
		fmt.Printf("Ключ успешно сгенерирован и сохранен в: %s\n", keyFilePath)
		fmt.Printf("Hex-представление: %016X\n", key)

	case "encrypt", "decrypt":
		if len(os.Args) < 5 {
			fmt.Printf("Использование: go run des.go %s <файл_ключа> <входной_файл> <выходной_файл>\n", mode)
			return
		}

		keyFilePath := os.Args[2]
		inputFile := os.Args[3]
		outputFile := os.Args[4]

		// Чтение ключа из файла
		key, err := readKeyFromFile(keyFilePath)
		if err != nil {
			fmt.Printf("Ошибка чтения ключа: %v\n", err)
			return
		}

		fmt.Printf("Загружен ключ: %016X\n", key)

		des := NewDES(key)

		// Проверка существования входного файла
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Printf("Входной файл не существует: %s\n", inputFile)
			return
		}

		switch mode {
		case "encrypt":
			err = des.encryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка шифрования: %v\n", err)
				return
			}
			fmt.Printf("Файл успешно зашифрован: %s -> %s\n", inputFile, outputFile)

		case "decrypt":
			err = des.decryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка дешифрования: %v\n", err)
				return
			}
			fmt.Printf("Файл успешно расшифрован: %s -> %s\n", inputFile, outputFile)
		}

	default:
		fmt.Println("Неизвестный режим. Используйте 'encrypt', 'decrypt' или 'genkey'")
		fmt.Println("Использование:")
		fmt.Println("  Шифрование: go run des.go encrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Дешифрование: go run des.go decrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Генерация ключа: go run des.go genkey <файл_ключа>")
	}
}
