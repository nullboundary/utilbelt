package utilbelt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"unicode"
)

//GenerateToken creates a urlEnocded hmac token
func GenerateToken(key []byte, seeds ...string) string {

	tokenSeed := strings.Join(seeds, "|")
	hmac := CalcHMAC(tokenSeed, key)

	return base64.URLEncoding.EncodeToString(hmac)

}

//VerifyToken returns true if token has a valid HMAC.
func VerifyToken(key []byte, authToken string, seeds ...string) (bool, error) {

	decodedMac, err := base64.URLEncoding.DecodeString(authToken)

	if err != nil {
		return false, fmt.Errorf("base64 Decode Error: %s", err)
	}
	tokenSeed := strings.Join(seeds, "|")
	return VerifyHMAC(tokenSeed, decodedMac, key), nil

}

//CalcHMAC makes an HMAC of a message and key
func CalcHMAC(message string, key []byte) []byte {

	mac := hmac.New(sha256.New, key)
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	return mac.Sum(nil)
}

//VerifyHMAC verifies a HMAC message
func VerifyHMAC(message string, macOfMessage []byte, key []byte) bool {

	mac := hmac.New(sha256.New, key)
	n, err := mac.Write([]byte(message))
	if n != len(message) || err != nil {
		panic(err)
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(macOfMessage, expectedMAC)
}

//EncryptAESCFB encrypts string to base64 crypto using AES
func EncryptAESCFB(key []byte, text string) string {

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

//DecryptAESCFB decrypts from base64 to decrypted string
func DecryptAESCFB(key []byte, cryptoText string) string {

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

//HashSha1Json computes a Sha1 hash of json data
func HashSha1Json(jsonData interface{}) []byte {

	//compute sha1 hash for json
	hash := sha1.New()
	enc := json.NewEncoder(hash) //json encode writes to the hash function
	enc.Encode(jsonData)
	return hash.Sum(nil)
}

//HashSha1Bytes computes a Sha1 hash of byte data
func HashSha1Bytes(hashBytes []byte) []byte {

	//compute sha1 hash of bytes
	hash := sha1.New()
	n, err := hash.Write(hashBytes)
	if n != len(hashBytes) || err != nil {
		panic(err)
	}
	return hash.Sum(nil)

}

//RandomStr generates a string of random letters and numbers with crypto/rand
func RandomStr(n int) string {
	g := big.NewInt(0)
	max := big.NewInt(130)
	bs := make([]byte, n)

	for i := range bs {
		g, _ = rand.Int(rand.Reader, max)
		r := rune(g.Int64())
		for !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			g, _ = rand.Int(rand.Reader, max)
			r = rune(g.Int64())
		}
		bs[i] = byte(g.Int64())
	}
	return string(bs)
}
