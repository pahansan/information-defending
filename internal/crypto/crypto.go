package crypto

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"
)

func ModExp(a, x, p int64) int64 {
	var y int64 = 1

	for x != 0 {
		mod_a := a % p
		if x&1 == 1 {
			y *= mod_a
		}
		a = mod_a * mod_a
		x >>= 1
	}

	return y % p
}

func testFerma(a, p int64) bool {
	if a >= p {
		return false
	}

	return ModExp(a, p-1, p) == 1
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandInt64(min, max int64) int64 {
	return min + r.Int63n(max-min)
}

func IsProbablyPrime(x int64) bool {

	if x <= 1 || x%2 == 0 {
		return false
	}

	if x == 3 {
		return true
	}

	iters := int64(100)

	for i := int64(0); i < iters && i < x; i++ {
		a := RandInt64(2, x-1)
		if a%x == 0 {
			continue
		}
		if !testFerma(a, x) {
			return false
		}
	}

	return true
}

func ExtendedGCD(a, b int64) (int64, int64, int64) {
	if a < 1 || b < 1 {
		return 0, 0, 0
	}

	u1, u2, u3 := a, int64(1), int64(0)
	v1, v2, v3 := b, int64(0), int64(1)

	if a < b {
		u1, v1 = v1, u1
		u2, v2 = v2, u2
		u3, v3 = v3, u3
	}

	for v1 != 0 {
		q := u1 / v1

		t1, t2, t3 := u1%v1, u2-q*v2, u3-q*v3
		u1, u2, u3 = v1, v2, v3
		v1, v2, v3 = t1, t2, t3
	}

	return u1, u2, u3
}

func ExtendedGCDRandoms() (int64, int64, int64, int64, int64) {
	var a int64
	b := RandInt64(1, 1000)

	for a < b {
		a = RandInt64(1, 1000)
	}

	u1, u2, u3 := ExtendedGCD(a, b)

	return a, b, u1, u2, u3
}

func GeneratePrime(lb, ub int64) int64 {
	if lb < 2 || ub < 3 {
		return 0
	}

	x := RandInt64(lb, ub)
	for !IsProbablyPrime(x) {
		x = RandInt64(lb, ub)
	}

	return x
}

func ExtendedGCDPrimes() (int64, int64, int64, int64, int64) {
	var a int64
	b := GeneratePrime(2, 1000)

	for a < b {
		a = GeneratePrime(2, 1000)
	}

	u1, u2, u3 := ExtendedGCD(a, b)

	return a, b, u1, u2, u3
}

func BSGS(a, y, p int64) []int64 {
	m := int64(math.Ceil(math.Sqrt(float64(p))))
	valueMap := make(map[int64]int64)

	y = y % p
	for j := int64(0); j < m; j++ {
		ayModP := (ModExp(a, j, p) * y) % p
		valueMap[ayModP] = j
	}

	answer := make([]int64, 0)

	for i := int64(1); i <= m; i++ {
		amModP := ModExp(a, i*m, p)
		value, ok := valueMap[amModP]
		if ok {
			answer = append(answer, i*m-value)
		}
	}

	return answer
}

func RandBSGS() ([]int64, int64, int64, int64) {
	a := GeneratePrime(2, 1000)
	p := GeneratePrime(2, 1000)
	for a >= p {
		a = GeneratePrime(2, 1000)
	}
	y := RandInt64(1, p-1)
	result := BSGS(a, y, p)
	return result, a, y, p
}

func GenerateP() int64 {
	q := GeneratePrime(257, 1000)
	p := 2*q + 1

	for !IsProbablyPrime(p) {
		q = GeneratePrime(257, 1000)
		p = 2*q + 1
	}

	return p
}

func GenerateG(p int64) int64 {
	q := (p - 1) / 2
	g := RandInt64(2, p-1)
	for ModExp(g, q, p) == 1 {
		g = RandInt64(2, p-1)
	}
	return g
}

func DiffieHellman(p, g, a, b int64) int64 {
	A := ModExp(g, a, p)
	B := ModExp(g, b, p)

	Ka := ModExp(B, a, p)
	Kb := ModExp(A, b, p)

	if Ka == Kb {
		return Ka
	}
	return -1
}

func DiffieHellmanByte(p, g, a, b int64) byte {
	return byte(DiffieHellman(p, g, a, b))
}

func RandDiffieHellman() (int64, int64, int64, int64, int64) {
	p := GenerateP()
	g := GenerateG(p)
	a := RandInt64(2, 100)
	b := RandInt64(2, 100)
	for b == a {
		b = RandInt64(2, 100)
	}

	K := DiffieHellman(p, g, a, b)

	return p, g, a, b, K
}

func RandDiffieHellmanByte() (int64, int64, int64, int64, byte) {
	p, g, a, b, k := RandDiffieHellman()
	return p, g, a, b, byte(k)
}

func GeneratePInBounds(lb, ub int64) int64 {
	return GeneratePrime(lb, ub)
}

func Gcd(a, b *big.Int) *big.Int {
	zero := big.NewInt(0)
	for b.Cmp(zero) != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}
	return a
}

func ModInverse(a, m *big.Int) *big.Int {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	g.GCD(x, y, a, m)
	if g.Cmp(big.NewInt(1)) != 0 {
		return nil
	}
	return x.Mod(x, m)
}

func ReadBigInt(prompt string) *big.Int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		val, ok := new(big.Int).SetString(text, 10)
		if ok {
			return val
		}
		fmt.Println("Ошибка: введите целое число")
	}
}
