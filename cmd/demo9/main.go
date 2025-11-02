package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"information-defending/internal/elgamal"
	"log"
	"math/big"
	"os"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)

	keyFile := generateCmd.String("key", "elgamal_keys", "File to save Elgamal keys")

	signInput := signCmd.String("input", "", "Input file to sign")
	signOutput := signCmd.String("output", "", "Output signature file (default: input.sig)")
	signKey := signCmd.String("key", "elgamal_keys", "Elgamal private key file")

	verifyInput := verifyCmd.String("input", "", "Input file to verify")
	verifySig := verifyCmd.String("signature", "", "Signature file")
	verifyKey := verifyCmd.String("key", "elgamal_keys", "Elgamal public key file")

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
	fmt.Println("  generate - generate Elgamal keys")
	fmt.Println("  sign     - sign a file")
	fmt.Println("  verify   - verify a file signature")
	fmt.Println("\nUse [command] -h for more information about a command")
}

func generateKeys(keyFile string) {
	fmt.Println("Generating Elgamal keys...")
	keys, _ := elgamal.GenerateKeys()

	err := saveKeys(keys, keyFile)
	if err != nil {
		log.Fatalf("Error saving keys: %v", err)
	}

	fmt.Printf("Keys saved to %s.pub and %s.priv\n", keyFile, keyFile)
	fmt.Printf("Public key (P): %s\n", keys.P.String())
	fmt.Printf("Public key (G): %s\n", keys.G.String())
	fmt.Printf("Public key (Y): %s\n", keys.Y.String())
	fmt.Printf("Private key (X): %s\n", keys.X.String())
}

func signFile(inputFile, outputFile, keyFile string) {
	fmt.Printf("Signing file: %s\n", inputFile)

	privKey, err := loadPrivateKey(keyFile + ".priv")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	pubKey, err := loadPublicKey(keyFile + ".pub")
	if err != nil {
		log.Fatalf("Error loading public key: %v", err)
	}

	keys := elgamal.Keys{P: pubKey.P, G: pubKey.G, Y: pubKey.Y, X: privKey.X}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	hash := sha256.Sum256(data)
	h := new(big.Int).SetBytes(hash[:])

	sign := elgamal.CountSign(&keys, h)

	signData := fmt.Sprintf("%s\n%s", sign.R.String(), sign.S.String())
	err = os.WriteFile(outputFile, []byte(signData), 0644)
	if err != nil {
		log.Fatalf("Error writing signature: %v", err)
	}

	fmt.Printf("Signature saved to: %s\n", outputFile)
	fmt.Printf("Signature:\nR: %s\nS: %s\n", sign.R.String(), sign.S.String())
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

	sign, err := loadSign(signatureFile)
	if err != nil {
		log.Fatalf("Error reading signature: %v", err)
	}

	hash := sha256.Sum256(data)
	h := big.NewInt(0).SetBytes(hash[:])

	ok := elgamal.CheckSign(sign, h, pubKey.Y, pubKey.G, pubKey.P)

	if ok {
		fmt.Println("✓ Signature is VALID")
	} else {
		fmt.Println("✗ Signature is INVALID")
	}
}

type PublicKey struct {
	P *big.Int
	G *big.Int
	Y *big.Int
}

type PrivateKey struct {
	X *big.Int
}

func saveKeys(keys *elgamal.Keys, baseName string) error {

	pubData := fmt.Sprintf("%s\n%s\n%s", keys.P.String(), keys.G.String(), keys.Y.String())
	err := os.WriteFile(baseName+".pub", []byte(pubData), 0644)
	if err != nil {
		return err
	}

	privData := fmt.Sprintf("%s", keys.X.String())
	err = os.WriteFile(baseName+".priv", []byte(privData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadPublicKey(filename string) (*PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var pStr, gStr, yStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s\n%s", &pStr, &gStr, &yStr)
	if err != nil {
		return nil, err
	}

	p, _ := new(big.Int).SetString(pStr, 10)
	g, _ := new(big.Int).SetString(gStr, 10)
	y, _ := new(big.Int).SetString(yStr, 10)

	return &PublicKey{P: p, G: g, Y: y}, nil
}

func loadPrivateKey(filename string) (*PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var xStr string
	_, err = fmt.Sscanf(string(data), "%s", &xStr)
	if err != nil {
		return nil, err
	}

	x, _ := new(big.Int).SetString(xStr, 10)

	return &PrivateKey{X: x}, nil
}

func loadSign(filename string) (*elgamal.Sign, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rStr, sStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s", &rStr, &sStr)
	if err != nil {
		return nil, err
	}

	r, _ := new(big.Int).SetString(rStr, 10)
	s, _ := new(big.Int).SetString(sStr, 10)

	return &elgamal.Sign{R: r, S: s}, nil
}
