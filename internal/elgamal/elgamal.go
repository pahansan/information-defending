package elgamal

import (
	"crypto/rand"
	"math/big"
	"os"
	"strings"
)

type Keys struct {
	P *big.Int
	G *big.Int
	X *big.Int
	Y *big.Int
}

type Sign struct {
	R *big.Int
	S *big.Int
}

func GenerateP() (*big.Int, error) {
	q, err := rand.Prime(rand.Reader, 256)
	if err != nil {
		return nil, err
	}
	doubleQ := q.Mul(q, big.NewInt(2))
	doubleQPlus1 := doubleQ.Add(doubleQ, big.NewInt(1))
	for !doubleQPlus1.ProbablyPrime(20) {
		q, err := rand.Prime(rand.Reader, 256)
		if err != nil {
			return nil, err
		}
		doubleQ = q.Mul(q, big.NewInt(2))
		doubleQPlus1 = doubleQ.Add(doubleQ, big.NewInt(1))
	}
	return doubleQPlus1, nil
}

func GenerateG(p *big.Int) (*big.Int, error) {
	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))
	q := new(big.Int).Set(pMinus1)
	q = q.Div(q, big.NewInt(2))
	tmp, err := rand.Int(rand.Reader, pMinus1)
	g := new(big.Int).Set(tmp)
	if err != nil {
		return nil, err
	}
	one := big.NewInt(1)
	// while g <= 1 && g >= p - 1 && g^^q mod p == 1
	for tmp.Cmp(one) != 1 && tmp.Exp(tmp, q, p).Cmp(one) == 0 {
		tmp, err := rand.Int(rand.Reader, pMinus1)
		g.Set(tmp)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

func GenerateX(p *big.Int) (*big.Int, error) {
	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))
	x, err := rand.Int(rand.Reader, pMinus1)
	if err != nil {
		return nil, err
	}
	for x.Cmp(big.NewInt(0)) == 0 {
		x, err = rand.Int(rand.Reader, pMinus1)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func GenerateY(g, x, p *big.Int) *big.Int {
	y := new(big.Int).Exp(g, x, p)
	return y
}

func generateK(p *big.Int) (*big.Int, error) {
	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))
	k, err := rand.Int(rand.Reader, pMinus1)
	if err != nil {
		return nil, err
	}
	for k.Cmp(big.NewInt(0)) == 0 {
		k, err = rand.Int(rand.Reader, pMinus1)
		if err != nil {
			return nil, err
		}
	}
	return k, nil
}

func modInverseK(k, p *big.Int) *big.Int {
	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))
	inverseK := new(big.Int).Set(k)
	inverseK = inverseK.ModInverse(inverseK, pMinus1)
	return inverseK
}

func countR(g, k, p *big.Int) *big.Int {
	r := new(big.Int).Exp(g, k, p)
	return r
}

func countU(h, x, r, p *big.Int) *big.Int {
	xr := new(big.Int).Set(x)
	xr = xr.Mul(xr, r)

	hMinusXR := new(big.Int).Set(h)
	hMinusXR = hMinusXR.Sub(hMinusXR, xr)

	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))

	u := new(big.Int).Set(hMinusXR)
	u = u.Mod(u, pMinus1)
	return u
}

func countS(k, u, p *big.Int) *big.Int {
	inverseK := modInverseK(k, p)

	ku := new(big.Int).Set(inverseK)
	ku = ku.Mul(ku, u)

	pMinus1 := new(big.Int).Set(p)
	pMinus1 = pMinus1.Sub(pMinus1, big.NewInt(1))

	s := new(big.Int).Set(ku)
	s = s.Mod(s, pMinus1)

	return s
}

func CountSign(keys *Keys, h *big.Int) *Sign {
	k, _ := generateK(keys.P)
	r := countR(keys.G, k, keys.P)
	u := countU(h, keys.X, r, keys.P)
	s := countS(k, u, keys.P)
	return &Sign{R: r, S: s}
}

func CheckSign(sign *Sign, h, y, g, p *big.Int) bool {
	left := new(big.Int).Exp(y, sign.R, p)
	tmp := new(big.Int).Exp(sign.R, sign.S, p)
	left.Mul(left, tmp)
	left.Mod(left, p)

	right := new(big.Int).Exp(g, h, p)

	return left.Cmp(right) == 0
}

func GenerateKeys() (*Keys, error) {
	p, err := GenerateP()
	if err != nil {
		return nil, err
	}
	g, err := GenerateG(p)
	if err != nil {
		return nil, err
	}
	x, err := GenerateX(p)
	if err != nil {
		return nil, err
	}
	y := GenerateY(g, x, p)
	return &Keys{P: p, G: g, X: x, Y: y}, nil
}

func RandElGamal(p, g *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))

	Ci, _ := rand.Int(rand.Reader, pMinus1)
	Ci = Ci.Add(Ci, big.NewInt(1))

	Di := new(big.Int).Exp(g, Ci, p)
	return Ci, Di
}

func ElGamalEncrypt(p, g, Db, k, m *big.Int) (*big.Int, *big.Int) {
	r := new(big.Int).Exp(g, k, p)
	e := new(big.Int).Mul(m, new(big.Int).Exp(Db, k, p))
	e = e.Mod(e, p)

	return r, e
}

func ElGamalDecrypt(e, r, p, Cb *big.Int) *big.Int {
	pMinus1Cb := new(big.Int).Sub(p, big.NewInt(1))
	pMinus1Cb = pMinus1Cb.Sub(pMinus1Cb, Cb)

	m := new(big.Int).Mul(e, new(big.Int).Exp(r, pMinus1Cb, p))
	return m.Mod(m, p)
}

func EncryptFile(inputFile, outputFile string, p, g, Db, k *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	result := []*big.Int{}

	for _, b := range data {
		m := big.NewInt(int64(b))
		r, e := ElGamalEncrypt(p, g, Db, k, m)
		result = append(result, r)
		result = append(result, e)
	}

	out := ""
	for _, val := range result {
		out += val.String() + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func DecryptFile(inputFile, outputFile string, p, Cb *big.Int) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []byte

	flag := true
	var r, e *big.Int
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if flag {
			r, _ = new(big.Int).SetString(line, 10)
		} else {
			e, _ = new(big.Int).SetString(line, 10)
			dec := ElGamalDecrypt(e, r, p, Cb)
			result = append(result, byte(dec.Int64()))
		}
		flag = !flag
	}

	return os.WriteFile(outputFile, result, 0644)
}
