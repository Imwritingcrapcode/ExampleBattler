package Site

import (
	"log"
	"time"
	. "../Abstract"
	"github.com/gorilla/websocket"
	"strconv"
	"net/http"
)

func AddUserQueue(user *ClientChannels) {
	QUEUE = append(QUEUE, user)
	log.Println("Added user:", user.UserID, "\tMain:", user.ChosenGirls[0], "\tSecondary:", user.ChosenGirls[1],
		"\tMainSkill:", user.SkillLevels[0], "\tSecSkill:", user.SkillLevels[1])
}

func RemoveUser(user *ClientChannels) bool {
	for k, other := range QUEUE {
		if other == user {
			QUEUE = append(QUEUE[:k], QUEUE[k+1:]...)
			return true
		}
	}
	return false
}

func EventOrganizer() {
	for {
		user := <-QUEUECHANNEL
		if user.IsTaken {
			log.Println("[Queue]", user.UserID, "is taken already and yet...", user.IsDesperate)
		} else if user.ShouldRemove {
			res := RemoveUser(user)
			if !res {
				log.Panic("[Queue] Not in queue " + strconv.FormatInt(user.UserID, 10))
			}
		} else if !user.IsDesperate {
			AddUserQueue(user)
			for j := 0; j < len(QUEUE); j++ {
				other := QUEUE[j]
				comp, char1, char2 := user.GetCompatibility(other)
				if user.UserID != other.UserID && !other.ShouldRemove && comp == 1 {
					user.PlayingAs = char1
					other.PlayingAs = char2
					user.IsTaken = true
					other.IsTaken = true
					RemoveUser(user)
					RemoveUser(other)
					ClientConnections[user.UserID] = user
					ClientConnections[other.UserID] = other
					go user.Take(other)
					go other.Take(user)
					info := "[Queue] MATCHED " + strconv.FormatInt(user.UserID, 10) + " WITH " + strconv.FormatInt(user.UserID, 10) + " AS IDEAL"
					log.Println(info)
					break
				}
			}
		} else {
			BestGrade := 1000
			BestOpponent := user
			for j := 0; j < len(QUEUE); j++ {
				other := QUEUE[j]
				otherGrade, _, _ := user.GetCompatibility(other)
				if user.UserID != other.UserID && !other.ShouldRemove && BestGrade > otherGrade {
					BestOpponent = other
					BestGrade = otherGrade
				}
			}
			info := "[Queue] MATCHED " + strconv.FormatInt(user.UserID, 10) + " WITH " + strconv.FormatInt(BestOpponent.UserID, 10) + " AS DESPERATE"
			if BestGrade == 1000 {
				user.IsTaken = true
				RemoveUser(user)
				ClientConnections[user.UserID] = user
				go user.Take(user)
			} else {
				user.IsTaken = true
				BestOpponent.IsTaken = true
				RemoveUser(user)
				RemoveUser(BestOpponent)
				comp, char1, char2 := user.GetCompatibility(BestOpponent)
				user.PlayingAs = char1
				BestOpponent.PlayingAs = char2
				ClientConnections[user.UserID] = user
				ClientConnections[BestOpponent.UserID] = BestOpponent
				go user.Take(BestOpponent)
				go BestOpponent.Take(user)
				info += " AND COMPATIBILITY " + strconv.Itoa(comp)
			}
			log.Println(info)
		}
	}
}

func KillIfTheyDisconnect(ws *websocket.Conn, user *ClientChannels) {
	for {
		var msg ClientMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			ws.Close()
			user.Disconnected <- "in WaitForClientMessages"
			return
		}
	}
}

func WaitForConnections(self *ClientChannels) {
	timer2 := time.NewTimer(CONNECTWAITTIME * time.Second)
	// check if connected and propagate forward
	select {
	case message := <-self.Time: //I've connected first.
	 timer2.Stop()
		if message {
			log.Println("[Queue]", FindBaseID(self.UserID).Username, "connected.")
		}
		self.Time <- message
		return
	case message := <-self.Opponent.Time: //Opponent's connected first.
		timer2.Stop()
		if message {
			log.Println("[Queue]", FindBaseID(self.Opponent.UserID).Username, "connected.")
		}
		self.Opponent.Time <- message
		return

	case <-timer2.C: //timed out
		timer2.Stop()
		user := FindBaseID(self.UserID)
		opp := FindBaseID(self.Opponent.UserID)
		log.Println("[Queue] Nobody connected of", user.Username, opp.Username, "killing their game.")
		EndGame(self)
		if user.CurrentActivity == Queuing {
			SetState(self.UserID, BrowsingCharacters)
		}
		if opp.CurrentActivity == Queuing {
			SetState(opp.UserID, BrowsingCharacters)
		}
		return
	}
}

