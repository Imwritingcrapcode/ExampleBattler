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
			w.WriteHeader(400)
			w.Write([]byte("you sent a bad request"))
			return
		}
		isConverting, secondsPassed, secondsLeft := GetConversionInfo(user.UserID)
		if convR.ReqType == "?" { // how much &state
			var response ConvResponse
			if isConverting {
				response = ConvResponse{
					ConversionRate:       nil,
					SecondsPerConversion: nil,
					IsConvertingRN: true,
					CurrentProgress:secondsPassed,
					Left : secondsLeft,
				}
			} else {
				response = ConvResponse{
					ConversionRate:       ConversionRate,
					SecondsPerConversion: SecondsPerConversion,
					IsConvertingRN: false,
					CurrentProgress: -1,
					Left : -1,
				}
			}

			res, err := json.Marshal(response)
			if err != nil {
				log.Println("[Conversion] for", user.Username, err)
				return
			}
			w.Write(res)
		} else if convR.ReqType == "!" { //order & claim
			if isConverting && secondsLeft > 0 {
				w.WriteHeader(400)
				w.Write([]byte("You can not order or claim atm"))
				return
			}
			if isConverting { //claim
				ClaimConversion(user.UserID)
			} else { //order
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
					newDustType = "g"
				case "g":
					newDustType = "s"
				default:
					log.Println("[Conversion] what is this type of dust D:")
				}
				array := make([]int, 3)
				//the amnt they'll get
				array[0] = int(math.Floor(rate * float64(amnt)))
				if array[0] < 1 {
					w.WriteHeader(400)
					w.Write([]byte("you will get less than 1 dust"))
					return
				}
				cost := int(float64(array[0])/rate)
				//how much is remaining
				array[1] = amnt - cost
				//how much time it will take in seconds!
				array[2] = array[0] * SecondsPerConversion[dust]
				//pay up
				user.SetDust(dust, user.GetDust(dust)-cost)
				//start converting
				StartConversion(user.UserID, array[2], array[0], newDustType)
				ready, err := json.Marshal(&array)
				if err != nil {
					log.Println("[Conversion] Marshalling error for", user.Username, err)
					w.WriteHeader(400)
					w.Write([]byte("You have sent a wrong request!"))
					return
				}
				w.WriteHeader(200)
				w.Write(ready)
			}
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte("You sent a wrong request"))
	 }
}

