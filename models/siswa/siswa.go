package siswa

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

type Siswa struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	NoInduk    string `json:"noinduk,omitempty" bson:"noinduk,omitempty"`
	Password   string `json:"password,omitempty" bson:"password,omitempty"`
	FotoProfil string `json:"fotoprofil,omitempty" bson:"fotoprofil,omitempty"` //simpan alamatnya saja
	Nama       string `json:"nama,omitempty" bson:"nama,omitempty"`
	TglLahir   string `json:"tgllahir,omitempty" bson:"tgllahir,omitempty"`
	Email      string `json:"email,omitempty" bson:"email,omitempty"`
	EmailOrtu  string `json:"emailortu,omitempty" bson:"emailortu,omitempty"`
	Gender     string `json:"gender,omitempty" bson:"gender,omitempty"`
	NoHp       string `json:"nohp,omitempty" bson:"nohp,omitempty"`
	Alamat     string `json:"alamat,omitempty" bson:"alamat,omitempty"`
	IdKelas    string `json:"idkelas,omitempty" bson:"idkelas,omitempty"`
}

/*Standar json pengembalian jika mengalami error*/
func ErrorReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	//fmt.Fprintf(w, "{error: %i, message: %q}", code, pesan)
	//return "{error: " + strconv.Itoa(code) + ", message: " + pesan + "}"
	return fmt.Sprintf("{\"error\": %d, \"message\": \"%s\"}", code, pesan)
}

/*Standar json pengembalian jika berhasil*/
func SuccessReturn(w http.ResponseWriter, pesan string, code int) string {
	w.WriteHeader(code)
	//return `{"success": " + strconv.Itoa(code) + ", message: " + pesan + "}`
	return fmt.Sprintf("{\"success\": %d, \"message\": \"%s\"}", code, pesan)
}

//Digunakan untuk membuka profil pribadi
//Cara menggunakan: http://linknya:9000/siswa/profil/
func GetProfile(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var siswa Siswa
	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get(konst.HeaderToken)
	sess := r.Header.Get(konst.HeaderSession)
	tokenSplit := strings.Split(token, ".")

	if stat, msg := auth.CheckToken(token); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	if stat, msg := auth.CheckSession(ses, sess, auth.Base64ToString(tokenSplit[1])); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	c := ses.DB(konst.DBName).C(konst.DBSiswa)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(auth.Base64ToString(tokenSplit[1]))}).Select(bson.M{"_id": 0, "password": 0}).One(&siswa)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	siswaJson, _ := json.Marshal(siswa)
	return string(siswaJson)
}

//Digunakan untuk mengakses profil Siswa yang lain
//Cara menggunakan: http://linknya.com:9000/siswa/nim/
func GetProfileOther(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var siswa Siswa

	ses := s.Copy()
	defer ses.Close()

	c := ses.DB(konst.DBName).C(konst.DBSiswa)

	err := c.Find(bson.M{"nip": path}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "emailortu": 0, "tgllahir": 0, "nohp": 0, "alamat": 0, "matapelajaran": 0}).One(&siswa)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	siswaJson, _ := json.Marshal(siswa)
	return string(siswaJson)
}

//Digunakan untuk mengedit data siswa yang terdapat pada database
//Format penggunaan http://linknya.com:9000/siswa/edit/
func EditSiswa(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var siswa Siswa
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
	c := s.DB(konst.DBName).C(konst.DBSiswa)

	err := json.NewDecoder(r.Body).Decode(&siswa)
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}

	messhex := auth.Base64ToString(tokenSplit[1])
	// messhex := hex.EncodeToString([]byte(mess))

	jsonsiswa, _ := json.Marshal(siswa)

	err = c.Find(bson.M{"_id": bson.ObjectId(messhex)}).One(&user)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	//Membatasi data yang dapat diedit
	if siswa.Id != "" {
		if user.LoginType != 4 {
			siswa.Id = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if siswa.NoInduk != "" {
		if user.LoginType != 4 {
			siswa.NoInduk = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if siswa.EmailOrtu != "" {
		if user.LoginType != 4 {
			siswa.EmailOrtu = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	if siswa.IdKelas != "" {
		if user.LoginType != 4 {
			siswa.IdKelas = ""
			return ErrorReturn(w, "Tidak Boleh Mengedit Data Ini", http.StatusForbidden)
		}
	}

	err = json.Unmarshal(jsonsiswa, &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Edit Request", http.StatusBadRequest)
	}

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(messhex)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Gagal Edit Data", http.StatusBadRequest)
	}

	return SuccessReturn(w, "Berhasil Edit Data", http.StatusOK)
}

func SiswaController(urle string, w http.ResponseWriter, r *http.Request) string {
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
			return EditSiswa(ses, w, r, pathe[1])
		} else if pathe[1] == "profil" {
			return GetProfile(ses, w, r)
		} else if pathe[1] == "" {
			return GetProfileOther(ses, w, r, pathe[1])
		}
	}
	
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
