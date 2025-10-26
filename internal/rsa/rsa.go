package rsa

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"strings"
)

type Keys struct {
	C *big.Int
	D *big.Int
	N *big.Int
}

func GenerateKeys() Keys {
	P, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Something went wrong:%s", err.Error())
	}
	Q, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Something went wrong:%s", err.Error())
	}
	d, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Something went wrong:%s", err.Error())
	}
	N := new(big.Int).Mul(P, Q)
	phi := new(big.Int).Mul(P.Sub(P, big.NewInt(1)), Q.Sub(Q, big.NewInt(1)))
	gcd := new(big.Int)
	for {
		gcd.GCD(nil, nil, d, phi)
		if gcd.Cmp(big.NewInt(1)) == 0 {
			break
		}
		d, err = rand.Prime(rand.Reader, 1024)
		if err != nil {
			log.Fatalf("Something went wrong:%s", err.Error())
		}
	}
	c := new(big.Int).ModInverse(d, phi)

	return Keys{c, d, N}
}

func Encrypt(m, d, N *big.Int) *big.Int {
	if m.Cmp(N) != -1 {
		return nil
	}

	e := new(big.Int).Exp(m, d, N)
	return e
}

func Decrypt(e, c, N *big.Int) *big.Int {
	if e.Cmp(N) != -1 {
		return nil
	}

	m := new(big.Int).Exp(e, c, N)
	return m
}

func EncryptFile(inputFile, outputFile string, d, N *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	result := []*big.Int{}

	for _, b := range data {
		m := big.NewInt(int64(b))
		e := Encrypt(m, d, N)
		result = append(result, e)
	}

	out := ""
	for _, val := range result {
		out += val.String() + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func DecryptFile(inputFile, outputFile string, c, N *big.Int) error {
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
		e, _ := new(big.Int).SetString(line, 10)
		m := Decrypt(e, c, N)
		result = append(result, byte(m.Int64()))
	}

	return os.WriteFile(outputFile, result, 0644)
}
