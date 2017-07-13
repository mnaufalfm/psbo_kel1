package files

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"os"

	"../../auth"
	"../../const"
	"../user"
)

type File struct {
	Id            string `json:"id,omitempty" bson:"_id,omitempty"`
	IdPengirim    string `json:"idpengirim,omitempty" bson:"idpengirim,omitempty"`
	NamaFile      string `json:"namafile,omitempty" bson:"namafile,omitempty"`
	AlamatFile    string `json:"alamatfile,omitempty" bson:"alamatfile,omitempty"`
	TanggalUpload int64  `json:"tglupload,omitempty" bson:"tglupload,omitempty"`
	IdKelas       string `json:"idkelas,omitempty" bson:"idkelas,omitempty"`
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

//Digunakan untuk upload file
//Cara menggunakan: http://linknya:9000/upload/idkelas/
func UploadFile(s *mgo.Session, w http.ResponseWriter, r *http.Request, pathe string) string {
	var dataFile File
	var user user.Pengguna
	returnFile := make(map[string]interface{})
	ses := s.Copy()
	defer s.Close()

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

	d := ses.DB(konst.DBName).C(konst.DBUser)
	err := d.Find(bson.M{"_id": bson.ObjectIdHex(messhex)}).One(&user)
	if err != nil {
		return ErrorReturn(w, "User Tidak Ditemukan", http.StatusBadRequest)
	}

	if user.LoginType != 2 && user.LoginType != 4 {
		fmt.Println(user.LoginType)
		return ErrorReturn(w, "Tidak Boleh Unggah File", http.StatusForbidden)
	}

	//Upload File
	file, header, err := r.FormFile("materi")
	if err != nil {
		return ErrorReturn(w, "Format Request Salah", http.StatusBadRequest)
	}
	defer file.Close()

	formatfile := header.Header.Get("Content-Type")

	fmt.Println(formatfile)

	if formatfile != "application/pdf" {
		return ErrorReturn(w, "Format File Tidak Diterima", http.StatusForbidden)
	}

	path := "files/materi/" + pathe

	dataFile.NamaFile = header.Filename
	dataFile.AlamatFile = path + "/"
	dataFile.IdKelas = pathe
	dataFile.IdPengirim = messhex
	dataFile.TanggalUpload = time.Now().Unix()

	c := ses.DB(konst.DBName).C(konst.DBFile)
	err = c.Insert(dataFile)
	if err != nil {
		return ErrorReturn(w, "Tambah Database Gagal", http.StatusBadGateway)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	f, err := os.OpenFile(path+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return ErrorReturn(w, "Gagal Upload File", http.StatusBadGateway)
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return ErrorReturn(w, "Gagal Upload File", http.StatusBadGateway)
	}

	returnFile["nama"] = dataFile.NamaFile
	_ = c.Find(bson.M{}).Sort("-_id").One(&dataFile)
	returnFile["alamat"] = "download/" + hex.EncodeToString([]byte(dataFile.Id))

	w.WriteHeader(http.StatusOK)
	fileJson, _ := json.Marshal(returnFile)
	return string(fileJson)
}

func FileController(urle string, w http.ResponseWriter, r *http.Request) string {
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
		if pathe[0] == "upload" && pathe[1] != "" {
			return UploadFile(ses, w, r, pathe[1])
		}
	}
	return ErrorReturn(w, "Path Tidak Ditemukan", http.StatusNotFound)
}
