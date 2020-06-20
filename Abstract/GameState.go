package Abstract

import (
	"log"
)

const (
	Offline            = iota
	MainPage
	BrowsingFriendList
	ConversionPage
	Shopping
	//the new ones go above this line -----
	BrowsingCharacters
	Queuing
	ReadyingForTheGame
	PlayingAs
	Disconnected
	GaveUp
	OpponentGaveUp
	JustFinishedTheGame
)

const TURNLENGTH = 70
const ACTUALTURNLENGTH = 60
const TICKEVERY = 10
const MAXDCTIME = 75
const CONNECTWAITTIME = 10

//SkillStates that are sent to front end
const (
	Available = 0
	//On CD = amount of turns left (> 0)
	NotUnlockedYet   = -1
	DisabledByEffect = -2
	Disabled         = -100
)

//EndGame states
const (
	GameContinue  = iota
	GameWon
	GameLost
	GameDraw
	GameOppGaveUp
	GameGaveUp
	GameCancelled
)

var ActivitiesToString = map[int]string{
	Offline:             "Offline",
	MainPage:            "Main Page",
	BrowsingFriendList:  "Browsing their friend list",
	ConversionPage:      "Converting currency",
	Shopping:            "Shopping",
	BrowsingCharacters:  "Browsing character gallery",
	Queuing:             "Queuing for a game",
	ReadyingForTheGame:  "Just found an opponent",
	PlayingAs:           "Playing as",
	JustFinishedTheGame: "Just finished a game",
	OpponentGaveUp:      "Just finished a game (opp gave up)",
	Disconnected:        "Disconnected from a game",
	GaveUp:              "Gave up in a game",
}

type GameState struct {
	Instruction     string            `json:"Instruction"`
	TurnNum         int               `json:"TurnNum"`
	TurnPlayer      int               `json:"TurnPlayer"`
	PlayerNum       int               `json:"PlayerNum"`
	OppNum          int               `json:"OppNum"`
	PlayerName      string            `json:"PlayerName"`
	OppName         string            `json:"OppName"`
	HP              int               `json:"HP"`
	MaxHP           int               `json:"MaxHP"`
	OppHP           int               `json:"OppHP"`
	OppMaxHP        int               `json:"OppMaxHP"`
	Effects         map[int]string    `json:"Effects"`
	OppEffects      map[int]string    `json:"OppEffects"`
	SkillState      map[string]int    `json:"SkillState"`
	OppSkillState   map[string]int    `json:"OppSkillState"`
	SkillNames      map[string]string `json:"SkillNames"`
	OppSkillNames   map[string]string `json:"OppSkillNames"`
	Defenses        map[Colour]int    `json:"Defenses"`
	OppDefenses     map[Colour]int    `json:"OppDefenses"`
	SkillColours    map[string]string `json:"SkillColours"`
	OppSkillColours map[string]string `json:"OppSkillColours"`
	EndState        int               `json:"EndState"`
}

func (this *GameState) Copy() *GameState {
	var other GameState
	//other.Instruction = this.Instruction
	other.TurnNum = this.TurnNum
	other.TurnPlayer = this.TurnPlayer
	other.PlayerNum = this.PlayerNum
	other.OppNum = this.OppNum
	other.PlayerName = this.PlayerName
	other.OppName = this.OppName
	other.HP = this.HP
	other.MaxHP = this.MaxHP
	other.OppHP = this.OppHP
	other.OppMaxHP = this.OppMaxHP
	other.EndState = this.EndState
	targetMap := make(map[int]string)
	target2 := make(map[int]string)
	target3 := make(map[string]int)
	target4 := make(map[string]int)
	target5 := make(map[string]string)
	target6 := make(map[string]string)

	target7 := make(map[Colour]int)
	target8 := make(map[Colour]int)
	target9 := make(map[string]string)
	target10 := make(map[string]string)

	for key, value := range this.Effects {
		targetMap[key] = value
	}
	other.Effects = targetMap

	for key, value := range this.OppEffects {
		target2[key] = value
	}
	other.OppEffects = target2

	for key, value := range this.SkillState {
		target3[key] = value
	}
	other.SkillState = target3

	for key, value := range this.OppSkillState {
		target4[key] = value
	}
	other.OppSkillState = target4

	for key, value := range this.SkillNames {
		target5[key] = value
	}
	other.SkillNames = target5
	for key, value := range this.OppSkillNames {
		target6[key] = value
	}
	other.OppSkillNames = target6

	for key, value := range this.Defenses {
		target7[key] = value
	}
	other.Defenses = target7
	for key, value := range this.OppDefenses {
		target8[key] = value
	}
	other.OppDefenses = target8

	for key, value := range this.SkillColours {
		target9[key] = value
	}
	other.SkillColours = target9
	for key, value := range this.OppSkillColours {
		target10[key] = value
	}
	other.OppSkillColours = target10

	return &other
}

