package shamir

import (
	"crypto/rand"
	"fmt"
	"information-defending/internal/crypto"
	"math/big"
	"os"
	"strings"
)

func GenerateKeys(p *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	var c *big.Int
	for {
		cand, _ := rand.Int(rand.Reader, pMinus1)
		if crypto.Gcd(cand, pMinus1).Cmp(big.NewInt(1)) == 0 {
			c = cand
			break
		}
	}
	d := crypto.ModInverse(c, pMinus1)
	return c, d
}

func Protocol(m, p, k1, k2 *big.Int) *big.Int {
	x1 := new(big.Int).Exp(m, k1, p)
	x2 := new(big.Int).Exp(x1, k2, p)
	return x2
}

func EncryptFile(inputFile, outputFile string, p, cA, cB *big.Int) error {
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
		enc := Protocol(m, p, cA, cB)
		result = append(result, enc)
	}

	out := ""
	for _, val := range result {
		out += val.String() + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func DecryptFile(inputFile, outputFile string, p, dA, dB *big.Int) error {
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
		dec := Protocol(val, p, dA, dB)
		result = append(result, byte(dec.Int64()))
	}

	return os.WriteFile(outputFile, result, 0644)
}
