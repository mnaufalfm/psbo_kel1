package ortu

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

type Ortu struct {
	Id            string   `json:"id,omitempty" bson:"_id,omitempty"`
	DataDiri      string   `json:"datadiri,omitempty" bson:"datadiri,omitempty"`
	JenisDataDiri int      `json:"jenisdatadiri,omitempty" bson:"jenisdatadiri,omitempty"` //1: KTP, 2: SIM, 3: Lainny
	Username      string   `json:"username,omitempty" bson:"username,omitempty"`
	Password      string   `json:"password,omitempty" bson:"password,omitempty"`
	FotoProfil    string   `json:"fotoprofil,omitempty" bson:"fotoprofil,omitempty"` //simpan alamatnya saja
	Nama          string   `json:"nama,omitempty" bson:"nama,omitempty"`
	TglLahir      string   `json:"tgllahir,omitempty" bson:"tgllahir,omitempty"`
	Email         string   `json:"email,omitempty" bson:"email,omitempty"`
	EmailSiswa    []string `json:"emailsiswa,omitempty" bson:"emailsiswa,omitempty"`
	Gender        string   `json:"gender,omitempty" bson:"gender,omitempty"`
	NoHp          string   `json:"nohp,omitempty" bson:"nohp,omitempty"`
	Alamat        string   `json:"alamat,omitempty" bson:"alamat,omitempty"`
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
//Cara menggunakan: http://linknya.com:9000/ortu/profil/
func GetProfile(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var ortu Ortu

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

	c := ses.DB(konst.DBName).C(konst.DBOrtu)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(auth.Base64ToString(tokenSplit[1]))}).Select(bson.M{"_id": 0, "password": 0}).One(&ortu)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	ortuJson, _ := json.Marshal(ortu)
	return string(ortuJson)
}

//Digunakan untuk mengakses profil Ortu yang lain
//Cara menggunakan: http://linknya.com:9000/ortu/id/
func GetProfileOther(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var ortu Ortu

	ses := s.Copy()
	defer ses.Close()

	c := ses.DB(konst.DBName).C(konst.DBOrtu)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(path)}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "tgllahir": 0, "nohp": 0, "alamat": 0}).One(&ortu)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	ortuJson, _ := json.Marshal(ortu)
	return string(ortuJson)
}

//Digunakan untuk mengedit data ortu yang terdapat pada database
//Format penggunaan http://linknya.com:9000/ortu/edit
func EditOrtu(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var ortu Ortu
	var bsonn map[string]interface{}
	var user user.Pengguna

	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get(konst.HeaderToken)
	sess := r.Header.Get(konst.HeaderSession)
	tokenSplit := strings.Split(token, ".")
	messhex := auth.Base64ToString(tokenSplit[1])

	if stat, msg := auth.CheckToken(token); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	if stat, msg := auth.CheckSession(ses, sess, messhex); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	c := ses.DB(konst.DBName).C(konst.DBOrtu)

	err := json.NewDecoder(r.Body).Decode(&ortu)
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}

	err = c.Find(bson.M{"_id": bson.ObjectIdHex(messhex)}).One(&user)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	if len(ortu.EmailSiswa) > 0 {
		if user.LoginType != 4 {
			ortu.EmailSiswa = []string{}
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Anak", http.StatusForbidden)
		}
	}

	jsonortu, _ := json.Marshal(ortu)
	err = json.Unmarshal(jsonortu, &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Edit Request", http.StatusBadRequest)
	}

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(messhex)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Gagal Edit Data", http.StatusBadRequest)
	}

	return SuccessReturn(w, "Berhasil Edit Data", http.StatusOK)
}

func OrtuController(urle string, w http.ResponseWriter, r *http.Request) string {
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
			return EditOrtu(ses, w, r)
		} else if pathe[1] == "profil" {
			return GetProfile(ses, w, r)
		} else if pathe[1] != "" {
			return GetProfileOther(ses, w, r, pathe[1])
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
