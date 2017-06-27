package session

import (
	"crypto/cipher"
	"crypto/des"
	"math/rand"
	"strconv"
	"time"

	"sync"

	"fmt"

	"gopkg.in/mgo.v2"
)

var Client struct {
	IdUser      string
	LogId       string
	LoggedIn    bool
	ExpiredTime int64
}

var SessionStore map[string]bool

func CheckSession(s *mgo.Session, token string) bool {

	return true
}

func CreateSession(s *mgo.Session, id string) string {
	if !CheckSession(s, id) {
		var klien Client
		var storageMutex sync.RWMutex

		klien.IdUser = id
		now := time.Now()
		dur, _ := time.ParseDuration("72h")
		expTime := now.Add(dur).Unix()
		klien.ExpiredTime = expTime

		klien.LoggedIn = true

		rand.Seed(now.Unix())
		plaintext := []byte(strconv.FormatInt(rand.Int63(), 10))
		triplekey := "studenth" + "studenth" + "studenth"
		block, _ := des.NewTripleDESCipher([]byte(triplekey))
		iv := []byte("cobasaja")
		mode := cipher.NewCBCEncrypter(block, iv)
		sessionidbyte := make([]byte, len(plaintext))
		mode.CryptBlocks(sessionidbyte, plaintext)
		sessionid := fmt.Sprintf("%x", sessionidbyte)
		klien.LogId = sessionid

		storageMutex.Lock()
		SessionStore[sessionid] = klien
		storageMutex.Unlock()

		return sessionid
	}
	return ""
}

func DeleteSession(s *mgo.Session, token string) bool {
	return true
}
