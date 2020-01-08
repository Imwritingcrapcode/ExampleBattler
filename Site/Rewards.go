package Site

import (
	"encoding/json"
	. "html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func AfterBattle(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[AfterBattle] Redirected to login, blah-blah")
		Redirect(w, r, "/login")
		return
	}
	client := FindBaseID(session.UserID)
	if r.Method == http.MethodGet {
		userfree := client.GatherFreeData()
		Path := "/Site/rewards.html"
		pwd, _ := os.Getwd()
		Path = strings.Replace(pwd+Path, "/", "\\", -1)
		log.Println("[rewards] " + Path)
		template, err := ParseFiles(Path)
		if err != nil {
			panic(err)
		}
		template.Execute(w, userfree)
	} else {

		rewards := *GetRewards(client)
		DeleteRewards(client)
		w.WriteHeader(200)
		res, err := json.Marshal(rewards)
		if err != nil {
			log.Println("[Rewards] for", client.Username, err)
		}
		w.Write(res)
	}

}
