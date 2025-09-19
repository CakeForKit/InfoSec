package rsa_alg

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
)

type RSA struct {
	publicKey  *PublicKey
	privateKey *PrivateKey
	keysFile   string
}

type PublicKey struct {
	E *big.Int
	N *big.Int
}

type PrivateKey struct {
	D *big.Int
	N *big.Int
}

func NewRSA(loadKeys bool, keysFileName string) (*RSA, error) {
	rsa := &RSA{
		keysFile: keysFileName,
	}

	if loadKeys && FileExists(rsa.keysFile) {
		rsa.loadKeys()
	} else {
		rsa.generateKeys()
		err := rsa.saveKeys()
		if err != nil {
			return nil, fmt.Errorf("NewRSA: %w", err)
		}
	}

	return rsa, nil
}

func (r *RSA) generatePrime(min, max int64) *big.Int {
	for {
		// Генерируем случайное число в диапазоне
		n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
		if err != nil {
			panic(err)
		}
		n.Add(n, big.NewInt(min))

		// Проверяем на простоту
		if n.ProbablyPrime(20) {
			return n
		}
	}
}

func (r *RSA) generateKeys() {
	// Генерация p и q
	p := r.generatePrime(100, 300)
	q := r.generatePrime(100, 300)
	for p.Cmp(q) == 0 {
		q = r.generatePrime(100, 300)
	}

	// Вычисляем n и φ(n)
	n := new(big.Int).Mul(p, q)
	phi := new(big.Int).Mul(
		new(big.Int).Sub(p, big.NewInt(1)),
		new(big.Int).Sub(q, big.NewInt(1)),
	)

	// Выбор открытой экспоненты e
	e := big.NewInt(65537)
	for new(big.Int).GCD(nil, nil, e, phi).Cmp(big.NewInt(1)) != 0 {
		e, _ = rand.Int(rand.Reader, new(big.Int).Sub(phi, big.NewInt(3)))
		e.Add(e, big.NewInt(3))
	}

	// Вычисление закрытой экспоненты d
	d := new(big.Int).ModInverse(e, phi)

	r.publicKey = &PublicKey{E: e, N: n}
	r.privateKey = &PrivateKey{D: d, N: n}
}

func (r *RSA) saveKeys() error {
	keys := map[string]interface{}{
		"public_key": map[string]string{
			"e": r.publicKey.E.String(),
			"n": r.publicKey.N.String(),
		},
		"private_key": map[string]string{
			"d": r.privateKey.D.String(),
			"n": r.privateKey.N.String(),
		},
	}

	jsonData, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return fmt.Errorf("saveKeys: %w", err)
	}

	err = os.WriteFile(r.keysFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("saveKeys: %w", err)
	}
	return nil
}

func (r *RSA) loadKeys() {
	data, err := os.ReadFile(r.keysFile)
	if err != nil {
		panic(err)
	}

	var keys map[string]interface{}
	err = json.Unmarshal(data, &keys)
	if err != nil {
		panic(err)
	}

	// Загрузка публичного ключа
	publicKey := keys["public_key"].(map[string]interface{})
	e := new(big.Int)
	e.SetString(publicKey["e"].(string), 10)
	n := new(big.Int)
	n.SetString(publicKey["n"].(string), 10)
	r.publicKey = &PublicKey{E: e, N: n}

	// Загрузка приватного ключа
	privateKey := keys["private_key"].(map[string]interface{})
	d := new(big.Int)
	d.SetString(privateKey["d"].(string), 10)
	n = new(big.Int)
	n.SetString(privateKey["n"].(string), 10)
	r.privateKey = &PrivateKey{D: d, N: n}
}

func (r *RSA) encryptByte(b byte) *big.Int {
	byteInt := big.NewInt(int64(b))
	return new(big.Int).Exp(byteInt, r.publicKey.E, r.publicKey.N)
}

func (r *RSA) decryptByte(encrypted *big.Int) byte {
	decrypted := new(big.Int).Exp(encrypted, r.privateKey.D, r.privateKey.N)
	return byte(decrypted.Int64())
}

func (r *RSA) EncryptFile(inputPath, outputPath string) error {
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

	buffer := make([]byte, 1)
	for {
		_, err := inputFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		encrypted := r.encryptByte(buffer[0])
		encryptedBytes := encrypted.Bytes()

		length := byte(len(encryptedBytes))
		if _, err := outputFile.Write([]byte{length}); err != nil {
			return err
		}
		if _, err := outputFile.Write(encryptedBytes); err != nil {
			return err
		}
	}

	return nil
}

func (r *RSA) DecryptFile(inputPath, outputPath string) error {
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

	for {
		// Читаем длину зашифрованного блока
		lengthBuf := make([]byte, 1)
		_, err := inputFile.Read(lengthBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		length := int(lengthBuf[0])
		if length == 0 {
			continue
		}

		// Читаем зашифрованные данные
		encryptedBytes := make([]byte, length)
		_, err = inputFile.Read(encryptedBytes)
		if err != nil {
			return err
		}

		encrypted := new(big.Int).SetBytes(encryptedBytes)
		decrypted := r.decryptByte(encrypted)

		if _, err := outputFile.Write([]byte{decrypted}); err != nil {
			return err
		}
	}

	return nil
}
