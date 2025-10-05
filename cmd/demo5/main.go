package main

import (
	"fmt"
	"information-defending/internal/crypto"
	"information-defending/internal/elgamal"
	"math/big"
)

func main() {
	// Исходное сообщение (байт)
	m := big.NewInt(127)

	// p и g генерятся как в системе Диффи-Хеллмана
	p := big.NewInt(crypto.GenerateP())
	g := big.NewInt(crypto.GenerateG(p.Int64()))

	// Секретный (Cb) и открытый (Db) ключи абонента B
	Cb, Db := elgamal.RandElGamal(p, g)

	fmt.Printf("p = %d g = %d\n", p, g)
	fmt.Printf("B: (cb=%d, db=%d)\n", Cb, Db)

	// Абонент A генерит случайно число k [2, p-1)
	k := big.NewInt(crypto.RandInt64(2, p.Int64()-1))
	fmt.Printf("k = %d\n", k)

	// Абонент A шифрует сообщение (1 байт) и получает пару чисел (r, e)
	r, e := elgamal.ElGamalEncrypt(p, g, Db, k, m)
	fmt.Printf("(r, e) = (%d, %d)\n", r, e)

	// Абонент B дешифрует сообщение (из 2 пар чисел)
	// и получает исходное сообщение (байт)
	m_result := elgamal.ElGamalDecrypt(e, r, p, Cb)
	fmt.Printf("m' = %d\n", m_result)
}
