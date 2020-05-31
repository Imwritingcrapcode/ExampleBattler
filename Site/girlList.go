package Site

import (
	. "../Abstract"
	. "../Characters"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func GirlListHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if AlrdyLoggedIn {
		log.Println("[GIRLLIST] "+"accessing GirlList for", session.UserID)
		if r.Method == http.MethodGet {
			//Path := "/Site/girllist.html" //old path lolololol
			Path := "/Site/girllist2.html"
			pwd, _ := os.Getwd()
			Path = strings.Replace(pwd+Path, "/", "\\", -1)
			log.Println("[GIRLLIST] " + Path)
			http.ServeFile(w, r, Path)

		} else {
			if GetState(session.UserID) <= Queuing { //TODO you know.
				//basically if they were queuing what disrupts the queue? and prompt to reconnect. ty
				SetState(session.UserID, BrowsingCharacters)
			}
			//1. Get a list (slice) of all the girls unlocked from the db
			userGirls := GetGirls(session.UserID)
			// 2. getinfo for each of those girls
			//2.5 put a GirlInfo inside
			info := make([]GirlInfo, len(userGirls))
			for i, _ := range userGirls {
				info[i] = *GetGirlInfo(userGirls[i].Number)
				/* all the fields below are taken from GirlInfo
				Name          string
				Rarity        string
				Tags          []string
				Skills        []string
				SkillColours  []string
				Description   string */
				userGirls[i].Name = info[i].Name
				userGirls[i].Rarity = info[i].Rarity
				userGirls[i].SkillColours = info[i].SkillColours
				userGirls[i].SkillColourCodes = info[i].SkillColourCodes
				userGirls[i].Tags = info[i].Tags
				userGirls[i].Skills = info[i].Skills
				userGirls[i].Description = info[i].Description
				userGirls[i].MainColour = info[i].MainColour
			}

			//3. Send it to the frontend
			var resline string
			standard := "Press \"Battle!\" when you are ready"
			/*if _, ok := ClientMap[session.Username]; ok {
				resline = "You are already in game. Queuing will terminate the current match (give up)"
			} else {
				resline = standard
			}*/
			resline = standard
			res1 := GirlListResponse{
				Girls:    userGirls,
				Response: resline,
			}

			res, err := json.Marshal(res1)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			log.Println("[GIRLLIST] "+"sending", len(userGirls), "girls to", session.UserID, ".")
			w.Write(res)

		}
	} else {
		log.Print("[GIRLLIST] " + "redirected to /login")
		Redirect(w, r, "/login")

	}
}
