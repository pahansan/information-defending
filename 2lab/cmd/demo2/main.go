package main

import (
	"2lab/internal/crypto"
	"fmt"
)

func main() {
	a, y, p := int64(5), int64(1), int64(7)
	fmt.Printf("%d^x %% %d = %d, x = %d\n", a, p, y, crypto.BSGS(a, y, p))
	a, y, p = int64(3), int64(3), int64(11)
	fmt.Printf("%d^x %% %d = %d, x = %d\n", a, p, y, crypto.BSGS(a, y, p))
	a, y, p = int64(7), int64(5), int64(17)
	fmt.Printf("%d^x %% %d = %d, x = %d\n", a, p, y, crypto.BSGS(a, y, p))

	result, a, y, p := crypto.RandBSGS()
	fmt.Printf("%d^x %% %d = %d, x = %d\n", a, p, y, result)

	fmt.Printf("Your a, y, p: ")
	fmt.Scan(&a)
	fmt.Scan(&y)
	fmt.Scan(&p)
	fmt.Printf("%d^x %% %d = %d, x = %d\n", a, p, y, crypto.BSGS(a, y, p))
}
