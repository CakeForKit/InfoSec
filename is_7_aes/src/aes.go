package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

// AES представляет реализацию алгоритма AES
type AES struct {
	key       []byte
	keySize   int
	rounds    int
	roundKeys [][]byte
}

// NewAES создает новый экземпляр AES с заданным ключом
func NewAES(key []byte) (*AES, error) {
	keySize := len(key)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("неверный размер ключа: %d байт, должен быть 16, 24 или 32", keySize)
	}

	aes := &AES{
		key:     key,
		keySize: keySize,
	}

	// Определяем количество раундов
	switch keySize {
	case 16: // 128 бит
		aes.rounds = 10
	case 24: // 192 бита
		aes.rounds = 12
	case 32: // 256 бит
		aes.rounds = 14
	}

	// Генерируем раундовые ключи
	if err := aes.keyExpansion(); err != nil {
		return nil, err
	}

	return aes, nil
}

// keyExpansion генерирует раундовые ключи
func (a *AES) keyExpansion() error {
	nk := a.keySize / 4
	nr := a.rounds

	// Количество слов в расширенном ключе
	keyWords := 4 * (nr + 1)
	a.roundKeys = make([][]byte, keyWords)

	// Инициализируем первые Nk слов исходным ключом
	for i := 0; i < nk; i++ {
		a.roundKeys[i] = make([]byte, 4)
		copy(a.roundKeys[i], a.key[i*4:(i+1)*4])
	}

	// Генерируем остальные слова
	for i := nk; i < keyWords; i++ {
		temp := make([]byte, 4)
		copy(temp, a.roundKeys[i-1])

		if i%nk == 0 {
			// RotWord + SubWord + Rcon
			temp = a.rotWord(temp)
			temp = a.subWord(temp)
			temp[0] ^= byte(rcon[i/nk] >> 24)
		} else if nk > 6 && i%nk == 4 {
			// SubWord для 256-битных ключей
			temp = a.subWord(temp)
		}

		a.roundKeys[i] = make([]byte, 4)
		for j := 0; j < 4; j++ {
			a.roundKeys[i][j] = a.roundKeys[i-nk][j] ^ temp[j]
		}
	}

	return nil
}

// rotWord выполняет циклический сдвиг влево на 1 байт
func (a *AES) rotWord(word []byte) []byte {
	return []byte{word[1], word[2], word[3], word[0]}
}

// subWord применяет S-Box к каждому байту слова
func (a *AES) subWord(word []byte) []byte {
	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = sBox[word[i]]
	}
	return result
}

// addRoundKey применяет XOR с раундовым ключом
func (a *AES) addRoundKey(state *[4][4]byte, round int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] ^= a.roundKeys[round*4+j][i]
		}
	}
}

// subBytes заменяет каждый байт состояния с помощью S-Box
func (a *AES) subBytes(state *[4][4]byte) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = sBox[state[i][j]]
		}
	}
}

// invSubBytes обратная операция subBytes
func (a *AES) invSubBytes(state *[4][4]byte) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = invSBox[state[i][j]]
		}
	}
}

// shiftRows выполняет сдвиг строк состояния
func (a *AES) shiftRows(state *[4][4]byte) {
	// Вторая строка: сдвиг на 1
	temp := state[1][0]
	state[1][0] = state[1][1]
	state[1][1] = state[1][2]
	state[1][2] = state[1][3]
	state[1][3] = temp

	// Третья строка: сдвиг на 2
	state[2][0], state[2][2] = state[2][2], state[2][0]
	state[2][1], state[2][3] = state[2][3], state[2][1]

	// Четвертая строка: сдвиг на 3
	temp = state[3][0]
	state[3][0] = state[3][1]
	state[3][1] = state[3][2]
	state[3][2] = state[3][3]
	state[3][3] = temp
}

// invShiftRows обратная операция shiftRows
func (a *AES) invShiftRows(state *[4][4]byte) {
	// Вторая строка: обратный сдвиг на 1
	temp := state[1][3]
	state[1][3] = state[1][2]
	state[1][2] = state[1][1]
	state[1][1] = state[1][0]
	state[1][0] = temp

	// Третья строка: сдвиг на 2 (обратный такой же)
	state[2][0], state[2][2] = state[2][2], state[2][0]
	state[2][1], state[2][3] = state[2][3], state[2][1]

	// Четвертая строка: обратный сдвиг на 3
	temp = state[3][0]
	state[3][0] = state[3][1]
	state[3][1] = state[3][2]
	state[3][2] = state[3][3]
	state[3][3] = temp
}

// xtime умножает байт на 2 в поле GF(2^8)
func (a *AES) xtime(b byte) byte {
	if b&0x80 != 0 {
		return (b << 1) ^ 0x1b
	}
	return b << 1
}

