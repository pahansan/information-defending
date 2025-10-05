package elgamal

import (
	"crypto/rand"
	"math/big"
)

func RandElGamal(p, g *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))

	Ci, _ := rand.Int(rand.Reader, pMinus1)
	Ci = Ci.Add(Ci, pMinus1)

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
