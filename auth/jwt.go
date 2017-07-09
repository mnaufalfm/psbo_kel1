package auth

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var Token struct {
	IdToken   string `json:"idtoken,omitempty"`
	JenisAkun int    `json:"jenisakun,omitempty"`
}

func StringToBase64(mes string) string {
	mess := []byte(mes)
	return base64.StdEncoding.EncodeToString(mess)
}

func Base64ToString(bes64 string) string {
	mess, _ := base64.StdEncoding.DecodeString(bes64)
	return hex.EncodeToString(mess)
}

//message dalam bentuk json
func TokenMaker(message string, key string) string {
	header := `{"alg":"HS256", "typ":"JWT"}`
	a := StringToBase64(header) + "." + StringToBase64(message)
	sign := ComputeHMAC256(a, key)
	token := a + "." + sign
	return token
}

func ComputeHMAC256(mes string, key string) string {
	//fmt.Println("ComputeHMAC256")
	//fmt.Println(mes)
	kun := []byte(key)
	pes := []byte(mes)
	h := hmac.New(sha256.New, kun)
	h.Write(pes)
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

func KeyTripleDES(key string) string {
	if len(key) > des.BlockSize {
		return key[:des.BlockSize] + key[:des.BlockSize] + key[:des.BlockSize]
	} else {
		padding := key[:(des.BlockSize - len(key))]
		return key + padding + key + padding + key + padding
	}
}

func EncryptTripleDES(mes string, key string) string {
	key = KeyTripleDES(key)
	block, _ := des.NewTripleDESCipher([]byte(key))
	iv := []byte("cobasaja")
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypbyte := make([]byte, len(mes))
	mode.CryptBlocks(encrypbyte, []byte(mes))
	encryp := fmt.Sprintf("%x", encrypbyte)
	return encryp
}

func DecryptTripleDES(cip string, key string) string {
	cc, _ := hex.DecodeString(cip)
	cip = string(cc)
	key = KeyTripleDES(key)
	block, _ := des.NewTripleDESCipher([]byte(key))
	iv := []byte("cobasaja")
	mode := cipher.NewCBCDecrypter(block, iv)
	textbyte := make([]byte, len(cip))
	mode.CryptBlocks(textbyte, []byte(cip))
	// text := fmt.Sprintf("%x", textbyte)
	return string(textbyte)
}
