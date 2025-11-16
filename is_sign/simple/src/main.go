package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func generateKeys(privateKeyPath, publicKeyPath string) error {
	// Генерация RSA ключа 2048 бит
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("ошибка генерации ключа: %v", err)
	}

	// Сохранение приватного ключа
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла приватного ключа: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("ошибка кодирования приватного ключа: %v", err)
	}

	// Сохранение публичного ключа
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла публичного ключа: %v", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга публичного ключа: %v", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return fmt.Errorf("ошибка кодирования публичного ключа: %v", err)
	}

	fmt.Printf("Ключи сгенерированы:\n - Приватный: %s\n - Публичный: %s\n", privateKeyPath, publicKeyPath)
	return nil
}

func signFile(inputFile, privateKeyPath, signatureFile string) error {
	// Чтение приватного ключа
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения приватного ключа: %v", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return fmt.Errorf("не удалось декодировать PEM блок приватного ключа")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("ошибка парсинга приватного ключа: %v", err)
	}

	// Чтение и хеширование файла
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	hash := sha256.Sum256(data)

	// Создание подписи
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("ошибка создания подписи: %v", err)
	}

	// Сохранение подписи
	if err := os.WriteFile(signatureFile, signature, 0644); err != nil {
		return fmt.Errorf("ошибка сохранения подписи: %v", err)
	}

	fmt.Printf("Подпись сохранена в %s\n", signatureFile)
	return nil
}

func verifySignature(inputFile, publicKeyPath, signatureFile string) error {
	// Чтение публичного ключа
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения публичного ключа: %v", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return fmt.Errorf("не удалось декодировать PEM блок публичного ключа")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("ошибка парсинга публичного ключа: %v", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("неверный тип публичного ключа")
	}

	// Чтение файла и вычисление хеша
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	hash := sha256.Sum256(data)

	// Чтение подписи
	signature, err := os.ReadFile(signatureFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения подписи: %v", err)
	}

	// Проверка подписи
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("подпись недействительна: %v", err)
	}

	fmt.Println("Подпись действительна.")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "gen-keys":
		privateKey := "private.pem"
		publicKey := "public.pem"

		// Парсинг аргументов для gen-keys
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "-priv", "--private":
				if i+1 < len(os.Args) {
					privateKey = os.Args[i+1]
					i++
				}
			case "-pub", "--public":
				if i+1 < len(os.Args) {
					publicKey = os.Args[i+1]
					i++
				}
			}
		}

		if err := generateKeys(privateKey, publicKey); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}

	case "sign":
		if len(os.Args) < 5 {
			fmt.Println("Ошибка: недостаточно аргументов для команды sign")
			fmt.Println("Использование: sign <input_file> -priv <private_key> [-o signature.bin]")
			os.Exit(1)
		}

		inputFile := os.Args[2]
		privateKey := ""
		output := "signature.bin"

		// Парсинг аргументов для sign
		for i := 3; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "-priv", "--private":
				if i+1 < len(os.Args) {
					privateKey = os.Args[i+1]
					i++
				}
			case "-o", "--output":
				if i+1 < len(os.Args) {
					output = os.Args[i+1]
					i++
				}
			}
		}

		if privateKey == "" {
			fmt.Println("Ошибка: необходимо указать приватный ключ с помощью -priv")
			os.Exit(1)
		}

		if err := signFile(inputFile, privateKey, output); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}

	case "verify":
		if len(os.Args) < 6 {
			fmt.Println("Ошибка: недостаточно аргументов для команды verify")
			fmt.Println("Использование: verify <input_file> -pub <public_key> -s <signature_file>")
			os.Exit(1)
		}

		inputFile := os.Args[2]
		publicKey := ""
		signature := ""

		// Парсинг аргументов для verify
		for i := 3; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "-pub", "--public":
				if i+1 < len(os.Args) {
					publicKey = os.Args[i+1]
					i++
				}
			case "-s", "--signature":
				if i+1 < len(os.Args) {
					signature = os.Args[i+1]
					i++
				}
			}
		}

		if publicKey == "" {
			fmt.Println("Ошибка: необходимо указать публичный ключ с помощью -pub")
			os.Exit(1)
		}
		if signature == "" {
			fmt.Println("Ошибка: необходимо указать файл подписи с помощью -s")
			os.Exit(1)
		}

		if err := verifySignature(inputFile, publicKey, signature); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Использование:")
	fmt.Println("  gen-keys [-priv private.pem] [-pub public.pem]")
	fmt.Println("  sign <input_file> -priv <private_key> [-o signature.bin]")
	fmt.Println("  verify <input_file> -pub <public_key> -s <signature_file>")
	fmt.Println("")
	fmt.Println("Примеры:")
	fmt.Println("  go run ./src/main.go gen-keys -priv private.pem -pub public.pem")
	fmt.Println("  go run ./src/main.go sign ./data/text.txt -priv private.pem -o signature.bin")
	fmt.Println("  go run ./src/main.go verify ./data/text.txt -pub public.pem -s signature.bin")
}
