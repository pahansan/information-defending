package gost

import (
	"crypto/rand"
	"log"
	"math/big"
)

type Keys struct {
	Q *big.Int // prime
	P *big.Int // p = bq + 1, prime
	A *big.Int // a^q mod p = 1

	X *big.Int // secret key
	Y *big.Int // public key
}

func generatePrime(bitSize int) (*big.Int, error) {
	for {
		prime, err := rand.Prime(rand.Reader, bitSize)
		if err != nil {
			return nil, err
		}

		if prime.Bit(bitSize-1) == 1 {
			return prime, nil
		}
	}
}

func generatePfromQ(q *big.Int, bitSize int) (*big.Int, error) {
	for {
		// generate b
		bBitSize := bitSize - q.BitLen() + 10
		b, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bBitSize)))
		if err != nil {
			return nil, err
		}

		// p = b*q + 1
		p := new(big.Int)
		p.Mul(b, q)
		p.Add(p, big.NewInt(1))

		if p.BitLen() != bitSize {
			continue
		}

		if p.Bit(bitSize-1) != 1 {
			continue
		}

		if p.ProbablyPrime(20) {
			return p, nil
		}
	}
}

func generateA(p, q *big.Int) (*big.Int, error) {
	for {
		// need: 1 < g < p
		// 0 <= q < p - 2
		g, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
		if err != nil {
			return nil, err
		}
		// 1 < g < p
		g.Add(g, big.NewInt(2))

		// (p - 1) / q
		t := new(big.Int).Sub(p, big.NewInt(1))
		t.Div(t, q)

		// a = g^((p-1)/q) mod p
		a := new(big.Int).Exp(g, t, p)
		if a.Cmp(big.NewInt(1)) == 1 {
			return a, nil
		}
	}
}

func GenerateLessThanNotZero(q *big.Int) (*big.Int, error) {
	// need: 0 < a < q
	// 0 <= a < q - 1
	a, err := rand.Int(rand.Reader, new(big.Int).Sub(q, big.NewInt(1)))
	if err != nil {
		return nil, err
	}
	// 0 < a < q
	a.Add(a, big.NewInt(1))

	return a, nil
}

func GenerateKeys() Keys {
	Q, err := generatePrime(256)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
	}
	P, err := generatePfromQ(Q, 1024)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
	}
	A, err := generateA(P, Q)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
	}
	X, err := GenerateLessThanNotZero(Q)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err.Error())
	}
	Y := new(big.Int).Exp(A, X, P)

	return Keys{Q: Q,
		P: P,
		A: A,
		X: X,
		Y: Y,
	}
}
