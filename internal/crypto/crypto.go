package crypto

import (
	"math"
	"math/rand"
	"os"
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

func randInt64(min, max int64) int64 {
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
		a := randInt64(2, x-1)
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
	b := randInt64(1, 1000)

	for a < b {
		a = randInt64(1, 1000)
	}

	u1, u2, u3 := ExtendedGCD(a, b)

	return a, b, u1, u2, u3
}

func generatePrime(lb, ub int64) int64 {
	if lb < 2 || ub < 3 {
		return 0
	}

	x := randInt64(lb, ub)
	for !IsProbablyPrime(x) {
		x = randInt64(lb, ub)
	}

	return x
}

func ExtendedGCDPrimes() (int64, int64, int64, int64, int64) {
	var a int64
	b := generatePrime(2, 1000)

	for a < b {
		a = generatePrime(2, 1000)
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
	a := generatePrime(2, 1000)
	p := generatePrime(2, 1000)
	for a >= p {
		a = generatePrime(2, 1000)
	}
	y := randInt64(1, p-1)
	result := BSGS(a, y, p)
	return result, a, y, p
}

func generateP() int64 {
	q := generatePrime(2, 1000)
	p := 2*q + 1

	for !IsProbablyPrime(p) {
		q = generatePrime(2, 1000)
		p = 2*q + 1
	}

	return p
}

func generateG(p int64) int64 {
	q := (p - 1) / 2
	g := randInt64(2, p-1)
	for ModExp(g, q, p) == 1 {
		g = randInt64(2, p-1)
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

func RandDiffieHellman() (int64, int64, int64, int64, int64) {
	p := generateP()
	g := generateG(p)
	a := randInt64(2, 100)
	b := randInt64(2, 100)
	for b == a {
		b = randInt64(2, 100)
	}

	K := DiffieHellman(p, g, a, b)

	return p, g, a, b, K
}

func ModExpArray(a []byte, x, p int64) []byte {
	result := make([]byte, len(a))
	for i, b := range a {
		result[i] = byte(ModExp(int64(b), x, p))
	}
	return result
}

func GeneratePInBounds(lb, ub int64) int64 {
	return generatePrime(lb, ub)
}

func GenerateShamirKeys(p int64) (int64, int64, int64, int64) {
	if p < 2 {
		return 0, 0, 0, 0
	}

	ca := generatePrime(2, 1000)
	cb := generatePrime(2, 1000)
	_, da, _ := ExtendedGCD(ca, p-1)
	_, db, _ := ExtendedGCD(cb, p-1)

	return ca, da, cb, db
}

func ShamirEnDeCrypt(k1, k2, p int64, m []byte) []byte {
	if p <= int64(^uint8(0)) {
		return nil
	}

	x1 := ModExpArray(m, k1, p)
	return ModExpArray(x1, k1, p)
}

func ShamirEnDeCryptFile(input, output string, k1, k2, p int64) error {
	fileBytes, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	data := ShamirEnDeCrypt(k1, k2, p, fileBytes)
	err = os.WriteFile(output, data, 0644)
	return err
}
