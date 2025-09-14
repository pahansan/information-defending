package main

import (
	"2lab/pkg/crypto"
	"fmt"
)

func main() {
	fmt.Println("aboba")
	fmt.Println(crypto.BSGS(7, 5, 17))
}
