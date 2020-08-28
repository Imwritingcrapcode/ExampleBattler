package Site

import (
	. "../Abstract"
	. "../Game"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
flavor: other important todos
• Mini skill buttons on girllist + stroke for the edge + other colours
• effect icons for Z89 & Euphoria, effect for Z89
• Please choose the girls first appears in red
• add friend sends a notification to the potential friend
• remove add friend button if already friends
//FLAVOR
//TODO cuter design of everything
//todo linux support, get a server
//TODO DDOS protection (not cloudflare apparently)
//TODO bot
//TODO random damage
//todo loading animation while queuing
//TODO global chat
//TODO dms chat
//TODO emojis (almost done but pasting in chat)
//TODO speech bubbles
//TODO skins
//TODO site themes
//TODO character wiki
//TODO news page
//TODO choose your pfp
//TODO battle your friends
//TODO ability draft, 2v2, ...
*/

func DistributeRewards(p1 *ClientChannels, won bool) {
	log.Println("[INGAME] gg distribute rewards", p1.UserID)
	if p1.State != GaveUp {
		log.Println("[ingame] state of ", p1.UserID, ActivitiesToString[p1.State])
		//1. Determine how much that girl should give
		//2. Add that much
		neededAmount := HowMuchDoesSheGive(p1.PlayingAs, won)
		//2.
		user := FindBaseID(p1.UserID)
		AddRewards(user, "w", neededAmount)
		AddRewards(user, strconv.Itoa(p1.PlayingAs), 1)
		user.SetDust("w", user.GetDust("w")+neededAmount)
	}
}

func Game(p1, p2 *ClientChannels) {
	var g1, g2 CharInt
	p1.Clock = &Clock{
		Client: p1,
	}
	p2.Clock = &Clock{
		Client: p2,
	}
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	PLAYER1 := p1.PlayingAs
	PLAYER2 := p2.PlayingAs

	coin := rand.Intn(2)
	if coin == 0 {
		InitAsNumber(gi1, PLAYER1)
		InitAsNumber(gi2, PLAYER2)
	} else {
		InitAsNumber(gi1, PLAYER2)
		InitAsNumber(gi2, PLAYER1)
	}
	SomeoneGaveUp := false
	var p1Won, p2Won bool
	var state1, state2 *GameState
	for i := 1; i < 21; i++ {
		if (i % 2) == 1 { //turn of girl1
			if PLAYER1 == gi1.Number { //this is p1's turn
				log.Println("[INGAME] turn of", gi1.Name, p1.UserID)
				SomeoneGaveUp = Turn2Channels(gi1, gi2, i, p1, p2)
			} else { //gi1==player2
				log.Println("[INGAME] turn of", gi1.Name, p2.UserID)
				SomeoneGaveUp = Turn2Channels(gi1, gi2, i, p2, p1)
			}
			if SomeoneGaveUp || p1.State == GaveUp || p2.State == GaveUp || !gi1.IsAlive() || !gi2.IsAlive() || i == 20 {
				if PLAYER1 == gi1.Number { //p1's turn
					state1 = GetGameStateChannels(gi1, gi2, i)
					state2 = GetGameStateChannelsOpp(gi1, gi2, i)
				} else { //gi1 = player2, gi1 & p2's turn.
					state1 = GetGameStateChannelsOpp(gi1, gi2, i)
					state2 = GetGameStateChannels(gi1, gi2, i)
				}
				break
			}
		} else { //turn of girl2
			if PLAYER1 == gi2.Number { //this is p1's turn
				log.Println("[INGAME] turn of", gi2.Name, p1.UserID)
				SomeoneGaveUp = Turn2Channels(gi2, gi1, i, p1, p2)
			} else { //gi2==player2
				log.Println("[INGAME] turn of", gi2.Name, p2.UserID)
				SomeoneGaveUp = Turn2Channels(gi2, gi1, i, p2, p1)
			}
			if SomeoneGaveUp || p1.State == GaveUp || p2.State == GaveUp || !gi1.IsAlive() || !gi2.IsAlive() || i == 20 {
				if PLAYER1 == gi2.Number { //p1's turn
					state1 = GetGameStateChannels(gi2, gi1, i)
					state2 = GetGameStateChannelsOpp(gi2, gi1, i)
				} else { //gi2 = player2, gi2 & p2's turn.
					state1 = GetGameStateChannelsOpp(gi2, gi1, i)
					state2 = GetGameStateChannels(gi2, gi1, i)
				}
				break
			}
		}
	}
	if SomeoneGaveUp || p1.State == GaveUp || p2.State == GaveUp {
		if p1.State == GaveUp {
			state1.EndState = GameGaveUp
			state2.EndState = GameOppGaveUp
			p1Won = false
			p2Won = true
		} else {
			state2.EndState = GameGaveUp
			state1.EndState = GameOppGaveUp
			p1Won = true
			p2Won = false
		}
	} else {
		if gi1.CurrHP > gi2.CurrHP && gi1.Number == PLAYER1 ||
			gi2.CurrHP > gi1.CurrHP && gi2.Number == PLAYER1 { //player1 won
			state1.EndState = GameWon
			state2.EndState = GameLost
			p1Won = true
			p2Won = false

		} else if gi1.CurrHP == gi2.CurrHP { //draw
			state1.EndState = GameDraw
			state2.EndState = GameDraw
			p1Won = true
			p2Won = true
		} else { //player2 won
			state2.EndState = GameWon
			state1.EndState = GameLost
			p1Won = false
			p2Won = true
		}
	}
	SetLastBattleResult(p1.UserID, state1.EndState)
	SetLastBattleResult(p2.UserID, state2.EndState)
	DistributeRewards(p1, p1Won)
	DistributeRewards(p2, p2Won)
	p1.Send(state1)
	p2.Send(state2)
	SetState(p1.UserID, JustFinishedTheGame)
	SetState(p2.UserID, JustFinishedTheGame)
	IncreaseBattles(p1.UserID, p1Won)
	IncreaseBattles(p2.UserID, p2Won)
	IncreaseMatchesAs(p1.UserID, p1.PlayingAs, p1Won, state1.EndState == GameGaveUp)
	IncreaseMatchesAs(p2.UserID, p2.PlayingAs, p2Won, state2.EndState == GameGaveUp)
	EndGame(p1)
}

func BotGame(p1 *ClientChannels, botChar int, DEPTH int) {
	if (p1.PlayingAs == 33 || botChar == 33) && DEPTH > 5 {
		DEPTH = 5
	} else if (p1.PlayingAs == 33 || botChar == 33) && DEPTH < 4 {
		DEPTH = 4
	}

	var g1, g2 CharInt
	p1.Clock = &Clock{
		Client: p1,
	}
	g1 = new(Girl)
	g2 = new(Girl)
	gi1 := g1.(*Girl)
	gi2 := g2.(*Girl)
	PLAYER1 := p1.PlayingAs
	PLAYER2 := botChar
	coin := rand.Intn(2)
	InitAsNumber(gi1, PLAYER1)
	InitAsNumber(gi2, PLAYER2)
	IGaveUp := false
	var p1Won bool
	var state1 *GameState
	botInput := make(chan string, 2)
	botClock := Clock{Client: p1}
	for i := 1; i < 21; i++ {
		if (i % 2) == coin {
			if gi1.HasEffect(ControlledByStT) {
				log.Println("[INGAME] turn of the bot (took control)", gi2.Name)
				go ThrowNextMoves(gi1.Copy(), gi2.Copy(), i, DEPTH, botInput)
			} else {
				log.Println("[INGAME] turn of", gi1.Name, p1.UserID)
			}
			IGaveUp = TurnChannels(gi1, gi2, i, p1, botInput, &botClock, false)

			if IGaveUp || p1.State == GaveUp || !gi1.IsAlive() || !gi2.IsAlive() || i == 20 {
				state1 = GetGameStateChannels(gi1, gi2, i)
				break
			}
		} else {
			if !gi2.HasEffect(ControlledByStT) {
				log.Println("[INGAME] turn of the bot", gi2.Name)
				go ThrowNextMoves(gi2.Copy(), gi1.Copy(), i, DEPTH, botInput)
			} else {
				log.Println("[INGAME] turn of (took control)", gi1.Name, p1.UserID)
			}
			IGaveUp = TurnChannels(gi2, gi1, i, p1, botInput, &botClock, true)

			if IGaveUp || p1.State == GaveUp || !gi1.IsAlive() || !gi2.IsAlive() || i == 20 {
				state1 = GetGameStateChannelsOpp(gi2, gi1, i)
				break
			}
		}
	}
	if IGaveUp || p1.State == GaveUp {
		state1.EndState = GameGaveUp
		p1Won = false
	} else {
		if gi1.CurrHP > gi2.CurrHP && gi1.Number == PLAYER1 ||
			gi2.CurrHP > gi1.CurrHP && gi2.Number == PLAYER1 { //player1 won
			state1.EndState = GameWon
			p1Won = true

		} else if gi1.CurrHP == gi2.CurrHP { //draw
			state1.EndState = GameDraw
			p1Won = true
		} else { //the bot won
			state1.EndState = GameLost
			p1Won = false
		}
	}
	SetLastBattleResult(p1.UserID, state1.EndState)
	DistributeRewards(p1, p1Won)
	p1.Send(state1)
	SetState(p1.UserID, JustFinishedTheGame)
	IncreaseBattles(p1.UserID, p1Won)
	IncreaseMatchesAs(p1.UserID, p1.PlayingAs, p1Won, state1.EndState == GameGaveUp)
	EndGame(p1)
}

func Standard(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Print("[INGAME] redirected to /login")
		Redirect(w, r, "/login")
		return
	}
	client := FindBaseID(session.UserID)
	channels, present := ClientConnections[client.UserID]
	if !present || channels.State <= Queuing {
		log.Print("[INGAME] " + client.Username + " redirected to /girllist bc their state is < queuing or they are not in the map")
		Redirect(w, r, "/girllist")
		return
	}
	state := channels.State
	//a valid registered user that is in game,
	//in queue and came after finding an opponent
	if state >= GaveUp {
		log.Print("[INGAME] " + client.Username + " redirected to /afterbattle since the game is over")
		Redirect(w, r, "/afterbattle")
		return
	}

	Path := "/html/interface3.html"
	pwd, _ := os.Getwd()
	Path = strings.Replace(pwd+Path, "/", "\\", -1)
	log.Println("[INGAME] " + Path)
	http.ServeFile(w, r, Path)
}

