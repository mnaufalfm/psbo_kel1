package call

import (
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"encoding/hex"
	"encoding/json"

	"../../../auth"
	"../../../const"
)

type Siswa struct {
	Id         string   `json:"id,omitempty" bson:"_id,omitempty"`
	NoInduk    string   `json:"noinduk,omitempty" bson:"noinduk,omitempty"`
	Password   string   `json:"password,omitempty" bson:"password,omitempty"`
	FotoProfil string   `json:"fotoprofil,omitempty" bson:"fotoprofil,omitempty"` //simpan alamatnya saja
	Nama       string   `json:"nama,omitempty" bson:"nama,omitempty"`
	TglLahir   string   `json:"tgllahir,omitempty" bson:"tgllahir,omitempty"`
	Email      string   `json:"email,omitempty" bson:"email,omitempty"`
	EmailOrtu  string   `json:"emailortu,omitempty" bson:"emailortu,omitempty"`
	Gender     string   `json:"gender,omitempty" bson:"gender,omitempty"`
	NoHp       string   `json:"nohp,omitempty" bson:"nohp,omitempty"`
	Alamat     string   `json:"alamat,omitempty" bson:"alamat,omitempty"`
	IdKelas    []string `json:"idkelas,omitempty" bson:"idkelas,omitempty"`
}

func GetSiswaByOrtu(s *mgo.Session, w http.ResponseWriter, r *http.Request, email []string) (bool, string) {
	var siswa Siswa
	var datasSiswa []string
	returnSiswa := make(map[string]interface{})

	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get(konst.HeaderToken)
	// fmt.Printf("token: %s\nsession: %s\n", token, sess)
	tokenSplit := strings.Split(token, ".")
	messhex := auth.Base64ToString(tokenSplit[1])

	c := ses.DB(konst.DBName).C(konst.DBSiswa)

	datasSiswa = []string{}
	for i := 0; i < len(email); i++ {
		err := c.Find(bson.M{"email": email[i]}).One(&siswa)
		if err != nil {
			return false, "Data Ortu Tidak Ditemukan"
		}
		link := ""
		if messhex == hex.EncodeToString([]byte(siswa.Id)) {
			link = "/siswa/profil"
		} else {
			link = "/siswa/" + messhex
		}
		dataSiswa := fmt.Sprintf("{\"nama\": %d, \"link\": \"%s\"}", siswa.Nama, link)
		datasSiswa = append(datasSiswa, dataSiswa)
	}

	returnSiswa["anak"] = datasSiswa
	anakJson, _ := json.Marshal(returnSiswa)
	return true, string(anakJson)
}
