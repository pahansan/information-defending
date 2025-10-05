package elgamal

import (
	"crypto/rand"
	"math/big"
	"os"
	"strings"
)

func RandElGamal(p, g *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))

	Ci, _ := rand.Int(rand.Reader, pMinus1)
	Ci = Ci.Add(Ci, big.NewInt(1))

	Di := new(big.Int).Exp(g, Ci, p)
	return Ci, Di
}

func ElGamalEncrypt(p, g, Db, k, m *big.Int) (*big.Int, *big.Int) {
	r := new(big.Int).Exp(g, k, p)
	e := new(big.Int).Mul(m, new(big.Int).Exp(Db, k, p))
	e = e.Mod(e, p)

	return r, e
}

func ElGamalDecrypt(e, r, p, Cb *big.Int) *big.Int {
	pMinus1Cb := new(big.Int).Sub(p, big.NewInt(1))
	pMinus1Cb = pMinus1Cb.Sub(pMinus1Cb, Cb)

	m := new(big.Int).Mul(e, new(big.Int).Exp(r, pMinus1Cb, p))
	return m.Mod(m, p)
}

func EncryptFile(inputFile, outputFile string, p, g, Db, k *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	result := []*big.Int{}

	for _, b := range data {
		m := big.NewInt(int64(b))
		r, e := ElGamalEncrypt(p, g, Db, k, m)
		result = append(result, r)
		result = append(result, e)
	}

	out := ""
	for _, val := range result {
		out += val.String() + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func DecryptFile(inputFile, outputFile string, p, Cb *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []byte

	flag := true
	var r, e *big.Int
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if flag {
			r, _ = new(big.Int).SetString(line, 10)
		} else {
			e, _ = new(big.Int).SetString(line, 10)
			dec := ElGamalDecrypt(e, r, p, Cb)
			result = append(result, byte(dec.Int64()))
		}
		flag = !flag
	}

	return os.WriteFile(outputFile, result, 0644)
}
