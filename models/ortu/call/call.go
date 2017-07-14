package call

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"../../../auth"
	"../../../const"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func GetOrtuBySiswa(s *mgo.Session, w http.ResponseWriter, r *http.Request, email string) (bool, string) {
	var ortu Ortu
	returnOrtu := make(map[string]interface{})

	ses := s.Copy()
	defer ses.Close()

	token := r.Header.Get(konst.HeaderToken)
	// fmt.Printf("token: %s\nsession: %s\n", token, sess)
	tokenSplit := strings.Split(token, ".")
	messhex := auth.Base64ToString(tokenSplit[1])

	c := ses.DB(konst.DBName).C(konst.DBOrtu)
	err := c.Find(bson.M{"email": email}).One(&ortu)
	if err != nil {
		return false, "Data Ortu Tidak Ditemukan"
	}

	returnOrtu["nama"] = ortu.Nama
	if messhex == hex.EncodeToString([]byte(ortu.Id)) {
		returnOrtu["linkprofil"] = "/ortu/profil"
	} else {
		returnOrtu["linkprofil"] = "/ortu/" + messhex
	}

	ortuJson, _ := json.Marshal(returnOrtu)
	return true, string(ortuJson)
}
