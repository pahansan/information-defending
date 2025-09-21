package main

import (
	"3lab/internal/crypto"
	"fmt"
)

func main() {
	p, g, a, b, K := crypto.RandDiffieHellman()

	fmt.Printf("p = %d, g = %d, a = %d, b = %d, K = %d\n", p, g, a, b, K)

	fmt.Println("Введите p, g, a, b")
	fmt.Scan(&p, &g, &a, &b)
	K = crypto.DiffieHellman(p, g, a, b)

	fmt.Printf("p = %d, g = %d, a = %d, b = %d, K = %d\n", p, g, a, b, K)
}
