package Site

import (
	. "../Abstract"
	. "../Characters"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// Users

func AddUser(reg_data *RegData) (bool, string) {
	rows, err := DATABASE.Query("SELECT username FROM users WHERE username = '" + reg_data.Username + "'")
	if err != nil {
		log.Println(err)
		return false, "Something(like selection) went horribly wrong"
	}
	defer rows.Close()
	if rows.Next() {
		return false, "There is already a user with that name"
	}

	// adding to users
	currTime := time.Now().UTC().UnixNano()
	reg_data.Password = EncryptPassword(reg_data.Password)
	statement, err := DATABASE.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
		return false, "Something (lke insertion) went horribly wrong"
	}
	_, err = statement.Exec(reg_data.Username, reg_data.Password)
	if err != nil {
		log.Println(err)
		return false, "Something (lke execution) went horribly wrong"
	}
	// and to userData

	var ID int64
	rows, err = DATABASE.Query("SELECT userID FROM users WHERE username = '" + reg_data.Username + "'")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&ID)
	}
	statement, err = DATABASE.Prepare("INSERT INTO userData (userID, lastLoginTime) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
	}
	statement.Exec(ID, currTime)

	//2 free girls
	rand.Seed(currTime)
	free_girl1 := ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	free_girl2 := ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	for free_girl2 == free_girl1 {
		free_girl2 = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
	}
	UnlockGirl(FindBase(reg_data.Username), free_girl1)
	UnlockGirl(FindBase(reg_data.Username), free_girl2)
	return true, "OK"
}

func PassCorrect(pass string, user *User) bool {
	var dbPass string
	if user == nil {
		return false
	}
	rows, err := DATABASE.Query("SELECT password FROM users WHERE username = '" + user.Username + "'")
	if err != nil {
		log.Println(err)
		return false
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&dbPass)
	} else {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass))
	if err != nil {
		log.Println("[Log In] Wrong pass for user", user.Username)
	}
	return err == nil
}

// UserData

func GetFriendData(userID int64, mutual bool) []string {
	var info []string
	var name string
	var activity int
	var duration int64
	if mutual { //get activity
		rows, err := DATABASE.Query("SELECT username, activity, lastActivityTime FROM users INNER JOIN userData ON users.userID = userData.userID WHERE users.userID = " + strconv.FormatInt(userID, 10))
		if err != nil {
			log.Println("[GetFriendData] " + err.Error())
			return nil
		}
		if rows.Next() {
			rows.Scan(&name, &activity, &duration)
			rows.Close()
			actString := ActivitiesToString[activity]
			info = append(info, name, actString)
			if activity == PlayingAs {
				charNum := GetLastPlayedAs(userID)
				info = append(info, ReleasedCharactersNames[charNum])
			} else if activity == Offline {
				curTime := time.Now().UTC().UnixNano()
				info = append(info, strconv.Itoa(int(time.Duration(curTime - duration).Seconds())))
			}
		} else {
			rows.Close()
			log.Panic("[GetFriendData] no such mutual friend found uwu")
		}
	} else { //don't get activity
		rows, err := DATABASE.Query("SELECT username FROM users WHERE userID = " + strconv.FormatInt(userID, 10))
		if err != nil {
			log.Println("[GetFriendData] " + err.Error())
			return nil
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&name)
			info = append(info, name)
		} else {
			log.Panic("[GetFriendData] no such friend found owo")
		}
	}
	return info
}

func FindBase(username string) *User {
	rows, err := DATABASE.Query("SELECT userData.userID, username, activity FROM users INNER JOIN userData ON users.userID = userData.userID WHERE username = '" + username + "'")
	if err != nil {
		log.Println("[FindBase] " + err.Error())
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		data := User{}
		rows.Scan(&data.UserID, &data.Username, &data.CurrentActivity)
		return &data
	}
	return nil
}

func FindBaseID(ID int64) *User {
	//log.Println("in FindBase ~!~")
	rows, err := DATABASE.Query("SELECT userData.userID, username, activity FROM users INNER JOIN userData ON users.userID = userData.userID WHERE users.userID = " + strconv.FormatInt(ID, 10))
	if err != nil {
		log.Println("[FindBaseID] " + err.Error())
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		data := User{}
		rows.Scan(&data.UserID, &data.Username, &data.CurrentActivity)
		return &data
	}
	return nil
}

