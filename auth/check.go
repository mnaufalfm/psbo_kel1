package auth

import (
	"crypto/hmac"
	"strings"
)

//untuk melakukan pengecekan keabsahan token. Return true jika token valid
func CheckToken(token string) (bool, string) {
	breakToken := strings.Split(token, ".")
	//fmt.Println("MasukToken")
	if len(breakToken) < 2 {
		//fmt.Println("Maksimal")
		return false, "Token Tidak Valid"
	}
	if breakToken[0] == "" || breakToken[1] == "" || breakToken[2] == "" {
		//fmt.Println("Lemah")
		return false, "Token Tidak Valid"
	}
	signSend := breakToken[2]
	signReal := ComputeHMAC256(breakToken[0]+"."+breakToken[1], "studenthack")
	if stat := hmac.Equal([]byte(signSend), []byte(signReal)); !stat {
		return false, "Token Tidak Valid"
	}
	return true, "Token Valid"
}
