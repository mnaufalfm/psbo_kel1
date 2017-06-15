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
)

type Pengguna struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string `json:"username,omitempty" bson:"username,omitempty"`
	Password   string `json:"password,omitempty" bson:"password,omitempty"`
	LoginType  int    `json:"logintype,omitempty" bson:"logintype,omitempty"` //1 = siswa, 2 = guru, 3 = ortu, 4 = school regulator
	StatusUser bool   `json:"statususer" bson:"statususer"`
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
		return fmt.Sprintf("{\"token\": \"%s\", \"access\": \"%s\"}", jwt.TokenMaker(log.Id, "studenthack"), jwt.StringToBase64(log.Username+" "+log.JenisAkun))
	}

	return ErrorReturn(w, "Password Salah", http.StatusForbidden)
}

//Digunakan untuk mengontrol path dari user (/user/...)
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
	}

	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
