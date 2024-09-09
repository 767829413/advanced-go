package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func main() {
	navicatPassword := NewNavicatPassword(12)
	encryptedPassword := ""
	decode, err := navicatPassword.Decrypt(encryptedPassword)
	if err != nil {
		fmt.Printf("Error decrypting password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(decode)
}

type NavicatPassword struct {
	version    int
	aesKey     []byte
	aesIv      []byte
	blowString string
	blowKey    []byte
	blowIv     []byte
}

func NewNavicatPassword(version int) *NavicatPassword {
	np := &NavicatPassword{
		version:    version,
		aesKey:     []byte("libcckeylibcckey"),
		aesIv:      []byte("libcciv libcciv "),
		blowString: "3DC5CA39",
	}
	np.blowKey = sha1.New().Sum([]byte(np.blowString))[:20]
	np.blowIv, _ = hex.DecodeString("d9c7c3c8870d64bd")
	return np
}

func (np *NavicatPassword) Decrypt(s string) (string, error) {
	switch np.version {
	case 11:
		return np.decryptEleven(s)
	case 12:
		return np.decryptTwelve(s)
	default:
		return "", fmt.Errorf("unsupported version: %d", np.version)
	}
}

func (np *NavicatPassword) decryptEleven(upperString string) (string, error) {
	str, err := hex.DecodeString(strings.ToLower(upperString))
	if err != nil {
		return "", err
	}

	round := len(str) / 8
	leftLength := len(str) % 8
	result := make([]byte, 0, len(str))
	currentVector := np.blowIv

	for i := 0; i < round; i++ {
		encryptedBlock := str[8*i : 8*(i+1)]
		temp, err := np.xorBytes(np.decryptBlock(encryptedBlock), currentVector)
		if err != nil {
			return "", err
		}
		currentVector, err = np.xorBytes(currentVector, encryptedBlock)
		if err != nil {
			return "", err
		}
		result = append(result, temp...)
	}

	if leftLength > 0 {
		currentVector = np.encryptBlock(currentVector)
		temp, err := np.xorBytes(str[8*round:], currentVector[:leftLength])
		if err != nil {
			return "", err
		}
		result = append(result, temp...)
	}

	return string(result), nil
}

func (np *NavicatPassword) decryptTwelve(upperString string) (string, error) {
	str, err := hex.DecodeString(strings.ToLower(upperString))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(np.aesKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, np.aesIv)
	plaintext := make([]byte, len(str))
	mode.CryptBlocks(plaintext, str)
	return string(plaintext), nil
}

func (np *NavicatPassword) encryptBlock(block []byte) []byte {
	cipher, _ := aes.NewCipher(np.blowKey)
	dest := make([]byte, len(block))
	cipher.Encrypt(dest, block)
	return dest
}

func (np *NavicatPassword) decryptBlock(block []byte) []byte {
	cipher, _ := aes.NewCipher(np.blowKey)
	dest := make([]byte, len(block))
	cipher.Decrypt(dest, block)
	return dest
}

func (np *NavicatPassword) xorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}

	return result, nil
}
