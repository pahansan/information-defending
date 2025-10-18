package main

import (
	"fmt"
	"information-defending/internal/crypto"
	"information-defending/internal/vernam"
	"log"
	"os"
)

func main() {
	original := []byte("Зашифрованное сообщение Вернама")
	err := os.WriteFile("input.txt", original, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p, g, a, b, k := crypto.RandDiffieHellmanByte()

	fmt.Printf("p = %d, g = %d, a = %d, b = %d, K = %d\n", p, g, a, b, k)

	err = vernam.EncryptFile("input.txt", "encrypted.txt", k)
	if err != nil {
		log.Fatal(err)
	}

	err = vernam.DecryptFile("encrypted.txt", "decrypted.txt", k)
	if err != nil {
		log.Fatal(err)
	}

}
