package base

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GetPlainPassword(pass string) string {
	if XConfig.Crypt.Key == "" {
		Warn.Print("crypt key empty.")
		return pass
	}

	txt, err := AesDecrypt(pass)
	if err != nil {
		Error.Printf("password [%s] decrypt error [%v].", pass, err)
	}

	return txt
}

func GetMaskPassword(pass string) string {
	if XConfig.Crypt.Key == "" {
		Warn.Print("crypt key empty.")
		return pass
	}

	tst, err := AesEncrypt(pass)
	if err != nil {
		Error.Printf("password [%s] encrypt error [%v].", pass, err)
	}

	return tst
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
