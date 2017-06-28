package main

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*type Rekening struct {
	NoRekening string `json:"norekening"`
	AtasNama   string `json:"atasnama"`
	Bank       string `json:"bank"`
}

type Pengguna struct {
	Id         string     `json:"id" bson:"_id,omitempty"`
	Username   string     `json:"username"`
	Password   string     `json:"pass"`
	FotoProfil string     `json:"fotoprofil"` //simpan alamatnya saja
	Nama       string     `json:"nama"`
	IdDiri     string     `json:"iddiri"`
	JenisID    int        `json:"jenisid"` //1=KTP, 2=SIM, 3=Paspor
	TglLahir   string     `json:"tgllahir"`
	Norek      []Rekening `json:"norek"`
	Email      string     `json:"email"`
	Gender     string     `json:"gender"`
	NoHp       string     `json:"nohp"`
	Alamat     string     `json:"alamat"`
}*/

type Rekening struct {
	NoRekening string `json:"norekening,omitempty"`
	AtasNama   string `json:"atasnama,omitempty"`
	Bank       string `json:"bank,omitempty"`
}

type Pengguna struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string        `json:"username,omitempty" bson:"username,omitempty"`
	Password   string        `json:"password,omitempty" bson:"password,omitempty"`
	FotoProfil string        `json:"fotoprofil,omitempty" bson:"fotoprofil,omitempty"` //simpan alamatnya saja
	Nama       string        `json:"nama,omitempty" bson:"nama,omitempty"`
	IdDiri     string        `json:"iddiri,omitempty" bson:"iddiri,omitempty"`
	JenisID    int           `json:"jenisid,omitempty" bson:"jenisid,omitempty"` //1=KTP, 2=SIM, 3=Paspor
	TglLahir   string        `json:"tgllahir,omitempty" bson:"tgllahir,omitempty"`
	Norek      []Rekening    `json:"norek,omitempty" bson:"norek,omitempty"`
	Email      string        `json:"email,omitempty" bson:"email,omitempty"`
	Gender     string        `json:"gender,omitempty" bson:"gender,omitempty"`
	NoHp       string        `json:"nohp,omitempty" bson:"nohp,omitempty"`
	Alamat     string        `json:"alamat,omitempty" bson:"alamat,omitempty"`
}

type Client struct {
	IdUser      string
	LogId       string
	LoggedIn    bool
	ExpiredTime int64
}

/*type Biasa struct {
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
}

type Berisik struct {
	Maklumi int     `json:"maklumi"`
	Haha    string  `json:"haha"`
	Hihi    []Biasa `json:"hihi"`
}*/

func ganteng(w http.ResponseWriter, r *http.Request, s *mgo.Session) []byte {
	ses := s.Copy()
	defer ses.Close()

	var log Pengguna

	c := ses.DB("coba").C("cobaa")

	err := json.NewDecoder(r.Body).Decode(&log)
	if err != nil {
		//fmt.Println("Cari data")
		fmt.Println("Error coy")
	}
	c.Insert(log)

	if log.Alamat == "" {
		fmt.Println("Kosong coy")
	}

	a, _ := json.Marshal(log)

	return a
}

type UserHandler int

func (u UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ses, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	defer ses.Close()
	ses.SetMode(mgo.Monotonic, true)

	io.WriteString(res, string(ganteng(res, req, ses)))
	fmt.Println(req.Header.Get("Auth"))
	//res.Write([]byte(`{"code":400, "error":"Hahaha iseng aja"}`))
}

