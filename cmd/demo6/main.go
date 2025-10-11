package main

import (
	"fmt"
	"information-defending/internal/rsa"
	"log"
	"os"
)

func main() {
	original := []byte("Зашифрованное сообщение RSA")
	err := os.WriteFile("input.txt", original, 0644)
	if err != nil {
		log.Fatal(err)
	}

	A := rsa.GenerateKeys()
	B := rsa.GenerateKeys()

	fmt.Printf("c_A = %d, d_A = %d, N_A = %d\n", A.C.Int64(), A.D.Int64(), A.N.Int64())
	fmt.Printf("c_B = %d, d_B = %d, N_B = %d\n", B.C.Int64(), B.D.Int64(), B.N.Int64())

	err = rsa.EncryptFile("input.txt", "encrypted.txt", B.D, B.N)
	if err != nil {
		log.Fatal(err)
	}

	err = rsa.DecryptFile("encrypted.txt", "decrypted.txt", B.C, B.N)
	if err != nil {
		log.Fatal(err)
	}
}
