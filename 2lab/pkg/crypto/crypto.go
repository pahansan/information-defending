package crypto

import (
	"fmt"
	"math"
	"math/rand"
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

func generatePrime() int64 {
	x := randInt64(2, 1000)

	for !IsProbablyPrime(x) {
		x = randInt64(2, 1000)
	}

	return x
}

func ExtendedGCDPrimes() (int64, int64, int64, int64, int64) {
	var a int64
	b := generatePrime()

	for a < b {
		a = generatePrime()
	}

	u1, u2, u3 := ExtendedGCD(a, b)

	return a, b, u1, u2, u3
}

func BSGS(a, y, p int64) int64 {
	m := int64(math.Ceil(math.Sqrt(float64(p))))
	k := m

	m_map := make(map[int64]int64)
	ay := y
	m_map[ay%p] = 0
	fmt.Printf("j = 0, ay_mod_p = %d\n", ay%p)

	for j := int64(1); j < m; j++ {
		ay *= a
		ay_mod_p := ay % p
		m_map[ay_mod_p] = j
		fmt.Printf("j = %d, ay_mod_p = %d\n", j, ay_mod_p)
	}

	am := ay * a / y
	fmt.Println()

	for i := int64(1); i <= k; i++ {
		am_mod_p := am % p
		am *= a
		_, ok := m_map[am_mod_p]
		fmt.Printf("i = %d, am_mod_p = %d\n", i, am_mod_p)
		if ok {
			return i*m - m_map[am_mod_p]
		}
	}
	return -1
}
