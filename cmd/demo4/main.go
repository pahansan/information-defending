package main

import (
	"fmt"
	"information-defending/internal/crypto"
	"information-defending/internal/shamir"
	"log"
	"math/big"
	"os"
)

func main() {
	original := []byte("Секретное сообщение по Шамиру")
	err := os.WriteFile("input.txt", original, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p := big.NewInt(crypto.GeneratePInBounds(1000, 5000))
	ca, da := shamir.GenerateKeys(p)
	cb, db := shamir.GenerateKeys(p)

	fmt.Printf("p = %d\n", p)
	fmt.Printf("A: (ca=%d, da=%d)\n", ca, da)
	fmt.Printf("B:  (cb=%d, db=%d)\n", cb, db)

	err = shamir.EncryptFile("input.txt", "encrypted.txt", p, ca, cb)
	if err != nil {
		log.Fatal(err)
	}

	err = shamir.DecryptFile("encrypted.txt", "decrypted.txt", p, da, db)
	if err != nil {
		log.Fatal(err)
	}
}
