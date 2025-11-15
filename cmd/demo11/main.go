package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"information-defending/internal/gost"
	"log"
	"math/big"
	"os"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)

	keyFile := generateCmd.String("key", "fips", "File to save FIPS keys")

	signInput := signCmd.String("input", "", "Input file to sign")
	signOutput := signCmd.String("output", "", "Output signature file (default: input.sig)")
	signKey := signCmd.String("key", "fips", "FIPS private key file")

	verifyInput := verifyCmd.String("input", "", "Input file to verify")
	verifySig := verifyCmd.String("signature", "", "Signature file")
	verifyKey := verifyCmd.String("key", "fips", "FIPS public key file")

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "generate":
		generateCmd.Parse(os.Args[2:])
		generateKeys(*keyFile)
	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signInput == "" {
			fmt.Println("Error: input file is required")
			signCmd.PrintDefaults()
			os.Exit(1)
		}
		if *signOutput == "" {
			*signOutput = *signInput + ".sig"
		}
		signFile(*signInput, *signOutput, *signKey)
	case "verify":
		verifyCmd.Parse(os.Args[2:])
		if *verifyInput == "" || *verifySig == "" {
			fmt.Println("Error: input file and signature file are required")
			verifyCmd.PrintDefaults()
			os.Exit(1)
		}
		verifySignature(*verifyInput, *verifySig, *verifyKey)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  generate - generate FIPS keys")
	fmt.Println("  sign     - sign a file")
	fmt.Println("  verify   - verify a file signature")
	fmt.Println("\nUse [command] -h for more information about a command")
}

func generateKeys(keyFile string) {
	fmt.Println("Generating FIPS keys...")
	keys := gost.GenerateKeys(160)

	err := saveKeys(keys, keyFile)
	if err != nil {
		log.Fatalf("Error saving keys: %v", err)
	}

	fmt.Printf("Keys saved to %s.pub and %s.priv\n", keyFile, keyFile)
	fmt.Printf("Prime Q (%d bits): %s\n", keys.Q.BitLen(), keys.Q.String())
	fmt.Printf("Prime P (%d bits): %s\n", keys.P.BitLen(), keys.P.String())
	fmt.Printf("A: %s\n", keys.A.String())
	fmt.Printf("Private key X: %s\n", keys.X.String())
	fmt.Printf("Public key Y: %s\n", keys.Y.String())
}

func signFile(inputFile, outputFile, keyFile string) {
	fmt.Printf("Signing file: %s\n", inputFile)

	privKey, err := loadPrivateKey(keyFile + ".priv")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	hash := sha1.Sum([]byte(data))
	h := new(big.Int).SetBytes(hash[:])

	// Gen (r, s)
	var r, s *big.Int
	for {
		k, err := gost.GenerateLessThanNotZero(privKey.Q)
		if err != nil {
			log.Fatalf("Something went wrong: %s", err.Error())
		}

		// r = (a^k mod p) mod q
		r = new(big.Int).Exp(privKey.A, k, privKey.P)
		r = r.Mod(r, privKey.Q)
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		// s = k^(-1) (h + xr) mod q
		k_inv := new(big.Int).ModInverse(k, privKey.Q)
		xr := new(big.Int).Mul(privKey.X, r)
		h_xr := new(big.Int).Add(h, xr)
		s = new(big.Int).Mul(k_inv, h_xr)
		s.Mod(s, privKey.Q)

		if s.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		break
	}

	// Save (r, s)
	signature := fmt.Sprintf("%s\n%s", r.String(), s.String())
	err = os.WriteFile(outputFile, []byte(signature), 0644)
	if err != nil {
		log.Fatalf("Error writing signature: %v", err)
	}

	fmt.Printf("Signature saved to: %s\n", outputFile)
	fmt.Printf("Signature (r, s):\n  r: %s\n  s: %s\n", r.String(), s.String())
}