func SetState(id int64, state int) {
	previousState := GetState(id)
	if previousState <= Offline && state > Offline {
		NotifyFriends(id)
	}
	log.Println("[SetState] updating lastactivity", id, ActivitiesToString[state], ActivitiesToString[previousState], previousState <= Offline, state > Offline)
	update, err := DATABASE.Prepare("UPDATE userData SET lastActivityTime = " + strconv.FormatInt(time.Now().UnixNano(), 10) + " WHERE  userID = " + strconv.FormatInt(id, 10))
	if err != nil {
		log.Println("[SetState] updating lastactivity", err)
		return
	}
	_, err = update.Exec()
	if err != nil {
		log.Println("[SetState] updating lastactivity", err)
		return
	}

	update, err = DATABASE.Prepare("UPDATE userData SET activity = " + strconv.Itoa(state) + " WHERE  userID = " + strconv.FormatInt(id, 10))
	if err != nil {
		log.Println("[SetState] updating activity", err)
		return
	}
	_, err = update.Exec()
	if err != nil {
		log.Println("[SetState] updating activity", err)
		return
	}
}

func GetState(id int64) int {
	rows, err := DATABASE.Query("SELECT activity FROM userData WHERE userID = " + strconv.FormatInt(id, 10))
	if err != nil {
		log.Println(err)
		return -1
	}
	defer rows.Close()
	if rows.Next() {
		var state int
		rows.Scan(&state)
		return state
	}
	log.Println("[GetState] Somehow no state has been found for user", id)
	return -1
}

func IncreaseBattles(userID int64, won bool) {
	var statement *sql.Stmt
	var err error
	if won {
		statement, err = DATABASE.Prepare("UPDATE userData SET battlesTotal = battlesTotal + 1, battlesWon = battlesWon + 1 " +
			"WHERE userID = " + strconv.FormatInt(userID, 10))
		if err != nil {
			log.Println(err)
		}
	} else {
		statement, err = DATABASE.Prepare("UPDATE userData SET battlesTotal = battlesTotal + 1 WHERE userID = " + strconv.FormatInt(userID, 10))
		if err != nil {
			log.Println(err)
		}
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}
}

func (user *User) GetBattles() int {
	var battles int
	rows, err := DATABASE.Query("SELECT battlesTotal FROM userData WHERE userID = " + strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println(err)
		return -1
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&battles)
	} else {
		log.Println(user.Username + "does not have any battles?!")
	}
	return battles
}

func (user *User) GetWonBattles() int {
	var battles int
	rows, err := DATABASE.Query("SELECT battlesWon FROM userData WHERE userID = " + strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println(err)
		return -1
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&battles)
	} else {
		log.Println(user.Username + "does not have any won battles?!")
	}
	return battles
}

func (user *User) GetDust(letter string) int {
	//log.Println("in GetDust ~!~")
	var dust int
	value, present := DustMap[letter]
	if !present {
		log.Println("[Database] wrong argument for GetDust")
		return -1
	}
	rows, err := DATABASE.Query("SELECT " + value + " FROM userData WHERE userID = " + strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println(err)
		return -1
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&dust)
	} else {
		log.Println(user.Username + "does not have any dust?!")
	}
	return dust
}

func (user *User) SetDust(letter string, amount int) {
	value, present := DustMap[letter]
	if !present {
		log.Println("[Database] wrong argument for SetDust")
	}
	statement, err := DATABASE.Prepare("UPDATE userData SET " + value + " = " + strconv.Itoa(amount) +
		" WHERE userID = " + strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println(err)
	}
	statement.Exec()
}

func GetLastPlayedAs(userID int64) int {
	var lastPlayedAs int
	rows, err := DATABASE.Query("SELECT lastPlayedAs FROM userData WHERE userID = " + strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println(err)
		return -1
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&lastPlayedAs)
	} else {
		log.Println(strconv.FormatInt(userID, 10) + "does not have lastPlayedAs?!")
	}
	return lastPlayedAs
}

func SetLastPlayedAs(userID int64, playedAS int) {
	statement, err := DATABASE.Prepare("UPDATE userData SET lastPlayedAs = " + strconv.Itoa(playedAS) +
		" WHERE userID = " + strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println(err)
	}
	statement.Exec()

}

// Girls

