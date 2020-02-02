package Site

import (
	"encoding/json"
	. "html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"math"
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
			log.Println(err.Error())
			w.WriteHeader(400)
			w.Write([]byte("you sent a bad request"))
			return
		}
		isConverting, secondsPassed, secondsLeft, dust, amnt := GetConversionInfo(user.UserID)
		if convR.ReqType == "?" || convR.ReqType == "!" && isConverting && secondsLeft < 1 { // how much & state & convert
			if convR.ReqType != "?" {
				ClaimConversion(user.UserID)
			}
			var response ConvResponse
			if convR.ReqType != "!" && isConverting {
				response = ConvResponse{
					ConversionRate:       nil,
					SecondsPerConversion: nil,
					IsConvertingRN:       true,
					CurrentProgress:      secondsPassed,
					Left:                 secondsLeft,
					DustType:             dust,
					Amount:               amnt,
				}
			} else {
				response = ConvResponse{
					ConversionRate:       ConversionRate,
					SecondsPerConversion: SecondsPerConversion,
					IsConvertingRN:       false,
					CurrentProgress:      -1,
					Left:                 -1,
				}
			}

			res, err := json.Marshal(response)
			if err != nil {
				log.Println("[Conversion] for", user.Username, err)
				return
			}
			w.Write(res)
		} else if convR.ReqType == "!" && !isConverting { //order
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
				w.Write([]byte("you are broke! You have " + strconv.Itoa(monies) + " of " + dust + " and you needed " + strconv.Itoa(amnt)))
				return
			}
			var newDustType string
			switch dust {
			case "w":
				newDustType = "b"
			case "b":
				newDustType = "y"
			case "y":
				newDustType = "p"
			case "p":
				newDustType = "s"
			default:
				log.Println("[Conversion] what is this type of dust D:")
			}
			response := ConvResponse{
				ConversionRate:       nil,
				SecondsPerConversion: nil,
				IsConvertingRN:       true,
				CurrentProgress:      0,
			}
			//the amnt they'll get
			get := int(math.Floor(rate * float64(amnt)))
			if get < 1 {
				w.WriteHeader(400)
				w.Write([]byte("you will get less than 1 dust"))
				return
			}
			cost := int(float64(get) / rate)
			//how much time it will take in seconds!
			time := get * SecondsPerConversion[dust]
			//pay up
			user.SetDust(dust, user.GetDust(dust)-cost)
			//start converting
			StartConversion(user.UserID, time, get, newDustType)
			//send Info
			response.Left = time
			response.DustType = newDustType
			response.Amount = get
			ready, err := json.Marshal(&response)
			if err != nil {
				log.Println("[Conversion] Marshalling error for", user.Username, err)
				w.WriteHeader(400)
				w.Write([]byte("You sent a wrong request!"))
				return
			}
			w.WriteHeader(200)
			w.Write(ready)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("You sent a wrong request. Idek what you meant."))
		}
	}
}
