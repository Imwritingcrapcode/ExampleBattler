package Site

import (
	. "../Abstract"
	. "../Characters"
	"database/sql"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//BattlerHandlers

var STLevels = []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 5, 5, 5, 5}
var ADLevels = []int{1, 1, 2, 2, 2, 2, 3, 3, 4, 4, 4, 4, 4, 5, 5, 6, 6, 6, 6}
var SPLevels = []int{1, 2, 4, 5, 5, 5, 6, 8, 9, 11, 11, 11, 11, 11, 13, 14, 16, 16, 16}
var RPLevels = []int{2, 6, 10, 13, 13, 14, 14, 16, 20, 23, 24, 24, 25, 27, 29, 31, 35, 37, 37}
var LFLevels = []int{4, 7, 12, 16, 17, 18, 20, 23, 27, 30, 31, 31, 31, 33, 34, 38, 41, 43, 44}

func GetLevelCaps(rarity string) []int {
	switch rarity {
	case "ST":
		return STLevels
	case "AD":
		return ADLevels
	case "SP":
		return SPLevels
	case "RP":
		return RPLevels
	case "LF":
		return LFLevels
	default:
		log.Println("unknown rarity: " + rarity)
		return nil
	}
}

var ClientConnections = map[int64]*ClientChannels{}

var DifferenceForWinLoss = map[string]int{
	Rarities[0]: 1,
	Rarities[1]: 2,
	Rarities[2]: 3,
	Rarities[3]: 4,
	Rarities[4]: 5,
}

func HowMuchDoesSheGive(girlNumber int, win bool) int {
	var num int
	rarity := GetGirlRarity(girlNumber)
	if win {
		num = 10 + DifferenceForWinLoss[rarity]
	} else {
		num = 10 - DifferenceForWinLoss[rarity]
	}

	return num
}

func EndGame(channels *ClientChannels) {
	log.Println("[GAME] Ended game as", channels.UserID, "against", channels.Opponent.UserID)
	_, stillOpen := <-channels.Input
	if stillOpen {
		close(channels.Input)
	}

	_, stillOpen = <-channels.KillConnection
	if stillOpen {
		close(channels.KillConnection)
	}

	close(channels.Time)

	if channels.Clock.State() {
		channels.Clock.Stop()
	}

	_, present := ClientConnections[channels.UserID]
	if present {
		delete(ClientConnections, channels.UserID)
		log.Println("[GAME] removed", channels.UserID, "from the game map")
	}
	if channels.Opponent != channels {
		_, stillOpen := <-channels.Opponent.Input
		if stillOpen {
			close(channels.Opponent.Input)
		}
		_, stillOpen = <-channels.Opponent.KillConnection
		if stillOpen {
			close(channels.Opponent.KillConnection)
		}

		close(channels.Opponent.Time)

		if channels.Opponent.Clock.State() {
			channels.Opponent.Clock.Stop()
		}
		_, present := ClientConnections[channels.Opponent.UserID]
		if present {
			delete(ClientConnections, channels.Opponent.UserID)
		}
		log.Println("[GAME] removed", channels.Opponent.UserID, "from the game map")
	}
}

//DB

var BotNames = []string{
	"BestHomie",
	"frozencalla",
	"IrJean",
	"Ron",
	"ukulele",
	"urkitten",
	"YouHaveTheHugeMegaGay",
}

var DustMap = map[string]string{
	"w": "wdust",
	"b": "bdust",
	"y": "ydust",
	"g": "gdust",
	"s": "sdust",
}


var DATABASE *sql.DB

type RewardsObj struct {
	BattleResult      int
	LastOpponentsName string
	Dusts             map[string]int
	Matches           []int
	Name              string
	TotalMatches      int
	ToAdd             int
}

type User struct {
	UserID          int64
	Username        string
	CurrentActivity int
	PlayingAs       int
}