/*func WaitForSoloConn(self *ClientChannels) {
	timer2 := time.NewTimer(CONNECTWAITTIME * time.Second)
	// check if connected
	select {
	case message := <-self.Time: //connected
		if message {
			log.Println("[Queue] I c u connected to a bot game,", self.UserID)
			timer2.Stop()
			return
		} else {
			log.Panic("[Queue] wtf.")
		}
	case <-timer2.C: //timed out
		log.Println("[Queue] Tsk tsk tsk, you queued and did not connect. Really?", self.UserID)
		channels, present := ClientConnections[self.UserID]
		if present {
			channels.State = GaveUp
		}
		user := FindBaseID(self.UserID)
		if user.CurrentActivity == Queuing {
			SetState(user.UserID, BrowsingCharacters)
		}
		EndGame(self)
	}
}*/

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Print("[Queue]" + " redirected to /login")
		Redirect(w, r, "/login")
		return
	}
	if r.Method == http.MethodGet {
		log.Println("[Queue] accessing queue for", session.UserID)
		if GetState(session.UserID) > Queuing && GetState(session.UserID) < JustFinishedTheGame {
			log.Println("[Queue] Terminating your game, ", session.UserID)
			channels, present := ClientConnections[session.UserID]
			if present {
				channels.GiveUp()
				delete(ClientConnections, session.UserID)
			} else {
				log.Println("[Queue] user not found in queue", session.UserID)
			}

		} else if GetState(session.UserID) == Queuing {
			http.Error(w, "You are already in queue", 400)
			return
		}
		// Upgrade initial GET request to a web socket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		//wait for the girls
		//after you receive them, search for a partner
		var userChannels ClientChannels
		var girls ClientMessage
		// Read in a new message as JSON and map it to a Message object
		err = ws.ReadJSON(&girls)
		if err != nil || girls.MainGirl == girls.SecondaryGirl || !HasGirl(session.UserID, girls.MainGirl) ||
			!HasGirl(session.UserID, girls.SecondaryGirl) {
			message := QueueResponse{
				Prompt:   "Error sending data",
				OK:       false,
				Location: "",
			}
			err := ws.WriteJSON(message)
			if err != nil {
				SetState(session.UserID, BrowsingCharacters)
				log.Println("[Queue] rror xc" + err.Error())
			}
			ws.Close()
			return
		}
		userChannels = ClientChannels{
			UserID:         session.UserID,
			Opponent:       nil,
			State:          Queuing,
			ChosenGirls:    []int{girls.MainGirl, girls.SecondaryGirl},
			SkillLevels:    []int{GetSkillLevel(session.UserID, girls.MainGirl), GetSkillLevel(session.UserID, girls.SecondaryGirl)},
			LastThing:      GameState{},
			Input:          make(chan string, 2),
			Output:         make(chan GameState, 4),
			HasGivenUp:     make(chan bool, 1),
			Time:           make(chan bool, 1),
			TimeOutput:     make(chan string, 1),
			KillConnection: make(chan struct{}, 1),
			Taken:          make(chan *ClientChannels, 1),
		}
		userChannels.Opponent = &userChannels
		log.Println("[Queue] Received girls", girls.MainGirl, girls.SecondaryGirl, "from", session.UserID)
		go KillIfTheyDisconnect(ws, &userChannels)
		SetState(session.UserID, Queuing)
		go ListenForOpp(ws, &userChannels)
		QUEUECHANNEL <- &userChannels
	}
}

func ListenForOpp(ws *websocket.Conn, user *ClientChannels) {
	timer := time.NewTimer(QUEUEWAITTIME)
	defer ws.Close()
	for {
		select {
		case why := <-user.Disconnected:
			timer.Stop()
			log.Println("[Queue]", user.UserID, "disconnected from queue with reason:", why)
			if !user.IsTaken {
				user.ShouldRemove = true
				QUEUECHANNEL <- user
			}
			close(user.Input)
			close(user.Output)
			close(user.HasGivenUp)
			close(user.Time)
			close(user.TimeOutput)
			close(user.KillConnection)
			close(user.Disconnected)
			close(user.Taken)
			SetState(user.UserID, BrowsingCharacters)
			return

		case Opp := <-user.Taken:
			timer.Stop()
			user.Opponent = Opp
			log.Println("Queue", user.UserID, "redirected to /game, opponnent - ", Opp.UserID)
			message := QueueResponse{
				Prompt:   "Your opponent is " + FindBaseID(user.Opponent.UserID).Username,
				Location: "/game",
				OK:       true,
			}
			user.State = ReadyingForTheGame

			go WaitForConnections(user)
			err := ws.WriteJSON(message)
			if err != nil {
				SetState(user.UserID, BrowsingCharacters)
				log.Println("[Queue] rror xc" + err.Error())
			} else if !message.OK {
				SetState(user.UserID, BrowsingCharacters)
				log.Println("[Queue] says not ok and", message)
			}
			return
		case <-timer.C:
			timer.Stop()
			if !user.IsTaken && !user.ShouldRemove {
				user.IsDesperate = true
				QUEUECHANNEL <- user
			}
		}
	}
}

