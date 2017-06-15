package post

import (
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"../../auth"

	"../comment"
)

var Post struct {
	Id            string            `json:"id,omitempty" bson:"_id,omitempty"`
	IdPengirim    string            `json:"idpengirim,omitempty" bson:"idpengirim,omitempty"`
	IsiPost       string            `json:"isipost,omitempty" bson:"isipost,omitempty"`
	TglPost       string            `json:"tglpost,omitempty" bson:"tglpost,omitempty"` //simpan alamatnya saja
	JumlahComment int               `json:"jumlahcomment,omitempty" bson:"jumlahcomment"`
	Comment       []comment.Comment `json:"comment,omitempty" bson:"comment,omitempty"`
	JumlahLike    int               `json:"jumlahlike,omitempty" bson:"jumlahlike"`
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

func CreatePost(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var post Post

	token := r.Header.Get("Auth")
	tokenSplit := strings.Split(token, ".")

	if !jwt.CheckToken(token) {
		return ErrorReturn(w, "Token yang Dikirimkan Invalid", http.StatusForbidden)
	}

	ses := s.Copy()
	defer s.Close()

}

func PostController(urle string, w http.ResponseWriter, r *http.Request) string {
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
			return GetSiswa(ses, w, r, pathe[1])
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