func UnlockGirl(user *User, girl int) {
	if HasGirl(user.UserID, girl) {
		log.Println("[UnlockGirl]", user.Username, "already has", girl)
		return
	}
	statement, err := DATABASE.Prepare("INSERT INTO girls (userID, girlNumber) VALUES (?, ?)")
	if err != nil {
		log.Println("[UnlockGirl]", err)
		return
	}
	_, err = statement.Exec(user.UserID, girl)
	if err != nil {
		log.Println(err)
		log.Println("[UnlockGirl]", err)
	}

}

func HasGirl(userID int64, girl int) bool {
	rows, err := DATABASE.Query("SELECT girlNumber FROM girls WHERE userID = " + strconv.FormatInt(userID, 10) +
		" AND girlNumber = " + strconv.Itoa(girl))
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	return rows.Next()
}

func GetGirls(id int64) []UserGirl {
	var girls []UserGirl
	rows, err := DATABASE.Query("SELECT girlNumber, girlLevel, matchesPlayed, matchesWon FROM girls WHERE userID = " + strconv.FormatInt(id, 10))
	if err != nil {
		log.Println("[FindBase] GetGirls: " + err.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		/// ???
		data := UserGirl{}
		rows.Scan(&data.Number, &data.Level, &data.MatchesPlayed, &data.MatchesWon)
		girls = append(girls, data)
	}
	//sort
	sort.Slice(girls, func(i, j int) bool {
		return girls[i].Number < girls[j].Number
	})

	if girls == nil {
		log.Println("[GetGirls] ", id, " does not have any girls D:")
	}
	return girls
}

func GetTotalMatches(userID int64, girl int) int {
	rows, err := DATABASE.Query("SELECT matchesPlayed FROM girls WHERE userID = " + strconv.FormatInt(userID, 10) +
		" AND girlNumber = " + strconv.Itoa(girl))
	if err != nil {
		log.Println("[GetTotalMatches] " + err.Error())
		return -1
	}
	defer rows.Close()
	var num int
	if rows.Next() {
		rows.Scan(&num)
		//log.Println(girl, "\t\t|\t", num, "\t\t|\t\t", level)
		return num
	} else {
		log.Println("[GetTotalMatches] girl not found for", userID, girl)
		return -1
	}
}

func GetMatchData(userID int64, girl int) []int {
	stuff := make([]int, 2)
	rows, err := DATABASE.Query("SELECT girlLevel, currMatches FROM girls WHERE userID = " + strconv.FormatInt(userID, 10) +
		" AND girlNumber = " + strconv.Itoa(girl))
	if err != nil {
		log.Println("[FindBase] GetMatchData: " + err.Error())
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&stuff[0], &stuff[1])
	}
	return stuff
}

func LevelUp(userID int64, as int, lt20 bool) {
	log.Println("[LevelUp] Levelled up", as, "for", userID)
	if lt20 {
		statement, err := DATABASE.Prepare("UPDATE girls SET girlLevel = girlLevel + 1" +
			" WHERE userID = " + strconv.FormatInt(userID, 10) + " AND girlNumber = " + strconv.Itoa(as))
		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Println(err)
		}
	}
	statement, err := DATABASE.Prepare("UPDATE girls SET currMatches = 0" +
		" WHERE userID = " + strconv.FormatInt(userID, 10) + " AND girlNumber = " + strconv.Itoa(as))
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}
}

func IncreaseMatchesAs(userID int64, as int, won, gaveUp bool) {
	var statement *sql.Stmt
	var err error
	//1, get the amnt of matches and the girl's lvl
	stuff := GetMatchData(userID, as) //0 - girllevel, 1 - curr matches
	level := stuff[0]
	curr := stuff[1] + 1
	//2. lvl up if needed
	if level < 20 && curr >= GetLevelCaps(GetGirlRarity(as))[level-1] {
		LevelUp(userID, as, level < 20)
	} else if !gaveUp {
		//2, no lvl up, give xp if we didn't give up
		statement, err = DATABASE.Prepare("UPDATE girls SET currMatches = currMatches + 1" +
			" WHERE userID = " + strconv.FormatInt(userID, 10) + " AND girlNumber = " + strconv.Itoa(as))
		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Println(err)
		}
	}
	//3. increase total play
	statement, err = DATABASE.Prepare("UPDATE girls SET matchesPlayed = matchesPlayed + 1" +
		" WHERE userID = " + strconv.FormatInt(userID, 10) + " AND girlNumber = " + strconv.Itoa(as))
	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}

	//4. if she won, +1 won match ez
	if won {
		statement, err = DATABASE.Prepare("UPDATE girls SET matchesWon = matchesWon + 1" +
			" WHERE userID = " + strconv.FormatInt(userID, 10) + " AND girlNumber = " + strconv.Itoa(as))
		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Println(err)
		}
	}
}

