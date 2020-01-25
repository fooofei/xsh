package base

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GetPlainPassword(pass string) (string, error) {
	if XConfig.Crypt.Key == "" {
		Warn.Print("crypt key empty.")
		return pass, nil
	}

	txt, err := AesDecrypt(pass)
	if err != nil {
		return pass, err
	}

	return txt, nil
}

func GetMaskPassword(pass string) (string, error) {
	if XConfig.Crypt.Key == "" {
		Warn.Print("crypt key empty.")
		return pass, nil
	}

	txt, err := AesEncrypt(pass)
	if err != nil {
		return pass, err
	}

	return txt, nil
}

func AesEncrypt(plaintext string) (string, error) {
	if XConfig.Crypt.Type != "aes" {
		return "", CryptTypeUnknown
	}

	if len(XConfig.Crypt.Key) != 32 {
		return "", CryptKeyIllegal
	}

	block, err := aes.NewCipher([]byte(XConfig.Crypt.Key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil

}

func AesDecrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	if XConfig.Crypt.Type != "aes" {
		return "", CryptTypeUnknown
	}

	if len(XConfig.Crypt.Key) != 32 {
		return "", CryptKeyIllegal
	}

	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(XConfig.Crypt.Key))
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("decrypt failed [%s]", "ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
