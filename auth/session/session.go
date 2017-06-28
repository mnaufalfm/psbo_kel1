package session

import (
	"math/rand"
	"strconv"
	"time"

	"sync"

	".."
	"gopkg.in/mgo.v2"
)

type Client struct {
	IdUser      string
	LogId       string
	LoggedIn    bool
	ExpiredTime int64
}

var SessionStore map[string]Client

//return true jika token valid, return false jika token invalid
func CheckSession(s *mgo.Session, token string, id string) (bool, string) {
	token = jwt.DecryptTripleDES(token, "bingunga")
	exp, ex := SessionStore[token]
	if !ex {
		return false, "Anda belum Login"
	}

	if exp.ExpiredTime >= time.Now().Unix() {
		DeleteSession(s, token)
		return false, "Sesi Anda Telah Habis"
	}

	if exp.IdUser != id {
		return false, "Sesi Tidak Valid"
	}

	return true, "Sesi Valid"
}

func CreateSession(s *mgo.Session, id string) (bool, string) {
	_, ex := SessionStore[id]
	if !ex {
		var klien Client
		var storageMutex sync.RWMutex

		klien.IdUser = id
		now := time.Now()
		dur, _ := time.ParseDuration("72h")
		expTime := now.Add(dur).Unix()
		klien.ExpiredTime = expTime

		klien.LoggedIn = true

		rand.Seed(now.Unix())
		randid := strconv.FormatInt(rand.Int63(), 10)
		sessionid := jwt.EncryptTripleDES(randid, "bingunga")
		klien.LogId = sessionid

		storageMutex.Lock()
		SessionStore[sessionid] = klien
		storageMutex.Unlock()

		return true, sessionid
	}
	return false, ""
}

func DeleteSession(s *mgo.Session, token string) bool {
	_, ex := SessionStore[token]
	if ex {
		delete(SessionStore, token)
		return true
	}
	return false
}