func BattlerHandler(w http.ResponseWriter, r *http.Request) {
	AlrdyLoggedIn, session := IsLoggedIn(r)
	if !AlrdyLoggedIn {
		log.Print("[INGAME] redirected to /login")
		Redirect(w, r, "/login")
		return
	}
	client := FindBaseID(session.UserID)
	channels, present := ClientConnections[client.UserID]
	if !present || channels.State <= Queuing {
		log.Print("[INGAME] "+client.Username+" redirected to /girllist", !present)
		Redirect(w, r, "/girllist")
		return
	}
	state := channels.State
	if state >= GaveUp {
		log.Print("[INGAME] " + client.Username + " redirected to /afterbattle since the game is over")
		Redirect(w, r, "/afterbattle")
		return
	}
	if state == Disconnected || state == ReadyingForTheGame {
		//upgrade their request.
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		if state == Disconnected { //reconnect
			log.Println("[INGAME] Reconnected! will show the last thing that happened for", client.Username)
			kill := make(chan struct{}, 1)
			channels.KillConnection = kill
			channels.Output <- channels.LastThing
			channels.Time <- true
			channels.State = PlayingAs
			SetState(client.UserID, PlayingAs)
			if channels.Opponent != channels {
				LS := channels.Opponent.LastThing
				LS.Instruction = "System:Opponent reconnected."
				channels.Opponent.Output <- LS
			}
		}

		defer ws.Close()
		defer close(channels.KillConnection)

		if state == ReadyingForTheGame && channels.Opponent != channels && channels.Opponent.State == PlayingAs {
			//we are the second one to connect
			log.Println("[INGAME] connected 2nd to the game as", client.Username)
			channels.Time <- true

		} else if state == ReadyingForTheGame && channels.Opponent != channels { //starting a game with a human
			channels.State = PlayingAs
			log.Println("[INGAME] started the game as", client.Username)
			channels.Time <- true
			go WaitForTheOther(channels)

		} else if state == ReadyingForTheGame { //starting a bot game
			log.Println("[INGAME] starting the game as", client.Username)
			channels.Time <- true
			p1 := channels
			p1.State = PlayingAs
			var char2 int
			botChosen := make([]int, 2)
			botChosen[0] = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
			botChosen[1] = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
			for botChosen[1] == botChosen[0] {
				botChosen[1] = ReleasedCharacters[rand.Intn(len(ReleasedCharacters))]
			}
			ch1 := p1.ChosenGirls
			ch2 := botChosen
			if ch1[0] != ch2[0] {
				p1.PlayingAs = ch1[0]
				char2 = ch2[0]
			} else if ch1[1] != ch2[1] {
				p1.PlayingAs = ch1[1]
				char2 = ch2[1]
			} else {
				coin := rand.Intn(2)
				if coin == 1 {
					p1.PlayingAs = ch1[1]
					char2 = ch2[0]
				} else {
					p1.PlayingAs = ch1[0]
					char2 = ch2[1]
				}
			}
			DEPTH := GetSkillLevel(p1.UserID, p1.PlayingAs) //3 to 6.
			log.Println("[BOTGAME] bot has difficulty: ", DEPTH)
			SetLastPlayedAs(p1.UserID, p1.PlayingAs)
			SetLastOpponentsName(p1.UserID, p1.UserID)
			SetState(p1.UserID, PlayingAs)
			go BotGame(p1, char2, DEPTH)
		}

		//Send messages to the client
		go WriteGameStates(ws, channels)

		//awaits for player messages and throws them into the input channel
		for {

			if channels.State >= JustFinishedTheGame {
				log.Println("[INGAME] returning")
				return
			}
			var msg string
			err := ws.ReadJSON(&msg)
			if err != nil {
				if channels.Opponent != channels && channels.Opponent.State == Disconnected &&
					channels.State < GaveUp {
					channels.State = JustFinishedTheGame
					channels.Opponent.State = JustFinishedTheGame
					SetState(channels.UserID, JustFinishedTheGame)
					SetState(channels.Opponent.UserID, JustFinishedTheGame)
					EndGame(channels)
				} else if channels.State < GaveUp {
					channels.State = Disconnected
					if channels.Clock != nil {
						go channels.Clock.Disconnected()
						if channels.Opponent != channels {
							LS := channels.Opponent.LastThing
							LS.Instruction = "System:Opponent disconnected."
							channels.Opponent.Output <- LS
						}
					}
					SetState(channels.UserID, BrowsingCharacters)
				}
				return
			} else {
				log.Println("[INGAME]", client.Username, "sent", msg)
				if msg == "GiveUp" {
					channels.GiveUp()
				} else if channels.State < GaveUp {
					channels.Input <- msg
				} else {
					_, stillOpen := <-channels.Input
					if stillOpen {
						close(channels.Input)
					}
				}
			}
		}
	}
}

