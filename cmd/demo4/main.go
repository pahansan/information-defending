package main

import (
	"fmt"
	"information-defending/internal/crypto"
	"log"
	"os"
)

func main() {
	original := []byte("Секретное сообщение по Шамиру")
	err := os.WriteFile("input.txt", original, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p := crypto.GeneratePInBounds(1000, 5000)
	ca, da, cb, db := crypto.GenerateShamirKeys(p)

	fmt.Printf("p = %d\n", p)
	fmt.Printf("A: (ca=%d, da=%d)\n", ca, da)
	fmt.Printf("B:  (cb=%d, db=%d)\n", cb, db)

	err = crypto.ShamirEnDeCryptFile("input.txt", "encrypted.txt", ca, cb, p)
	if err != nil {
		log.Fatal(err)
	}

	err = crypto.ShamirEnDeCryptFile("encrypted.txt", "result.txt", da, db, p)
	if err != nil {
		log.Fatal(err)
	}
}