// Rewards

func AddRewards(user *User, kind string, amount int) {
	//log.Println("in AddRewards ~!~")
	//1. does that user already have that type of reward?
	//2. if so, add some more
	//3. if nah, add a new one
	rows, err := DATABASE.Query("SELECT rewardType FROM rewards WHERE userID = " +
		strconv.FormatInt(user.UserID, 10) + " AND rewardType = '" + kind + "'")

	if err != nil {
		log.Println("[Database] AddRewards: " + err.Error())
		return
	}

	//dust gets replaced, not dust is ok
	_, isDust := DustMap[kind]
	if rows.Next() && !isDust {
		rows.Close()
		statement, err := DATABASE.Prepare("UPDATE rewards SET amount = amount + " +
			strconv.Itoa(amount) + " WHERE userID = " + strconv.FormatInt(user.UserID, 10) +
			" AND rewardType = '" + kind + "'")

		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Println(err)
		}

	} else if rows.Next() && isDust {
		rows.Close()
		statement, err := DATABASE.Prepare("UPDATE rewards SET amount = " +
			strconv.Itoa(amount) + " WHERE userID = " + strconv.FormatInt(user.UserID, 10) +
			" AND rewardType = '" + kind + "'")

		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Println(err)
		}
	} else {
		rows.Close()
		statement, err := DATABASE.Prepare("INSERT INTO rewards (userID, rewardType, amount) VALUES (?, ?, ?)")
		if err != nil {
			log.Println(err)
		}
		_, err = statement.Exec(user.UserID, kind, amount)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetRewards(user *User) *RewardsObj {
	rows, err := DATABASE.Query("SELECT rewardType, amount FROM rewards WHERE userID = " +
		strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println("[GetRewards] GetRewards: " + err.Error())
		return nil
	}
	defer rows.Close()
	number := GetLastPlayedAs(user.UserID)
	playedAs := strconv.Itoa(number)
	rewards := RewardsObj{
		BattleResult:      0,
		LastOpponentsName: "",
		Dusts:             make(map[string]int, 5),
		Name:              ReleasedCharactersNames[number],
		Matches:           GetLevelCaps(GetGirlRarity(number)),
		TotalMatches:      GetTotalMatches(user.UserID, number), //TODO change to level +currMAtches
		ToAdd:             0,
	}
	var rtype string
	var amnt int
	for rows.Next() {
		rows.Scan(&rtype, &amnt)
		_, isDust := DustMap[rtype]
		if isDust {
			rewards.Dusts[rtype] += amnt
		} else if playedAs == rtype { //found this gurl
			rewards.ToAdd = amnt
		}
	}

	if rewards.ToAdd > 0 {
		rewards.LastOpponentsName = GetLastOpponentsName(user.UserID)
		rewards.BattleResult = GetLastBattleResult(user.UserID)
	}

	return &rewards
}

func DeleteRewards(user *User) {
	del, err := DATABASE.Prepare("DELETE FROM rewards WHERE userID = " + strconv.FormatInt(user.UserID, 10))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = del.Exec()
	if err != nil {
		log.Println(err)
	}
	log.Println("[Database] rewards deleted for " + strconv.FormatInt(user.UserID, 10))
}

func SetLastBattleResult(UserID int64, EndState int) {
	update, err := DATABASE.Prepare("UPDATE userData SET lastBattleResult = " + strconv.Itoa(EndState) + " WHERE  userID = " + strconv.FormatInt(UserID, 10))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = update.Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func GetLastBattleResult(UserID int64) int {
	rows, err := DATABASE.Query("SELECT lastBattleResult FROM userData WHERE userID = " + strconv.FormatInt(UserID, 10))
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()
	if rows.Next() {
		var state int
		rows.Scan(&state)
		return state
	}
	log.Println("[GetLastBattleResult] Somehow no res has been found for user" + strconv.FormatInt(UserID, 10))
	return 0
}

func SetLastOpponentsName(UserID1 int64, UserID2 int64) {
	var name string
	if UserID1 != UserID2 {
		rows, err := DATABASE.Query("SELECT username FROM users WHERE userID = " + strconv.FormatInt(UserID2, 10))
		if err != nil {
			log.Println("[GetLastOpponentsName] " + err.Error())
		}
		if rows.Next() {
			rows.Scan(&name)
		} else {
			log.Panic("[GetLastOpponentsName] no such person found owo")
		}
		rows.Close()
	} else {
		index := rand.Intn(len(BotNames))
		name = BotNames[index]
		for name == FindBaseID(UserID1).Username {
			index = rand.Intn(len(BotNames))
			name = BotNames[index]
		}
	}
	update, err := DATABASE.Prepare("UPDATE userData SET lastOpponentsName = '" + name + "' WHERE  userID = " + strconv.FormatInt(UserID1, 10))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = update.Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func GetLastOpponentsName(UserID int64) string {
	rows, err := DATABASE.Query("SELECT lastOpponentsName FROM userData WHERE userID = " + strconv.FormatInt(UserID, 10))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer rows.Close()
	if rows.Next() {
		var state string
		rows.Scan(&state)
		return state
	}
	log.Println("[GetLastOpponentsName] no last name??? " + strconv.FormatInt(UserID, 10))
	return ""
}

// Sessions

func AddSession(session *Session) {
	//delete other sessions for the same user
	del, err := DATABASE.Prepare("DELETE FROM sessions WHERE userID = " + strconv.FormatInt(session.UserID, 10))
	if err != nil {
		log.Println(err)
	}
	_, err = del.Exec()
	if err != nil {
		log.Println(err)
	}
	statement, err := DATABASE.Prepare("INSERT INTO sessions (sessionID, userID, expires) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = statement.Exec(session.ID, session.UserID, session.Expires)
	if err != nil {
		log.Println(err)
	}

}

func IsActiveSession(ID string) *Session {
	var session Session
	rows, err := DATABASE.Query("SELECT sessionID, userID, expires FROM sessions WHERE sessionID = '" + ID + "'")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&session.ID, &session.UserID, &session.Expires)
		return &session

	} else if rows.Err() != nil {
		log.Println("[isActiveSession] Something is horribly wrong with that row(s): " + rows.Err().Error())
	}
	return nil
}

func ClearSessions() time.Time {
	currTime := time.Now()
	statement, err := DATABASE.Prepare("UPDATE userData SET activity = 0 WHERE userData.userID IN" +
		"(SELECT sessions.userID FROM sessions WHERE sessions.expires < " + strconv.FormatInt(currTime.UnixNano(), 10) + ")")

	if err != nil {
		log.Println(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println(err)
	}

	del, err := DATABASE.Prepare("DELETE FROM sessions WHERE expires < " + strconv.FormatInt(currTime.UnixNano(), 10))
	if err != nil {
		log.Println(err)
	}
	_, err = del.Exec()
	if err != nil {
		log.Println(err)
	}
	log.Println("[Sessions] Cleared sessions!!")
	return currTime
}

func DeleteSession(session *Session) {
	del, err := DATABASE.Prepare("DELETE FROM sessions WHERE sessionID = '" + session.ID + "'")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = del.Exec()
	if err != nil {
		log.Println(err)
	}
	log.Println("[Sessions] session deleted for " + strconv.FormatInt(session.UserID, 10))
}

//Misc

func (user *User) GatherFreeData() *UserFree {
	data := UserFree{}
	moneyI := MoneyInfo{}
	moneyI.W = user.GetDust("w")
	moneyI.B = user.GetDust("b")
	moneyI.Y = user.GetDust("y")
	moneyI.P = user.GetDust("p")
	moneyI.S = user.GetDust("s")
	data.Monies = moneyI
	data.Username = user.Username
	data.BattlesTotal = user.GetBattles()
	data.BattlesWon = user.GetWonBattles()
	return &data
}

//Friends

func AddFriend(userID int64, friendID int64) {
	if IsFriend(userID, friendID) {
		log.Println("[AddFriend]", userID, "and", friendID, "are already friends!")
		return
	}
	statement, err := DATABASE.Prepare("INSERT INTO friends (userID, friendID) VALUES (?, ?)")
	if err != nil {
		log.Println("[AddFriend]", err)
		return
	}
	_, err = statement.Exec(userID, friendID)
	if err != nil {
		log.Println(err)
		log.Println("[AddFirned]", err)
	}
}

func RemoveFriend(userID int64, friendID int64) {
	if !IsFriend(userID, friendID) {
		log.Println("[RemoveFriend]", userID, "and", friendID, "were not friends!")
		return
	}

	del, err := DATABASE.Prepare("DELETE FROM friends WHERE userID = " + strconv.FormatInt(userID, 10) +
		" AND friendID = " + strconv.FormatInt(friendID, 10))
	if err != nil {
		log.Println("[RemoveFriend]", err)
		return
	}
	_, err = del.Exec()
	if err != nil {
		log.Println("[RemoveFriend]", err)
	}
	log.Println("[RemovedFriend] for", userID, "and", friendID)
}

func IsFriend(userID, friendID int64) bool {
	rows, err := DATABASE.Query("SELECT friendID FROM friends WHERE userID = " + strconv.FormatInt(userID, 10) +
		" AND friendID = " + strconv.FormatInt(friendID, 10))
	if err != nil {
		log.Println("[IsFriend]", userID, friendID, err)
	}
	defer rows.Close()
	return rows.Next()
}

func GetFriendLists(userID int64) *FriendList {
	var from []int64
	//1. select all that are from me
	rows, err := DATABASE.Query("SELECT friendID FROM friends WHERE userID = " + strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println("[GetFriendLists]", userID, err)
	}
	var friendID int64
	for rows.Next() {
		rows.Scan(&friendID)
		from = append(from, friendID)
	}
	rows.Close()
	//2. select all that are to me
	var to []int64
	rows, err = DATABASE.Query("SELECT userID FROM friends WHERE friendID = " + strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println("[GetFriendLists2]", userID, err)
	}
	for rows.Next() {
		rows.Scan(&friendID)
		to = append(to, friendID)
	}
	rows.Close()
	//3. find the intersection (take out pending and all that)
	var incoming, pending []string
	var friends [][]string
	for _, id := range from { //for every id that is in from, get info and check if they are in to
		i := Contains(to, id)
		mutual := i != -1
		info := GetFriendData(id, mutual)
		if mutual { //mutual friends
			to = append(to[:i], to[i+1:]...)
			friends = append(friends, info)
		} else { //they did not accept(yet?)
			pending = append(pending, info[0])
		}
	}
	//now the other way around: go through the to-requests that are left
	for _, id := range to {
		info := GetFriendData(id, false)[0]
		incoming = append(incoming, info)
	}
	//4. replace ids with names and add statuses where applicable
	res := FriendList{
		Friends:  friends,
		Incoming: incoming,
		Pending:  pending,
	}
	return &res

}

//CONVERSION
func GetConversionInfo(userID int64) (bool, int, int, string, int, int) {
	rows, err := DATABASE.Query("SELECT begins, duration, type, give, notified FROM conversions WHERE userID = " + strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println("[GetConversionInfo]", userID, err)
	}
	if rows.Next() {
		var begins int64
		var dtype string
		var amnt, duration, notified int
		rows.Scan(&begins, &duration, &dtype, &amnt, &notified)
		rows.Close()
		curTime := time.Now().UnixNano()
		secondsPassed := int(time.Duration(curTime - begins).Seconds())
		if secondsPassed >= duration {
			return true, duration, 0, dtype, amnt, notified
		} else {
			return true, secondsPassed, duration - secondsPassed, dtype, amnt, notified
		}
	} else {
		return false, -1, -1, "", -1, 0
	}
}

func ClaimConversion(UserID int64) {
	isConverting, _, secondsLeft, dustType, amnt, _ := GetConversionInfo(UserID)
	if !isConverting || secondsLeft > 0 {
		log.Println("[Claim Conversion] called too early")
		return
	}

	//2. Delet the old conversion
	del, err := DATABASE.Prepare("DELETE FROM conversions WHERE userID = " + strconv.FormatInt(UserID, 10))
	if err != nil {
		log.Println("[ClaimConversion]", err)
		return
	}
	_, err = del.Exec()
	if err != nil {
		log.Println("[ClaimConversion]", err)
	}
	//3. Add the moneys
	user := User{UserID: UserID}
	user.SetDust(dustType, user.GetDust(dustType)+amnt)

}

func StartConversion(UserID int64, duration, amnt int, dustType string) {
	begins := time.Now().UnixNano()
	statement, err := DATABASE.Prepare("INSERT INTO conversions (userID, begins, duration, give, type, notified) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("[StartConversion]", err)
		return
	}
	_, err = statement.Exec(UserID, begins, duration, amnt, dustType, 0)
	if err != nil {
		log.Println(err)
		log.Println("[StartConversion]", err)
	}
}

func NotifiedConversion(UserID int64) {
	statement, err := DATABASE.Prepare("UPDATE conversions SET notified = 1 WHERE userID = " + strconv.FormatInt(UserID, 10))
	if err != nil {
		return
	}
	statement.Exec()
}

//notifications

func AddNotification(UserID int64, text, redirect string) {
	statement, err := DATABASE.Prepare("INSERT INTO notifications (userID, text, redirect, seen) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("[AddNotification]", err)
		return
	}
	_, err = statement.Exec(UserID, text, redirect, 0)
	if err != nil {
		log.Println(err)
		log.Println("[AddNotification]", err)
	} else {
		log.Println("[AddNotification] Notification added", UserID, text)
	}
}

func GetNotifications(UserID int64) [][]string {
	var notifications [][]string
	var ids []int
	notifications = make([][]string, 0)
	ids = make([]int, 0)
	rows, err := DATABASE.Query("SELECT text, redirect, NotifID FROM notifications WHERE userID = " + strconv.FormatInt(UserID, 10) +
		" AND seen = 0")
	if err != nil {
		log.Println("[GetNotifications]", UserID, err)
	}
	for rows.Next() {
		var text, redirect string
		var id int
		rows.Scan(&text, &redirect, &id)
		notifications = append(notifications, []string{text, redirect})
		ids = append(ids, id)
	}
	rows.Close()

	for _, v := range ids {
		statement, _ := DATABASE.Prepare("UPDATE notifications SET seen = 1 WHERE NotifID = " + strconv.Itoa(v))
		statement.Exec()
	}
	return notifications
}

func DeleteNotifications(UserID int64, kind string) {
	if kind == "all" {
		del, err := DATABASE.Prepare("DELETE FROM notifications WHERE userID = " + strconv.FormatInt(UserID, 10))
		if err != nil {
			log.Println("[DeleteNotifications]", err)
			return
		}
		_, err = del.Exec()
		if err != nil {
			log.Println("[DeleteNotifications]", err)
		}
	} else if kind == "seen" {
		del, err := DATABASE.Prepare("DELETE FROM notifications WHERE userID = " + strconv.FormatInt(UserID, 10) +
			" AND seen = 1")
		if err != nil {
			log.Println("[DeleteNotifications]", err)
			return
		}
		_, err = del.Exec()
		if err != nil {
			log.Println("[DeleteNotifications]", err)
		}
	} else {
		del, err := DATABASE.Prepare("DELETE FROM notifications WHERE userID = " + strconv.FormatInt(UserID, 10) +
			" AND redirect = " + kind)
		if err != nil {
			log.Println("[DeleteNotifications]", err)
			return
		}
		_, err = del.Exec()
		if err != nil {
			log.Println("[DeleteNotifications]", err)
		}
	}
}

func SeeNotifications(UserID int64, kind string) {
	if kind == "all" {
		statement, err := DATABASE.Prepare("UPDATE notifications SET seen = 1 WHERE userID = " + strconv.FormatInt(UserID, 10))
		if err != nil {
			return
		}
		statement.Exec()
	} else {
		statement, err := DATABASE.Prepare("UPDATE notifications SET seen = 1 WHERE userID = " + strconv.FormatInt(UserID, 10) +
			" AND redirect = " + kind)
		if err != nil {
			return
		}
		statement.Exec()
	}
}

func NotifyFriends(userID int64) {
	friends := GetFriendLists(userID)
	us := FindBaseID(userID)
	for _, v := range friends.Friends {
		friend := FindBase(v[0])
		if friend.CurrentActivity > 0 { //is online
			AddNotification(friend.UserID, "Your friend <b>"+us.Username+"</b> has just come online!", "friends")
			log.Println("[Notifications] added friend notif for", friend.Username, us.Username)
		}
	}
}