func WriteGameStates(ws *websocket.Conn, channels *ClientChannels) {
	for {
		//they are connected, the game is still going
		select {
		case <-channels.KillConnection:
			log.Println("[INGAME] killed the connection", ws.RemoteAddr().String())
			return
		case Message, stillOpen := <-channels.Output:
			if channels.State >= ReadyingForTheGame && channels.State <= OpponentGaveUp && stillOpen {
				err := ws.WriteJSON(Message)
				//log.Println("[INGAME] Sending a game state: turn", strconv.Itoa(Message.TurnNum)+",",
				//	Message.EndState,
				//	"to ws", ws.RemoteAddr().String()+",",
				//	"to", FindBaseID(channels.UserID).Username)
				if err != nil {
					return
				}
				if Message.EndState > 0 {
					channels.State = JustFinishedTheGame
					close(channels.Output)
					return
				}
			} else {
				return
			}
		case Message, stillOpen := <-channels.TimeOutput:
			if channels.State >= ReadyingForTheGame && channels.State <= OpponentGaveUp && stillOpen {
				err := ws.WriteJSON(Message)
				log.Println("[INGAME] Sending a time state,",
					"to ws", ws.RemoteAddr().String(),
					Message, "to", FindBaseID(channels.UserID).Username)
				if err != nil {
					return
				}
			} else {
				return
			}

		}
	}
}

