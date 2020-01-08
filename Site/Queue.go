package Site

import (
	. "../Abstract"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Print("[Queue]" + " redirected to /login")
		Redirect(w, r, "/login")
	} else if r.Method == http.MethodGet {
		user := FindBaseID(session.UserID)
		log.Println("[Queue] accessing queue for", user.Username)
		if user.GetState() > Queuing && user.GetState() < JustFinishedTheGame {
			log.Println("[Queue] Terminating your game, " + user.Username)
			channels, present := ClientConnections[user.UserID]
			if present {
				channels.Input <- "GiveUp"
			} else {
				log.Println("[Queue] user not found in queue", user.Username)
			}

		} else if user.GetState() == Queuing {
			http.Error(w, "You are already in queue", 400)
			return
		}
		SetState(user.UserID, Queuing)

		// Upgrade initial GET request to a web socket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()
		//Search for a partner while hoping the client will not d/c
		outputChannel := make(chan QueueResponse, 2)
		inputChannel := make(chan ClientMessage, 1)
		var userChannels ClientChannels
		go WaitForClientMessages(ws, outputChannel, inputChannel)

		for {
			select {
			case girls, ok := <-inputChannel:
				{
					if ok {
						//configure client channels
						if err != nil || girls.MainGirl == girls.SecondaryGirl || !HasGirl(user.UserID, girls.MainGirl) ||
							!HasGirl(user.UserID, girls.SecondaryGirl) {
							message := QueueResponse{
								Prompt:   "Error sending data",
								OK:       false,
								Location: "",
							}
							outputChannel <- message
							log.Println("[Queue] Error sending data", user.Username)
						} else {
							input := make(chan string, 2)
							output := make(chan GameState, 4)
							gg := make(chan bool, 1)
							Time := make(chan bool, 1)
							TimeO := make(chan string, 1)
							kill := make(chan struct{}, 1)
							userChannels = ClientChannels{
								UserID:         user.UserID,
								Opponent:       nil,
								State:          Queuing,
								ChosenGirls:    []int{girls.MainGirl, girls.SecondaryGirl},
								LastThing:      GameState{},
								Input:          input,
								Output:         output,
								HasGivenUp:     gg,
								Time:           Time,
								TimeOutput:     TimeO,
								KillConnection: kill,
							}
							log.Println("[Queue] Received girls", girls.MainGirl, girls.SecondaryGirl, "from", user.Username)

							//add to the client map
							QueueClients[user.UserID] = &userChannels

							//push to the queue
							queueID := UserQueue.Len()
							UserQueue.Push(&Item{
								UserID:   user.UserID,
								Priority: time.Now().Unix(),
								Index:    queueID,
							})

							//wait for the girls
							//after you receive them, search for a partner
							SearchForPartner(&userChannels, user.UserID, outputChannel)
						}
					}
				}
			case message := <-outputChannel:
				{
					if message == (QueueResponse{}) {
						log.Println("[Queue] Client DCed")
						//client d/ced
						UserQueue.Remove(user.UserID)
						delete(QueueClients, user.UserID)
						SetState(user.UserID, BrowsingCharacters)
						return
					} else {
						//send the opp + location message to the client
						if message.OK {
							ClientConnections[user.UserID] = &userChannels
							log.Println("[Queue] added", userChannels.UserID, "to cli con. deleted from qu ch.")
						}
						err := ws.WriteJSON(message)
						if err != nil {
							delete(QueueClients, user.UserID)
							log.Println("[Queue] rror xc" + err.Error())
						} else if !message.OK {
							delete(QueueClients, user.UserID)
							SetState(user.UserID, BrowsingCharacters)
							log.Println("[Queue] says not ok and", message)
						}
						delete(QueueClients, user.UserID) //???
						return

					}
				}
			}
		}
	}
}