//Client for the gameServer
type ClientChannels struct {
	UserID         int64
	Opponent       *ClientChannels
	State          int
	ChosenGirls    []int
	SkillLevels		[]int
	PlayingAs      int
	LastThing      GameState
	Clock          *Clock
	Input          chan string
	Output         chan GameState
	HasGivenUp     chan bool
	Time           chan bool
	TimeOutput     chan string
	KillConnection chan struct{}
	Taken               chan *ClientChannels
	Disconnected        chan string
	IsTaken             bool
	IsDesperate         bool
	ShouldRemove        bool
}

func (user *ClientChannels) GetCompatibility(other *ClientChannels) (int, int, int) {
	//Let's say I am always more important than the other.
	myMain := user.ChosenGirls[0]
	myMainSkill := user.SkillLevels[0]
	mySec := user.ChosenGirls[1]
	mySecSkill := user.SkillLevels[1]
	theirMain := other.ChosenGirls[0]
	theirMainSkill := other.SkillLevels[0]
	theirSec := other.ChosenGirls[1]
	theirSecSkill := other.SkillLevels[1]

	switch {
	case myMain != theirMain && myMainSkill == theirMainSkill:
		return 1, myMain, theirMain
	case myMain != theirSec && myMainSkill == theirSecSkill:
		return 2, myMain, theirSec
	case mySec != theirMain && mySecSkill == theirMainSkill:
		return 3, mySec, theirMain
	case mySec != theirSec && mySecSkill == theirSecSkill:
		return 4, mySec, theirSec
	case myMain != theirMain && (-1 == myMainSkill-theirMainSkill || myMainSkill-theirMainSkill == 1):
		return 5, myMain, theirMain
	case myMain != theirSec && (-1 == myMainSkill-theirSecSkill || myMainSkill-theirSecSkill == 1):
		return 6, myMain, theirSec
	case mySec != theirMain && (-1 == myMainSkill-theirSecSkill || myMainSkill-theirSecSkill == 1):
		return 7, mySec, theirMain
	case mySec != theirSec && (-1 == mySecSkill-theirSecSkill || mySecSkill-theirSecSkill == 1):
		return 8, mySec, theirSec
	case myMain != theirMain && (-2 == myMainSkill-theirMainSkill || myMainSkill-theirMainSkill == 2):
		return 9, myMain, theirMain
	case myMain != theirSec && (-2 == myMainSkill-theirSecSkill || myMainSkill-theirSecSkill == 2):
		return 10, myMain, theirSec
	case mySec != theirMain && (-2 == myMainSkill-theirSecSkill || myMainSkill-theirSecSkill == 2):
		return 11, mySec, theirMain
	case mySec != theirSec && (-2 == mySecSkill-theirSecSkill || mySecSkill-theirSecSkill == 2):
		return 12, mySec, theirSec
	case myMain != theirMain:
		return 13, myMain, theirMain
	case myMain != theirSec:
		return 14, myMain, theirSec
	case mySec != theirMain:
		return 15, mySec, theirMain
	case mySec != theirSec:
		return 16, mySec, theirSec
	default:
		panic("oh lol")
		return 1000, myMain, theirMain
	}
}

func (user ClientChannels) Take(other *ClientChannels) {
	user.Taken <- other
}

func (i *ClientChannels) GiveUp() {
	log.Println(i.UserID, "has given up!")
	if !(i.State >= GaveUp) {
		i.State = GaveUp
		if i.Opponent != i && i.Opponent.State != GaveUp {
			i.Opponent.State = OpponentGaveUp
		}
		i.HasGivenUp <- true
		close(i.Input)
	}
}

func (this *ClientChannels) Send(state *GameState) {
	this.LastThing = *state.Copy()
	if this.State != Disconnected {
		this.Output <- *state
	}
}
