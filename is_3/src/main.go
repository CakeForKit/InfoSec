package main

import (
	"fmt"
	"is_3/src/rsa_alg"
)

const (
	keysFileName string = "rsa_keys.json"
)

func main2() {
	rsa1, err := rsa_alg.NewRSA(false, keysFileName)
	_ = rsa1
	if err != nil {
		fmt.Println(err.Error())
	}

	rsa, err := rsa_alg.NewRSA(true, keysFileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = rsa.EncryptFile("./data/input.txt", "./data/outputEncr.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = rsa.DecryptFile("./data/outputEncr.txt", "./data/outputDecr.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	fmt.Println("\nВыберите действие:")
	fmt.Println("1. Сгенерировать новые ключи")
	fmt.Println("2. Зашифровать файл")
	fmt.Println("3. Расшифровать файл")
	fmt.Println("4. Выход")
	for {
		var choice int
		fmt.Print("Ваш выбор: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case 1:
			rsa, err := rsa_alg.NewRSA(false, keysFileName)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Новые ключи сгенерированы и сохранены")
			_ = rsa

		case 2:
			rsa, err := rsa_alg.NewRSA(true, keysFileName)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			var inputFile string
			fmt.Print("Введите путь к файлу для шифрования: ")
			fmt.Scan(&inputFile)

			if !rsa_alg.FileExists(inputFile) {
				fmt.Println("Файл не найден!")
				continue
			}

			outputFile := inputFile + ".encrypted.zip"
			err = rsa.EncryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка шифрования: %v\n", err)
				continue
			}
			fmt.Printf("Файл зашифрован и сохранен как %s\n", outputFile)

		case 3:
			rsa, _ := rsa_alg.NewRSA(true, keysFileName)

			var inputFile string
			fmt.Print("Введите путь к зашифрованному файлу: ")
			fmt.Scan(&inputFile)

			if !rsa_alg.FileExists(inputFile) {
				fmt.Println("Файл не найден!")
				continue
			}

			outputFile := inputFile + ".decrypted.zip"
			err := rsa.DecryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("Ошибка расшифровки: %v\n", err)
				continue
			}
			fmt.Printf("Файл расшифрован и сохранен как %s\n", outputFile)

		case 4:
			fmt.Println("Выход...")
			return

		default:
			fmt.Println("Неверный выбор!")
		}
	}
}
