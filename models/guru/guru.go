package guru

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../auth"
	"../../const"
	"../user"
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

//Digunakan untuk membuka profil pribadi
//Cara menggunakan: http://linknya.com:9000/guru/profil/
func GetProfile(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	// fmt.Println("GetProfil")
	var guru Guru

	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get(konst.HeaderToken)
	sess := r.Header.Get(konst.HeaderSession)
	// fmt.Printf("token: %s\nsession: %s\n", token, sess)
	tokenSplit := strings.Split(token, ".")
	// fmt.Println(auth.Base64ToString(tokenSplit[1]))

	if stat, msg := auth.CheckToken(token); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	if stat, msg := auth.CheckSession(ses, sess, auth.Base64ToString(tokenSplit[1])); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	c := ses.DB(konst.DBName).C(konst.DBGuru)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(auth.Base64ToString(tokenSplit[1]))}).Select(bson.M{"_id": 0, "password": 0}).One(&guru)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	guruJson, _ := json.Marshal(guru)
	return string(guruJson)
}

//Digunakan untuk mengakses profil Guru yang lain
//Cara menggunakan: http://linknya.com:9000/guru/id/
func GetProfileOther(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var guru Guru

	ses := s.Copy()
	defer ses.Close()

	c := ses.DB(konst.DBName).C(konst.DBGuru)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(path)}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "tgllahir": 0, "nohp": 0, "alamat": 0, "matapelajaran": 0}).One(&guru)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	guruJson, _ := json.Marshal(guru)
	return string(guruJson)
}

//Digunakan untuk mengedit data guru yang terdapat pada database
//Format penggunaan http://linknya.com:9000/guru/edit
func EditGuru(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var guru Guru
	var bsonn map[string]interface{}
	var user user.Pengguna

	token := r.Header.Get(konst.HeaderToken)
	sesi := r.Header.Get(konst.HeaderSession)
	tokenSplit := strings.Split(token, ".")

	if stat, msg := auth.CheckToken(token); !stat {
		return ErrorReturn(w, msg, http.StatusForbidden)
	}

	if stat, msg := auth.CheckSession(s, sesi, auth.Base64ToString(tokenSplit[1])); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	ses := s.Copy()
	defer ses.Close()
	c := ses.DB(konst.DBName).C(konst.DBGuru)

	err := json.NewDecoder(r.Body).Decode(&guru)
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}

	messhex := auth.Base64ToString(tokenSplit[1])
	// messhex := hex.EncodeToString([]byte(mess))

	fmt.Println(guru.MataPelajaran)

	err = c.Find(bson.M{"_id": bson.ObjectId(messhex)}).One(&user)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	//Membatasi data yang dapat diedit
	if guru.Id != "" {
		if user.LoginType != 4 {
			guru.Id = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if guru.NIP != "" {
		if user.LoginType != 4 {
			guru.NIP = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if guru.IdKelas != "" {
		if user.LoginType != 4 {
			guru.IdKelas = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if len(guru.MataPelajaran) > 0 {
		if user.LoginType != 4 {
			guru.MataPelajaran = []Pelajaran{}
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	jsonguru, _ := json.Marshal(guru)
	err = json.Unmarshal(jsonguru, &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Edit Request", http.StatusBadRequest)
	}

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(messhex)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Gagal Edit Data", http.StatusBadRequest)
	}

	return SuccessReturn(w, "Berhasil Edit Data", http.StatusOK)
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
		if pathe[1] == "edit" {
			return EditGuru(ses, w, r)
		} else if pathe[1] == "profil" {
			return GetProfile(ses, w, r)
		} else if pathe[1] != "" {
			return GetProfileOther(ses, w, r, pathe[1])
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
