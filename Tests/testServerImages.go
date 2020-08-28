package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.Println("READY TO TEST!")
	http.HandleFunc("/favicon.ico", FavIcoFix)
	http.HandleFunc("/images/", ImageHandler2)
	http.HandleFunc("/scripts/", ScriptsHandler2)
	http.HandleFunc("/html/", ScriptsHandler3)
	err := http.ListenAndServe(":1117", nil)
	if err != nil {
		panic(err)
	}
}

func ImageHandler2(w http.ResponseWriter, r *http.Request) {
	var Path = r.URL.Path
	pwd, _ := os.Getwd()
	Path = strings.Replace(Path, "/", "\\", -1)
	log.Println("[IMAGE] Accessing " + Path + " as C.")
	img, err := ioutil.ReadFile(pwd + Path)
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func ScriptsHandler2(w http.ResponseWriter, r *http.Request) {
	log.Println("[SCRIPTS] access by C!")
	pwd, _ := os.Getwd()
	fs := http.FileServer(http.Dir(pwd + "\\scripts"))
	realhandler := http.StripPrefix("/scripts/", fs).ServeHTTP
	realhandler(w, r)
}
func ScriptsHandler3(w http.ResponseWriter, r *http.Request) {
	log.Println("[html] access by C!")
	pwd, _ := os.Getwd()
	fs := http.FileServer(http.Dir(pwd + "\\html"))
	realhandler := http.StripPrefix("/html/", fs).ServeHTTP
	realhandler(w, r)
}

func FavIcoFix(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{})
}
