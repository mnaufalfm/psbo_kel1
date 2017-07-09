package guru

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../auth"
)

type Pelajaran struct {
	Matpel  string `json:"matpel,omitempty" bson:"matpel,omitempty"`
	IdKelas string `json:"idkelas,omitempty" bson:"idkelas,omitempty"`
}

type Guru struct {
	Id            string      `json:"id,omitempty" bson:"_id,omitempty"`
	NIP           string      `json:"nip,omitempty" bson:"nip,omitempty"`
	Password      string      `json:"password,omitempty" bson:"password,omitempty"`
	FotoProfil    string      `json:"fotoprofil,omitempty" bson:"fotoprofil,omitempty"` //simpan alamatnya saja
	Nama          string      `json:"nama,omitempty" bson:"nama,omitempty"`
	TglLahir      string      `json:"tgllahir,omitempty" bson:"tgllahir,omitempty"`
	Email         string      `json:"email,omitempty" bson:"email,omitempty"`
	MataPelajaran []Pelajaran `json:"matapelajaran,omitempty" bson:"matapelajaran,omitempty"`
	Gender        string      `json:"gender,omitempty" bson:"gender,omitempty"`
	NoHp          string      `json:"nohp,omitempty" bson:"nohp,omitempty"`
	Alamat        string      `json:"alamat,omitempty" bson:"alamat,omitempty"`
	IdKelas       string      `json:"idkelas,omitempty" bson:"idkelas,omitempty"` //id kelas yang diwalikan olehnya
}

/*Standar json pengembalian jika mengalami error*/
func ErrorReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	return fmt.Sprintf("{\"error\": %d, \"message\": \"%s\"}", code, pesan)
}

/*Standar json pengembalian jika berhasil*/
func SuccessReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	return fmt.Sprintf("{\"success\": %d, \"message\": \"%s\"}", code, pesan)
}

func GetGuru(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) (bool, string) {
	//Jika membuka profil user lain dan milik sendiri
	//linknya:9000/siswa/noinduk

	//Format Return: bool, string
	//bool: true: menandakan user mengakses datanya sendiri, false: menandakan user mengakses data orang lain
	//string: json pengembalian
	var guru Guru
	ses := s.Copy()
	defer ses.Close()

	c := ses.DB("studenthack").C("teachers")

	//resBody, err := ioutil.ReadAll(r.Body)
	//token := string(resBody)
	token := r.Header.Get("Auth")
	sess := r.Header.Get("Session")

	tokenSplit := strings.Split(token, ".")

	if !auth.CheckToken(token) {
		err := c.Find(bson.M{"nip": path}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "emailortu": 0, "tgllahir": 0, "nohp": 0, "alamat": 0, "matapelajaran": 0}).One(&guru)
		if err != nil {
			return false, ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
		}
	}

	if sess == "" || !auth.CheckSession(sess) {
		err := c.Find(bson.M{"nip": path}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "emailortu": 0, "tgllahir": 0, "nohp": 0, "alamat": 0, "matapelajaran": 0}).One(&guru)
		if err != nil {
			return false, ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
		}
	}

	if jwt.CheckToken(token) {
		idaccess := strings.Split(token, ".")[1]
		idaccesss := jwt.Base64ToString(idaccess)
		idhex := hex.EncodeToString([]byte(idaccesss))
		err := c.Find(bson.M{"nip": path}).One(&guru)
		if err != nil {
			return false, ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
		}
		if idhex != guru.Id {
			err = c.Find(bson.M{"nip": path}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "emailortu": 0, "tgllahir": 0, "nohp": 0, "alamat": 0, "matapelajaran": 0}).One(&guru)
		}
	} else {
		_ = c.Find(bson.M{"noinduk": path}).Select(bson.M{"_id": 0, "password": 0}).One(&guru)
	}

	//Pengaturan return untuk mengatur pengembalian data berdasarkan siapa yang membuka dan profil siapa yang dibuka (belum dilakukan)
	w.WriteHeader(http.StatusOK)
	us, _ := json.Marshal(guru)
	return string(us)
}

//Digunakan untuk mengedit data guru yang terdapat pada database
func EditGuru(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	//Format penggunaan http://localhost:9000/edit/nip
	var guru Guru
	var bsonn map[string]interface{}

	token := r.Header.Get("Auth")
	sesi := r.Header.Get("Session")
	tokenSplit := strings.Split(token, ".")

	if !jwt.CheckToken(token) {
		return ErrorReturn(w, "Token yang Dikirimkan Invalid", http.StatusForbidden)
	}

	if stat, msg := session.CheckSession(s, sesi, tokenSplit[1]); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	ses := s.Copy()
	defer ses.Close()

	err := json.NewDecoder(r.Body).Decode(&guru)
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}

	mess := jwt.Base64ToString(tokenSplit[1])
	messhex := hex.EncodeToString([]byte(mess))

	jsonguru, _ := json.Marshal(guru)
	err = json.Unmarshal(jsonguru, &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Edit Request", http.StatusBadRequest)
	}

	c := s.DB("studenthack").C("teachers")

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(messhex)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Gagal Edit Data", http.StatusBadRequest)
	}

	return SuccessReturn(w, "Berhasil Edit Data", http.StatusOK)
}

func CheckDupGuru(s *mgo.Session, p Guru) string {
	var ret string

	ses := s.Copy()
	defer ses.Close()

	c := ses.DB("studenthack").C("teachers")

	d, _ := c.Find(bson.M{"nip": p.NIP}).Count()
	if d > 0 {
		ret = "NIP"
	}

	d, _ = c.Find(bson.M{"email": p.Email}).Count()
	if d > 0 {
		if ret != "" {
			ret = ret + ", Email"
		} else {
			ret = "Email"
		}
	}

	d, _ = c.Find(bson.M{"nohp": p.NoHp}).Count()
	if d > 0 {
		if ret != "" {
			ret = ret + ", Nomor Handphone"
		} else {
			ret = "Nomor Handphone"
		}
	}

	return ret
}

func GuruController(urle string, w http.ResponseWriter, r *http.Request) string {
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

	if len(pathe) >= 2 {
		if pathe[0] == "edit" {
			return EditGuru(ses, w, r, pathe[1])
		} else if pathe[1] == "" {
			return GetGuru(ses, w, r, pathe[1])
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
