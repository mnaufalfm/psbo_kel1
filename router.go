package main

import (
	"fmt"
	"net/http"

	"io"
	"sync"

	"./auth"
	"./models/comment"
	"./models/guru"
	"./models/ortu"
	"./models/post"
	"./models/user"
	"./models/files"
)

type UserHandler int
type PostHandler int
type GuruHandler int
type KomenHandler int
type OrtuHandler int
type FileHandler int

func (u UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//io.WriteString(res,req.RequestURI)
	io.WriteString(res, user.UserController(req.RequestURI, res, req))
}

func (p PostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, post.PostController(req.RequestURI, res, req))
}

func (g GuruHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, guru.GuruController(req.RequestURI, res, req))
}

func (c KomenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, comment.CommentController(req.RequestURI, res, req))
}

func (o OrtuHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, ortu.OrtuController(req.RequestURI, res, req))
}

func (f FileHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, files.FileController(req.RequestURI, res, req))
}

func DefaultServe(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	io.WriteString(res, fmt.Sprintf("{\"error\": %d, \"message\": \"%s\"}", http.StatusNotFound, "Path Tidak Ditemukan"))
}

func main() {
	var pengg UserHandler
	var post PostHandler
	var guru GuruHandler
	var komen KomenHandler
	var ortu OrtuHandler
	var file FileHandler

	auth.SessionStore = make(map[string]auth.Client)
	auth.StorageMutex = sync.RWMutex{}

	mux := http.NewServeMux()
	mux.Handle("/guru/", guru)
	mux.Handle("/login/", pengg)
	mux.Handle("/logout/", pengg)
	mux.Handle("/post/", post)
	mux.Handle("/like/", post)
	mux.Handle("/comment/", komen)
	mux.Handle("/ortu/", ortu)
	mux.Handle("/upload/",file)
	mux.HandleFunc("/", DefaultServe)

	http.ListenAndServe(":9000", mux)
}
