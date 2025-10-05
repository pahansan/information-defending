package main

import (
	"fmt"
	"information-defending/internal/crypto"
	"information-defending/internal/elgamal"
	"log"
	"math/big"
	"os"
)

func main() {
	original := []byte("Зашифрованное сообщение Эль-Гамаля")
	err := os.WriteFile("input.txt", original, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// p и g генерятся как в системе Диффи-Хеллмана
	p := big.NewInt(crypto.GenerateP())
	g := big.NewInt(crypto.GenerateG(p.Int64()))

	// Секретный (Cb) и открытый (Db) ключи абонента B
	Cb, Db := elgamal.RandElGamal(p, g)

	fmt.Printf("p = %d g = %d\n", p, g)
	fmt.Printf("B: (cb=%d, db=%d)\n", Cb, Db)

	// Абонент A генерит случайное число k [2, p-1)
	k := big.NewInt(crypto.RandInt64(2, p.Int64()-1))
	fmt.Printf("k = %d\n", k)

	err = elgamal.EncryptFile("input.txt", "encrypted.txt", p, g, Db, k)
	if err != nil {
		log.Fatal(err)
	}

	err = elgamal.DecryptFile("encrypted.txt", "decrypted.txt", p, Cb)
	if err != nil {
		log.Fatal(err)
	}
}
