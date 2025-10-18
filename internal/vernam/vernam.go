package vernam

import (
	"fmt"
	"os"
	"strings"
)

func Encrypt(m, k byte) byte {
	return m ^ k
}

func Decrypt(m, k byte) byte {
	return m ^ k
}

func EncryptFile(inputFile, outputFile string, k byte) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	result := []byte{}

	for _, b := range data {
		e := Encrypt(b, k)
		result = append(result, e)
	}

	out := ""
	for _, val := range result {
		out += fmt.Sprintf("%d", val) + "\n"
	}
	return os.WriteFile(outputFile, []byte(out), 0644)
}

func DecryptFile(inputFile, outputFile string, k byte) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []byte
	var e byte

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Sscanf(line, "%d", &e)
		m := Decrypt(e, k)
		result = append(result, m)
	}

	return os.WriteFile(outputFile, result, 0644)
}