func main() {
	//var berisik Berisik

	// a := time.Now()
	// sub, _ := time.ParseDuration("1s")
	// c := a.Add(sub)
	// fmt.Println(a.Unix())
	// fmt.Println(c.Unix())
	// if a.Unix() > c.Unix() {
	// 	fmt.Println("Hai")
	// } else if a.Unix() < c.Unix() {
	// 	fmt.Println("Hello")
	// }

	var coba map[string]Client

	_, ok := coba["syalala"]
	if !ok {
		fmt.Println("Hai")
	}

	// rand.Seed(time.Now().Unix())
	// fmt.Println(rand.Intn(10-0) + 0)

	// skrg := time.Now()
	// fmt.Println(skrg)
	// sub, _ := time.ParseDuration("72h")
	// fmt.Println(sub)
	// hm := skrg.Add(sub)
	// fmt.Println(hm)
	// ca := hm.String()
	// fmt.Println(ca)
	// ba, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", ca)
	// fmt.Println(ba)
	// if ba == hm {
	// 	fmt.Println("Yes")
	// }
	// fmt.Println(hm.Unix())

	// a := "Hai"
	// fmt.Println(a)
	// b := a + "Nama"
	// fmt.Println(b)
	// a = a + a + a
	// fmt.Println(a)
	// c := b - "Nama"
	// fmt.Println(c)

	//jsonhaha := []byte(`{"maklumi":5,"haha":"blabla"}`)
	//jsonblob := "{username: 'williamhanugra', pass: 'ganteng123', fotoprofil: 'blabla.jpg',nama: 'Lu William Hanugra',iddiri: '135060700111084',jenisid: 1,tgllahir: '14 April 2017',norek: [{norekening:'44444',atasnama:'William Hanugra',bank:'IPB Syariah'}],email: 'cipatonthesky@gmail.com',gender: 'L',nohp: '087873766464',alamat: 'Pondok Bu Sri'}"

	//fmt.Println(json)
	//err := json.Unmarshal(jsonhaha, &berisik)
	//if err != nil {
	//	fmt.Println("Gagal coy")
	//}
	//fmt.Printf("%+v", berisik)

	//bb := reflect.ValueOf(&berisik).Elem()
	//for i := 0; i < bb.NumField(); i++ {
	//	fmt.Printf("%s %s %v\n", bb.Type().Field(i).Name, bb.Type(), bb.Field(i).Interface())
	//}
	//fmt.Println()
	//fmt.Printf("%+v", berisik.Hihi[0])

	/*type ColorGroup struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	//var b Pengguna
	//err := json.NewDecoder([]byte(jsonblob)).Decode(&b)
	var cg ColorGroup
	a := []byte(`{"status": 1, "message": "Reds"}`)
	err := json.Unmarshal(a, cg)
	if err != nil {
		fmt.Println("error:", err)
	}
	bb := reflect.ValueOf(&cg).Elem()
	for i := 0; i < bb.NumField(); i++ {
		fmt.Println(bb.Field(i).Interface())
	}*/
	//os.Stdout.Write(b)

	/*s := "/guru/edit"
	s = s[1:]
	fmt.Println(s)*/

	/*pesan := 5
	haha := "Walah to"
	fmt.Printf("{message: %d, haha: %q}", pesan, haha)*/

	//sum := sha256.Sum256([]byte("hello world\n"))
	//fmt.Println(fmt.Sprintf("%x", sum))
	//fmt.Printf("%x", sum)

	/*var sum string
	sum = fmt.Sprintf("%x", sha256.Sum256([]byte("hello world\n")))
	fmt.Println(sum)*/
	// var user Pengguna
	// objectid := bson.NewObjectId().Hex()
	// a := 5
	// b := "bodo amat"
	// fmt.Printf("hai: %d \"%s\" %s", a, b, objectid)
	//var uu UserHandler
	//http.ListenAndServe(":9000", uu)

	// ses, err := mgo.Dial("localhost:27017")
	// if err != nil {
	// 	panic(err)
	// }

	// defer ses.Close()
	// ses.SetMode(mgo.Monotonic, true)

	// c := ses.DB("coba").C("cobaa")

	// err = c.Find(bson.M{"username": "apadeh"}).One(&user)
	// if err != nil {
	// 	panic(err)
	// }
	// hee := string(user.Id)
	// fmt.Printf(" Idnya = %s, Username = %s", hex(hee), user.Username)

	//gan, erro := json.Marshal(user.Id)
	//if erro != nil {
	//	panic(err)
	//}
	//fmt.Println(string(gan))
	//a, _ := hex.DecodeString(user.Id.Hex())
	/*b := base64.StdEncoding.EncodeToString([]byte(user.Id))
	fmt.Println(b)
	d, _ := base64.StdEncoding.DecodeString("WQl3aW8xBav4+UsI")
	e := string(d)
	f := hex.EncodeToString([]byte(e))
	fmt.Println(f)
	//fmt.Println(jwt.TokenMaker(string([]byte(user.Id)), "anggunauranaufalwilliam"))
	//fmt.Printf("%+v", user)
	//if user.Id.Hex() == f {
	//		fmt.Println("Bacot")
	//	}

	fmt.Printf("%s %s", reflect.TypeOf(user.Id), hex.EncodeToString([]byte(user.Id)))*/
	//	fmt.Println(reflect.TypeOf(user.Id.Hex()))

	//data := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9")
	//a := make([]byte, base64.RawStdEncoding.EncodedLen(len(data)))
	//str := base64.StdEncoding.EncodeToString(data)
	//base64.RawStdEncoding.Encode(a, data)
	//fmt.Println(hex.EncodeToString(a))
	//fmt.Printf("%x", a)
	//fmt.Println(str)

	//kun := []byte("berisikamatlu")
	//pes := []byte(mes)
	//h := hmac.New(sha256.New, kun)
	//h.Write(data)
	//fmt.Println(base64.RawStdEncoding.EncodeToString(h.Sum(nil)))
	//pass := "ganteng123"

	//fmt.Println(fmt.Sprintf("%x", sha256.Sum256([]byte(pass))))

	//fmt.Println(string("184054b78da172c42e37015fb66dd6968b582846f4226c9edfad9da80dc2bf22"))
	//s := append([]string{"1", "2"}, []string{"3", "4"}...)
	//fmt.Println(s)

	//var bsonn map[string]interface{}
	//err = json.Unmarshal([]byte(s), &bsonn)
	//fmt.Println(bsonn)
}
