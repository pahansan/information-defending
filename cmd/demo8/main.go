package main

import (
	"crypto/sha256"
	"fmt"
	"information-defending/internal/rsa"
	"math/big"
)

func main() {
	m := []byte("Текст, который необходимо подписать")

	// Подписываем
	A := rsa.GenerateKeys()
	hash := sha256.Sum256(m)
	y := big.NewInt(0).SetBytes(hash[:])
	s := rsa.Decrypt(y, A.C, A.N)
	fmt.Println("Подпись:", s)

	// Проверяем подпись
	w := rsa.Encrypt(s, A.D, A.N)
	if w.Cmp(y) == 0 {
		fmt.Println("Подпись правильная")
	} else {
		fmt.Println("Подпись неправильная")
	}
}