func verifySignature(inputFile, signatureFile, keyFile string) {
	fmt.Printf("Verifying file: %s\n", inputFile)

	pubKey, err := loadPublicKey(keyFile + ".pub")
	if err != nil {
		log.Fatalf("Error loading public key: %v", err)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	signatureData, err := os.ReadFile(signatureFile)
	if err != nil {
		log.Fatalf("Error reading signature: %v", err)
	}

	// Get (r, s)
	var rStr, sStr string
	_, err = fmt.Sscanf(string(signatureData), "%s\n%s", &rStr, &sStr)
	if err != nil {
		log.Fatalf("Error parsing signature: %v", err)
	}

	r := big.NewInt(0)
	r.SetString(rStr, 10)
	s := big.NewInt(0)
	s.SetString(sStr, 10)

	// 2.
	// 	0 < r < q
	// 	0 < s < q
	if r.Cmp(pubKey.Q) >= 0 || r.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("✗ Signature is INVALID: r is out of range")
		return
	}
	if s.Cmp(pubKey.Q) >= 0 || s.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("✗ Signature is INVALID: s is out of range")
		return
	}

	// Hash
	hash := sha1.Sum([]byte(data))
	h := new(big.Int).SetBytes(hash[:])

	// 3.
	// 	u1 = h * s^-1 mod q
	//  u2 = r * s^-1 mod q
	sInv := new(big.Int).ModInverse(s, pubKey.Q)

	u1 := new(big.Int).Mul(h, sInv)
	u1.Mod(u1, pubKey.Q)

	u2 := new(big.Int).Mul(r, sInv)
	u2.Mod(u2, pubKey.Q)

	// 4.
	// 	v = ((a^u1 * y^u2) mod p) mod q
	au1 := new(big.Int).Exp(pubKey.A, u1, pubKey.P)
	yu2 := new(big.Int).Exp(pubKey.Y, u2, pubKey.P)

	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, pubKey.P)
	v.Mod(v, pubKey.Q)

	// 5.
	// 	v == r
	if v.Cmp(r) == 0 {
		fmt.Println("✓ Signature is VALID")
	} else {
		fmt.Println("✗ Signature is INVALID")
		fmt.Printf("Computed v: %s\n", v.String())
		fmt.Printf("Signature r: %s\n", r.String())
	}
}

type GOSTPublicKey struct {
	Q *big.Int
	P *big.Int
	A *big.Int
	Y *big.Int
}

type GOSTPrivateKey struct {
	Q *big.Int
	P *big.Int
	A *big.Int
	X *big.Int
	Y *big.Int
}

func saveKeys(keys gost.Keys, baseName string) error {
	pubData := fmt.Sprintf("%s\n%s\n%s\n%s",
		keys.Q.String(),
		keys.P.String(),
		keys.A.String(),
		keys.Y.String())
	err := os.WriteFile(baseName+".pub", []byte(pubData), 0644)
	if err != nil {
		return err
	}

	privData := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		keys.Q.String(),
		keys.P.String(),
		keys.A.String(),
		keys.X.String(),
		keys.Y.String())
	err = os.WriteFile(baseName+".priv", []byte(privData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadPublicKey(filename string) (*GOSTPublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var qStr, pStr, aStr, yStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s\n%s\n%s", &qStr, &pStr, &aStr, &yStr)
	if err != nil {
		return nil, err
	}

	Q := big.NewInt(0)
	Q.SetString(qStr, 10)

	P := big.NewInt(0)
	P.SetString(pStr, 10)

	A := big.NewInt(0)
	A.SetString(aStr, 10)

	Y := big.NewInt(0)
	Y.SetString(yStr, 10)

	return &GOSTPublicKey{Q: Q, P: P, A: A, Y: Y}, nil
}

func loadPrivateKey(filename string) (*GOSTPrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var qStr, pStr, aStr, xStr, yStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s\n%s\n%s\n%s", &qStr, &pStr, &aStr, &xStr, &yStr)
	if err != nil {
		return nil, err
	}

	Q := big.NewInt(0)
	Q.SetString(qStr, 10)

	P := big.NewInt(0)
	P.SetString(pStr, 10)

	A := big.NewInt(0)
	A.SetString(aStr, 10)

	X := big.NewInt(0)
	X.SetString(xStr, 10)

	Y := big.NewInt(0)
	Y.SetString(yStr, 10)

	return &GOSTPrivateKey{Q: Q, P: P, A: A, X: X, Y: Y}, nil
}
