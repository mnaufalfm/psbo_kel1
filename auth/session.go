package auth

import (
	"crypto/des"
	"fmt"
	"math/rand"
	"time"

	"sync"

	konst "../const"
	"gopkg.in/mgo.v2"
)

type Client struct {
	IdUser      string
	LogId       string
	LoggedIn    bool
	ExpiredTime int64
}

var SessionStore map[string]Client
var StorageMutex sync.RWMutex

//return true jika token valid, return false jika token invalid
func CheckSession(s *mgo.Session, token string, id string) (bool, string) {
	// fmt.Println(token)
	tokenn := DecryptTripleDES(token, "12345678")
	StorageMutex.RLock()
	exp, ex := SessionStore[id]
	StorageMutex.RUnlock()
	// fmt.Println(token)
	fmt.Println(SessionStore)
	if !ex {
		return false, "Anda Belum Login"
	}

	if exp.ExpiredTime <= time.Now().Unix() {
		DeleteSession(s, id)
		return false, "Sesi Anda Telah Habis"
	}

	if exp.LogId != tokenn {
		// fmt.Println(id)
		return false, "Sesi Tidak Valid"
	}

	return true, "Sesi Valid"
}

func SessionIdRand(seed int64) string {
	rand.Seed(seed)
	b := ""
	for i := 0; i < des.BlockSize; i++ {
		b = b + string(konst.Letter[rand.Int63()%int64(len(konst.Letter))])
	}
	return b
}

func CreateSession(s *mgo.Session, id string) (bool, string) {
	StorageMutex.RLock()
	_, ex := SessionStore[id]
	StorageMutex.RUnlock()
	if !ex {
		var klien Client

		klien.IdUser = id
		now := time.Now()
		dur, _ := time.ParseDuration("72h")
		expTime := now.Add(dur).Unix()
		klien.ExpiredTime = expTime

		klien.LoggedIn = true

		// rand.Seed(now.Unix())
		// randid := strconv.FormatInt(rand.Int63(), 10)
		randid := SessionIdRand(now.Unix())
		StorageMutex.RLock()
		_, ch := SessionStore[randid]
		StorageMutex.RUnlock()
		for ch {
			randid = SessionIdRand(now.Unix())
		}
		// fmt.Println(randid)
		sessionid := EncryptTripleDES(randid, "12345678")
		klien.LogId = randid

		StorageMutex.Lock()
		// fmt.Println("CreateSession Lock")
		SessionStore[id] = klien
		StorageMutex.Unlock()
		// fmt.Println("CreateSession Unlock")

		fmt.Println(SessionStore)

		return true, sessionid
	}
	return false, ""
}

func DeleteSession(s *mgo.Session, id string) (bool, string) {
	// token = DecryptTripleDES(token, "12345678")
	StorageMutex.RLock()
	_, ex := SessionStore[id]
	StorageMutex.RUnlock()
	// fmt.Println("DeleteSession")
	if ex {
		// fmt.Println("right")
		StorageMutex.Lock()
		// fmt.Println("Lock")
		delete(SessionStore, id)
		StorageMutex.Unlock()
		// fmt.Println("Unlock")
		fmt.Println(SessionStore)
		return true, "Session Berhasil Dihapus"
	}
	return false, "Session Gagal Dihapus"
}