//finds a partner for me. is supposed to set my opponent to the correct one in my client channels. self as opponent = bot
//everyone removes THEMSELVES from the queue, but not the other person
//but the one who found the opponent first sets their state to readying
func SearchForPartner(self *ClientChannels, queueID int64, msgchnl chan QueueResponse) {
	log.Println("[Queue] Searching for partner as " + strconv.FormatInt(self.UserID, 10))
	timer := time.NewTimer(WAITTIME)
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-timer.C:
			{ //checking one last time
				if self.State >= ReadyingForTheGame {
					UserQueue.Remove(queueID)
					msgchnl <- QueueResponse{
						Prompt:   "Your opponent is " + FindBaseID(self.UserID).Username,
						Location: "/game",
						OK:       true,
					}
					return
				}
				//start a bot game bc foreveralone
				log.Println("[Queue] couldn't find player opponent for ", FindBaseID(self.UserID).Username)
				log.Println("[Queue] " + FindBaseID(self.UserID).Username + "'s state now is " + "Found opponent")
				self.Opponent = self
				self.State = ReadyingForTheGame
				UserQueue.Remove(queueID)
				msgchnl <- QueueResponse{
					Prompt:   "Your opponent is " + FindBaseID(self.UserID).Username,
					Location: "/game",
					OK:       true,
				}
				go WaitForSoloConn(self)
				ticker.Stop()
				timer.Stop()
				return
			}
		case <-ticker.C:
			if self.State >= ReadyingForTheGame {
				UserQueue.Remove(queueID)
				msgchnl <- QueueResponse{
					Prompt:   "Your opponent is " + FindBaseID(self.Opponent.UserID).Username,
					Location: "/game",
					OK:       true,
				}
				ticker.Stop()
				timer.Stop()
				return
			}
			queue := make(PriorityQueue, len(UserQueue))
			copy(queue, UserQueue)
			for _, item := range queue {
				other, present := QueueClients[item.UserID]
				if present && other.State == Queuing && self.State == Queuing && self.UserID != other.UserID {
					UserQueue.Remove(queueID)
					self.Opponent = other
					other.Opponent = self
					other.State = ReadyingForTheGame
					log.Println("[Queue] " + FindBaseID(other.Opponent.UserID).Username + "'s state now is " + "Found opponent")
					self.State = ReadyingForTheGame
					msgchnl <- QueueResponse{
						Prompt:   "Your opponent is " + FindBaseID(self.Opponent.UserID).Username,
						Location: "/game",
						OK:       true,
					}
					log.Println("[Queue] " + FindBaseID(self.UserID).Username + "'s state now is " + "Found opponent")
					log.Println("[Queue] "+"found an opponent for", FindBaseID(self.UserID).Username, "and", FindBaseID(other.UserID).Username)
					log.Println("[Queue] Waiting for one of them to connect...")
					ticker.Stop()
					timer.Stop()
					go WaitForConnections(self)
					return
				}
			}
		}
	}
}

func WaitForClientMessages(ws *websocket.Conn, msgchnl chan QueueResponse, clientchnl chan ClientMessage) {
	stillGoing := true
	for {
		var msg ClientMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			msgchnl <- QueueResponse{}
			return
		} else if stillGoing {
			clientchnl <- msg
			stillGoing = false
		} else {
			close(clientchnl)
		}
	}
}

func WaitForConnections(self *ClientChannels) {
	timer2 := time.NewTimer(CONNECTWAITTIME * time.Second)
	// check if connected
	select {
	case message := <-self.Time: //I've connected first.
		if message {
			log.Println("[Queue]", FindBaseID(self.UserID).Username, "connected first.")
		}
		timer2.Stop()
		return
	case message := <-self.Opponent.Time: //Opponent's connected first.
		if message {
			log.Println("[Queue]", FindBaseID(self.Opponent.UserID).Username, "connected first.")
		}
		timer2.Stop()
		return

	case <-timer2.C: //timed out
		log.Println("[Queue] Nobody connected of", FindBaseID(self.UserID).Username, FindBaseID(self.Opponent.UserID).Username, "killing their game.")
		EndGame(self)
		SetState(self.UserID, BrowsingCharacters)
		SetState(self.Opponent.UserID, BrowsingCharacters)
	}
	return
}

func WaitForSoloConn(self *ClientChannels) {
	timer2 := time.NewTimer(CONNECTWAITTIME * time.Second)
	// check if connected
	select {
	case message := <-self.Time: //connected
		if message {
			log.Println("[Queue] I c u connected to a bot game,", self.UserID)
			timer2.Stop()
			return
		} else {
			log.Panic("[queue] wtf.")
		}
	case <-timer2.C: //timed out
		log.Println("[Queue] Tsk tsk tsk, you queued and did not connect. Really?")
		channels, present := ClientConnections[self.UserID]
		if present {
			channels.State = GaveUp
		}
		SetState(self.UserID, BrowsingCharacters)
		EndGame(self)
	}
}
