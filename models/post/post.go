package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../auth"
	"../../const"

	"io/ioutil"

	"encoding/hex"

	"../comment"
	"../user"
)

type Post struct {
	Id            string            `json:"id,omitempty" bson:"_id,omitempty"`
	IdPengirim    string            `json:"idpengirim,omitempty" bson:"idpengirim,omitempty"`
	IsiGambar     string            `json:"isigambar,omitempty" bson:"isigambar,omitempty"`
	IsiPost       string            `json:"isipost,omitempty" bson:"isipost,omitempty"`
	TglPost       int64             `json:"tglpost,omitempty" bson:"tglpost,omitempty"` //simpan alamatnya saja
	JumlahComment int               `json:"jumlahcomment" bson:"jumlahcomment"`
	Comment       []comment.Comment `json:"comment" bson:"comment,omitempty"`
	JumlahLike    int               `json:"jumlahlike" bson:"jumlahlike"`
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

//Digunakan untuk membuat post baru
//Cara menggunakan: http://linknya:9000/post/create/
//Format json pengiriman : {"post":"isipost"}
func CreatePost(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var post Post
	var postt map[string]interface{}
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

	_ = json.Unmarshal(resBody, &postt)

	//Fungsi utama membuat post
	post.IdPengirim = auth.Base64ToString(tokenSplit[1])
	post.IsiPost = postt["post"].(string)
	post.TglPost = time.Now().Unix()

	c := ses.DB(konst.DBName).C(konst.DBPost)
	err = c.Insert(post)
	if err != nil {
		return ErrorReturn(w, "Membuat Post Gagal", http.StatusBadGateway)
	}

	return SuccessReturn(w, "Membuat Post Berhasil", http.StatusOK)
}

//Digunakan untuk mendapatkan post dengan idpost tertentu
//Cara menggunakan: http://linknya:9000/post/id/
/*Format yang dikembalikan:
  {"id":"idpost","post":"isipostnya","profilpemilik":"/guru/id(misalkan)","jumlahkomen":0, "komen":[]Komen,"jumlahlike":0}.*/
func GetOnePost(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var post Post
	returnPost := make(map[string]interface{})

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

	c := ses.DB(konst.DBName).C(konst.DBPost)

	err := c.Find(bson.M{"_id": bson.ObjectId(path)}).One(&post)
	if err != nil {
		return ErrorReturn(w, "Post Tidak Ditemukan", http.StatusBadRequest)
	}

	// post.Id = hex.EncodeToString([]byte(post.Id))
	returnPost["id"] = hex.EncodeToString([]byte(post.Id))
	returnPost["post"] = post.IsiPost
	returnPost["jumlahkomen"] = post.JumlahComment
	// returnPost["komen"] = post.Comment
	returnPost["jumlahlike"] = post.JumlahLike

	//Menentukan isi dari returnPost["profilpemilik"]
	var user user.Pengguna
	d := ses.DB(konst.DBName).C(konst.DBUser)
	err = d.Find(bson.M{"_id": bson.ObjectId(post.IdPengirim)}).One(&user)
	if err != nil {
		return ErrorReturn(w, "Akun Pemilik Tidak Ditemukan", http.StatusBadRequest)
	}
	stat, role := konst.GetRoleString(user.LoginType)
	if !stat {
		return ErrorReturn(w, role, http.StatusBadRequest)
	}

	if auth.Base64ToString(tokenSplit[1]) == post.IdPengirim {
		returnPost["profilpemilik"] = "/" + role + "/profil/"
	} else {
		returnPost["profilpemilik"] = "/" + role + "/" + post.IdPengirim
	}

	//Mengurusi komentar
	komens := []string{}
	// returnPost["komen"] = []string{}
	for i := 0; i < len(post.Comment); i++ {
		stat, msg := comment.GetComment(ses, w, r, post.Comment[i])
		if !stat {
			return ErrorReturn(w, msg, http.StatusBadRequest)
		}
		komens = append(komens, msg)
	}
	returnPost["komen"] = komens

	w.WriteHeader(http.StatusOK)
	postJson, _ := json.Marshal(returnPost)
	return string(postJson)
}

//Digunakan untuk mendapatkan post dengan idpost tertentu
//Cara menggunakan: http://linknya:9000/post/all
/*Format pengembalian (return): {"posts":[{}{}{}]}*/
func GetAllPost(s *mgo.Session, w http.ResponseWriter, r *http.Request) string {
	var idAll []interface{}
	var posts []string
	returnPost := make(map[string]interface{})

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

	c := ses.DB(konst.DBName).C(konst.DBPost)
	err := c.Find(bson.M{}).Select(bson.M{"_id": 1}).All(&idAll)
	if err != nil {
		return ErrorReturn(w, "Tidak Ada Post", http.StatusBadRequest)
	}
	for i := 0; i < len(idAll); i++ {
		id, _ := json.Marshal(idAll[i].(bson.M))
		id = id[8 : len(id)-2]
		post := GetOnePost(ses, w, r, string(id))
		posts = append(posts, post)
	}
	returnPost["posts"] = posts

	w.WriteHeader(http.StatusOK)
	postJson, _ := json.Marshal(returnPost)
	return string(postJson)
}

//Digunakan untuk meng-like post tertentu
//Cara menggunakan: http://linknya:9000/like/idpost/
func LikePost(s *mgo.Session, w http.ResponseWriter, r *http.Request, path string) string {
	var post Post
	bsonn := make(map[string]interface{})

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

	c := ses.DB(konst.DBName).C(konst.DBPost)

	err := c.Find(bson.M{"_id": bson.ObjectIdHex(path)}).One(&post)
	if err != nil {
		return ErrorReturn(w, "Post Tidak Ditemukan", http.StatusBadRequest)
	}

	bsonn["jumlahlike"] = post.JumlahLike + 1
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(path)}, bson.M{"$set": bsonn})
	if err != nil {
		return ErrorReturn(w, "Like Post Gagal", http.StatusBadRequest)
	}

	return SuccessReturn(w, "Like Post Berhasil", http.StatusOK)
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
		if pathe[0] == "like" && pathe[1] != "" {
			return LikePost(ses, w, r, pathe[1])
		} else if pathe[1] == "create" {
			return CreatePost(ses, w, r)
		} else if pathe[2] == "" {
			return GetOnePost(ses, w, r, pathe[2])
		} else if pathe[2] == "all" {
			return GetAllPost(ses, w, r)
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
