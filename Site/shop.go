package Site

import (
	. "../Abstract"
	"encoding/json"
	. "html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"math"
	"math/rand"
)

func Shop(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Shop] Redirected to login")
		Redirect(w, r, "/login")
		return
	}
	client := FindBaseID(session.UserID)
	if r.Method == http.MethodGet {
		userfree := client.GatherFreeData()
		Path := "/Site/shop.html"
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

func ShopItems(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[ShopItems] Redirected to login")
		Redirect(w, r, "/login")
		return
	}
	//client := FindBaseID(session.UserID)
	if r.Method == http.MethodGet {
		items := *GetPurchaseableItems()
		w.WriteHeader(200)
		res, err := json.Marshal(items)
		if err != nil {
			log.Println("[ShopItems] for", session.UserID, err)
		}
		w.Write(res)
	} else if r.Method == http.MethodPost {
		var item string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "You sent an invalid purchase request.", 400)
			return
		}
		Purchase(session.UserID, item, w)
	}
}

func Purchase(userID int64, ID string, w http.ResponseWriter) {
	exists, item := GetItemByID(ID)
	if !exists {
		http.Error(w, "You sent a wrong purchase request", 400)
		return
	}
	user := FindBaseID(userID)
	dust := user.GetDust(item.Dust)
	if dust >= item.Cost {
		user.SetDust(item.Dust, dust-item.Cost)
		if item.Type == "pack" {
			value := ReleasedCharactersPacks[ID]
			index := rand.Intn(len(value))
			gotGirl := value[index]
			if !HasGirl(userID, gotGirl) {
				UnlockGirl(user, gotGirl)
				w.WriteHeader(200)
				w.Write([]byte("Enjoy your purchase - " + ReleasedCharactersNames[gotGirl]))
			} else {
				dust = user.GetDust(item.Dust)
				user.SetDust(item.Dust, dust+int(math.Floor(float64(item.Cost/2.0))))
				w.WriteHeader(200)
				w.Write([]byte("You've got a duplicate - " + ReleasedCharactersNames[gotGirl]))
			}
		}
	} else {
		http.Error(w, "You don't have enough dust!~.", 400)
		return
	}
}

func GetPurchaseableItems() *[]ShopItem {
	items := make([]ShopItem, 0)
	for _, rarity := range Rarities {
		if len(ReleasedCharactersPacks[rarity]) != 0 {
			_, item := GetItemByID(rarity)
			items = append(items, item)
		}
	}
	return &items
}