// mixColumns перемешивает столбцы состояния
func (a *AES) mixColumns(state *[4][4]byte) {
	for i := 0; i < 4; i++ {
		s0 := state[0][i]
		s1 := state[1][i]
		s2 := state[2][i]
		s3 := state[3][i]

		h := s0 ^ s1 ^ s2 ^ s3
		state[0][i] ^= h ^ a.xtime(s0^s1)
		state[1][i] ^= h ^ a.xtime(s1^s2)
		state[2][i] ^= h ^ a.xtime(s2^s3)
		state[3][i] ^= h ^ a.xtime(s3^s0)
	}
}

// invMixColumns обратная операция mixColumns
func (a *AES) invMixColumns(state *[4][4]byte) {
	for i := 0; i < 4; i++ {
		a := state[0][i]
		b := state[1][i]
		c := state[2][i]
		d := state[3][i]

		state[0][i] = a.xtime(a.xtime(a.xtime(a^b^c^d)^a^b^d) ^ a ^ b ^ c)
		state[1][i] = a.xtime(a.xtime(a.xtime(a^b^c^d)^b^c^a) ^ b ^ c ^ d)
		state[2][i] = a.xtime(a.xtime(a.xtime(a^b^c^d)^c^d^b) ^ c ^ d ^ a)
		state[3][i] = a.xtime(a.xtime(a.xtime(a^b^c^d)^d^a^c) ^ d ^ a ^ b)
	}
}

// bytesToState преобразует массив байт в состояние 4x4
func (a *AES) bytesToState(input []byte) [4][4]byte {
	var state [4][4]byte
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = input[i+4*j]
		}
	}
	return state
}

// stateToBytes преобразует состояние 4x4 в массив байт
func (a *AES) stateToBytes(state [4][4]byte) []byte {
	output := make([]byte, 16)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			output[i+4*j] = state[i][j]
		}
	}
	return output
}

// EncryptBlock шифрует один блок (16 байт)
func (a *AES) EncryptBlock(block []byte) []byte {
	if len(block) != 16 {
		panic("размер блока должен быть 16 байт")
	}

	state := a.bytesToState(block)

	// Начальный раунд
	a.addRoundKey(&state, 0)

	// Основные раунды
	for round := 1; round < a.rounds; round++ {
		a.subBytes(&state)
		a.shiftRows(&state)
		a.mixColumns(&state)
		a.addRoundKey(&state, round)
	}

	// Финальный раунд
	a.subBytes(&state)
	a.shiftRows(&state)
	a.addRoundKey(&state, a.rounds)

	return a.stateToBytes(state)
}

// DecryptBlock расшифровывает один блок (16 байт)
func (a *AES) DecryptBlock(block []byte) []byte {
	if len(block) != 16 {
		panic("размер блока должен быть 16 байт")
	}

	state := a.bytesToState(block)

	// Финальный раунд (обратный)
	a.addRoundKey(&state, a.rounds)
	a.invShiftRows(&state)
	a.invSubBytes(&state)

	// Основные раунды (обратные)
	for round := a.rounds - 1; round > 0; round-- {
		a.addRoundKey(&state, round)
		a.invMixColumns(&state)
		a.invShiftRows(&state)
		a.invSubBytes(&state)
	}

	// Начальный раунд
	a.addRoundKey(&state, 0)

	return a.stateToBytes(state)
}

// TestBasicOperations тестирует базовые операции AES
func (a *AES) TestBasicOperations() {
	fmt.Println("\n=== ТЕСТ БАЗОВЫХ ОПЕРАЦИЙ ===")

	// Тестовый блок из стандарта AES
	testBlock := []byte{
		0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d,
		0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x07, 0x34,
	}

	fmt.Printf("Исходный блок: %x\n", testBlock)

	encrypted := a.EncryptBlock(testBlock)
	fmt.Printf("Зашифрованный блок: %x\n", encrypted)

	decrypted := a.DecryptBlock(encrypted)
	fmt.Printf("Расшифрованный блок: %x\n", decrypted)

	// Проверяем совпадение
	match := true
	for i := 0; i < 16; i++ {
		if testBlock[i] != decrypted[i] {
			match = false
			break
		}
	}

	if match {
		fmt.Println("✓ Базовые операции работают корректно")
	} else {
		fmt.Println("✗ Ошибка в базовых операциях AES!")
		fmt.Printf("Ожидалось: %x\n", testBlock)
		fmt.Printf("Получено:  %x\n", decrypted)
	}
}

