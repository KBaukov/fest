package handlers

import (
	"fest/config"
	"log"
	"net/http"
)

var (
	//sessStore      = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))
	cfg, _ = config.LoadConfig("config.json")
	webres = cfg.FrontRoute.WebResFolder
)

func GetWebResources(w http.ResponseWriter, r *http.Request) {
	log.Println("###: ", r.URL.Path)
	http.ServeFile(w, r, "./"+webres+r.URL.Path)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	log.Println("###: ", r.URL.Path)
	http.ServeFile(w, r, "./"+webres+"/index.html")
}

func GetMap(w http.ResponseWriter, r *http.Request) {
	log.Println("###: ", r.URL.Path)
	http.ServeFile(w, r, "./"+webres+"/map.html")
}

func UserReg(w http.ResponseWriter, r *http.Request) {
	log.Println("###: ", r.URL.Path)
	http.ServeFile(w, r, "./"+webres+"/index.html")
}
