package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func gcd(a, b *big.Int) *big.Int {
	zero := big.NewInt(0)
	for b.Cmp(zero) != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}
	return a
}

func modInverse(a, m *big.Int) *big.Int {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	g.GCD(x, y, a, m)
	if g.Cmp(big.NewInt(1)) != 0 {
		return nil
	}
	return x.Mod(x, m)
}

func generateKeys(p *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	var c *big.Int
	for {
		cand, _ := rand.Int(rand.Reader, pMinus1)
		if gcd(cand, pMinus1).Cmp(big.NewInt(1)) == 0 {
			c = cand
			break
		}
	}
	d := modInverse(c, pMinus1)
	return c, d
}

func shamirProtocol(m, p, cA, dA, cB, dB *big.Int) *big.Int {
	x1 := new(big.Int).Exp(m, cA, p)
	x2 := new(big.Int).Exp(x1, cB, p)
	x3 := new(big.Int).Exp(x2, dA, p)
	x4 := new(big.Int).Exp(x3, dB, p)
	return x4
}

func encryptFile(inputFile, outputFile string, p, cA, dA, cB, dB *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	result := []*big.Int{}

	for _, b := range data {
		m := big.NewInt(int64(b))
		if m.Cmp(pMinus1) >= 0 {
			return fmt.Errorf("байт %d >= p, увеличьте простое число p", b)
		}
		enc := shamirProtocol(m, p, cA, dA, cB, dB)
		result = append(result, enc)
	}

	out := ""
	for _, val := range result {
		out += val.String() + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func decryptFile(inputFile, outputFile string, p, cA, dA, cB, dB *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []byte

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		val, ok := new(big.Int).SetString(line, 10)
		if !ok {
			return fmt.Errorf("ошибка парсинга числа: %s", line)
		}
		dec := shamirProtocol(val, p, cA, dA, cB, dB)
		result = append(result, byte(dec.Int64()))
	}

	return os.WriteFile(outputFile, result, 0644)
}

func readBigInt(prompt string) *big.Int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		val, ok := new(big.Int).SetString(text, 10)
		if ok {
			return val
		}
		fmt.Println("Ошибка: введите целое число")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Шифр Шамира для файлов (Go)")
	fmt.Println("1 - Ввести параметры вручную")
	fmt.Println("2 - Сгенерировать автоматически")
	fmt.Print("Выбор: ")
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, _ := strconv.Atoi(choiceStr)

	var p, cA, dA, cB, dB *big.Int

	if choice == 1 {
		p = readBigInt("Введите простое число p (> 255): ")
		cA = readBigInt("Введите cA: ")
		dA = readBigInt("Введите dA: ")
		cB = readBigInt("Введите cB: ")
		dB = readBigInt("Введите dB: ")
	} else {
		p = big.NewInt(257)
		cA, dA = generateKeys(p)
		cB, dB = generateKeys(p)
		fmt.Println("Автоматически сгенерированные параметры:")
		fmt.Println("p =", p)
		fmt.Println("A: cA =", cA, "dA =", dA)
		fmt.Println("B: cB =", cB, "dB =", dB)
	}

	fmt.Println("\nРежим работы:")
	fmt.Println("1 - Шифрование")
	fmt.Println("2 - Расшифровка")
	fmt.Print("Выбор: ")
	modeStr, _ := reader.ReadString('\n')
	modeStr = strings.TrimSpace(modeStr)
	mode, _ := strconv.Atoi(modeStr)

	fmt.Print("Введите имя входного файла: ")
	inFile, _ := reader.ReadString('\n')
	inFile = strings.TrimSpace(inFile)

	fmt.Print("Введите имя выходного файла: ")
	outFile, _ := reader.ReadString('\n')
	outFile = strings.TrimSpace(outFile)

	if mode == 1 {
		if err := encryptFile(inFile, outFile, p, cA, dA, cB, dB); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		fmt.Println("Файл зашифрован:", outFile)
	} else {
		if err := decryptFile(inFile, outFile, p, cA, dA, cB, dB); err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		fmt.Println("Файл расшифрован:", outFile)
	}
}