//Throws bot moves into the given bot input channel.
func ThrowNextMoves(player, opp CharInt, i, DEPTH int, botInput chan string) {
	_, _, moves := MiniMax(player, opp, i, DEPTH, true, []int{})
	//dur := 7
	//time.Sleep(time.Duration(dur) * time.Second)
	if (player.(*Girl)).HasEffect(SpedUp) {
		botInput <- ToStringStrat(moves[0:1])
		botInput <- ToStringStrat(moves[1:2])
	} else {
		botInput <- ToStringStrat(moves[0:1])
	}
	//log.Println("figured out a strat for the turn", i, prediction, gameends, ToStringStrat(moves))
}

func WaitForTheOther(channels *ClientChannels) {
	p1 := channels
	p2 := channels.Opponent
	timer := time.NewTimer(CONNECTWAITTIME * time.Second)
	// check if connected
	select {
	case message := <-channels.Opponent.Time: //connected
		if message {
			SetLastPlayedAs(p1.UserID, p1.PlayingAs)
			SetLastOpponentsName(p1.UserID, p2.UserID)
			SetLastOpponentsName(p2.UserID, p1.UserID)
			SetLastPlayedAs(p2.UserID, p2.PlayingAs)
			SetState(p1.UserID, PlayingAs)
			SetState(p2.UserID, PlayingAs)
			Game(channels, channels.Opponent)
			return
		}
	case <-timer.C: //timed out
		log.Println("[INGAME]", channels.UserID, channels.Opponent.UserID, "failed to connect.")
		SetState(channels.UserID, BrowsingCharacters)
		channels.Send(&GameState{
			Instruction: "Error:Opponent failed to connect",
			EndState:    GameCancelled,
		})
		return
	}
}