type Session struct {
	ID      string
	UserID  int64
	Expires int64
}

type UserGirl struct {
	Number        int
	Level         int
	MatchesPlayed int
	MatchesWon    int
	//all the fields below are taken from GirlInfo
	Name         string
	Rarity       string
	Tags         []string
	Skills       []string
	SkillColours []string
	SkillColourCodes []string
	Description  string
	MainColour   string
}

type FriendList struct {
	Friends  [][]string `json:"friends"`
	Incoming []string   `json:"incoming"`
	Pending  []string   `json:"pending"`
}

//GameServer

const WAITTIME = 15 * time.Nanosecond

var QueueClients = map[int64]*ClientChannels{}

var upgrader = websocket.Upgrader{}

var UserQueue = make(PriorityQueue, 0)

type QueueResponse struct {
	OK       bool   `json:"OK"`
	Prompt   string `json:"Prompt"`
	Location string `json:"Location"`
}

type ClientMessage struct {
	MainGirl      int `json:"MainGirl"`
	SecondaryGirl int `json:"SecondaryGirl"`
}

//GirlList

type GirlListResponse struct {
	Girls    []UserGirl `json:"girls"`
	Response string     `json:"response"`
}

//LoginRegister

const SessionsClearPeriod = time.Hour * 24

var LastSessionsCleared time.Time

type RegData struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberme"`
}

type UserFree struct {
	Username     string `json:"Username"`
	BattlesTotal int    `json:"BattlesTotal"`
	BattlesWon   int    `json:"BattlesWon"`
	WhiteDust    int    `json:"WhiteDust"`
	BlueDust     int    `json:"BlueDust"`
	YellowDust   int    `json:"YellowDust"`
	GreenDust    int    `json:"GreenDust"`
	StarDust     int    `json:"StarDust"`
}

//servmain

const OfflineEvery = 15 * time.Minute

//shop

type ShopItem struct {
	Name string
	ID   string
	Type string
	Dust string
	Cost int
}

func GetItemByID(ID string) (bool, ShopItem) {
	item := ShopItem{}
	exists := false
	switch ID {
	case Rarities[0]:
		exists = len(ReleasedCharactersPacks[Rarities[0]]) != 0
		item.Name = Rarities[0] + " pack"
		item.Cost = 300
		item.Dust = "w"
		item.Type = "pack"
		item.ID = ID
	case Rarities[1]:
		exists = len(ReleasedCharactersPacks[Rarities[1]]) != 0
		item.Name = Rarities[1] + " pack"
		item.Cost = 185
		item.Dust = "b"
		item.Type = "pack"
		item.ID = ID
	case Rarities[2]:
		exists = len(ReleasedCharactersPacks[Rarities[2]]) != 0
		item.Name = Rarities[2] + " pack"
		item.Cost = 124
		item.Dust = "y"
		item.Type = "pack"
		item.ID = ID
	case Rarities[3]:
		exists = len(ReleasedCharactersPacks[Rarities[3]]) != 0
		item.Name = Rarities[3] + " pack"
		item.Cost = 72
		item.Dust = "g"
		item.Type = "pack"
		item.ID = ID
	case Rarities[4]:
		exists = len(ReleasedCharactersPacks[Rarities[4]]) != 0
		item.Name = Rarities[4] + " pack"
		item.Cost = 15
		item.Dust = "s"
		item.Type = "pack"
		item.ID = ID
	}
	return exists, item
}

//conversion

type ConvRequest struct {
	ReqType string
	Amount int
	DustType string
}

type ConvResponse struct {
	ConversionRate map [string]float64
	SecondsPerConversion map [string]int
}

var ConversionRate = map[string]float64{
	"w": 0.5,
	"b": 0.4,
	"y": 0.25,
	"g": 0.2,
}

var SecondsPerConversion = map [string]int{
	"w": 24,
	"b": 30,
	"y": 45,
	"g": 60,
}