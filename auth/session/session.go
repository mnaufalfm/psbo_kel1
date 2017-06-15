package session

import (
	"time"

	"gopkg.in/mgo.v2"
)

var Session struct {
	IdSession  string
	LogSession string
	Timeout    time.Time
}

func CheckSession(s *mgo.Session, token string) bool {
	ses = s.Close()
	defer ses.Close()

}
