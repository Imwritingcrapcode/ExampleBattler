package Site

import (
	"encoding/json"
	. "html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
)

func Conversion(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Shop] Redirected to login")
		Redirect(w, r, "/login")
		return
	}
	user := FindBaseID(session.UserID)
	if r.Method == http.MethodGet {
		userfree := user.GatherFreeData()
		Path := "/Site/conversion.html"
		pwd, _ := os.Getwd()
		Path = strings.Replace(pwd+Path, "/", "\\", -1)
		log.Println("[Conversion] " + Path)
		template, err := ParseFiles(Path)
		if err != nil {
			panic(err)
		}
		template.Execute(w, userfree)
	} else if r.Method == http.MethodPost {
		convR := ConvRequest{}
		//type, amount, dusttype
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&convR)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("you sent a bad request"))
			return
		}
		if convR.ReqType == "!" {
			var amnt = convR.Amount
			var dust = convR.DustType
			rate, present := ConversionRate[dust]
			if !present {
				w.WriteHeader(400)
				w.Write([]byte("this is a wrong type of dust!"))
				return
			}
			monies := user.GetDust(dust)
			if monies < amnt {
				w.WriteHeader(400)
				w.Write([]byte("you are broke! You have " + strconv.Itoa(monies) + " of " + dust))
				return
			}
			array := make([]int, 3)
			//the amnt they'll get
			array[0] = int(rate * float64(amnt))
			//how much is remaining
			array[1] = amnt - int(float64(array[0])/rate)
			//how much time it will take in seconds!
			array[2] = array[0] * SecondsPerConversion[dust]
			ready, err := json.Marshal(&array)
			if err != nil {
				log.Println("[Conversion] Marshalling error for", user.Username, err)
				w.WriteHeader(400)
				w.Write([]byte("You have sent a wrong request!"))
				return
			}
			w.WriteHeader(200)
			w.Write(ready)
		} else if convR.ReqType == "?" {
			response := ConvResponse{
				ConversionRate:       ConversionRate,
				SecondsPerConversion: SecondsPerConversion,
			}
			res, err := json.Marshal(response)
			if err != nil {
				log.Println("[Conversion] for", user.Username, err)
				return
			}
			w.Write(res)

		} else {
			w.WriteHeader(400)
			w.Write([]byte("You sent a wrong request"))
		}
	}

}
