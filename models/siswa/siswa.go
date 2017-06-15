package siswa

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../auth"
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

func CheckDupSiswa(s *mgo.Session, p Siswa) string {
	var ret string

	ses := s.Copy()
	defer ses.Close()

	c := ses.DB("studenthack").C("students")

	d, _ := c.Find(bson.M{"noinduk": p.NoInduk}).Count()
	if d > 0 {
		ret = "NoInduk"
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

func GetSiswa(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	//Jika membuka profil user lain dan milik sendiri
	//linknya:9000/siswa/noinduk
	var siswa Siswa
	ses := s.Copy()
	defer ses.Close()

	c := ses.DB("studenthack").C("students")

	//resBody, err := ioutil.ReadAll(r.Body)
	//token := string(resBody)
	token := r.Header.Get("Auth")
	if jwt.CheckToken(token) {
		idaccess := strings.Split(token, ".")[1]
		idaccesss := jwt.Base64ToString(idaccess)
		idhex := hex.EncodeToString([]byte(idaccesss))
		err := c.Find(bson.M{"noinduk": path}).One(&siswa)
		if err != nil {
			return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
		}
		if idhex != siswa.Id {
			err = c.Find(bson.M{"noinduk": path}).Select(bson.M{"_id": 0, "password": 0, "email": 0, "emailortu": 0, "tgllahir": 0, "nohp": 0, "alamat": 0}).One(&siswa)
		}
	} else {
		_ = c.Find(bson.M{"noinduk": path}).Select(bson.M{"_id": 0, "password": 0}).One(&siswa)
	}

	//Pengaturan return untuk mengatur pengembalian data berdasarkan siapa yang membuka dan profil siapa yang dibuka (belum dilakukan)
	w.WriteHeader(http.StatusOK)
	us, _ := json.Marshal(siswa)
	return string(us)
}

func EditSiswa(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	//penyesuaian sedikit
	ses := s.Copy()
	defer ses.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	req := string(reqBody)
	token := r.Header.Get("Auth")
	tokenSplit := strings.Split(token, ".")
	if req == "" {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}
	//fmt.Println(tokenSplit[0] + "." + tokenSplit[1] + "." + tokenSplit[2])
	if !jwt.CheckToken(token) {
		return ErrorReturn(w, "Token yang Dikirimkan Invalid", http.StatusForbidden)
	}
	mess := jwt.Base64ToString(tokenSplit[1])
	messhex := hex.EncodeToString([]byte(mess))
	//fmt.Println(mess)

	//kk, _ := json.Marshal(mess)
	//err := json.Unmarshal([]byte(mess), &sebelumEdit)
	//if err != nil {
	//	panic(err)
	//}

	var bsonn map[string]interface{}
	err := json.Unmarshal([]byte(req), &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Edit Request", http.StatusBadRequest)
	}
	//fmt.Println(sebelumEdit.Username)
	//fmt.Println(bsonn)

	c := s.DB("studenthack").C("students")

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
		} else if pathe[1] == "" {
			return GetSiswa(ses, w, r, pathe[1])
		}
	}

	// if pathe[0] == "login" {
	// 	return LoginUser(ses, w, r)
	// } else if pathe[0] == "registrasi" {
	// 	return RegistrasiUser(ses, w, r)
	// } else if pathe[0] == "edit" {
	// 	return EditUser(ses, w, r, pathe[1])
	// }

	// if len(pathe) >= 2 {
	// 	if pathe[1] != "" {
	// 		return GetUser(ses, w, r, pathe[1])
	// 	}
	// }
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
