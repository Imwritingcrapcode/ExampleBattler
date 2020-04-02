package main

import (
	. "../Abstract"
	. "../Site"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func FavIcoFix(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{})
}

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = f
}

func main() {
	outputToLogTxt := false
	if outputToLogTxt {

		f, err1 := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err1 != nil {
			fmt.Println("error opening file: %v", err1)
			return

		}
		log.SetOutput(f)
		redirectStderr(f)
	}

	var err error

	DATABASE, err = sql.Open("sqlite3", "Server\\Battler.db")
	if err != nil {
		panic(err)
	}
	statement, err := DATABASE.Prepare("CREATE TABLE IF NOT EXISTS userData (userID INTEGER PRIMARY KEY, battlesTotal INTEGER DEFAULT 0, battlesWon INTEGER DEFAULT 0, lastBattleResult INTEGER DEFAULT 0, lastLoginTime INTEGER NOT NULL, wdust INTEGER DEFAULT 0, bdust INTEGER DEFAULT 0, ydust INTEGER DEFAULT 0, pdust INTEGER DEFAULT 0, sdust INTEGER DEFAULT 0, activity INTEGER DEFAULT 0, lastActivityTime INTEGER DEFAULT 0, lastPlayedAs INTEGER DEFAULT 0, lastOpponentsName TEXT DEFAULT \"\")")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	/*statement, err = DATABASE.Prepare("ALTER TABLE userData RENAME COLUMN gdust TO pdust")
	if err != nil {
		panic(err)
	}
	statement.Exec()*/

	statement, err = DATABASE.Prepare("UPDATE userData SET activity = 0")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS users (userID INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT,  password TEXT NOT NULL, FOREIGN KEY (userID) REFERENCES userData(userID))")
	if err != nil {
		panic(err)

	}
	statement.Exec()
	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS girls (userID INTEGER, girlNumber INTEGER NOT NULL, girlLevel INTEGER DEFAULT 1, currMatches INTEGER DEFAULT 0, matchesPlayed INTEGER DEFAULT 0, matchesWon INTEGER DEFAULT 0, FOREIGN KEY (userID) REFERENCES userData(userID), PRIMARY KEY (userID, girlNumber))")
	if err != nil {
		panic(err)

	}
	statement.Exec()

	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS rewards (rewardID INTEGER PRIMARY KEY AUTOINCREMENT, userID INTEGER NOT NULL, rewardType TEXT NOT NULL, amount INTEGER DEFAULT 1, FOREIGN KEY(userID) REFERENCES userData(userID))")
	if err != nil {
		panic(err)

	}
	statement.Exec()
	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS sessions (sessionID TEXT PRIMARY KEY, userID INTEGER NOT NULL, expires INTEGER NOT NULL, FOREIGN KEY(userID) REFERENCES userData(userID))")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS friends (userID INTEGER NOT NULL, friendID INTEGER NOT NULL, FOREIGN KEY(userID) REFERENCES userData(userID), PRIMARY KEY (userID, friendID))")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	/*statement, err = DATABASE.Prepare("DROP TABLE conversions")
	if err != nil {
		panic(err)
	}
	statement.Exec()*/

	statement, err = DATABASE.Prepare("CREATE TABLE IF NOT EXISTS conversions (userID INTEGER PRIMARY KEY, begins INTEGER NOT NULL, duration INTEGER NOT NULL, give INTEGER NOT NULL, type TEXT NOT NULL, FOREIGN KEY(userID) REFERENCES userData(userID))")
	if err != nil {
		panic(err)
	}
	statement.Exec()

	/*urkitten := FindBase("urkitten")
	for number, _ := range ReleasedCharactersNames {
		UnlockGirl(urkitten, number)
	}*/

	rand.Seed(time.Now().UTC().UnixNano())
	LastSessionsCleared = ClearSessions()

	http.HandleFunc("/game", Standard)
	http.HandleFunc("/battler", BattlerHandler)
	http.HandleFunc("/images/", ImageHandler)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/girllist", GirlListHandler)
	http.HandleFunc("/favicon.ico", FavIcoFix)
	http.HandleFunc("/friends", FriendsHandler)
	http.HandleFunc("/friendlist", FriendListHandler)
	http.HandleFunc("/afterbattle", AfterBattle)
	http.HandleFunc("/queue", HandleConnections)
	http.HandleFunc("/scripts/", ScriptsHandler)
	http.HandleFunc("/shop", Shop)
	http.HandleFunc("/shopitems", ShopItems)
	http.HandleFunc("/conversion", Conversion)
	http.HandleFunc("/freeinfo", FreeInfo)

	go OfflinePeople()
	err = http.ListenAndServe(":1119", nil)
	if err != nil {
		panic(err)
	}

}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	var Path = r.URL.Path
	pwd, _ := os.Getwd()
	Path = strings.Replace(Path, "/", "\\", -1)
	if !AlrdyLoggedIn && !strings.HasPrefix(Path, "\\images\\locked\\") ||
		AlrdyLoggedIn {

		if AlrdyLoggedIn {
			log.Println("[IMAGE] Accessing " + Path + " as " + strconv.FormatInt(session.UserID, 10))
		} else {
			log.Println("[IMAGE] Accessing " + Path + " as guest.")
		}
		img, err := ioutil.ReadFile(pwd + Path)
		if err != nil {
			log.Println(err.Error())
			Redirect(w, r, "/")
		}
		if strings.HasSuffix(Path, ".gif") {
			w.Header().Set("Content-Type", "image/gif")
		} else if strings.HasSuffix(Path, ".png") {
			w.Header().Set("Content-Type", "image/png")
		}
		w.Header().Set("Cache-Control", "max-age=31536000")
		w.Write(img)
	} else {
		if AlrdyLoggedIn {
			log.Println("[IMAGE] Access denied of " + Path + " as " + strconv.FormatInt(session.UserID, 10))
		} else {
			log.Println("[IMAGE] Access denied of " + Path)
		}
		Redirect(w, r, "/login")
	}
}

func ScriptsHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, _ := IsLoggedIn(r)
	if !strings.HasSuffix(r.URL.Path, ".css") && (!AlrdyLoggedIn || r.URL.Path == "/scripts" || r.URL.Path == "/scripts/") {
		log.Println("[SCRIPTS] access denied!")
		Redirect(w, r, "/login")
	} else {
		//user := FindBaseID(session.UserID)
		//log.Println("[SCRIPTS] access by", user.Username)
		pwd, _ := os.Getwd()
		fs := http.FileServer(http.Dir(pwd + "\\scripts"))
		realhandler := http.StripPrefix("/scripts/", fs).ServeHTTP
		realhandler(w, r)
	}
}

func OfflinePeople() {
	ticker := time.NewTicker(OfflineEvery)
	for {
		<-ticker.C
		log.Println("[OfflinePeople] ticked")
		idsToOffline := make([]int64, 0)
		rows, err := DATABASE.Query("SELECT userID, lastActivityTime FROM userData WHERE activity > 0")
		if err != nil {
			log.Println("[OfflinePeople]", err)
		}
		for rows.Next() {
			var userID, lastActivityTime int64
			rows.Scan(&userID, &lastActivityTime)
			_, present := ClientConnections[userID]
			timePassed := time.Now().UnixNano() - lastActivityTime
			if timePassed > (OfflineEvery).Nanoseconds() && !present {
				idsToOffline = append(idsToOffline, userID)
			}
		}
		rows.Close()
		for _, id := range idsToOffline {
			SetState(id, Offline)
			log.Println("[OfflinePeople] set", id, "offline")
		}

	}
}
