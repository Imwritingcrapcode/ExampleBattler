package Site

import (
	. "../Abstract"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Redirect(w http.ResponseWriter, r *http.Request, where string) {
	http.Redirect(w, r, where, 303)
}

func IsLoggedIn(r *http.Request) (bool, *Session) {
	var AlrdyLoggedIn bool
	var session *Session
	for _, cookie := range r.Cookies() {
		session = IsActiveSession(cookie.Value)
		if cookie.Name == "BattlerCookie" &&
			session != nil &&
			time.Now().UnixNano() <= session.Expires {
			AlrdyLoggedIn = true
			break
		}
	}
	return AlrdyLoggedIn, session
}

func IsLegit(reg_data *RegData) bool {
	var verdict bool
	re := regexp.MustCompile(`[^ぁ-ンa-zA-Zа-яА-Я0-9\-_.•ёЁ]+`)
	re2 := regexp.MustCompile(`\s`)
	if re.MatchString(reg_data.Username) ||
		re2.MatchString(reg_data.Password) ||
		len(reg_data.Username) > 32 ||
		len(reg_data.Password) > 32 ||
		len(reg_data.Password) < 6 ||
		len(reg_data.Username) < 3 {
		verdict = false
	} else {
		verdict = true
	}
	return verdict
}

func EncryptPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 9)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func Register(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if AlrdyLoggedIn {
		log.Print("[Register] " + strconv.FormatInt(session.UserID, 10) + " redirected to /")
		Redirect(w, r, "/")
	} else {
		log.Println("[Register] access by a new user")
		if r.Method == http.MethodPost {
			reg_data := RegData{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reg_data)
			if err != nil {
				panic(err)
			}
			if IsLegit(&reg_data) {
				ok, err := AddUser(&reg_data)
				if !ok {
					log.Println("[Register] did not register:", err)
					w.WriteHeader(400)
					w.Write([]byte("Sorry, your reg data is wrong"))
				} else {
					log.Println("[Register] registered", reg_data.Username)
					w.Write([]byte("/login"))
				}
			} else {
				w.WriteHeader(400)
				w.Write([]byte("Sorry, your reg data is wrong"))
				log.Println("[Register] Sorry, your reg data is wrong")
			}
		} else {
			Path := "/Site/register.html"
			pwd, _ := os.Getwd()
			Path = strings.Replace(pwd+Path, "/", "\\", -1)
			log.Println("[Register] " + Path)
			http.ServeFile(w, r, Path)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if AlrdyLoggedIn {
		log.Print("[Log In] " + strconv.FormatInt(session.UserID, 10) + " redirected to /")
		Redirect(w, r, "/")
	} else {
		if r.Method == http.MethodPost {
			reg_data := RegData{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reg_data)
			if err != nil {
				panic(err)
			}
			log.Println("[Log In] Requesting log in as", reg_data.Username)
			if IsLegit(&reg_data) {
				user := FindBase(reg_data.Username)
				if user != nil && PassCorrect(reg_data.Password, user) {
					curTime := time.Now()
					ID, err := bcrypt.GenerateFromPassword([]byte(user.Username+strconv.FormatInt(curTime.UnixNano(), 10)), 5)
					if err != nil {
						log.Println("[Log In] Can't encrypt session/cookie id")
						return
					}

					var expiration time.Time
					if reg_data.RememberMe {
						expiration = curTime.Add(24 * 14 * time.Hour)
					} else {
						expiration = curTime.Add(2 * time.Hour)
					}
					//remember the session
					session := Session{
						ID:      string(ID),
						Expires: expiration.UnixNano(),
						UserID:  user.UserID,
					}
					if time.Now().Sub(LastSessionsCleared) > SessionsClearPeriod {
						log.Println("[Log In] Clearing Sessions!")
						ClearSessions()
						LastSessionsCleared = time.Now()
					}
					AddSession(&session)
					//Set the client's cookie
					cookie := http.Cookie{Name: "BattlerCookie",
						Value: string(ID),
						Expires: expiration}
					http.SetCookie(w, &cookie)
					w.WriteHeader(200)
					w.Write([]byte("/"))
				} else if user == nil {
					log.Println("[Log In] Wrong username")
					w.WriteHeader(400)
					w.Write([]byte("Incorrect username or password!"))
				}
			} else {
				w.WriteHeader(400)
				log.Println("[Log In] Wrong password for", reg_data.Username)
				w.Write([]byte("Incorrect username or password!"))

			}
		} else {
			Path := "/Site/login.html"
			pwd, _ := os.Getwd()
			Path = strings.Replace(pwd+Path, "/", "\\", -1)
			log.Println("[Log In]", Path)
			http.ServeFile(w, r, Path)
		}
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	loggedIn, _ := IsLoggedIn(r)
	if !loggedIn {
		log.Println("[Welcome] UNWELCOMED")
		Redirect(w, r, "/login")
		return
	}
	if r.Method == http.MethodGet {
		log.Println("[Welcome] GET WELCOME")
		//userfree := *user.GatherFreeData()
		Path := "/Site/main.html"
		pwd, _ := os.Getwd()
		Path = strings.Replace(pwd+Path, "/", "\\", -1)
		log.Println(Path)
		/*template, err := ParseFiles(Path)
		if err != nil {
			panic(err)
		}
		template.Execute(w, userfree)*/
		http.ServeFile(w, r, Path)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Logout] Redirected to login")
		Redirect(w, r, "/login")
		return
	}
	user := FindBaseID(session.UserID)
	_, present := ClientConnections[user.UserID]
	if present {
		delete(ClientConnections, user.UserID)
	}
	for _, other := range UserQueue {
		if other.UserID == user.UserID {
			UserQueue.Remove(other.UserID)
			break
		}
	}
	SetState(session.UserID, Offline)
	DeleteNotifications(user.UserID, "all")
	log.Println("[Logout] see you,", session.UserID)
	DeleteSession(session)
	log.Println("[Logout] Redirected to login")
	Redirect(w, r, "/login")
}
