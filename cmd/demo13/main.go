package main

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"information-defending/internal/rsa"
	"math/big"
)

type VoitingResult int

const (
	Yes VoitingResult = iota
	No
	Forgo
)

func main() {
	serverKeys := rsa.GenerateKeys()

	// 1. A voiting
	rnd, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 512))
	if err != nil {
		panic(err)
	}

	voitResult := big.NewInt(int64(Yes))
	n := new(big.Int).Lsh(rnd, 2)
	n.Or(n, voitResult)

	// 2. Generate r
	r, err := rand.Int(rand.Reader, serverKeys.N)
	if err != nil {
		panic(err)
	}

	gcd := new(big.Int).GCD(nil, nil, r, serverKeys.N)
	if gcd.Cmp(big.NewInt(1)) != 0 {
		panic("r and N are not coprime")
	}

	// 3. h = SHA3(n)
	hash := sha512.New()
	hash.Write(n.Bytes())
	h := new(big.Int).SetBytes(hash.Sum(nil))
	h.Mod(h, serverKeys.N) // Приводим к модулю N

	// 4. h' = h * r^d mod N
	rE := new(big.Int).Exp(r, serverKeys.D, serverKeys.N)
	hBlinded := new(big.Int).Mul(h, rE)
	hBlinded.Mod(hBlinded, serverKeys.N)

	// 5. Server помечает, что выдал биллютень
	// s' = h'^C mod N
	sBlinded := new(big.Int).Exp(hBlinded, serverKeys.C, serverKeys.N)

	// 6. s = s'*r-1 mod N
	rInv := new(big.Int).ModInverse(r, serverKeys.N)
	s := new(big.Int).Mul(sBlinded, rInv)
	s.Mod(s, serverKeys.N)

	// 7. SHA3(n) == s^d mod N
	check := new(big.Int).Exp(s, serverKeys.D, serverKeys.N)
	if check.Cmp(h) == 0 {
		fmt.Println("Valid blind signature!!")
	} else {
		fmt.Println("Invalid signature!!")
		fmt.Printf("Original hash: %s\n", h.String())
		fmt.Printf("Verified hash: %s\n", check.String())
	}
}
