package Site

import (
	"net/http"
	"log"
	"encoding/json"
)

func FreeInfo(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[FreeInfo] Redirected to login")
		Redirect(w, r, "/login")
		return
	}
	user := FindBaseID(session.UserID)
	data := user.GatherFreeData()
	res, err := json.Marshal(data)
	if err != nil {
		log.Println("[FreeInfo] for", user.Username, err)
	}
	w.WriteHeader(200)
	w.Write(res)
}
