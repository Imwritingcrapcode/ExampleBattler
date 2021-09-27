package Site

import (
	"net/http"
	"log"
	"os"
	"strings"
	"github.com/gorilla/websocket"
	"strconv"
	. "../Abstract"
)

// lobby HERE FOR NOW BC IT'S NOT DONE YET
var LobbyClients = map[int64]*LobbyClient{}

const (
	LeftTheLobby           = iota
	Error
	LobbyList
	ConnectedBlue
	ConnectedRed
	SpectatorConnectedBlue
	SpectatorConnectedRed
	SentGirls
	Ready
	Unready
	GameStart
	LeaderChange
	//--------------- they can't see past this point but LobbyHan can
	ExplainYourself
	ChangeState
	ConnectedFirst
)

type LobbyResponse struct {
	Event int    `json:"Event"`
	Who   string `json:"Who"`
}

type LobbyClient struct {
	UserID        int64
	Username      string
	InputChannel  chan string        //what the player says
	OutputChannel chan LobbyResponse //what we say
	Girls         []int
	State         int
	Team          bool
	IsLeader      bool
	LobbyHandler  *LobbyHandler
}

type LobbyHandler struct {
	Guests         []*LobbyClient
	Length         int
	Capacity       int
	SpecCapacity   int
	ReadyAmount    int
	BlueAmount     int
	RedAmount      int
	BlueSpecAmount int
	RedSpecAmount  int
	Leader         string
	Requests       chan LobbyResponse
}

func (this *LobbyHandler) ExplainYourself() string {
	var d string
	d += "Length: " + strconv.Itoa(this.Length) + "\n"
	d += "Capacity: " + strconv.Itoa(this.Capacity) + "\n"
	d += "SpecCapacity: " + strconv.Itoa(this.SpecCapacity) + "\n"
	d += "ReadyAmount: " + strconv.Itoa(this.ReadyAmount) + "\n"
	d += "BlueAmount: " + strconv.Itoa(this.BlueAmount) + "\n"
	d += "RedAmount: " + strconv.Itoa(this.RedAmount) + "\n"
	d += "BlueSpecAmount: " + strconv.Itoa(this.BlueSpecAmount) + "\n"
	d += "RedSpecAmount: " + strconv.Itoa(this.RedSpecAmount) + "\n"
	d += "Leader: " + this.Leader + "\n"
	return d
}

