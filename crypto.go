package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func GenerateToken(key string, seeds ...string) string {

	tokenSeed := strings.Join(seeds, "|")
	hmac := CalcHMAC(tokenSeed, key)
	return base64.URLEncoding.EncodeToString(hmac)

}

//////////////////////////////////////////////////////////////////////////
//
// verifyToken returns true if messageMAC is a valid HMAC tag for message.
//
//
//////////////////////////////////////////////////////////////////////////
func VerifyToken(key string, authToken string, seeds ...string) (bool, error) {

	decodedMac, err := base64.URLEncoding.DecodeString(authToken)
	if err != nil {
		return false, fmt.Errorf("base64 Decode Error: %s", err)
	}
	tokenSeed := strings.Join(seeds, "|")
	return VerifyHMAC(tokenSeed, decodedMac, key), nil

}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func CalcHMAC(message string, key string) []byte {

	mac := hmac.New(sha256.New, []byte(key))
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	return mac.Sum(nil)
}

//////////////////////////////////////////////////////////////////////////
//
//
//
//
//////////////////////////////////////////////////////////////////////////
func VerifyHMAC(message string, macOfMessage []byte, key string) bool {

	mac := hmac.New(sha256.New, []byte(key))
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(macOfMessage, expectedMAC)
}

//////////////////////////////////////////////////////////////////////////
//
// encrypt string to base64 crypto using AES
//
//
//////////////////////////////////////////////////////////////////////////
func encryptAESCFB(key []byte, text string) string {

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

//////////////////////////////////////////////////////////////////////////
//
// decrypt from base64 to decrypted string
//
//
//////////////////////////////////////////////////////////////////////////
func decryptAESCFB(key []byte, cryptoText string) string {

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
