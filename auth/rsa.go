package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"strconv"
)

func generateKeyPair() (*rsa.PrivateKey, error) {
	bits, err := strconv.Atoi(os.Getenv("RSA_KEY_SIZE"))
	if err != nil {
		return nil, err
	}
	keyPair, err := rsa.GenerateKey(rand.Reader, bits)
	return keyPair, err
}
