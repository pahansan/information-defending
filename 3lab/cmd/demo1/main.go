package main

import (
	"3lab/internal/crypto"
	"fmt"
)

func main() {
	fmt.Printf("1. Функция быстрого возведения числа в степень по модулю:\n")
	fmt.Printf("\t%d^%d mod %2d = %d\n", 5, 12, 7, crypto.ModExp(5, 12, 7))
	fmt.Printf("\t%d^%d mod %2d = %d\n", 3, 21, 11, crypto.ModExp(3, 21, 11))
	fmt.Printf("\t%d^%d mod %2d = %d\n", 7, 31, 17, crypto.ModExp(7, 31, 17))

	fmt.Printf("\n2. Тест ферма:\n")
	fmt.Printf("\t%4d is probably prime: %t\n", 3, crypto.IsProbablyPrime(3))
	fmt.Printf("\t%4d is probably prime: %t\n", 2377, crypto.IsProbablyPrime(2377))
	fmt.Printf("\t%4d is probably prime: %t\n", 10, crypto.IsProbablyPrime(10))
	fmt.Printf("\t%4d is probably prime: %t\n", 11, crypto.IsProbablyPrime(11))

	gcd, x, y := crypto.ExtendedGCD(10, 35)
	fmt.Printf("\n3. Расширенный алгоритм Евклида:\n")
	fmt.Printf("\ta = 10, b = 35:\n")
	fmt.Printf("\tgcd(a, b) = %d, x = %d, y = %d\n", gcd, x, y)

	a, b, gcd, x, y := crypto.ExtendedGCDRandoms()
	fmt.Printf("\n\tRandom numbers:\n")
	fmt.Printf("\ta = %d, b = %d:\n", a, b)
	fmt.Printf("\tgcd(a, b) = %d, x = %d, y = %d\n", gcd, x, y)

	a, b, gcd, x, y = crypto.ExtendedGCDPrimes()
	fmt.Printf("\n\tProbably prime numbers:\n")
	fmt.Printf("\ta = %d, b = %d:\n", a, b)
	fmt.Printf("\tgcd(a, b) = %d, x = %d, y = %d\n", gcd, x, y)

	fmt.Printf("\n\tYour numbers:\n\tEnter a, b: ")
	fmt.Scan(&a)
	fmt.Scan(&b)
	gcd, x, y = crypto.ExtendedGCD(a, b)
	fmt.Printf("\tgcd(a, b) = %d, x = %d, y = %d\n", gcd, x, y)
}
