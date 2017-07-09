package main

import (
	"net/http"

	"io"

	//"./models/projek"
	"sync"

	"./auth"
	"./models/user"
)

type UserHandler int
type ProjekHandler int

func (u UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//io.WriteString(res,req.RequestURI)
	io.WriteString(res, user.UserController(req.RequestURI, res, req))
}

func (p ProjekHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//io.WriteString(res, projek.ProjekController(req.RequestURI, res, req))
}

func main() {
	var pengg UserHandler
	var proj ProjekHandler

	auth.SessionStore = make(map[string]auth.Client)
	auth.StorageMutex = sync.RWMutex{}

	mux := http.NewServeMux()
	mux.Handle("/login/", pengg)
	mux.Handle("/logout/", pengg)
	mux.Handle("/projek/", proj)
	mux.Handle("/like/", proj)
	mux.Handle("/comment/", proj)

	http.ListenAndServe(":9000", mux)
}
