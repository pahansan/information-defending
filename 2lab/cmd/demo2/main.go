package main

import (
	"2lab/internal/crypto"
	"fmt"
)

func main() {
	fmt.Println(crypto.BSGS(5, 1, 7))
	fmt.Println(crypto.BSGS(3, 3, 11))
	fmt.Println(crypto.BSGS(7, 5, 17))
}
