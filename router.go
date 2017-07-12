package main

import (
	"net/http"

	"io"
	"sync"

	"./auth"
	"./models/guru"
	"./models/user"
	"./models/projek"
)

type UserHandler int
type ProjekHandler int
type GuruHandler int

func (u UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//io.WriteString(res,req.RequestURI)
	io.WriteString(res, user.UserController(req.RequestURI, res, req))
}

func (p ProjekHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, projek.ProjekController(req.RequestURI, res, req))
}

func (g GuruHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, guru.GuruController(req.RequestURI, res, req))
}

func main() {
	var pengg UserHandler
	var proj ProjekHandler
	var guru GuruHandler

	auth.SessionStore = make(map[string]auth.Client)
	auth.StorageMutex = sync.RWMutex{}

	mux := http.NewServeMux()
	mux.Handle("/guru/", guru)
	mux.Handle("/login/", pengg)
	mux.Handle("/logout/", pengg)
	mux.Handle("/projek/", proj)
	mux.Handle("/like/", proj)
	mux.Handle("/comment/", proj)

	http.ListenAndServe(":9000", mux)
}
