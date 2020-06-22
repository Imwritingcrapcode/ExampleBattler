package Site

import (
	. "../Abstract"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func FriendsHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn || r.Method != http.MethodGet {
		log.Println("[Friends] " + "redirected to /login")
		Redirect(w, r, "/login")
		return
	}
	Path := "/html/friends.html"
	pwd, _ := os.Getwd()
	Path = strings.Replace(pwd+Path, "/", "\\", -1)
	log.Println("[Friends] " + Path, session.UserID)
	http.ServeFile(w, r, Path)}

func FriendListHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		Redirect(w, r, "/login")
		log.Println("[FriendList] redirected to /")
		return
	}
	if r.Method == http.MethodGet {
		SetState(session.UserID, BrowsingFriendList)
		res1 := GetFriendLists(session.UserID)
		res, err := json.Marshal(res1)
		if err != nil {
			log.Panic(err)
		}
		w.WriteHeader(200)
		//user := FindBaseID(session.UserID)
		//log.Println("[FriendList] sending", len(res1.Friends), "friends to "+strconv.FormatInt(user.UserID, 10)+".")
		//log.Println("and", len(res1.Incoming), "requests as well as", len(res1.Pending), "pending ones.")
		w.Write(res)
	} else if r.Method == http.MethodPost {
		friendReq := make([]string, 2)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&friendReq)
		if err != nil {
			http.Error(w, "You sent an invalid friend request.", 400)
			log.Println("[FriendListPost]", err.Error(), session.UserID)
			return
		}
		if friendReq[0] == "Remove" {
			other := FindBase(friendReq[1])
			if other != nil { //a valid user
				if !IsFriend(session.UserID, other.UserID) {
					if IsFriend(other.UserID, session.UserID) {
						RemoveFriend(other.UserID, session.UserID)
						w.WriteHeader(200)
						w.Write([]byte("Rejected a friend request from " + friendReq[1] + "."))
						log.Println("[FriendListPost] rejected a friend request from", other.UserID, "to", session.UserID)
					} else {
						http.Error(w, "You sent an invalid friend request.", 400)
						log.Println("[FriendListPost] invalid request", session.UserID, "removing someone who's not a friend")
					}
				} else {
					RemoveFriend(session.UserID, other.UserID)
					w.WriteHeader(200)
					if IsFriend(other.UserID, session.UserID) {
						w.Write([]byte("Removed " + friendReq[1] + " from your friend list."))
						log.Println("[FriendListPost] removed a friend", other.UserID, "for", session.UserID)
					} else {
						w.Write([]byte("Revoked friend request to " + friendReq[1] + "."))
						log.Println("[FriendListPost] revoked a friend request", other.UserID, "for", session.UserID)
					}
				}
			} else { //an invalid user
				http.Error(w, "You sent an invalid friend request.", 400)
				log.Println("[FriendListPost] invalid request", session.UserID, "trying to remove", friendReq[1])
			}
		} else if friendReq[0] == "Add" {
			other := FindBase(friendReq[1])
			if other != nil { //a valid user
				if other.UserID == session.UserID {
					http.Error(w, "You are already your own friend. Hopefully.", 400)
					log.Println("[FriendListPost] invalid request", session.UserID, "adding yourself.")
					return
				}
				fromthem := IsFriend(other.UserID, session.UserID)
				fromme := IsFriend(session.UserID, other.UserID)
				if fromthem && fromme{
					http.Error(w, "You are already friends with this user.", 400)
					log.Println("[FriendListPost] invalid request", session.UserID, "adding someone who's a friend already")
				} else if fromme { //they were already a friend
					http.Error(w, "You have already sent a friend request to this user.", 400)
					log.Println("[FriendListPost] invalid request", session.UserID, "adding someone who you've already sent a req to")
				} else {
					AddFriend(session.UserID, other.UserID)
					w.WriteHeader(200)
					if fromthem { //accepting
						w.Write([]byte("Accepted friend request from " + friendReq[1] + "."))
						log.Println("[FriendListPost] added a friend", other.UserID, "for", session.UserID)
					} else { //sending first
						w.Write([]byte("Sent friend request to " + friendReq[1] + "."))
						log.Println("[FriendListPost] sent request to ", other.UserID, "from", session.UserID)
					}
				}
			} else { //an invalid user
				http.Error(w, "You sent an invalid friend request.", 400)
				log.Println("[FriendListPost] invalid request", session.UserID, "trying to add", friendReq[1])
			}
		} else {
			http.Error(w, "You sent an invalid friend request.", 400)
			log.Println("[FriendListPost] invalid request", session.UserID, "not a valid action")
		}
	}
}
