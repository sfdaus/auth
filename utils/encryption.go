package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"prakarsa-app/config"
)

func generateKey() []byte {
	hash := sha256.Sum256([]byte(config.LoadConfig().EncryptionSecretKey))
	return hash[:]
}

func EncryptDeterministic(plainText string) (string, error) {
	key := generateKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCFBEncrypter(block, iv)

	plainBytes := []byte(plainText)
	cipherBytes := make([]byte, len(plainBytes))
	stream.XORKeyStream(cipherBytes, plainBytes)

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func DecryptDeterministic(cipherText string) (string, error) {
	key := generateKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCFBDecrypter(block, iv)

	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	plainBytes := make([]byte, len(cipherBytes))
	stream.XORKeyStream(plainBytes, cipherBytes)

	return string(plainBytes), nil
}
