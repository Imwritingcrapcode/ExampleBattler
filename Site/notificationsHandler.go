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
		Redirect(w, r, "/login")
		return
	}
	user := FindBaseID(session.UserID)
	activity := user.CurrentActivity
	if r.Method == http.MethodGet {
		if activity != ConversionPage {
			over, _, left, dtype, amnt, notified := GetConversionInfo(user.UserID)
			if over && left < 1 && notified < 1 {
				AddNotification(user.UserID,
					"Your conversion for "+strconv.Itoa(amnt)+" :"+DustTypesToFilesMap[dtype]+"_dust_small: has finished!",
					"conversion")
				NotifiedConversion(user.UserID)
			}
		}
		client, present := ClientConnections[user.UserID]
		if present && client.State == Disconnected && activity != PlayingAs {
			AddNotification(user.UserID,
				"Your game as <b>" + ReleasedCharactersNames[client.PlayingAs] + "</b> is still going! Would you like to reconnect?","game")
		}
		switch {
		case activity == ConversionPage:
			SeeNotifications(user.UserID, "conversion")
		case activity == BrowsingFriendList:
			SeeNotifications(user.UserID, "friends")
		}
		notifications := GetNotifications(user.UserID)
		res, err := json.Marshal(notifications)
		if err != nil {
			log.Println("[Notifications] for", user.Username, err)
			return
		}
		w.Write(res)
		log.Println("[Notifications] Got notifications for", user.Username, notifications)
		DeleteNotifications(user.UserID, "seen")
	}
}
