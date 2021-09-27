package Site

import (
	"net/http"
	"log"
	"encoding/json"
	"strconv"
	. "../Abstract"
)

func Notifications(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Notifications] Redirected to login")
		w.WriteHeader(400)
		w.Write([]byte("/login"))
		return
	}
	SetState(session.UserID, GetState(session.UserID))
	user := FindBaseID(session.UserID)
	activity := user.CurrentActivity
	if r.Method == http.MethodGet {
		if activity != ConversionPage {
			over, _, left, dtype, amnt, notified := GetConversionInfo(user.UserID)
			if over && left < 1 && notified < 1 {
				AddNotification(user.UserID,
					"Your conversion for "+strconv.Itoa(amnt)+" :"+DustTypesToFilesMap[dtype]+"_dust_small: is over!",
					"conversion", false)
				NotifiedConversion(user.UserID)
			}
		}
		client, present := GetBattle(user.UserID)
		if present && client.State == Disconnected && activity != PlayingAs {
			AddNotification(user.UserID,
				"Your game as <b>" + ReleasedCharactersNames[client.PlayingAs] + "</b> is still going! Would you like to reconnect?","game", false)
		}
		switch {
		case activity == ConversionPage:
			SeeNotifications(user.UserID, "conversion")
		/*case activity == BrowsingFriendList:
			SeeNotifications(user.UserID, "friends")*/
		}
		notifications := GetNotifications(user.UserID, true)
		res, err := json.Marshal(notifications)
		if err != nil {
			log.Println("[Notifications] for", user.Username, err)
			return
		}
		w.Write(res)
		if len(notifications) > 0 {
			log.Println("[Notifications] Got notifications for", user.Username, notifications, ActivitiesToString[activity])
		}
		DeleteNotifications(user.UserID, "seen")
	}
}
