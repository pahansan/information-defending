package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"information-defending/internal/rsa"
	"log"
	"math/big"
	"os"
)

func main() {
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)

	keyFile := generateCmd.String("key", "rsa_keys", "File to save RSA keys")

	signInput := signCmd.String("input", "", "Input file to sign")
	signOutput := signCmd.String("output", "", "Output signature file (default: input.sig)")
	signKey := signCmd.String("key", "rsa_keys", "RSA private key file")

	verifyInput := verifyCmd.String("input", "", "Input file to verify")
	verifySig := verifyCmd.String("signature", "", "Signature file")
	verifyKey := verifyCmd.String("key", "rsa_keys", "RSA public key file")

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
	fmt.Println("  generate - generate RSA keys")
	fmt.Println("  sign     - sign a file")
	fmt.Println("  verify   - verify a file signature")
	fmt.Println("\nUse [command] -h for more information about a command")
}

func generateKeys(keyFile string) {
	fmt.Println("Generating RSA keys...")
	keys := rsa.GenerateKeys()

	err := saveKeys(keys, keyFile)
	if err != nil {
		log.Fatalf("Error saving keys: %v", err)
	}

	fmt.Printf("Keys saved to %s.pub and %s.priv\n", keyFile, keyFile)
	fmt.Printf("Public key (N): %s\n", keys.N.String())
	fmt.Printf("Public exponent (D): %s\n", keys.D.String())
	fmt.Printf("Private exponent (C): %s\n", keys.C.String())
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

	hash := sha256.Sum256(data)
	y := big.NewInt(0).SetBytes(hash[:])

	s := rsa.Decrypt(y, privKey.C, privKey.N)

	signatureBytes := s.Bytes()
	err = os.WriteFile(outputFile, signatureBytes, 0644)
	if err != nil {
		log.Fatalf("Error writing signature: %v", err)
	}

	fmt.Printf("Signature saved to: %s\n", outputFile)
	fmt.Printf("Signature: %x\n", signatureBytes)
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

	signatureBytes, err := os.ReadFile(signatureFile)
	if err != nil {
		log.Fatalf("Error reading signature: %v", err)
	}

	hash := sha256.Sum256(data)
	y := big.NewInt(0).SetBytes(hash[:])

	s := big.NewInt(0).SetBytes(signatureBytes)

	w := rsa.Encrypt(s, pubKey.D, pubKey.N)

	if w.Cmp(y) == 0 {
		fmt.Println("✓ Signature is VALID")
	} else {
		fmt.Println("✗ Signature is INVALID")
		fmt.Printf("Expected hash: %x\n", y.Bytes())
		fmt.Printf("Recovered hash: %x\n", w.Bytes())
	}
}

type PublicKey struct {
	N *big.Int
	D *big.Int
}

type PrivateKey struct {
	N *big.Int
	C *big.Int
}

func saveKeys(keys rsa.Keys, baseName string) error {

	pubData := fmt.Sprintf("%s\n%s", keys.N.String(), keys.D.String())
	err := os.WriteFile(baseName+".pub", []byte(pubData), 0644)
	if err != nil {
		return err
	}

	privData := fmt.Sprintf("%s\n%s", keys.N.String(), keys.C.String())
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

	var nStr, dStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s", &nStr, &dStr)
	if err != nil {
		return nil, err
	}

	N := big.NewInt(0)
	N.SetString(nStr, 10)

	D := big.NewInt(0)
	D.SetString(dStr, 10)

	return &PublicKey{N: N, D: D}, nil
}

func loadPrivateKey(filename string) (*PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var nStr, cStr string
	_, err = fmt.Sscanf(string(data), "%s\n%s", &nStr, &cStr)
	if err != nil {
		return nil, err
	}

	N := big.NewInt(0)
	N.SetString(nStr, 10)

	C := big.NewInt(0)
	C.SetString(cStr, 10)

	return &PrivateKey{N: N, C: C}, nil
}