func (this *LobbyHandler) ListenAndManageTheLobby() {
	for {
		req := <-this.Requests
		if req.Event == LeftTheLobby {
			remove := -1
			team := true
			state := SpectatorConnectedBlue
			var person *LobbyClient
			for i, guest := range this.Guests {
				if guest.Username == req.Who {
					person = guest
					remove = i
					team = guest.Team
					state = guest.State
					guest.State = LeftTheLobby
				} else {
					guest.OutputChannel <- req
				}
			}
			if remove == -1 {
				continue
			}
			this.Guests = append(this.Guests[:remove], this.Guests[remove+1:]...)
			for key, _ := range LobbyClients {
				if key == person.UserID {
					delete(LobbyClients, key)
					break
				}
			}
			if state == Ready {
				this.ReadyAmount--
			}
			if team { //true => Blue
				if state == SpectatorConnectedBlue {
					this.BlueSpecAmount--
				} else {
					this.BlueAmount--
				}
			} else { //false => Red
				if state == SpectatorConnectedRed {
					this.RedSpecAmount--
				} else {
					this.RedAmount--
				}
			}
			this.Length--
			//log.Println("[LobbyHan] ppl0", this.BlueAmount, this.RedAmount, this.BlueSpecAmount, this.RedSpecAmount)
			if this.Length < 1 {
				log.Println("[LobbyHan] Deleted lobby for", this.Leader)
				return
			}
			if this.Leader == req.Who {
				this.Leader = this.Guests[0].Username
				message := LobbyResponse{
					Event: LeaderChange,
					Who:   this.Leader,
				}
				for _, guest := range this.Guests {
					guest.OutputChannel <- message
				}
			}
		} else if req.Event == GameStart {
			BluePlayers := make([]*LobbyClient, 0)
			RedPlayers := make([]*LobbyClient, 0)
			found := 0
			for _, guest := range this.Guests {
				if guest.State == Ready {
					if guest.Team {
						BluePlayers = append(BluePlayers, guest)
					} else {
						RedPlayers = append(RedPlayers, guest)
					}
					found++
					if found == this.Capacity {
						break
					}
				}
			}
			blueChannels := make([]*ClientChannels, this.Capacity/2)
			redChannels := make([]*ClientChannels, this.Capacity/2)
			for i, bplayer := range BluePlayers {
				rplayer := RedPlayers[i]
				userBlueChannels := ClientChannels{
					UserID:         bplayer.UserID,
					Opponent:       nil,
					State:          ReadyingForTheGame,
					ChosenGirls:    []int{bplayer.Girls[0], bplayer.Girls[1]},
					SkillLevels:    []int{GetSkillLevel(bplayer.UserID, bplayer.Girls[0]), GetSkillLevel(bplayer.UserID, bplayer.Girls[1])},
					LastThing:      GameState{},
					Input:          make(chan string, 2),
					Output:         make(chan GameState, 4),
					HasGivenUp:     make(chan bool, 1),
					Time:           make(chan bool, 1),
					TimeOutput:     make(chan string, 1),
					Taken:          make(chan *ClientChannels, 1),
					Disconnected:   make(chan string, 2),
				}
				userRedChannels := ClientChannels{
					UserID:         rplayer.UserID,
					Opponent:       nil,
					State:          ReadyingForTheGame,
					ChosenGirls:    []int{rplayer.Girls[0], rplayer.Girls[1]},
					SkillLevels:    []int{GetSkillLevel(rplayer.UserID, rplayer.Girls[0]), GetSkillLevel(rplayer.UserID, rplayer.Girls[1])},
					LastThing:      GameState{},
					Input:          make(chan string, 2),
					Output:         make(chan GameState, 4),
					HasGivenUp:     make(chan bool, 1),
					Time:           make(chan bool, 1),
					TimeOutput:     make(chan string, 1),
					Taken:          make(chan *ClientChannels, 1),
					Disconnected:   make(chan string, 2),
				}
				_, girl1Blue, girlRed := userBlueChannels.GetCompatibility(&userRedChannels)
				userBlueChannels.PlayingAs = girl1Blue
				userRedChannels.PlayingAs = girlRed
				userBlueChannels.Opponent = &userRedChannels
				userRedChannels.Opponent = &userBlueChannels

				AddBattle(bplayer.UserID, &userBlueChannels)
				AddBattle(rplayer.UserID, &userRedChannels)

				blueChannels[i] = &userBlueChannels
				redChannels[i] = &userRedChannels
				log.Println("[Lobby]", bplayer.Username, "playing as", userBlueChannels.PlayingAs, rplayer.Username, "playing as", userRedChannels.PlayingAs)
			}

			for _, guest := range this.Guests {
				guest.OutputChannel <- req
			}

			for _, channel := range blueChannels {
				go WaitForConnection(channel)
			}
			for _, channel := range redChannels {
				go WaitForConnection(channel)
			}

		} else if req.Event == Ready {
			for _, guest := range this.Guests {
				guest.OutputChannel <- req
			}
			this.ReadyAmount++
			if this.ReadyAmount != this.Capacity {
				continue
			}

			this.Requests <- LobbyResponse{
				Event: GameStart,
				Who:   "/game",
			}
		} else if req.Event == Unready {
			for _, guest := range this.Guests {
				guest.OutputChannel <- req
				if guest.Username == req.Who {
					if guest.Team {
						guest.State = ConnectedBlue
					} else {
						guest.State = ConnectedRed
					}
				}
			}
			this.ReadyAmount--
		} else if req.Event == LeaderChange {
			if req.Who != this.Leader {
				found1 := false
				found2 := false
				var leader *LobbyClient
				for _, guest := range this.Guests {
					if guest.Username == req.Who {
						log.Println("[LobbyHan] Changed the leader of", this.Leader, "to", req.Who)
						this.Leader = guest.Username
						found1 = true
					} else if guest.Username == this.Leader {
						leader = guest
						found2 = true
					}
					if found1 && found2 {
						break
					}
				}
				if found1 {
					for _, guest := range this.Guests {
						guest.OutputChannel <- req
					}
				} else {
					log.Println("[LobbyHan] User not found to change the leader of", this.Leader, "to", req.Who)
					leader.OutputChannel <- LobbyResponse{
						Event: Error,
						Who:   "Invalid user leader change.",
					}
				}
			} else {
				var leader *LobbyClient
				for _, guest := range this.Guests {
					if guest.Username == this.Leader {
						leader = guest
						break
					}
				}
				leader.OutputChannel <- LobbyResponse{
					Event: Error,
					Who:   "Can't change leader to yourself.",
				}
			}
		} else if req.Event == ConnectedFirst { //someone connected
			user := FindBase(req.Who)
			if user == nil {
				log.Fatal("[LobbyHandler] why did you do this to me??? " + req.Who)
				return
			}
			var person *LobbyClient
			found := false
			for ID, value := range LobbyClients {
				if ID == user.UserID {
					person = value
					found = true
					break
				}
			}
			if !found {
				continue
			}
			this.Length++
			log.Println("[LobbyHandler] Added", user.Username, this.Length, "to", this.Leader)
			designatedState := ConnectedBlue
			if this.BlueAmount+this.RedAmount < this.Capacity {
				if this.BlueAmount > this.RedAmount {
					designatedState = ConnectedRed
					this.RedAmount++
					person.Team = false
				} else {
					this.BlueAmount++
					person.Team = true
				}
			} else {
				if this.BlueSpecAmount > this.RedSpecAmount {
					designatedState = SpectatorConnectedRed
					this.RedSpecAmount++
					person.Team = false
				} else {
					designatedState = SpectatorConnectedBlue
					this.BlueSpecAmount++
					person.Team = true
				}
			}
			person.State = designatedState
			req.Event = designatedState
			guestNames := ""
			//log.Println(len(this.Guests), this.Guests)
			for _, guest := range this.Guests {
				guest.OutputChannel <- req
				var state int
				if guest.State == Ready {
					if guest.Team {
						state = ConnectedBlue
					} else {
						state = ConnectedRed
					}
				} else {
					state = guest.State
				}
				guestNames += " " + guest.Username + " " + strconv.Itoa(state)
			}
			if this.Length > 1 {
				person.OutputChannel <- LobbyResponse{
					Event: LobbyList,
					Who:   guestNames,
				}
			}
			person.OutputChannel <- LobbyResponse{
				Event: LeaderChange,
				Who:   this.Leader,
			}
			this.Guests = append(this.Guests, person)
			person.OutputChannel <- req
		} else { //someone is trying to change their state
			info := strings.Split(req.Who, " ")
			success := false
			name := info[0]
			var person *LobbyClient
			for _, guest := range this.Guests {
				if name == guest.Username {
					person = guest
					break
				}
			}
			prevState := person.State
			prevTeam := person.Team
			if info[1] == "connect" && info[2] == "blue" && prevState == ConnectedBlue ||
				info[1] == "connect" && info[2] == "red" && prevState == ConnectedRed ||
				info[1] == "spectate" && info[2] == "blue" && prevState == SpectatorConnectedBlue ||
				info[1] == "spectate" && info[2] == "red" && prevState == SpectatorConnectedRed {
				person.OutputChannel <- LobbyResponse{
					Error,
					"You're already in this state",
				}
				continue
			}
			//log.Println("[LobbyHan] ppl2", this.BlueAmount, this.RedAmount, this.BlueSpecAmount, this.RedSpecAmount)
			if info[1] == "connect" && this.BlueAmount+this.RedAmount < this.Capacity {
				if info[2] == "blue" && this.BlueAmount < this.Capacity/2 {
					success = true
					person.State = ConnectedBlue
					this.BlueAmount++
					person.Team = true
				} else if info[2] == "red" && this.RedAmount < this.Capacity/2 {
					success = true
					person.State = ConnectedRed
					this.RedAmount++
					person.Team = false
				}
			} else if info[1] == "spectate" { //info[1] == "spectate"
				if info[2] == "blue" && this.BlueSpecAmount < this.SpecCapacity/2 {
					success = true
					person.State = SpectatorConnectedBlue
					this.BlueSpecAmount++
					person.Team = true
				} else if info[2] == "red" && this.RedAmount < this.SpecCapacity/2 {
					success = true
					person.State = SpectatorConnectedRed
					this.RedSpecAmount++
					person.Team = false
				}
			}

			req.Event = person.State
			req.Who = person.Username
			if success {
				if prevState == Ready {
					this.ReadyAmount--
				}
				if prevTeam { //was blue
					if prevState == SpectatorConnectedBlue {
						this.BlueSpecAmount--
					} else {
						this.BlueAmount--
					}
				} else { //was red
					if prevState == SpectatorConnectedRed {
						this.RedSpecAmount--
					} else {
						this.RedAmount--
					}
				}
				for _, guest := range this.Guests {
					guest.OutputChannel <- req
				}
			} else {
				person.OutputChannel <- LobbyResponse{
					Event: Error,
					Who:   "Nope.",
				}
			}
			//log.Println("[LobbyHan] ppl2", this.BlueAmount, this.RedAmount, this.BlueSpecAmount, this.RedSpecAmount)
		}
		/*for _, guest := range this.Guests {
			if guest.Username == this.Leader {
				guest.OutputChannel <- LobbyResponse{
					ExplainYourself,
					this.ExplainYourself(),
				}
				break
			}
		}*/
	}
}

