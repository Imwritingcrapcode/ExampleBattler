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
	PlayingAs      int
	LastThing      GameState
	Clock          *Clock
	Input          chan string
	Output         chan GameState
	HasGivenUp     chan bool
	Time           chan bool
	TimeOutput     chan string
	KillConnection chan struct{}
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