// Encrypt шифрует данные с использованием режима CBC
func (a *AES) Encrypt(plaintext []byte) []byte {
	// Добавление padding по PKCS#7
	blockSize := 16
	padding := blockSize - (len(plaintext) % blockSize)
	if padding == 0 {
		padding = blockSize
	}

	padded := make([]byte, len(plaintext)+padding)
	copy(padded, plaintext)
	for i := len(plaintext); i < len(padded); i++ {
		padded[i] = byte(padding)
	}

	// Генерация IV (Initialization Vector)
	iv := make([]byte, blockSize)
	rand.Read(iv)

	// Шифрование в режиме CBC
	ciphertext := make([]byte, len(padded)+blockSize)
	copy(ciphertext[:blockSize], iv)

	currentIV := iv
	for i := 0; i < len(padded); i += blockSize {
		block := make([]byte, blockSize)
		copy(block, padded[i:i+blockSize])

		// XOR с текущим IV
		for j := 0; j < blockSize; j++ {
			block[j] ^= currentIV[j]
		}

		encryptedBlock := a.EncryptBlock(block)
		copy(ciphertext[blockSize+i:blockSize+i+blockSize], encryptedBlock)

		// Обновляем IV для следующего блока
		currentIV = encryptedBlock
	}

	return ciphertext
}

// Decrypt расшифровывает данные с использованием режима CBC
func (a *AES) Decrypt(ciphertext []byte) []byte {
	blockSize := 16
	if len(ciphertext) < blockSize*2 || len(ciphertext)%blockSize != 0 {
		panic(fmt.Sprintf("неверный размер зашифрованных данных: %d байт", len(ciphertext)))
	}

	// Извлечение IV
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// Расшифровка в режиме CBC
	plaintext := make([]byte, len(ciphertext))
	currentIV := iv

	for i := 0; i < len(ciphertext); i += blockSize {
		block := make([]byte, blockSize)
		copy(block, ciphertext[i:i+blockSize])

		decryptedBlock := a.DecryptBlock(block)

		// XOR с текущим IV
		for j := 0; j < blockSize; j++ {
			plaintext[i+j] = decryptedBlock[j] ^ currentIV[j]
		}

		// Обновляем IV для следующего блока
		currentIV = block
	}

	// Удаление padding
	if len(plaintext) == 0 {
		return plaintext
	}

	padding := int(plaintext[len(plaintext)-1])
	if padding < 1 || padding > blockSize {
		return plaintext
	}

	// Проверяем корректность padding
	validPadding := true
	for i := len(plaintext) - padding; i < len(plaintext); i++ {
		if plaintext[i] != byte(padding) {
			validPadding = false
			break
		}
	}

	if validPadding {
		return plaintext[:len(plaintext)-padding]
	}

	return plaintext
}