func LobbyTestPage(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, _ := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Lobby Test Page]")
		Redirect(w, r, "/login")
		return
	}
	if r.Method == http.MethodGet {
		Path := "/html/lobbytest.html"
		pwd, _ := os.Getwd()
		Path = strings.Replace(pwd+Path, "/", "\\", -1)
		log.Println("[lobby] " + Path)
		http.ServeFile(w, r, Path)
	}
}

func IsInALobby(ID int64) bool {
	for key, _ := range LobbyClients {
		if key == ID {
			return true
		}
	}
	return false
}

func SendMessages(ws *websocket.Conn, client *LobbyClient) {
	for {
		select {
		case Message, stillOpen := <-client.OutputChannel:
			if stillOpen && client.State > LeftTheLobby {
				err := ws.WriteJSON(Message)
				log.Println("[Lobby] Sending a message", Message.Event, Message.Who,
					"to ws", ws.RemoteAddr().String()+",",
					"to", FindBaseID(client.UserID).Username)
				if err != nil {
					return
				}
			} else {
				return
			}
		}

	}
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Println("[Lobby Test Page] back to login")
		Redirect(w, r, "/login")
		return
	}
	log.Println("[Lobby] accessing lobby for", session.UserID)
	channels, present := GetBattle(session.UserID)
	if present {
		log.Println("[Lobby] Terminating your game, ", session.UserID)
		channels.GiveUp()
	} else {
		log.Println("[Lobby] user not found in game", session.UserID)
	}
	if IsInALobby(session.UserID) || IsInQueue(session.UserID) {
		log.Println("[Lobby Test Page] You're already in a lobby, baka", session.UserID)
		http.Error(w, "You're already in a lobby or in a queue, baka.", 400)
		return
	}

	user := FindBaseID(session.UserID)

	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	var msg string
	err = ws.ReadJSON(&msg)
	if err != nil {
		log.Println("[Lobby] that's what she said:", user.UserID, msg, err.Error())
		return
	}
	split := strings.Split(msg, " ")
	created := false
	var LobbyHan *LobbyHandler
	if len(split) > 2 {
		ws.WriteJSON(LobbyResponse{
			Event: Error,
			Who:   "Invalid command",
		})
		log.Println("[Lobby] that's what she said:", user.UserID, msg)
		return
	} else if len(split) == 1 {
		if split[0] != "create" {
			ws.WriteJSON(LobbyResponse{
				Event: Error,
				Who:   "Invalid command",
			})
			log.Println("[Lobby] that's what she said:", user.UserID, msg)
			return
		}
		created = true
		LobbyHan = &LobbyHandler{
			Guests:         make([]*LobbyClient, 0),
			Capacity:       2,
			SpecCapacity:   12,
			ReadyAmount:    0,
			Length:         0,
			BlueAmount:     0,
			RedAmount:      0,
			BlueSpecAmount: 0,
			RedSpecAmount:  0,
			Leader:         user.Username,
			Requests:       make(chan LobbyResponse, 10),
		}
		go LobbyHan.ListenAndManageTheLobby()
	} else {
		if split[0] != "connect" {
			ws.WriteJSON(LobbyResponse{
				Event: Error,
				Who:   "Invalid command",
			})
			log.Println("[Lobby] that's what she said:", user.UserID, msg)
			return
		}
		leader := split[1]
		found := false
		for _, person := range LobbyClients {
			if person.Username == leader {
				LobbyHan = person.LobbyHandler
				found = true
				break
			}
		}
		if !found {
			ws.WriteJSON(LobbyResponse{
				Event: Error,
				Who:   "This lobby doesn't exist or the invitation has expired.",
			})
			log.Println("[Lobby] couldn't find:", user.UserID, msg)
			return
		} else if LobbyHan.Length == LobbyHan.SpecCapacity+LobbyHan.Capacity {
			ws.WriteJSON(LobbyResponse{
				Event: Error,
				Who:   "This lobby is already full.",
			})
			log.Println("[Lobby] lobby full:", user.UserID, msg)
			return
		}

	}
	client := LobbyClient{
		UserID:        session.UserID,
		Username:      user.Username,
		InputChannel:  make(chan string, 3),
		OutputChannel: make(chan LobbyResponse, 10),
		Girls:         make([]int, 2),
		State:         ConnectedFirst,
		Team:          created,
		IsLeader:      created,
		LobbyHandler:  LobbyHan,
	}
	LobbyClients[client.UserID] = &client

	//Send messages to the client
	go SendMessages(ws, &client)

	//try to connect
	client.LobbyHandler.Requests <- LobbyResponse{
		Event: ConnectedFirst,
		Who:   client.Username,
	}
	SetState(client.UserID, InLobby)
	//awaits for player messages and throws them into the input channel
	for {
		var msg2 string
		err := ws.ReadJSON(&msg2)
		if err != nil {
			log.Println("[LobbyHan] bye bye or error", err)
			client.LobbyHandler.Requests <- LobbyResponse{
				Event: LeftTheLobby,
				Who:   client.Username,
			}
			return
		} else {
			log.Println("[Lobby]", client.Username, "sent", msg2)
			msg := strings.Split(msg2, " ")
			if len(msg) == 1 && msg[0] == "ready" { //ready
				if client.State != ConnectedRed && client.State != ConnectedBlue {
					if client.State == SpectatorConnectedRed || client.State == SpectatorConnectedBlue {
						client.OutputChannel <- LobbyResponse{
							Event: Error,
							Who:   "You're not to be ready, you're a spectator.",
						}
					} else if client.State == Ready {
						client.OutputChannel <- LobbyResponse{
							Event: Error,
							Who:   "You're already ready.",
						}

					} else if client.Girls[0] == 0 {
						client.OutputChannel <- LobbyResponse{
							Event: Error,
							Who:   "Choose your characters first.",
						}
					} else {
						client.State = Ready
						client.LobbyHandler.Requests <- LobbyResponse{
							Event: Ready,
							Who:   client.Username,
						}
					}
				} else {
					client.State = Ready
					client.LobbyHandler.Requests <- LobbyResponse{
						Event: Ready,
						Who:   client.Username,
					}
				}
			} else if len(msg) == 2 && msg[0] == "not" && msg[1] == "ready" { //unready
				if client.State != Ready {
					client.OutputChannel <- LobbyResponse{
						Event: Error,
						Who:   "You weren't ready in the first place.",
					}
				} else {
					client.State = Unready
					client.LobbyHandler.Requests <- LobbyResponse{
						Event: Unready,
						Who:   client.Username,
					}
				}
			} else if len(msg) == 3 && msg[0] == "change" && msg[1] == "leader" {
				if client.LobbyHandler.Leader != client.Username {
					client.OutputChannel <- LobbyResponse{
						Event: Error,
						Who:   "You aren't the lobby leader.",
					}
				} else {
					who := msg[2]
					log.Println("[Lobby] Tried to change the leader of", client.LobbyHandler.Leader, "to", who)
					client.LobbyHandler.Requests <- LobbyResponse{
						Event: LeaderChange,
						Who:   who,
					}
				}
			} else if len(msg) == 4 && msg[0] == "send" && msg[1] == "girls" && len(msg[2]) < 4 && len(msg[3]) < 4 {
				if client.State == ConnectedRed || client.State == ConnectedBlue {
					girl1, err1 := strconv.Atoi(msg[2])
					girl2, err2 := strconv.Atoi(msg[3])
					if err1 != nil || err2 != nil || girl1 == girl2 {
						client.OutputChannel <- LobbyResponse{
							Event: Error,
							Who:   "These aren't proper girl numbers.",
						}
					} else if !HasGirl(client.UserID, girl1) || !HasGirl(client.UserID, girl2) {
						client.OutputChannel <- LobbyResponse{
							Event: Error,
							Who:   "You don't have this/these girls unlocked.",
						}
					} else {
						client.Girls[0] = girl1
						client.Girls[1] = girl2
						client.OutputChannel <- LobbyResponse{
							Event: SentGirls,
							Who:   "Your girls have been accepted.",
						}
					}
				} else if client.State == SpectatorConnectedRed || client.State == SpectatorConnectedBlue {
					client.OutputChannel <- LobbyResponse{
						Event: Error,
						Who:   "You're not to send girls, you're a spectator.",
					}
				} else if client.State == Ready {
					client.OutputChannel <- LobbyResponse{
						Event: Error,
						Who:   "You can't change your girls once you're ready.",
					}
				}
			} else if len(msg) == 3 && msg[0] == "change" && (msg[1] == "spectate" || msg[1] == "connect") && (msg[2] == "red" || msg[2] == "blue") {
				client.LobbyHandler.Requests <- LobbyResponse{
					Event: ChangeState,
					Who:   client.Username + " " + msg[1] + " " + msg[2],
				}
			} else { //TODO kick function
				//TODO invite function
				client.OutputChannel <- LobbyResponse{
					Event: Error,
					Who:   "Invalid command.",
				}
			}
		}
	}
}
