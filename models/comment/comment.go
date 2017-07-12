package comment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"../../auth"
	"../../const"
	"../post"
	"../user"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	IdPembuat  string `json:"idpembuat,omitempty" bson:"idpembuat,omitempty"`
	IsiComment string `json:"isicomment,omitempty" bson:"isicomment,omitempty"`
	TglComment int64  `json:"tglcomment,omitempty" bson:"tglcomment,omitempty"`
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

//Digunakan untuk membuat komen dan memasukkannya ke post
//Cara menggunakan: http://linknya:9000/comment/id
//Format pengiriman: {"komen": "isikomennya"}
func CreateComment(s *mgo.Session, w http.ResponseWriter, r *http.Request, id string) string {
	var komen Comment
	var komenPost []Comment
	var post post.Post
	komennSend, bsonn := make(map[string]interface{})

	ses := s.Copy()
	defer s.Close()

	token := r.Header.Get(konst.HeaderToken)
	sess := r.Header.Get(konst.HeaderSession)
	tokenSplit := strings.Split(token, ".")

	if stat, msg := auth.CheckToken(token); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	if stat, msg := auth.CheckSession(ses, sess, auth.Base64ToString(tokenSplit[1])); !stat {
		return ErrorReturn(w, msg, http.StatusBadRequest)
	}

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}

	_ = json.Unmarshal(resBody, &komennSend)

	komen.Id = bson.NewObjectId().Hex()
	komen.IdPembuat = auth.Base64ToString(tokenSplit[1])
	komen.IsiComment = komennSend["komen"].(string)
	komen.TglComment = time.Now().Unix()

	jsonKomen, _ := json.Marshal(komen)

	d := ses.DB(konst.DBName).C(konst.DBPost)
	err = d.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&post)
	if err != nil {
		return ErrorReturn(w, "Post Tidak Ditemukan", http.StatusBadRequest)
	}
	post.Comment = append(post.Comment, jsonKomen)

	jsonpost, _ := json.Marshal(post)
	err = json.Unmarshal(jsonpost, &bsonn)
	if err != nil {
		return ErrorReturn(w, "Tambah Komentar Gagal", http.StatusBadRequest)
	}

	err = d.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Tambah Komentar Gagal", http.StatusBadGateway)
	}

	return SuccessReturn(w, "Tambah Komentar Berhasil", http.StatusOK)
}

//Digunakan untuk membuat return json Komentar yang dipanggil oleh post.go
//Format pengembalian: {"id":"idkomen","komen":"isikomennya","profilpemilik":"/guru/id(misalkan)"}
func GetComment(s *mgo.Session, w http.ResponseWriter, r *http.Request, komen Comment) (bool, string) {
	// var komen Comment
	returnComment := make(map[string]interface{})

	ses := s.Copy()
	defer s.Close()

	token := r.Header.Get(konst.HeaderToken)
	sess := r.Header.Get(konst.HeaderSession)
	tokenSplit := strings.Split(token, ".")

	if stat, msg := auth.CheckToken(token); !stat {
		return stat, msg
	}

	if stat, msg := auth.CheckSession(ses, sess, auth.Base64ToString(tokenSplit[1])); !stat {
		return stat, msg
	}

	// err := json.Unmarshal([]byte(jsonKomen), &komen)
	// if err != nil {
	// 	return false, "Format Komentar Salah"
	// }

	returnComment["id"] = komen.Id
	returnComment["komen"] = komen.IsiComment

	var user user.Pengguna
	d := ses.DB(konst.DBName).C(konst.DBUser)
	err := d.Find(bson.M{"_id": bson.ObjectId(komen.IdPembuat)}).One(&user)
	if err != nil {
		return false, ErrorReturn(w, "Akun Pemilik Tidak Ditemukan", http.StatusBadRequest)
	}
	stat, role := konst.GetRoleString(user.LoginType)
	if !stat {
		return false, ErrorReturn(w, role, http.StatusBadRequest)
	}

	if auth.Base64ToString(tokenSplit[1]) == komen.IdPembuat {
		returnComment["profilpemilik"] = "/" + role + "/profil/"
	} else {
		returnComment["profilpemilik"] = "/" + role + "/" + komen.IdPembuat
	}

	w.WriteHeader(http.StatusOK)
	komenJson, _ := json.Marshal(returnComment)
	return string(komenJson)
}
