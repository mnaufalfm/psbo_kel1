package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"crypto/sha256"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../auth"
	"../../auth/session"
)

type Pengguna struct {
	Id        string `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string `json:"username,omitempty" bson:"username,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
	LoginType int    `json:"logintype,omitempty" bson:"logintype,omitempty"` //1 = siswa, 2 = guru, 3 = ortu, 4 = school regulator
}

func ErrorReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	//fmt.Fprintf(w, "{error: %i, message: %q}", code, pesan)
	//return "{error: " + strconv.Itoa(code) + ", message: " + pesan + "}"
	return fmt.Sprintf("{\"error\": %d, \"message\": \"%s\"}", code, pesan)
}

/*func SuccessReturn(w http.ResponseWriter, json []byte, pesan string, code int) string {
	w.WriteHeader(code)
	//fmt.Fprintf(w, "{message: %q}", pesan)
	return string(json)
}*/

func SuccessReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	//return `{"success": " + strconv.Itoa(code) + ", message: " + pesan + "}`
	return fmt.Sprintf("{\"success\": %d, \"message\": \"%s\"}", code, pesan)
}

//Fungsi untuk login
func LoginUser(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	//Digunakan untuk login ke halaman user
	var log Pengguna
	ses := s.Copy()
	defer ses.Close()

	//fmt.Println("Login User")

	err := json.NewDecoder(r.Body).Decode(&log)
	if err != nil {
		//fmt.Println("Cari data")
		return ErrorReturn(w, "Login Gagal", http.StatusBadRequest)
	}

	c := ses.DB("studenthack").C("users")

	encryptPassLogin := fmt.Sprintf("%x", sha256.Sum256([]byte(log.Password)))

	err = c.Find(bson.M{"username": log.Username}).One(&log)
	if err != nil {
		//fmt.Println("User Hilang")
		return ErrorReturn(w, "Anda Belum Registrasi", http.StatusBadRequest)
	}

	encryptPass := log.Password
	if encryptPass == encryptPassLogin {
		w.WriteHeader(http.StatusOK)
		stat, msg := session.CreateSession(ses, log.Id)
		if !stat {
			return ErrorReturn(w, "Login Gagal", http.StatusBadRequest)
		}
		// return fmt.Sprintf("{\"token\": \"%s\", \"access\": \"%s\", \"session\": \"%s\"}", jwt.TokenMaker(log.Id, "studenthack"), jwt.StringToBase64(log.Username+" "+strconv.Itoa(log.LoginType)), msg)
		return fmt.Sprintf("{\"token\": \"%s\", \"session\": \"%s\"}", jwt.TokenMaker(log.Id, "studenthack"), msg)
	}

	return ErrorReturn(w, "Password Salah", http.StatusForbidden)
}

//Fungsi untuk logout
func LogoutUser(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get("Auth")
	sess := r.Header.Get("Session")

	token = strings.Split(token, ".")[1]

	if jwt.CheckToken(token) {
		if stat, msg := session.CheckSession(ses, sess, jwt.Base64ToString(token)); stat {
			if session.DeleteSession(ses, sess) {
				return SuccessReturn(w, "Logout Sukses", http.StatusOK)
			}
		} else {
			return ErrorReturn(w, msg, http.StatusBadRequest)
		}
	}

	return ErrorReturn(w, "Logout Gagal", http.StatusBadRequest)
}

//Digunakan untuk mengontrol login/logout
func UserController(urle string, w http.ResponseWriter, r *http.Request) string {

	urle = urle[1:]
	pathe := strings.Split(urle, "/")
	fmt.Println(pathe[0] + " " + pathe[1])

	ses, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer ses.Close()
	ses.SetMode(mgo.Monotonic, true)
	//IndexCreating(ses)

	if pathe[0] == "login" {
		return LoginUser(ses, w, r)
	} else if pathe[0] == "logout" {
		return LogoutUser(ses, w, r)
	}

	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