// FileEncrypt шифрует файл
func (a *AES) FileEncrypt(inputFile, outputFile string) error {
	// Чтение входного файла
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	fmt.Printf("=== ШИФРОВАНИЕ ===\n")
	fmt.Printf("Исходный файл: %d байт\n", len(input))
	fmt.Printf("Содержимое: %s\n", string(input))
	fmt.Printf("Hex: %x\n", input)

	// Шифрование
	ciphertext := a.Encrypt(input)

	fmt.Printf("Результат шифрования: %d байт\n", len(ciphertext))
	fmt.Printf("Hex зашифрованных данных: %x\n", ciphertext)

	// Запись зашифрованных данных
	err = os.WriteFile(outputFile, ciphertext, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

// FileDecrypt расшифровывает файл
func (a *AES) FileDecrypt(inputFile, outputFile string) error {
	// Чтение зашифрованного файла
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	fmt.Printf("=== РАСШИФРОВАНИЕ ===\n")
	fmt.Printf("Зашифрованный файл: %d байт\n", len(ciphertext))
	fmt.Printf("Hex зашифрованных данных: %x\n", ciphertext)

	if len(ciphertext) < 32 || len(ciphertext)%16 != 0 {
		return fmt.Errorf("неверный размер зашифрованного файла: %d байт", len(ciphertext))
	}

	// Расшифровка
	plaintext := a.Decrypt(ciphertext)

	fmt.Printf("Результат расшифровки: %d байт\n", len(plaintext))
	fmt.Printf("Содержимое после расшифровки: %s\n", string(plaintext))
	fmt.Printf("Hex после расшифровки: %x\n", plaintext)

	// Запись расшифрованных данных
	err = os.WriteFile(outputFile, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

// TestSmallFile тестирует работу с маленьким файлом
func (a *AES) TestSmallFile() {
	fmt.Println("\n=== ТЕСТ МАЛЕНЬКОГО ФАЙЛА ===")

	// Тестовые данные разного размера
	testCases := [][]byte{
		[]byte("A"),           // 1 байт
		[]byte("Hello"),       // 5 байт
		[]byte("Hello World"), // 11 байт
		[]byte(""),            // 0 байт
	}

	for i, testData := range testCases {
		fmt.Printf("\nТест %d: размер %d байт\n", i+1, len(testData))
		fmt.Printf("Исходные данные: '%s'\n", testData)

		// Шифрование
		ciphertext := a.Encrypt(testData)
		fmt.Printf("Зашифровано: %d байт\n", len(ciphertext))

		// Расшифровка
		plaintext := a.Decrypt(ciphertext)
		fmt.Printf("Расшифровано: %d байт, данные: '%s'\n", len(plaintext), plaintext)

		if string(testData) == string(plaintext) {
			fmt.Printf("✓ Тест %d пройден\n", i+1)
		} else {
			fmt.Printf("✗ Тест %d не пройден\n", i+1)
			fmt.Printf("  Ожидалось: '%s' (%x)\n", testData, testData)
			fmt.Printf("  Получено:  '%s' (%x)\n", plaintext, plaintext)
		}
	}
}

// generateAndSaveKey генерирует случайный ключ и сохраняет его в файл
func generateAndSaveKey(filename string, keySize int) ([]byte, error) {
	var key []byte
	switch keySize {
	case 128:
		key = make([]byte, 16)
	case 192:
		key = make([]byte, 24)
	case 256:
		key = make([]byte, 32)
	default:
		return nil, fmt.Errorf("неподдерживаемый размер ключа: %d", keySize)
	}

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filename, key, 0644)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// readKeyFromFile читает ключ из файла
func readKeyFromFile(filename string) ([]byte, error) {
	key, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Проверяем размер ключа
	switch len(key) {
	case 16, 24, 32:
		return key, nil
	default:
		return nil, fmt.Errorf("неверный размер ключа: %d байт. Должен быть 16, 24 или 32 байта", len(key))
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование:")
		fmt.Println("  Шифрование: go run aes.go encrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Дешифрование: go run aes.go decrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Генерация ключа: go run aes.go genkey <файл_ключа> [128|192|256]")
		fmt.Println("")
		fmt.Println("Размеры ключа: 128, 192 или 256 бит (по умолчанию: 128)")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "genkey":
		if len(os.Args) < 3 {
			fmt.Println("Использование: go run aes.go genkey <файл_ключа> [128|192|256]")
			return
		}
		keyFilePath := os.Args[2]

		keySize := 128
		if len(os.Args) >= 4 {
			switch os.Args[3] {
			case "128":
				keySize = 128
			case "192":
				keySize = 192
			case "256":
				keySize = 256
			default:
				fmt.Printf("Неверный размер ключа: %s. Используйте 128, 192 или 256\n", os.Args[3])
				return
			}
		}

		key, err := generateAndSaveKey(keyFilePath, keySize)
		if err != nil {
			fmt.Printf("Ошибка генерации ключа: %v\n", err)
			return
		}
		fmt.Printf("Ключ успешно сгенерирован и сохранен в: %s\n", keyFilePath)
		fmt.Printf("Размер ключа: %d бит\n", keySize)
		fmt.Printf("Hex-представление: %s\n", hex.EncodeToString(key))

	case "encrypt", "decrypt":
		if len(os.Args) < 5 {
			fmt.Printf("Использование: go run aes.go %s <файл_ключа> <входной_файл> <выходной_файл>\n", mode)
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

		fmt.Printf("Загружен ключ: %s\n", hex.EncodeToString(key))
		fmt.Printf("Размер ключа: %d бит\n", len(key)*8)

		// Создание экземпляра AES
		aes, err := NewAES(key)
		if err != nil {
			fmt.Printf("Ошибка создания AES: %v\n", err)
			return
		}

		// Тестирование базовых операций
		aes.TestBasicOperations()

		// Проверка существования входного файла
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Printf("Входной файл не существует: %s\n", inputFile)
			return
		}

		switch mode {
		case "encrypt":
			err = aes.FileEncrypt(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка шифрования: %v\n", err)
				return
			}

		case "decrypt":
			err = aes.FileDecrypt(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка дешифрования: %v\n", err)
				return
			}
		}

	default:
		fmt.Println("Неизвестный режим. Используйте 'encrypt', 'decrypt' или 'genkey'")
		fmt.Println("Использование:")
		fmt.Println("  Шифрование: go run aes.go encrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Дешифрование: go run aes.go decrypt <файл_ключа> <входной_файл> <выходной_файл>")
		fmt.Println("  Генерация ключа: go run aes.go genkey <файл_ключа> [128|192|256]")
	}
}
